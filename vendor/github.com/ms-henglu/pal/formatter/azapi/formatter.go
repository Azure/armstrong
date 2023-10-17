package azapi

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/ms-henglu/pal/formatter"
	"github.com/ms-henglu/pal/types"
)

var _ formatter.Formatter = &AzapiFormatter{}

var pluralizeClient = pluralize.NewClient()

type AzapiFormatter struct {
	existingResourceSet  map[string]bool
	labels               map[string]bool
	valueToKeyMap        map[string]string
	dataSourceActionUrls map[string]bool
}

func ignoreKeywords() []string {
	return []string{
		"operationresults",
		"asyncoperations",
		"operationstatuses",
		"operationsStatus",
		"operations",
	}
}

func (formatter *AzapiFormatter) Format(r types.RequestTrace) string {
	if formatter.existingResourceSet == nil {
		formatter.existingResourceSet = make(map[string]bool)
		formatter.labels = make(map[string]bool)
		formatter.valueToKeyMap = make(map[string]string)
		formatter.dataSourceActionUrls = make(map[string]bool)
	}

	if r.Host != "management.azure.com" {
		// ignore the request to other hosts
		return ""
	}

	if shouldIgnore(r.Url) {
		return ""
	}

	resourceId := GetId(r.Url)
	resourceType, apiVersion := parseResourceTypeApiVersion(r.Url)

	// ignore the request to fetch the resource group's resources
	if resourceType == "" && strings.Contains(resourceId, "/resources") {
		return ""
	}

	var requestBody interface{}
	_ = json.Unmarshal([]byte(r.Request.Body), &requestBody)
	var responseBody interface{}
	if r.Response != nil && r.Response.Body != "" {
		_ = json.Unmarshal([]byte(r.Response.Body), &responseBody)
	}
	def := AzapiDefinition{
		AzureResourceType: resourceType,
		ApiVersion:        apiVersion,
		Body:              requestBody,
		Output:            responseBody,
		ResourceId:        resourceId,
		Method:            r.Method,
		AdditionalFields:  make(map[string]Value),
	}

	switch r.Method {
	case "PUT":
		switch {
		case strings.EqualFold(def.AzureResourceType, "Microsoft.KeyVault/vaults/accessPolicies"):
			def = formatter.formatAsAzapiActionResource(def)
		case IsResourceAction(def.ResourceId):
			def = formatter.formatAsAzapiActionResource(def)
		case formatter.existingResourceSet[def.ResourceId]:
			def = formatter.formatAsAzapiUpdateResource(def)
		default:
			def = formatter.formatAsAzapiResource(def)
		}
	case "GET":
		if r.StatusCode == 200 && !formatter.existingResourceSet[def.ResourceId] {
			if IsResourceAction(def.ResourceId) {
				// resource action data source
				check := fmt.Sprintf("%s.%s", def.Method, def.ResourceId)
				if _, ok := formatter.dataSourceActionUrls[check]; !ok {
					def = formatter.formatAsAzapiActionDataSource(def)
					formatter.dataSourceActionUrls[check] = true
				} else {
					return ""
				}
			} else {
				// reading a resource which is created by the service
				check := fmt.Sprintf("%s.%s", def.Method, def.ResourceId)
				if _, ok := formatter.dataSourceActionUrls[check]; !ok {
					def = formatter.formatAsAzapiDataSource(def)
					formatter.dataSourceActionUrls[check] = true
					formatter.existingResourceSet[def.ResourceId] = true
				} else {
					return ""
				}
			}
		} else {
			return ""
		}
	case "POST":
		if r.Request.Body == "" {
			check := fmt.Sprintf("%s.%s", def.Method, def.ResourceId)
			if _, ok := formatter.dataSourceActionUrls[check]; !ok {
				def = formatter.formatAsAzapiActionDataSource(def)
				formatter.dataSourceActionUrls[check] = true
			} else {
				return ""
			}
		} else {
			def = formatter.formatAsAzapiActionResource(def)
		}
	case "PATCH":
		if !formatter.existingResourceSet[def.ResourceId] {
			log.Printf("[WARN] PATCH %s is not a created resource", def.ResourceId)
			return ""
		} else {
			def = formatter.formatAsAzapiActionResource(def)
		}
	case "DELETE":
		if !formatter.existingResourceSet[def.ResourceId] {
			log.Printf("[WARN] DELETE %s is not a created resource", def.ResourceId)
		}
		return ""
	default:
		return ""
	}

	def.Body = formatter.injectReference(def.Body)

	address := fmt.Sprintf("%s.%s.%s", def.Kind, def.ResourceName, def.Label)
	if def.Kind == "resource" {
		address = fmt.Sprintf("%s.%s", def.ResourceName, def.Label)
	}
	prefix := fmt.Sprintf(`jsondecode(%s.output)`, address)
	formatter.populateReference(prefix, def.Output)

	return def.String()
}

func shouldIgnore(url string) bool {
	resourceType := GetResourceType(url)
	if strings.EqualFold(resourceType, "Microsoft.ApiManagement/service/apis/operations") || strings.EqualFold(resourceType, "Microsoft.ApiManagement/service/apis/operations/tags") {
		return false
	}
	for _, v := range ignoreKeywords() {
		if strings.Contains(url, v) {
			return true
		}
	}
	return false
}

func (formatter *AzapiFormatter) formatAsAzapiResource(def AzapiDefinition) AzapiDefinition {
	formatter.existingResourceSet[def.ResourceId] = true
	def.Kind = "resource"
	def.ResourceName = "azapi_resource"
	def.Label = newUniqueLabel(def.ResourceName, defaultLabel(def.AzureResourceType), &formatter.labels)

	def.AdditionalFields["parent_id"] = formatter.tryAddressOrLiteral(GetParentId(def.ResourceId))
	def.AdditionalFields["name"] = NewLiteralValue(GetName(def.ResourceId))
	if def.Body != nil {
		if requestBody, ok := def.Body.(map[string]interface{}); ok && requestBody != nil {
			if location := requestBody["location"]; location != nil {
				def.AdditionalFields["location"] = NewLiteralValue(location.(string))
				delete(requestBody, "location")
			}
			if name := requestBody["name"]; name != nil {
				delete(requestBody, "name")
			}
			if strings.EqualFold(def.AzureResourceType, "Microsoft.ManagedIdentity/userAssignedIdentities") {
				delete(requestBody, "properties")
			}
			def.Body = requestBody
		}
	}
	formatter.valueToKeyMap[strings.ToUpper(def.ResourceId)] = fmt.Sprintf("%s.%s.id", def.ResourceName, def.Label)

	return def
}

func (formatter *AzapiFormatter) formatAsAzapiDataSource(def AzapiDefinition) AzapiDefinition {
	def.Kind = "data"
	def.ResourceName = "azapi_resource"
	def.Label = newUniqueLabel(def.ResourceName, defaultLabel(def.AzureResourceType), &formatter.labels)

	parentId := GetParentId(def.ResourceId)
	if parentAddress := formatter.valueToKeyMap[strings.ToUpper(parentId)]; parentAddress != "" {
		def.AdditionalFields["parent_id"] = NewReferenceValue(parentAddress)
		def.AdditionalFields["name"] = NewLiteralValue(GetName(def.ResourceId))
	} else {
		def.AdditionalFields["resource_id"] = NewLiteralValue(def.ResourceId)
	}
	def.AdditionalFields["response_export_values"] = NewRawValue(`["*"]`)
	formatter.valueToKeyMap[strings.ToUpper(def.ResourceId)] = fmt.Sprintf("data.%s.%s.id", def.ResourceName, def.Label)

	return def
}

func (formatter *AzapiFormatter) formatAsAzapiActionResource(def AzapiDefinition) AzapiDefinition {
	def.Kind = "resource"
	def.ResourceName = "azapi_resource_action"

	action := ""
	if parts := strings.Split(def.ResourceId, "/"); len(parts)%2 == 0 {
		action = parts[len(parts)-1]
		def.ResourceId = strings.Join(parts[:len(parts)-1], "/")
		def.AdditionalFields["action"] = NewLiteralValue(action)
	}

	label := action
	if label == "" {
		label = fmt.Sprintf("%s_%s", strings.ToLower(def.Method), defaultLabel(def.AzureResourceType))
	}
	def.Label = newUniqueLabel(def.ResourceName, label, &formatter.labels)

	def.AdditionalFields["resource_id"] = formatter.tryAddressOrLiteral(def.ResourceId)
	if def.Method != "POST" {
		def.AdditionalFields["method"] = NewLiteralValue(def.Method)
	}

	return def
}

func (formatter *AzapiFormatter) formatAsAzapiActionDataSource(def AzapiDefinition) AzapiDefinition {
	def.Kind = "data"
	def.ResourceName = "azapi_resource_action"

	action := ""
	if parts := strings.Split(def.ResourceId, "/"); len(parts)%2 == 0 {
		action = parts[len(parts)-1]
		def.ResourceId = strings.Join(parts[:len(parts)-1], "/")
		def.AdditionalFields["action"] = NewLiteralValue(action)
	}

	label := action
	if label == "" {
		label = fmt.Sprintf("%s_%s", strings.ToLower(def.Method), defaultLabel(def.AzureResourceType))
	}
	def.Label = newUniqueLabel(def.ResourceName, label, &formatter.labels)

	def.AdditionalFields["resource_id"] = formatter.tryAddressOrLiteral(def.ResourceId)
	if def.Method != "POST" {
		def.AdditionalFields["method"] = NewLiteralValue(def.Method)
	}

	return def
}

func (formatter *AzapiFormatter) formatAsAzapiUpdateResource(def AzapiDefinition) AzapiDefinition {
	def.Kind = "resource"
	def.ResourceName = "azapi_update_resource"

	def.Label = newUniqueLabel(def.ResourceName, fmt.Sprintf("update_%s", defaultLabel(def.AzureResourceType)), &formatter.labels)

	parentId := GetParentId(def.ResourceId)
	if resourceIdAddress := formatter.valueToKeyMap[strings.ToUpper(def.ResourceId)]; resourceIdAddress != "" {
		def.AdditionalFields["resource_id"] = NewReferenceValue(resourceIdAddress)
	} else if parentAddress := formatter.valueToKeyMap[strings.ToUpper(parentId)]; parentAddress != "" {
		def.AdditionalFields["parent_id"] = NewReferenceValue(parentAddress)
		def.AdditionalFields["name"] = NewLiteralValue(GetName(def.ResourceId))
	} else {
		def.AdditionalFields["resource_id"] = NewLiteralValue(def.ResourceId)
	}

	if def.Body != nil {
		bodyMap, ok := def.Body.(map[string]interface{})
		if ok && bodyMap != nil {
			delete(bodyMap, "id")
			delete(bodyMap, "name")
			delete(bodyMap, "type")
		}
	}

	return def
}

func (formatter *AzapiFormatter) tryAddressOrLiteral(resourceId string) Value {
	if address := formatter.valueToKeyMap[strings.ToUpper(resourceId)]; address != "" {
		return NewReferenceValue(address)
	}
	return NewLiteralValue(resourceId)
}

func (formatter *AzapiFormatter) injectReference(raw interface{}) interface{} {
	if raw == nil {
		return nil
	}
	switch value := raw.(type) {
	case map[string]interface{}:
		out := make(map[string]interface{})
		for k, v := range value {
			key := formatter.injectReference(k)
			out[key.(string)] = formatter.injectReference(v)
		}
		return out
	case []interface{}:
		for i, v := range value {
			value[i] = formatter.injectReference(v)
		}
		return value
	case string:
		if address := formatter.valueToKeyMap[strings.ToUpper(value)]; address != "" {
			return fmt.Sprintf("${%s}", address)
		} else {
			return value
		}
	default:
		return value
	}
}

func (formatter *AzapiFormatter) populateReference(prefix string, raw interface{}) {
	if raw == nil {
		return
	}
	switch value := raw.(type) {
	case map[string]interface{}:
		for k, v := range value {
			formatter.populateReference(fmt.Sprintf("%s.%s", prefix, k), v)
		}
	case []interface{}:
		for i, v := range value {
			formatter.populateReference(fmt.Sprintf("%s[%d]", prefix, i), v)
		}
	case string:
		if len(value) < 32 {
			// if the string is too short, it's probably not a resource id
			return
		}
		if _, ok := formatter.valueToKeyMap[strings.ToUpper(value)]; !ok {
			formatter.valueToKeyMap[strings.ToUpper(value)] = prefix
		}
	}

}
