package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/ms-henglu/armstrong/resource/types"
	"github.com/ms-henglu/armstrong/swagger"
	"github.com/ms-henglu/armstrong/utils"
	"github.com/sirupsen/logrus"
	_ "golang.org/x/text/cases"
)

var pluralizeClient = pluralize.NewClient()

func NewAzapiDefinitionsFromSwagger(apiPath swagger.ApiPath) []types.AzapiDefinition {
	methodMap := make(map[string]bool)
	for _, method := range apiPath.Methods {
		methodMap[method] = true
	}

	def := types.AzapiDefinition{
		Id:                apiPath.Path,
		AzureResourceType: apiPath.ResourceType,
		ApiVersion:        apiPath.ApiVersion,
		BodyFormat:        types.BodyFormatHcl,
		AdditionalFields:  make(map[string]types.Value),
	}

	label := defaultLabel(apiPath.ResourceType)

	switch {
	case len(methodMap) == 1 && methodMap[http.MethodGet]:
		def.Kind = types.KindDataSource
		switch apiPath.ApiType {
		case swagger.ApiTypeList:
			def.ResourceName = "azapi_resource_list"
			parentId := utils.ScopeOfListAction(apiPath.Path)
			def.AdditionalFields["parent_id"] = types.NewStringLiteralValue(parentId)
			resourceName := def.AzureResourceType[strings.LastIndex(def.AzureResourceType, "/")+1:]
			scope := strings.Title(defaultLabel(utils.ResourceTypeOfResourceId(parentId)))
			def.Label = fmt.Sprintf("list%sBy%s", strings.Title(resourceName), scope)
		case swagger.ApiTypeResource:
			def.ResourceName = "azapi_resource"
			def.AdditionalFields["parent_id"] = types.NewStringLiteralValue(utils.ParentIdOfResourceId(apiPath.Path))
			def.AdditionalFields["name"] = types.NewStringLiteralValue(utils.LastSegment(apiPath.Path))
			def.Label = label
		case swagger.ApiTypeResourceAction, swagger.ApiTypeProviderAction:
			def.ResourceName = "azapi_resource_action"
			def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodGet)
			def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(ResourceIdFromActionPath(apiPath.Path))
			action := utils.ActionName(apiPath.Path)
			def.AdditionalFields["action"] = types.NewStringLiteralValue(action)
			def.Label = action
			if def.Label == "" {
				def.Label = utils.LastSegment(apiPath.Path)
			}
		}
	case len(methodMap) == 1 && methodMap[http.MethodPost]:
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodPost)
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(ResourceIdFromActionPath(apiPath.Path))
		action := utils.ActionName(apiPath.Path)
		def.AdditionalFields["action"] = types.NewStringLiteralValue(action)
		def.Label = action
		if def.Label == "" {
			def.Label = utils.LastSegment(apiPath.Path)
		}
		// defaults to resource, but if the request body is nil, it's a data source
		// TODO: check the swagger spec to see if it's a data source
		def.Kind = types.KindResource
		examplePath := apiPath.ExampleMap[http.MethodPost]
		if requestBody, err := RequestBodyFromExample(examplePath); err == nil {
			if requestBody == nil {
				def.Kind = types.KindDataSource
			}
			def.Body = requestBody
		} else {
			logrus.Warnf("failed to get request body from example, "+
				"this usually means that the `x-ms-examples` extension is not set correctly for %s in the swagger spec. %v", apiPath.Path, err)
		}

	case methodMap[http.MethodGet] && methodMap[http.MethodPut] && methodMap[http.MethodDelete]:
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource"
		def.AdditionalFields["parent_id"] = types.NewStringLiteralValue(utils.ParentIdOfResourceId(apiPath.Path))
		def.AdditionalFields["name"] = types.NewStringLiteralValue(utils.LastSegment(apiPath.Path))
		def.AdditionalFields["schema_validation_enabled"] = types.NewRawValue("false")
		def.Label = label
		examplePath := apiPath.ExampleMap[http.MethodPut]
		if requestBody, err := RequestBodyFromExample(examplePath); err == nil {
			def.Body = requestBody
			if requestBody != nil {
				if requestBodyMap, ok := def.Body.(map[string]interface{}); ok && requestBody != nil {
					if location := requestBodyMap["location"]; location != nil {
						def.AdditionalFields["location"] = types.NewStringLiteralValue(location.(string))
						delete(requestBodyMap, "location")
					}
					if name := requestBodyMap["name"]; name != nil {
						delete(requestBodyMap, "name")
					}
					delete(requestBodyMap, "id")
					def.Body = requestBodyMap
				}
			}
		} else {
			logrus.Warnf("failed to get request body from example, "+
				"this usually means that the `x-ms-examples` extension is not set correctly for %s in the swagger spec. %v", apiPath.Path, err)
		}

	case methodMap[http.MethodPut] && methodMap[http.MethodGet]:
		def.Kind = types.KindResource
		def.ResourceName = "azapi_update_resource"
		def.AdditionalFields["parent_id"] = types.NewStringLiteralValue(utils.ParentIdOfResourceId(apiPath.Path))
		def.AdditionalFields["name"] = types.NewStringLiteralValue(utils.LastSegment(apiPath.Path))
		def.Label = label
		examplePath := apiPath.ExampleMap[http.MethodPut]
		if requestBody, err := RequestBodyFromExample(examplePath); err == nil {
			def.Body = requestBody
		} else {
			logrus.Warnf("failed to get request body from example, "+
				"this usually means that the `x-ms-examples` extension is not set correctly for %s in the swagger spec. %v", apiPath.Path, err)
		}

	case methodMap[http.MethodPut]:
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(ResourceIdFromActionPath(apiPath.Path))
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodPut)
		action := utils.ActionName(apiPath.Path)
		def.AdditionalFields["action"] = types.NewStringLiteralValue(action)
		if action != "" {
			def.Label = action
		} else {
			def.Label = fmt.Sprintf("put_%s", label)
		}
		examplePath := apiPath.ExampleMap[http.MethodPut]
		if requestBody, err := RequestBodyFromExample(examplePath); err == nil {
			def.Body = requestBody
		} else {
			logrus.Warnf("failed to get request body from example, "+
				"this usually means that the `x-ms-examples` extension is not set correctly for %s in the swagger spec. %v", apiPath.Path, err)
		}
	}

	res := make([]types.AzapiDefinition, 0)
	if def.Label != "" {
		res = append(res, def)
	}

	if methodMap[http.MethodPatch] {
		def := types.AzapiDefinition{
			AzureResourceType: apiPath.ResourceType,
			ApiVersion:        apiPath.ApiVersion,
			BodyFormat:        types.BodyFormatHcl,
			AdditionalFields:  make(map[string]types.Value),
		}
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodPatch)
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(ResourceIdFromActionPath(apiPath.Path))
		action := utils.ActionName(apiPath.Path)
		def.AdditionalFields["action"] = types.NewStringLiteralValue(action)
		if action != "" {
			def.Label = action
		} else {
			def.Label = fmt.Sprintf("patch_%s", label)
		}
		examplePath := apiPath.ExampleMap[http.MethodPatch]
		if requestBody, err := RequestBodyFromExample(examplePath); err == nil {
			def.Body = requestBody
			if requestBody != nil {
				if requestBodyMap, ok := def.Body.(map[string]interface{}); ok && requestBody != nil {
					if location := requestBodyMap["location"]; location != nil {
						delete(requestBodyMap, "location")
					}
					if name := requestBodyMap["name"]; name != nil {
						delete(requestBodyMap, "name")
					}
					delete(requestBodyMap, "id")
					def.Body = requestBodyMap
				}
			}
		} else {
			logrus.Warnf("failed to get request body from example, "+
				"this usually means that the `x-ms-examples` extension is not set correctly for %s in the swagger spec. %v", apiPath.Path, err)
		}
		res = append(res, def)
	}

	return res
}

func ResourceIdFromActionPath(input string) string {
	id := strings.TrimPrefix(strings.TrimSuffix(input, "/"), "/")
	parts := strings.Split(id, "/")
	if len(parts)%2 == 0 {
		return input
	}
	return "/" + strings.Join(parts[:len(parts)-1], "/")
}

func defaultLabel(resourceType string) string {
	parts := strings.Split(resourceType, "/")
	label := "test"
	if len(parts) != 0 {
		label = parts[len(parts)-1]
		label = pluralizeClient.Singular(label)
	}
	return label
}

func RequestBodyFromExample(examplePath string) (interface{}, error) {
	data, err := os.ReadFile(examplePath)
	if err != nil {
		return nil, err
	}
	var example struct {
		Parameters map[string]interface{} `json:"parameters"`
	}
	err = json.Unmarshal(data, &example)
	if err != nil {
		return nil, err
	}

	var body interface{}
	for _, value := range example.Parameters {
		if bodyMap, ok := value.(map[string]interface{}); ok {
			body = bodyMap
		}
	}
	return body, nil
}
