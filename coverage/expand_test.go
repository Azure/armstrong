package coverage_test

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

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
							t.Logf("%s %s %s %s %s %s", resourceProvider, version, operationKind, resourceType, action, pathPatternStr)
							refs = append(refs, &ref)
						}
					}
					for pathPatternStr, ref := range operationInfo.OperationRefs {
						t.Logf("%s %s %s %s %s", resourceProvider, version, operationKind, resourceType, pathPatternStr)
						refs = append(refs, &ref)
					}
				}
			}
		}
	}

	t.Logf("refs: %d", len(refs))
	t.Skip()

	for _, ref := range refs {
		azureRepoUrl := "/Users/wangtao/go/src/github.com/Azure/azure-rest-api-specs/specification/"
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

		if modelName == "" {
			panic("modelName is empty")
		}

		swaggerPath = strings.Replace(swaggerPath, "https:/", "https://", 1)

		_, err = coverage.Expand(modelName, swaggerPath)
		if err != nil {
			panic(err)
		}

	}

}
