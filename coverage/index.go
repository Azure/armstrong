package coverage

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	openapispec "github.com/go-openapi/spec"
	"github.com/magodo/azure-rest-api-index/azidx"
)

const (
	indexFileURL = "https://raw.githubusercontent.com/teowa/azure-rest-api-index-file/main/index.json.zip"
	azureRepoURL = "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/"
)

var indexCache *azidx.Index

func GetIndex() (*azidx.Index, error) {
	if indexCache != nil {
		return indexCache, nil
	}

	resp, err := http.Get(indexFileURL)
	if err != nil {
		return nil, fmt.Errorf("get index file from %v: %+v", indexFileURL, err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("download index file zip: %+v", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, fmt.Errorf("read index file zip: %+v", err)
	}

	var unzippedIndexBytes []byte
	for _, zipFile := range zipReader.File {
		if strings.EqualFold(zipFile.Name, "index.json") {
			unzippedIndexBytes, err = readZipFile(zipFile)
			if err != nil {
				return nil, fmt.Errorf("unzip index file: %+v", err)
			}
			break
		}
	}

	if len(unzippedIndexBytes) == 0 {
		return nil, fmt.Errorf("index file not found in zip")
	}

	var index azidx.Index
	if err := json.Unmarshal(unzippedIndexBytes, &index); err != nil {
		return nil, fmt.Errorf("unmarshal index file: %+v", err)
	}
	indexCache = &index

	log.Printf("[INFO] load index based commit: https://github.com/Azure/azure-rest-api-specs/tree/%s", index.Commit)
	return indexCache, nil
}

type SwaggerModel struct {
	ApiPath     string
	ModelName   string
	SwaggerPath string
}

func GetModelInfoFromIndex(resourceId, apiVersion string) (*SwaggerModel, error) {
	index, err := GetIndex()
	if err != nil {
		return nil, err
	}

	resourceURL := fmt.Sprintf("https://management.azure.com%s?api-version=%s", resourceId, apiVersion)
	uRL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL %s: %+v", resourceURL, err)
	}
	ref, err := index.Lookup("PUT", *uRL)
	if err != nil {
		return nil, err
	}

	model, err := GetModelInfoFromIndexRef(openapispec.Ref{Ref: *ref}, azureRepoURL)
	if err != nil {
		return nil, err
	}
	if model.ModelName == "" {
		return nil, fmt.Errorf("PUT model not found for %s", ref.String())
	}

	return model, nil
}

func GetModelInfoFromIndexRef(ref openapispec.Ref, swaggerRepo string) (*SwaggerModel, error) {
	_, swaggerPath := SchemaNamePathFromRef(swaggerRepo, ref)

	relativeBase := swaggerRepo + strings.Split(ref.GetURL().Path, "/")[0]
	operation, err := openapispec.ResolvePathItemWithBase(nil, ref, &openapispec.ExpandOptions{RelativeBase: relativeBase})

	if err != nil {
		return nil, err
	}

	pointerTokens := ref.GetPointer().DecodedTokens()
	apiPath := pointerTokens[1]

	var modelName string
	for _, param := range operation.Parameters {
		paramRef := param.Ref
		if paramRef.String() != "" {
			refParam, err := openapispec.ResolveParameterWithBase(nil, param.Ref, &openapispec.ExpandOptions{RelativeBase: swaggerPath})
			if err != nil {
				return nil, fmt.Errorf("resolve param ref %q: %+v", param.Ref.String(), err)
			}

			// Update the param
			param = *refParam
		}
		if param.In == "body" {
			if paramRef.String() != "" {
				modelName, swaggerPath = SchemaNamePathFromRef(swaggerPath, paramRef)
			}

			if param.Schema.Ref.String() != "" {
				modelName, swaggerPath = SchemaNamePathFromRef(swaggerPath, param.Schema.Ref)
			}
			break
		}
	}
	return &SwaggerModel{
		ApiPath:     apiPath,
		ModelName:   modelName,
		SwaggerPath: swaggerPath,
	}, nil
}

func readZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}
