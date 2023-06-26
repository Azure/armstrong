package coverage_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/jsonreference"
	openapispec "github.com/go-openapi/spec"
	"github.com/ms-henglu/armstrong/coverage"
)

func TestGetModelInfoFromSingleSwaggerFile(t *testing.T) {
	resourceId := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ex-resources/providers/Microsoft.Media/mediaServices/mediatest/transforms/transform1"
	swaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"

	apiPath, modelName, modelSwaggerPath, err := coverage.GetModelInfoFromSingleSwaggerFile(resourceId, swaggerPath)
	if err != nil {
		t.Errorf("get model info error:%+v", err)
	}

	expectedApiPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}"
	if *apiPath != expectedApiPath {
		t.Errorf("apiPath should be %s, but got %s", expectedApiPath, *apiPath)
	}

	expectedModelName := "Transform"
	if *modelName != expectedModelName {
		t.Errorf("modelName should be %s, but got %s", expectedModelName, *modelName)
	}

	expectedModelSwaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"
	if *modelSwaggerPath != expectedModelSwaggerPath {
		t.Errorf("modelSwaggerPath should be %s, but got %s", expectedModelSwaggerPath, *modelSwaggerPath)
	}
}

func TestExpand(t *testing.T) {
	modelName := "Transform"
	modelSwaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"
	model, err := coverage.Expand(modelName, modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("expanded model %s", string(out))
}

// try to expand all models
func TestExpandAll(t *testing.T) {
	azureRepoDir := os.Getenv("AZURE_REPO_DIR")
	if azureRepoDir == "" {
		t.Skip("AZURE_REPO_DIR is not set")
	}
	t.Logf("azure repo dir: %s", azureRepoDir)

	index, err := coverage.GetIndex()
	if err != nil {
		panic(err)
	}

	refMaps := make(map[string]*jsonreference.Ref, 0)
	for resourceProvider, versionRaw := range index.ResourceProviders {
		for version, methodRaw := range versionRaw {
			for operationKind, resourceTypeRaw := range methodRaw {
				if operationKind != "PUT" && operationKind != "POST" {
					continue
				}
				for resourceType, operationInfo := range resourceTypeRaw {
					if operationInfo == nil {
						continue
					}
					for action, operationRefs := range operationInfo.Actions {
						for pathPatternStr, ref := range operationRefs {
							if false {
								t.Logf("%s %s %s %s %s %s", resourceProvider, version, operationKind, resourceType, action, pathPatternStr)
							}
							refMaps[ref.String()] = &ref
						}
					}
					for pathPatternStr, ref := range operationInfo.OperationRefs {
						if false {
							t.Logf("%s %s %s %s %s", resourceProvider, version, operationKind, resourceType, pathPatternStr)
						}
						refMaps[ref.String()] = &ref
					}
				}
			}
		}
	}
	index = nil

	t.Logf("refs numbers: %d", len(refMaps))

	refChan := make(chan *jsonreference.Ref)

	var waitGroup sync.WaitGroup
	waitGroup.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			for ref := range refChan {
				t.Logf("%v ref: %v", i, ref.String())
				swaggerPath := filepath.Join(azureRepoDir, ref.GetURL().Path)
				operation, err := openapispec.ResolvePathItemWithBase(nil, openapispec.Ref{Ref: *ref}, &openapispec.ExpandOptions{RelativeBase: azureRepoDir + "/" + strings.Split(ref.GetURL().Path, "/")[0]})
				if err != nil {
					panic(err)
				}

				var modelName string
				for _, param := range operation.Parameters {
					if param.In == "body" {
						var modelRelativePath string
						modelName, modelRelativePath = coverage.SchemaNamePathFromRef(param.Schema.Ref)
						if modelRelativePath != "" {
							swaggerPath = filepath.Join(filepath.Dir(swaggerPath), modelRelativePath)
						}
					}
				}

				// post may have no model
				if operation.Put != nil && modelName == "" {
					panic("modelName is empty")
				}

				model, err := coverage.Expand(modelName, swaggerPath)
				if err != nil {
					panic(err)
				}

				model = model
				model = nil
				operation = nil
				ref = nil
			}

			waitGroup.Done()
		}(i)
	}

	for _, ref := range refMaps {
		refChan <- ref
	}
}
