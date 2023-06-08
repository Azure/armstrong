package coverage_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-openapi/jsonreference"
	openapispec "github.com/go-openapi/spec"
	"github.com/ms-henglu/armstrong/coverage"
)

func TestLoad(t *testing.T) {
	resourceId := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ex-resources/providers/Microsoft.Media/mediaServices/mediatest/transforms/transform1"
	swaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"

	apiPath, modelName, modelSwaggerPath, err := coverage.PathPatternFromId(resourceId, swaggerPath)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(*apiPath, *modelName, *modelSwaggerPath)

	model, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("model", string(out))
}

func TestExpand(t *testing.T) {
	modelName := "CustomPersistentDiskProperties"
	modelSwaggerPath := "/home/wangta/go/src/github.com/azure/azure-rest-api-specs/specification/appplatform/resource-manager/Microsoft.AppPlatform/stable/2022-12-01/appplatform.json"
	model, err := coverage.Expand(modelName, modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("model", string(out))
}

func TestExpandAll(t *testing.T) {
	index, err := coverage.GetIndex()
	if err != nil {
		t.Error(err)
	}

	refs := make([]*jsonreference.Ref, 0)
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
							resourceProvider = fmt.Sprintf("%s %s %s %s %s %s", resourceProvider, version, operationKind, resourceType, action, pathPatternStr)
							refs = append(refs, &ref)
						}
					}
					for pathPatternStr, ref := range operationInfo.OperationRefs {
						resourceProvider = fmt.Sprintf("%s %s %s %s %s", resourceProvider, version, operationKind, resourceType, pathPatternStr)
						refs = append(refs, &ref)
					}
				}
			}
		}
	}
	index = nil

	t.Logf("refs: %d", len(refs))

	refChans := make(chan *jsonreference.Ref)

	var waitGroup sync.WaitGroup
	waitGroup.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func(i int) {
			for ref := range refChans {
				t.Logf("%v ref: %v", i, ref.String())
				azureRepoUrl := "/home/wangta/go/src/github.com/azure/azure-rest-api-specs/specification/"
				swaggerPath := filepath.Join(azureRepoUrl, ref.GetURL().Path)
				operation, err := openapispec.ResolvePathItemWithBase(nil, openapispec.Ref{Ref: *ref}, &openapispec.ExpandOptions{RelativeBase: azureRepoUrl + "/" + strings.Split(ref.GetURL().Path, "/")[0]})
				if err != nil {
					t.Error(err)
				}

				var modelName string
				for _, param := range operation.Parameters {
					if param.In == "body" {
						var modelRelativePath string
						modelName, modelRelativePath = coverage.SchemaInfoFromRef(param.Schema.Ref)
						if modelRelativePath != "" {
							swaggerPath = filepath.Join(filepath.Dir(swaggerPath), modelRelativePath)
						}
					}
				}

				// post may have no model
				if operation.Put != nil && modelName == "" {
					panic("modelName is empty")
				}

				swaggerPath = strings.Replace(swaggerPath, "https:/", "https://", 1)

				model, err := coverage.Expand(modelName, swaggerPath)
				if err != nil {
					panic(err)
				}

				model = model
				model = nil
				//operation = nil
				//ref = nil
			}

			waitGroup.Done()
		}(i)
	}

	for _, ref := range refs {
		refChans <- ref
	}
}

func TestGR(t *testing.T) {
	c := make(chan int)
	var w sync.WaitGroup
	w.Add(5)

	for i := 1; i <= 5; i++ {
		go func(i int, ci <-chan int) {
			j := 1
			for v := range ci {
				time.Sleep(time.Millisecond)
				fmt.Printf("%d.%d got %d\n", i, j, v)
				j += 1
			}
			w.Done()
		}(i, c)
	}

	for i := 1; i <= 25; i++ {
		c <- i
	}
	close(c)
	w.Wait()
}
