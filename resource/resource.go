package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/ms-henglu/armstrong/hcl"
	"github.com/ms-henglu/armstrong/helper"
	"github.com/ms-henglu/armstrong/types"
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

func NewResourceFromExample(filepath string) (*Resource, error) {
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
				for _, value := range parametersMap {
					if bodyMap, ok := value.(map[string]interface{}); ok {
						body = bodyMap
					}
				}
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
				} else if id := GetIdFromResponseExample(responseMap["202"]); len(id) > 0 {
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

func (r Resource) GetHcl(dependencyHcl string, useRawJsonPayload bool) string {
	body := ""
	if useRawJsonPayload {
		jsonBody, _ := json.MarshalIndent(r.GetBody(dependencyHcl), "", "    ")
		body = fmt.Sprintf(`<<BODY
%s
BODY`, jsonBody)
	} else {
		hclBody := hcl.MarshalIndent(r.GetBody(dependencyHcl), "", "  ")
		body = fmt.Sprintf(`jsonencode(%s)`, hclBody)
	}
	return fmt.Sprintf(`
resource "azapi_resource" "test" {
    name = "%s"
	parent_id = %s
	type = "%s@%s"
 	body = %s
    schema_validation_enabled = false
}
`, helper.GetRandomResourceName(), r.GetParentReference(dependencyHcl), helper.GetResourceType(r.ExampleId), r.ApiVersion, body)
}

func (r Resource) GetBody(dependencyHcl string) interface{} {
	replacements := make(map[string]string)
	for _, mapping := range r.PropertyDependencyMappings {
		if mapping.ValuePath != "parent" && len(mapping.Reference) > 0 {
			parts := strings.Split(mapping.Reference, ".")
			resourceType := parts[0]
			propertyName := parts[1]
			if target := helper.GetResourceFromHcl(dependencyHcl, resourceType); len(target) > 0 {
				ref := target + "." + propertyName
				replacements[mapping.ValuePath] = "${" + ref + "}"
			} else {
				log.Printf("[WARN] dependency not found, resource type: %s", resourceType)
			}
		}
	}
	replacements[".location"] = "westeurope"
	removes := []string{".name"}
	return GetUpdatedBody(r.ExampleBody, replacements, removes, "")
}

func (r Resource) GetParentReference(dependencyHcl string) string {
	for _, mapping := range r.PropertyDependencyMappings {
		if mapping.ValuePath == "parent" && len(mapping.Reference) > 0 {
			parts := strings.Split(mapping.Reference, ".")
			resourceType := parts[0]
			propertyName := parts[1]
			if target := helper.GetResourceFromHcl(dependencyHcl, resourceType); len(target) > 0 {
				ref := target + "." + propertyName
				return ref
			} else {
				log.Printf("[WARN] dependency not found, resource type: %s", resourceType)
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
				log.Printf("[INFO] found dependency: %s", dep.ResourceType)
				r.PropertyDependencyMappings[index].Reference = dep.ResourceType + "." + dep.ReferredProperty
				dependencyHcl = helper.GetCombinedHcl(dependencyHcl, helper.GetRenamedHcl(dep.ExampleConfiguration))
				break // take the first match
			}
		}
	}
	return dependencyHcl
}
