package coverage

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	openapispec "github.com/go-openapi/spec"
	"github.com/magodo/azure-rest-api-index/azidx"
)

func getIndex(azureRepoDir string, refreshIndex bool) (*azidx.Index, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	indexFilePath := filepath.Join(cacheDir, "armstrong", "index.json")
	if _, err := os.Stat(indexFilePath); err != nil || refreshIndex {
		log.Printf("[INFO]azure-rest-api-specs root dir: %s\n", azureRepoDir)
		index, err := azidx.BuildIndex(azureRepoDir, "")
		if err != nil {
			return nil, err
		}

		b, err := json.MarshalIndent(index, "", "  ")
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(indexFilePath, b, 0644)
		if err != nil {
			return nil, err
		}
	}

	b, err := os.ReadFile(indexFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading index file %s: %v", indexFilePath, err)
	}
	var index azidx.Index
	if err := json.Unmarshal(b, &index); err != nil {
		return nil, fmt.Errorf("unmarshal index file: %v", err)
	}
	return &index, nil
}

func PathPatternFromIdFromIndex(resourceId, apiVersion, azureRepoDir string, refreshIndex bool) (*string, *string, *string, error) {
	index, err := getIndex(azureRepoDir, refreshIndex)
	if err != nil {
		return nil, nil, nil, err
	}

	resourceURL := fmt.Sprintf("https://management.azure.com%s?api-version=%s", resourceId, apiVersion)
	uRL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing URL %s: %v", resourceURL, err)
	}
	ref, err := index.Lookup("PUT", *uRL)
	if err != nil {
		return nil, nil, nil, err
	}

	swaggerPath := filepath.Join(azureRepoDir, ref.GetURL().Path)
	operation, err := openapispec.ResolvePathItemWithBase(nil, openapispec.Ref{Ref: *ref}, &openapispec.ExpandOptions{RelativeBase: azureRepoDir + "/" + strings.Split(ref.GetURL().Path, "/")[0]})

	if err != nil {
		return nil, nil, nil, err
	}

	pointerTokens := ref.GetPointer().DecodedTokens()
	apiPath := pointerTokens[1]

	var modelName string
	for _, param := range operation.Parameters {
		if param.In == "body" {
			var modelRelativePath string
			modelName, modelRelativePath = SchemaInfoFromRef(param.Schema.Ref)
			if modelRelativePath != "" {
				swaggerPath = filepath.Join(filepath.Dir(swaggerPath), modelRelativePath)
			}
		}
	}

	return &apiPath, &modelName, &swaggerPath, nil
}
