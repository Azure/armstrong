package swagger

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"
	"github.com/ms-henglu/armstrong/utils"
	"golang.org/x/exp/slices"
)

// Load loads the swagger spec from the given path
func Load(swaggerPath string) ([]ApiPath, error) {
	swaggerSpec, err := loads.JSONSpec(swaggerPath)
	if err != nil {
		return nil, err
	}

	paths := swaggerSpec.Spec().Paths
	if paths == nil {
		return nil, fmt.Errorf("paths is nil, swagger path: %v", swaggerPath)
	}

	apiVersion := swaggerSpec.Spec().Info.Version

	apiPaths := make([]ApiPath, 0)
	for pathKey, pathItem := range paths.Paths {
		apiPath := ApiPath{
			Path:         pathKey,
			ApiVersion:   apiVersion,
			ResourceType: utils.ResourceTypeOfResourceId(pathKey),
			ApiType:      ApiTypeUnknown,
		}

		operationMap := make(map[string]spec.Operation)
		if pathItem.Get != nil {
			operationMap[http.MethodGet] = *pathItem.Get
		}
		if pathItem.Put != nil {
			operationMap[http.MethodPut] = *pathItem.Put
		}
		if pathItem.Post != nil {
			operationMap[http.MethodPost] = *pathItem.Post
		}
		if pathItem.Delete != nil {
			operationMap[http.MethodDelete] = *pathItem.Delete
		}
		if pathItem.Patch != nil {
			operationMap[http.MethodPatch] = *pathItem.Patch
		}
		if pathItem.Head != nil {
			operationMap[http.MethodHead] = *pathItem.Head
		}

		methods := make([]string, 0)
		exampleMap := make(map[string]string)
		operationIdMap := make(map[string]string)
		for method, operation := range operationMap {
			methods = append(methods, method)
			// the example is in the Extensions["x-ms-examples"], here's an example:
			// "x-ms-examples": {
			//   "Get lists of an automation account": {
			//     "$ref": "./examples/listAutomationAccountKeys.json"
			//   }
			// },
			exampleList := make([]string, 0)
			if operation.Extensions != nil && operation.Extensions["x-ms-examples"] != nil {
				if examples, ok := operation.Extensions["x-ms-examples"].(map[string]interface{}); ok && examples != nil {
					for _, exampleItem := range examples {
						if exampleItemMap, ok := exampleItem.(map[string]interface{}); ok && exampleItemMap != nil {
							exampleRef := exampleItemMap["$ref"]
							if exampleRef == nil {
								continue
							}
							exampleList = append(exampleList, path.Clean(path.Join(filepath.ToSlash(swaggerPath), "..", exampleRef.(string))))
						}
					}
				}
			}
			sort.Strings(exampleList)
			if len(exampleList) > 0 {
				exampleMap[method] = exampleList[0]
			}
			operationIdMap[method] = operation.ID
		}

		sort.Strings(methods)
		apiPath.Methods = methods
		apiPath.ExampleMap = exampleMap
		apiPath.OperationIdMap = operationIdMap
		apiPaths = append(apiPaths, apiPath)
	}

	// check the operation whether it's a list operation
	// for example:
	// /subscriptions/{subscriptionId}/providers/Microsoft.Automation/automationAccounts
	// /.../providers/Microsoft.Automation/automationAccounts/{automationAccountName}/modules
	resourceTypes := make(map[string]bool)
	// add a special resource type "operations"
	resourceTypes["operations"] = true
	for _, apiPath := range apiPaths {
		if lastSegment := utils.LastSegment(apiPath.ResourceType); lastSegment != "" {
			resourceTypes[lastSegment] = true
		}
	}
	for i, apiPath := range apiPaths {
		if len(apiPath.Methods) != 1 || apiPath.Methods[0] != http.MethodGet {
			continue
		}
		lastSegment := utils.LastSegment(apiPath.Path)
		if lastSegment != "" && resourceTypes[lastSegment] {
			apiPaths[i].ResourceType += "/" + lastSegment
			apiPaths[i].ApiType = ApiTypeList
		}
	}

	// decide the api type
	for i, apiPath := range apiPaths {
		if apiPath.ApiType != ApiTypeUnknown {
			continue
		}
		switch {
		case !utils.IsAction(apiPath.Path):
			apiPaths[i].ApiType = ApiTypeResource
		case strings.Contains(apiPath.ResourceType, "/"):
			apiPaths[i].ApiType = ApiTypeResourceAction
		default:
			apiPaths[i].ApiType = ApiTypeProviderAction
		}
	}

	slices.SortFunc(apiPaths, func(i, j ApiPath) int {
		return strings.Compare(i.Path, j.Path)
	})
	return apiPaths, nil
}
