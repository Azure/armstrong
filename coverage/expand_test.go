package coverage_test

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/go-openapi/jsonreference"
	openapispec "github.com/go-openapi/spec"
	"github.com/ms-henglu/armstrong/coverage"
)

func TestExpand_MediaTransform(t *testing.T) {
	modelName := "Transform"
	modelSwaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"
	model, err := coverage.Expand(modelName, modelSwaggerPath)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("expanded model %s", string(out))
}

// try to expand all PUT and POST models
func TestExpandAll(t *testing.T) {
	azureRepoDir := os.Getenv("AZURE_REST_REPO_DIR")
	if azureRepoDir == "" {
		t.Skip("AZURE_REST_REPO_DIR is not set")
	}
	if !strings.HasSuffix(azureRepoDir, "specification/") {
		t.Fatalf("AZURE_REST_REPO_DIR must specify the specification folder, e.g., AZURE_REST_REPO_DIR=\"/home/test/go/src/github.com/azure/azure-rest-api-specs/specification/\"")
	}

	t.Logf("azure repo dir: %s", azureRepoDir)
	index, err := coverage.GetIndex()
	if err != nil {
		t.Fatal(err)
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
						for pathPattern, ref := range operationRefs {
							t.Logf("%s %s %s %s %s %s", resourceProvider, version, operationKind, resourceType, action, pathPattern)
							refMaps[ref.String()] = &ref
						}
					}
					for pathPattern, ref := range operationInfo.OperationRefs {
						t.Logf("%s %s %s %s %s", resourceProvider, version, operationKind, resourceType, pathPattern)
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

				model, err := coverage.GetModelInfoFromIndexRef(openapispec.Ref{Ref: *ref}, azureRepoDir)
				if err != nil {
					panic(fmt.Errorf("get model info from index ref %s: %+v", ref.String(), err))
				}

				_, err = coverage.Expand(model.ModelName, model.SwaggerPath)
				if err != nil {
					panic(fmt.Errorf("process %s, expand %s from %s: %+v", ref.String(), model.ModelName, model.SwaggerPath, err))
				}

				// clean up
				model = nil
				ref = nil
			}

			waitGroup.Done()
		}(i)
	}

	for _, ref := range refMaps {
		refChan <- ref
	}
}
