package coverage

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	openapispec "github.com/go-openapi/spec"
	"github.com/magodo/azure-rest-api-index/azidx"
)

func getIndex() (*azidx.Index, error) {
	indexUrl := "https://raw.githubusercontent.com/teowa/azure-rest-api-index-file/main/index.json"
	resp, err := http.Get(indexUrl)
	if err != nil {
		return nil, fmt.Errorf("get index file (%v): %v", indexUrl, err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read index file: %v", err)
	}

	var index azidx.Index
	if err := json.Unmarshal(b, &index); err != nil {
		return nil, fmt.Errorf("unmarshal index file: %v", err)
	}
	return &index, nil
}

func PathPatternFromIdFromIndex(resourceId, apiVersion string) (*string, *string, *string, *string, error) {
	index, err := getIndex()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	log.Printf("[INFO] load index based commit: https://github.com/Azure/azure-rest-api-specs/tree/%s", index.Commit)

	resourceURL := fmt.Sprintf("https://management.azure.com%s?api-version=%s", resourceId, apiVersion)
	uRL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("parsing URL %s: %v", resourceURL, err)
	}
	ref, err := index.Lookup("PUT", *uRL)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	azureRepoUrl := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/"
	swaggerPath := filepath.Join(azureRepoUrl, ref.GetURL().Path)
	operation, err := openapispec.ResolvePathItemWithBase(nil, openapispec.Ref{Ref: *ref}, &openapispec.ExpandOptions{RelativeBase: azureRepoUrl + "/" + strings.Split(ref.GetURL().Path, "/")[0]})

	if err != nil {
		return nil, nil, nil, nil, err
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

	swaggerPath = strings.Replace(swaggerPath, "https:/", "https://", 1)

	return &apiPath, &modelName, &swaggerPath, &index.Commit, nil
}
