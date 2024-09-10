package coverage

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	openapispec "github.com/go-openapi/spec"
	"github.com/magodo/azure-rest-api-index/azidx"
	"github.com/sirupsen/logrus"
)

const (
	indexFileURL = "https://raw.githubusercontent.com/teowa/azure-rest-api-index-file/main/index.json.zip"
	azureRepoURL = "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/"
)

var indexCache *azidx.Index

func GetIndexFromLocalDir(swaggerRepo, indexFilePath string) (*azidx.Index, error) {
	if indexCache != nil {
		return indexCache, nil
	}

	if indexFilePath != "" {
		if _, err := os.Stat(indexFilePath); err == nil {
			byteValue, _ := os.ReadFile(indexFilePath)

			var index azidx.Index
			if err := json.Unmarshal(byteValue, &index); err != nil {
				return nil, fmt.Errorf("unmarshal index file: %+v", err)
			}
			indexCache = &index

			logrus.Infof("load index from cache file %s", indexFilePath)

			return indexCache, nil
		}
	}

	logrus.Infof("building index from from local swagger %s, it might take several minutes", swaggerRepo)
	index, err := azidx.BuildIndex(swaggerRepo, "")
	if err != nil {
		logrus.Error(fmt.Sprintf("failed to build index: %+v", err))
		return nil, err
	}
	logrus.Infof("index successfully built on commit %+v", index.Commit)

	indexCache = index

	if indexFilePath != "" {
		jsonBytes, err := json.Marshal(&index)
		if err != nil {
			logrus.Warningf("failed to marshal index: %+v", err)
			return index, nil
		}

		err = os.WriteFile(indexFilePath, jsonBytes, 0644)
		if err != nil {
			logrus.Warningf("failed to write index cache file %s: %+v", indexFilePath, err)
			return index, nil
		}

		logrus.Infof("index successfully saved to cache file %s", indexFilePath)
	}

	return index, nil
}

func GetIndex(indexFilePath string) (*azidx.Index, error) {
	if indexCache != nil {
		return indexCache, nil
	}

	if indexFilePath != "" {
		if _, err := os.Stat(indexFilePath); err == nil {
			byteValue, _ := os.ReadFile(indexFilePath)

			var index azidx.Index
			if err := json.Unmarshal(byteValue, &index); err != nil {
				return nil, fmt.Errorf("unmarshal index file: %+v", err)
			}
			indexCache = &index

			logrus.Infof("load index from cache file %s", indexFilePath)

			return indexCache, nil
		}
	}

	resp, err := http.Get(indexFileURL)
	if err != nil {
		return nil, fmt.Errorf("get index file from %v: %+v", indexFileURL, err)
	}

	logrus.Infof("downloading index file from %s", indexFileURL)

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

	logrus.Infof("load index based commit: https://github.com/Azure/azure-rest-api-specs/tree/%s", index.Commit)

	if indexFilePath != "" {
		jsonBytes, err := json.Marshal(&index)
		if err != nil {
			logrus.Warningf("failed to marshal index: %+v", err)
			return indexCache, nil
		}

		err = os.WriteFile(indexFilePath, jsonBytes, 0644)
		if err != nil {
			logrus.Warningf("failed to write index cache file %s: %+v", indexFilePath, err)
			return indexCache, nil
		}

		logrus.Infof("index successfully saved to cache file %s", indexFilePath)
	}

	return indexCache, nil
}

type SwaggerModel struct {
	ApiPath     string
	ModelName   string
	SwaggerPath string
	OperationID string
}

// GetModelInfoFromIndex will try to download online index from https://github.com/teowa/azure-rest-api-index-file, and get model info from it
// if the index is already downloaded as in {indexFilePath}, it will use the cached index
func GetModelInfoFromIndex(resourceId, apiVersion, method, indexFilePath string) (*SwaggerModel, error) {
	index, err := GetIndex(indexFilePath)
	if err != nil {
		return nil, err
	}

	resourceURL := fmt.Sprintf("https://management.azure.com%s?api-version=%s", resourceId, apiVersion)
	uRL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL %s: %+v", resourceURL, err)
	}
	ref, err := index.Lookup(method, *uRL)
	if err != nil {
		return nil, fmt.Errorf("lookup PUT URL %s in index: %+v", resourceURL, err)
	}

	model, err := GetModelInfoFromIndexRef(openapispec.Ref{Ref: *ref}, azureRepoURL)
	if err != nil {
		return nil, fmt.Errorf("get model %s: %+v", ref, err)
	}
	if model.ModelName == "" {
		return nil, fmt.Errorf("PUT model not found for %s", ref.String())
	}

	return model, nil
}

// GetModelInfoFromLocalIndex tries to build index from local swagger repo and get model info from it
func GetModelInfoFromLocalIndex(resourceId, apiVersion, method, swaggerRepo, indexCacheFile string) (*SwaggerModel, error) {
	swaggerRepo, err := filepath.Abs(swaggerRepo)
	if err != nil {
		return nil, fmt.Errorf("swagger repo path %q is invalid: %+v", swaggerRepo, err)
	}

	if _, err := os.Stat(swaggerRepo); os.IsNotExist(err) {
		return nil, fmt.Errorf("swagger repo path %q is invalid: path does not exist", swaggerRepo)
	}

	swaggerRepo = strings.TrimSuffix(swaggerRepo, "/")

	if !strings.HasSuffix(swaggerRepo, "specification") {
		return nil, fmt.Errorf("swagger repo path %q is invalid: must point to \"specification\", e.g., /home/projects/azure-rest-api-specs/specification", swaggerRepo)
	}

	swaggerRepo += "/"

	index, err := GetIndexFromLocalDir(swaggerRepo, indexCacheFile)
	if err != nil {
		return nil, fmt.Errorf("build index from local dir %s: %+v", swaggerRepo, err)
	}

	resourceURL := fmt.Sprintf("https://management.azure.com%s?api-version=%s", resourceId, apiVersion)
	uRL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, fmt.Errorf("parsing URL %s: %+v", resourceURL, err)
	}
	ref, err := index.Lookup(method, *uRL)
	if err != nil {
		return nil, fmt.Errorf("lookup PUT URL %s in index: %+v", resourceURL, err)
	}

	model, err := GetModelInfoFromIndexRef(openapispec.Ref{Ref: *ref}, swaggerRepo)
	if err != nil {
		return nil, fmt.Errorf("get model %s: %+v", ref, err)
	}
	if model.ModelName == "" {
		return nil, fmt.Errorf("PUT model not found for %s", ref.String())
	}

	return model, nil
}

func GetModelInfoFromIndexRef(ref openapispec.Ref, swaggerRepo string) (*SwaggerModel, error) {
	_, swaggerPath := SchemaNamePathFromRef(swaggerRepo, ref)

	seperator := "/"
	// in windows the ref might use backslashes
	if strings.Contains(ref.GetURL().Path, string(os.PathSeparator)) {
		seperator = string(os.PathSeparator)
	}

	relativeBase := swaggerRepo + strings.Split(ref.GetURL().Path, seperator)[0]
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

func MockResourceIDFromType(azapiResourceType string) (string, string) {
	const (
		managementGroupId = "/providers/Microsoft.Management/managementGroups/group1"
		subscritionSeg    = "/subscriptions/00000000-0000-0000-0000-000000000000"
		resourceGroupSeg  = "resourceGroups/rg"
	)
	resourceType := strings.Split(azapiResourceType, "@")[0]
	apiVersion := strings.Split(azapiResourceType, "@")[1]
	resourceProvider := strings.Split(resourceType, "/")[0]
	rTypes := strings.Split(resourceType, "/")[1:]
	typeIds := strings.Join(rTypes, "/xxx/") + "/xxx"

	return fmt.Sprintf("%s/%s/providers/%s/%s", subscritionSeg, resourceGroupSeg, resourceProvider, typeIds), apiVersion
}

func GetModelInfoFromIndexWithType(azapiResourceType, method, indexCacheFile string) (*SwaggerModel, error) {
	resourceId, apiVersion := MockResourceIDFromType(azapiResourceType)
	return GetModelInfoFromIndex(resourceId, apiVersion, method, indexCacheFile)
}
