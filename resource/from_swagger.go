package resource

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"net/http"
	"os"
	"strings"

	"github.com/azure/armstrong/resource/types"
	"github.com/azure/armstrong/swagger"
	"github.com/azure/armstrong/utils"
	pluralize "github.com/gertd/go-pluralize"
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
	caser := cases.Title(language.Und, cases.NoLower)
	res := make([]types.AzapiDefinition, 0)
	coveredMethodMap := make(map[string]bool)
	switch {
	case methodMap[http.MethodGet] && methodMap[http.MethodPut] && methodMap[http.MethodDelete]:
		coveredMethodMap[http.MethodGet] = true
		coveredMethodMap[http.MethodPut] = true
		coveredMethodMap[http.MethodDelete] = true
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s, %s, %s", apiPath.OperationIdMap[http.MethodPut], apiPath.OperationIdMap[http.MethodGet], apiPath.OperationIdMap[http.MethodDelete]),
			fmt.Sprintf("PUT GET DELETE %s", apiPath.Path),
		}
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
		if def.Label != "" {
			res = append(res, def)
		}
	case methodMap[http.MethodPut] && methodMap[http.MethodGet]:
		coveredMethodMap[http.MethodGet] = true
		coveredMethodMap[http.MethodPut] = true
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodPut]),
			fmt.Sprintf("PUT %s", apiPath.Path),
		}
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(apiPath.Path)
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodPut)
		examplePath := apiPath.ExampleMap[http.MethodPut]
		if requestBody, err := RequestBodyFromExample(examplePath); err == nil {
			def.Body = requestBody
		} else {
			logrus.Warnf("failed to get request body from example, "+
				"this usually means that the `x-ms-examples` extension is not set correctly for %s in the swagger spec. %v", apiPath.Path, err)
		}
		def.Label = fmt.Sprintf("put_%s", label)
		res = append(res, def)

		def = types.AzapiDefinition{
			Id:                apiPath.Path,
			AzureResourceType: apiPath.ResourceType,
			ApiVersion:        apiPath.ApiVersion,
			BodyFormat:        types.BodyFormatHcl,
			AdditionalFields:  make(map[string]types.Value),
		}
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodGet]),
			fmt.Sprintf("GET %s", apiPath.Path),
		}
		def.Kind = types.KindDataSource
		def.ResourceName = "azapi_resource"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(apiPath.Path)
		def.AdditionalFields["depends_on"] = types.NewRawValue(fmt.Sprintf("[ azapi_resource_action.put_%s ]", label))
		def.Label = label
		res = append(res, def)
	case methodMap[http.MethodPut]:
		coveredMethodMap[http.MethodPut] = true
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodPut]),
			fmt.Sprintf("PUT %s", apiPath.Path),
		}
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
		if def.Label != "" {
			res = append(res, def)
		}
	case len(methodMap) == 1 && methodMap[http.MethodGet]:
		coveredMethodMap[http.MethodGet] = true
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodGet]),
			fmt.Sprintf("%s %s", http.MethodGet, apiPath.Path),
		}
		def.Kind = types.KindDataSource
		switch apiPath.ApiType {
		case swagger.ApiTypeList:
			def.ResourceName = "azapi_resource_list"
			parentId := utils.ScopeOfListAction(apiPath.Path)
			def.AdditionalFields["parent_id"] = types.NewStringLiteralValue(parentId)
			resourceName := def.AzureResourceType[strings.LastIndex(def.AzureResourceType, "/")+1:]
			scope := caser.String(defaultLabel(utils.ResourceTypeOfResourceId(parentId)))
			def.Label = fmt.Sprintf("list%sBy%s", caser.String(resourceName), scope)
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
		if def.Label != "" {
			res = append(res, def)
		}
	case len(methodMap) == 1 && methodMap[http.MethodPost]:
		coveredMethodMap[http.MethodPost] = true
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodPost]),
			fmt.Sprintf("%s %s", http.MethodPost, apiPath.Path),
		}
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
			def.Body = requestBody
		} else {
			logrus.Warnf("failed to get request body from example, "+
				"this usually means that the `x-ms-examples` extension is not set correctly for %s in the swagger spec. %v", apiPath.Path, err)
		}
		if def.Label != "" {
			res = append(res, def)
		}
	case methodMap[http.MethodGet] && methodMap[http.MethodDelete]:
		coveredMethodMap[http.MethodGet] = true
		coveredMethodMap[http.MethodDelete] = true
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodGet]),
			fmt.Sprintf("GET %s", apiPath.Path),
		}
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(apiPath.Path)
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodGet)
		def.Label = fmt.Sprintf("get_%s", label)
		res = append(res, def)

		def = types.AzapiDefinition{
			Id:                apiPath.Path,
			AzureResourceType: apiPath.ResourceType,
			ApiVersion:        apiPath.ApiVersion,
			BodyFormat:        types.BodyFormatHcl,
			AdditionalFields:  make(map[string]types.Value),
		}
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodDelete]),
			fmt.Sprintf("DELETE %s", apiPath.Path),
		}
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(apiPath.Path)
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodDelete)
		def.AdditionalFields["depends_on"] = types.NewRawValue(fmt.Sprintf("[ azapi_resource_action.get_%s ]", label))
		def.Label = fmt.Sprintf("delete_%s", label)
		res = append(res, def)
	case methodMap[http.MethodGet] && methodMap[http.MethodPatch]:
		coveredMethodMap[http.MethodGet] = true
		coveredMethodMap[http.MethodPatch] = true

		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodPatch]),
			fmt.Sprintf("PATCH %s", apiPath.Path),
		}
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(apiPath.Path)
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodPatch)
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
		def.Label = fmt.Sprintf("patch_%s", label)
		res = append(res, def)

		def = types.AzapiDefinition{
			Id:                apiPath.Path,
			AzureResourceType: apiPath.ResourceType,
			ApiVersion:        apiPath.ApiVersion,
			BodyFormat:        types.BodyFormatHcl,
			AdditionalFields:  make(map[string]types.Value),
		}
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodGet]),
			fmt.Sprintf("GET %s", apiPath.Path),
		}
		def.Kind = types.KindDataSource
		def.ResourceName = "azapi_resource"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(apiPath.Path)
		def.AdditionalFields["depends_on"] = types.NewRawValue(fmt.Sprintf("[ azapi_resource_action.patch_%s ]", label))
		def.Label = label
		res = append(res, def)
	}

	if methodMap[http.MethodPatch] && !coveredMethodMap[http.MethodPatch] {
		coveredMethodMap[http.MethodPatch] = true
		def := types.AzapiDefinition{
			AzureResourceType: apiPath.ResourceType,
			ApiVersion:        apiPath.ApiVersion,
			BodyFormat:        types.BodyFormatHcl,
			AdditionalFields:  make(map[string]types.Value),
		}
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodPatch]),
			fmt.Sprintf("PATCH %s", apiPath.Path),
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

	if methodMap[http.MethodHead] {
		coveredMethodMap[http.MethodHead] = true
		def := types.AzapiDefinition{
			AzureResourceType: apiPath.ResourceType,
			ApiVersion:        apiPath.ApiVersion,
			BodyFormat:        types.BodyFormatHcl,
			AdditionalFields:  make(map[string]types.Value),
		}
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodHead]),
			fmt.Sprintf("HEAD %s", apiPath.Path),
		}
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodHead)
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(ResourceIdFromActionPath(apiPath.Path))
		action := utils.ActionName(apiPath.Path)
		def.AdditionalFields["action"] = types.NewStringLiteralValue(action)
		if action != "" {
			def.Label = action
		} else {
			def.Label = fmt.Sprintf("head_%s", label)
		}
		res = append(res, def)
	}

	if methodMap[http.MethodDelete] && !coveredMethodMap[http.MethodDelete] {
		coveredMethodMap[http.MethodDelete] = true
		def.LeadingComments = []string{
			fmt.Sprintf("OperationId: %s", apiPath.OperationIdMap[http.MethodDelete]),
			fmt.Sprintf("DELETE %s", apiPath.Path),
		}
		def.Kind = types.KindResource
		def.ResourceName = "azapi_resource_action"
		def.AdditionalFields["resource_id"] = types.NewStringLiteralValue(apiPath.Path)
		def.AdditionalFields["method"] = types.NewStringLiteralValue(http.MethodDelete)
		def.Label = fmt.Sprintf("delete_%s", label)
		res = append(res, def)
	}

	notCoveredMethods := make([]string, 0)
	for method := range methodMap {
		if !coveredMethodMap[method] {
			notCoveredMethods = append(notCoveredMethods, method)
		}
	}
	if len(notCoveredMethods) != 0 {
		// TODO: GET and POST on a collection URL are not supported
		logrus.Errorf("there are methods not covered: %v for API path %s", notCoveredMethods, apiPath.Path)
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
