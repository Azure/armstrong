package resource

import (
	"encoding/json"
	"os"

	"github.com/ms-henglu/armstrong/resource/types"
	"github.com/ms-henglu/armstrong/utils"
	"github.com/sirupsen/logrus"
)

func NewAzapiDefinitionFromExample(exampleFilepath string, kind string) (types.AzapiDefinition, error) {
	data, err := os.ReadFile(exampleFilepath)
	if err != nil {
		return types.AzapiDefinition{}, err
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
		return types.AzapiDefinition{}, err
	}

	var body interface{}
	locationValue := ""
	if kind == "resource" && example.Parameters != nil {
		for _, value := range example.Parameters {
			if bodyMap, ok := value.(map[string]interface{}); ok {
				logrus.Debugf("found request body from example: %v", bodyMap)
				if location := bodyMap["location"]; location != nil {
					locationValue = location.(string)
					delete(bodyMap, "location")
				}
				if name := bodyMap["name"]; name != nil {
					delete(bodyMap, "name")
				}
				delete(bodyMap, "id")
				body = bodyMap
				break
			}
		}
		if body == nil {
			logrus.Warnf("found no request body from example")
		}
	}

	var id string
	for _, statusCode := range []string{"200", "201", "202"} {
		if response, ok := example.Responses[statusCode]; ok && response.Body.Id != "" {
			logrus.Debugf("found id from %s response: %s", statusCode, response.Body.Id)
			id = response.Body.Id
			break
		}
	}
	if id == "" {
		logrus.Warnf("found no id from example")
	}

	resourceType := utils.ResourceTypeOfResourceId(id)
	logrus.Debugf("resource type of %s is %s", id, resourceType)

	var apiVersion string
	if example.Parameters != nil && example.Parameters["api-version"] != nil {
		apiVersion = example.Parameters["api-version"].(string)
	}
	if apiVersion == "" {
		apiVersion = "TODO"
		logrus.Warnf("found no api-version from example, please specify it manually")
	}

	out := types.AzapiDefinition{
		Id:                id,
		Kind:              types.Kind(kind),
		ResourceName:      "azapi_resource",
		Label:             defaultLabel(resourceType),
		AzureResourceType: resourceType,
		ApiVersion:        apiVersion,
		Body:              body,
		BodyFormat:        types.BodyFormatHcl,
		AdditionalFields: map[string]types.Value{
			"parent_id": types.NewStringLiteralValue(utils.ParentIdOfResourceId(id)),
			"name":      types.NewStringLiteralValue(utils.LastSegment(id)),
		},
	}

	if kind == "resource" {
		out.AdditionalFields["schema_validation_enabled"] = types.NewRawValue("false")
		out.AdditionalFields["ignore_missing_property"] = types.NewRawValue("false")
		out.AdditionalFields["ignore_casing"] = types.NewRawValue("false")
		if locationValue != "" {
			out.AdditionalFields["location"] = types.NewStringLiteralValue(locationValue)
		}
	}
	return out, nil
}
