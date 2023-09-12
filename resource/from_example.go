package res

import (
	"os"

	"encoding/json"
	"github.com/ms-henglu/armstrong/utils"
)

func NewAzapiDefinitionFromExample(exampleFilepath string, kind string) (AzapiDefinition, error) {
	data, err := os.ReadFile(exampleFilepath)
	if err != nil {
		return AzapiDefinition{}, err
	}
	var example struct {
		Parameters map[string]interface{} `json:"parameters"`
		Responses  map[string]struct {
			Body struct {
				Id string `json:"id"`
			} `json:"body"`
		} `json:"responses"`
	}
	err = json.Unmarshal(data, &example)
	if err != nil {
		return AzapiDefinition{}, err
	}

	var body interface{}
	if kind == "resource" && example.Parameters != nil {
		for _, value := range example.Parameters {
			if bodyMap, ok := value.(map[string]interface{}); ok {
				body = bodyMap
			}
		}
	}

	var id string
	for _, statusCode := range []string{"200", "201", "202"} {
		if response, ok := example.Responses[statusCode]; ok && response.Body.Id != "" {
			id = response.Body.Id
		}
	}

	resourceType := utils.ResourceTypeOfResourceId(id)

	var apiVersion string
	if example.Parameters != nil && example.Parameters["api-version"] != nil {
		apiVersion = example.Parameters["api-version"].(string)
	}

	return AzapiDefinition{
		Id:                id,
		Kind:              Kind(kind),
		ResourceName:      "azapi_resource",
		Label:             defaultLabel(resourceType),
		AzureResourceType: resourceType,
		ApiVersion:        apiVersion,
		Body:              body,
		AdditionalFields: map[string]Value{
			"parent_id":                 NewStringLiteralValue(utils.ParentIdOfResourceId(id)),
			"name":                      NewStringLiteralValue(utils.LastSegment(id)),
			"schema_validation_enabled": NewRawValue("false"),
		},
	}, nil
}
