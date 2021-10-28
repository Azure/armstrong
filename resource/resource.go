package resource

import (
	"encoding/json"
	"fmt"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/types"
	"io/ioutil"
	"strings"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/helper"
)

type Resource struct {
	ApiVersion                 string
	ExampleId                  string
	ExampleBody                interface{}
	PropertyDependencyMappings []PropertyDependencyMapping
}

type PropertyDependencyMapping struct {
	ValuePath string
	Value     string
	Reference string
}

func NewResourceFromExample(filepath, bodyPropertyName string) (*Resource, error) {
	exampleData, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var example interface{}
	err = json.Unmarshal(exampleData, &example)
	if err != nil {
		return nil, err
	}

	var body interface{}
	exampleId := ""
	apiVersion := ""
	mappings := make([]PropertyDependencyMapping, 0)
	if exampleMap, ok := example.(map[string]interface{}); ok {
		if exampleMap["parameters"] != nil {
			if parametersMap, ok := exampleMap["parameters"].(map[string]interface{}); ok {
				body = parametersMap[bodyPropertyName]
				mappings = append(mappings, GetKeyValueMappings(body, "")...)

				apiVersion = parametersMap["api-version"].(string)
			}
		}

		if exampleMap["responses"] != nil {
			if responseMap, ok := exampleMap["responses"].(map[string]interface{}); ok {
				if id := GetIdFromResponseExample(responseMap["200"]); len(id) > 0 {
					exampleId = id
				} else if id := GetIdFromResponseExample(responseMap["201"]); len(id) > 0 {
					exampleId = id
				}
				if len(exampleId) > 0 {
					mappings = append(mappings, PropertyDependencyMapping{
						ValuePath: "parent",
						Value:     GetParentIdFromId(exampleId),
					})
				}
			}
		}
	}

	return &Resource{
		ApiVersion:                 apiVersion,
		ExampleId:                  exampleId,
		ExampleBody:                body,
		PropertyDependencyMappings: mappings,
	}, nil
}

func (r Resource) GetHcl(dependencyHcl string) string {
	body, _ := json.MarshalIndent(r.GetBody(dependencyHcl), "", "    ")
	return fmt.Sprintf(`
resource "azurermg_resource" "test" {
	url = "%s"
	api_version = "%s"
 	body = <<BODY
%s
BODY
}
`, r.GetUrl(dependencyHcl), r.ApiVersion, body)
}

func (r Resource) GetBody(dependencyHcl string) interface{} {
	replacements := make(map[string]string, 0)
	for _, mapping := range r.PropertyDependencyMappings {
		if mapping.ValuePath != "parent" && len(mapping.Reference) > 0 {
			parts := strings.Split(mapping.Reference, ".")
			resourceType := parts[0]
			propertyName := parts[1]
			if target := helper.GetResourceFromHcl(dependencyHcl, resourceType); len(target) > 0 {
				ref := target + "." + propertyName
				replacements[mapping.ValuePath] = "${" + ref + "}"
			} else {
				fmt.Printf("[WARN] dependency not found, resource type: %s", resourceType)
			}
		}
	}
	return GetUpdatedBody(r.ExampleBody, replacements, "")
}

func (r Resource) GetUrl(dependencyHcl string) string {
	for _, mapping := range r.PropertyDependencyMappings {
		if mapping.ValuePath == "parent" && len(mapping.Reference) > 0 {
			parts := strings.Split(mapping.Reference, ".")
			resourceType := parts[0]
			propertyName := parts[1]
			if target := helper.GetResourceFromHcl(dependencyHcl, resourceType); len(target) > 0 {
				ref := target + "." + propertyName
				return strings.ReplaceAll(r.ExampleId, mapping.Value, "${"+ref+"}")
			} else {
				fmt.Printf("[WARN] dependency not found, resource type: %s", resourceType)
			}
		}
	}
	return r.ExampleId
}

func (r Resource) GetDependencyHcl(deps []types.Dependency) string {
	dependencyHcl := ""
	for index, mapping := range r.PropertyDependencyMappings {
		for _, dep := range deps {
			if helper.IsValueMatchPattern(mapping.Value, dep.Pattern) {
				r.PropertyDependencyMappings[index].Reference = dep.ResourceType + "." + dep.ReferredProperty
				dependencyHcl = helper.GetCombinedHcl(dependencyHcl, helper.GetRenamedHcl(dep.ExampleConfiguration))
				break // take the first match
			}
		}
	}
	return dependencyHcl
}
