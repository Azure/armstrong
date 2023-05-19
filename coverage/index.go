package coverage

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"

	"github.com/go-openapi/loads"
	"github.com/magodo/azure-rest-api-index/azidx"
)

func PathPatternFromIdFromIndex(resourceId, apiVersion string) (*string, *string, *string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
		return nil, nil, nil, err
	}

	baseDir := filepath.Join(cacheDir, "armstrong")
	azureRepo, err := NewAzureRepo(baseDir)
	if err != nil {
		return nil, nil, nil, err
	}

	log.Println(azureRepo.SpecRootDir)
	index, err := azidx.BuildIndex(azureRepo.SpecRootDir, "")
	if err != nil {
		return nil, nil, nil, err
	}
	//b, err := json.MarshalIndent(index, "", "  ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//outputFile := filepath.Join(cacheDir, "armstrong", "index.json")
	//err = os.WriteFile(outputFile, b, 0644)
	//if err != nil {
	//	return nil, nil, nil, err
	//}
	//b, err := os.ReadFile(outputFile)
	//if err != nil {
	//	return fmt.Errorf("reading index file %s: %v", outputFile, err)
	//}
	//var index azidx.Index
	//if err := json.Unmarshal(b, &index); err != nil {
	//	return fmt.Errorf("unmarshal index file: %v", err)
	//}
	resourceURL := fmt.Sprintf("https://management.azure.com%s?api-version=%s", resourceId, apiVersion)
	uRL, err := url.Parse(resourceURL)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parsing URL %s: %v", resourceURL, err)
	}
	ref, err := index.Lookup("PUT", *uRL)
	if err != nil {
		return nil, nil, nil, err
	}

	swaggerPath := ref.GetURL().Path
	doc, err := loads.JSONSpec(swaggerPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("loading swagger spec: %+v", err)
	}

	spec := doc.Spec()

	pointerTokens := ref.GetPointer().DecodedTokens()
	apiPath := pointerTokens[1]
	path := spec.Paths.Paths[apiPath]
	operation := path.Put
	var modelName string
	for _, param := range operation.Parameters {
		if param.In == "body" {
			var modelRelativePath string
			modelName, modelRelativePath = SchemaInfoFromRef(param.Schema.Ref)
			if modelRelativePath != "" {
				//fmt.Println("modelRelativePath", modelRelativePath)
				swaggerPath = filepath.Join(filepath.Dir(swaggerPath), modelRelativePath)
			}
		}
	}

	return &apiPath, &modelName, &swaggerPath, nil
}
