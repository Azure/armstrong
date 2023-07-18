package coverage_test

import (
	"encoding/json"
	"fmt"
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

func TestExpand_MediaTranform(t *testing.T) {
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
	// e.g., AZURE_REST_REPO_DIR="/home/test/go/src/github.com/azure/azure-rest-api-specs/specification/"
	azureRepoDir := os.Getenv("AZURE_REST_REPO_DIR")
	if azureRepoDir == "" {
		t.Skip("AZURE_REST_REPO_DIR is not set")
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
				swaggerPath := filepath.Join(azureRepoDir, ref.GetURL().Path)
				operation, err := openapispec.ResolvePathItemWithBase(nil, openapispec.Ref{Ref: *ref}, &openapispec.ExpandOptions{RelativeBase: azureRepoDir + "/" + strings.Split(ref.GetURL().Path, "/")[0]})
				if err != nil {
					panic(fmt.Errorf("resolve operation %q from %s: %v", ref.String(), swaggerPath, err))
				}

				var modelName string
				for _, param := range operation.Parameters {
					paramRef := param.Ref
					if paramRef.String() != "" {
						refParam, err := openapispec.ResolveParameterWithBase(nil, param.Ref, &openapispec.ExpandOptions{RelativeBase: swaggerPath})
						if err != nil {
							panic(fmt.Errorf("resolve parameter %q from %s: %v", param.Ref.String(), swaggerPath, err))
						}

						// Update the param
						param = *refParam
					}
					if param.In == "body" {
						if paramRef.String() != "" {
							_, swaggerPath = coverage.SchemaNamePathFromRef(swaggerPath, paramRef)
						}

						modelName, swaggerPath = coverage.SchemaNamePathFromRef(swaggerPath, param.Schema.Ref)
						break
					}
				}

				// post may have no model
				if operation.Put != nil && modelName == "" {
					panic(fmt.Errorf("resolve %s from %s: modelName is empty", ref.String(), swaggerPath))
				}

				_, err = coverage.Expand(modelName, swaggerPath)
				if err != nil {
					panic(fmt.Errorf("expand %s from %s: %+v", modelName, swaggerPath, err))
				}

				// clean up
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
