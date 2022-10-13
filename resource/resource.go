package resource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/ms-henglu/armstrong/hcl"
	"github.com/ms-henglu/armstrong/resource/utils"
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
				if id := utils.GetIdFromResponseExample(responseMap["200"]); len(id) > 0 {
					exampleId = id
				} else if id := utils.GetIdFromResponseExample(responseMap["201"]); len(id) > 0 {
					exampleId = id
				} else if id := utils.GetIdFromResponseExample(responseMap["202"]); len(id) > 0 {
					exampleId = id
				}
				if len(exampleId) > 0 {
					mappings = append(mappings, PropertyDependencyMapping{
						ValuePath: "parent",
						Value:     utils.GetParentIdFromId(exampleId),
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

func (r Resource) Hcl(dependencyHcl string, useRawJsonPayload bool) string {
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
resource "azapi_resource" "%[1]s" {
	type = "%[3]s@%[4]s"
    name = "%[1]s"
	parent_id = %[2]s

 	body = %[5]s

    schema_validation_enabled = false
}
`, hcl.RandomName(), r.FindParentReference(dependencyHcl), utils.GetResourceType(r.ExampleId), r.ApiVersion, body)
}

func (r Resource) GetBody(dependencyHcl string) interface{} {
	replacements := make(map[string]string)
	for _, mapping := range r.PropertyDependencyMappings {
		if mapping.ValuePath != "parent" && len(mapping.Reference) > 0 {
			parts := strings.Split(mapping.Reference, ".")
			if len(parts) == 3 {
				replacements[mapping.ValuePath] = "${" + mapping.Reference + "}"
				continue
			}
			resourceType := parts[0]
			propertyName := parts[1]
			if target := hcl.FindResourceAddress(dependencyHcl, resourceType); len(target) > 0 {
				ref := target + "." + propertyName
				replacements[mapping.ValuePath] = "${" + ref + "}"
			} else {
				log.Printf("[WARN] dependency not found, resource type: %s", resourceType)
			}
		}
	}
	replacements[".location"] = "westeurope"
	removes := []string{".name"}
	return utils.GetUpdatedBody(r.ExampleBody, replacements, removes, "")
}

func (r Resource) FindParentReference(dependencyHcl string) string {
	for _, mapping := range r.PropertyDependencyMappings {
		if mapping.ValuePath == "parent" && len(mapping.Reference) > 0 {
			parts := strings.Split(mapping.Reference, ".")
			if len(parts) == 3 {
				return mapping.Reference
			}
			resourceType := parts[0]
			propertyName := parts[1]
			if target := hcl.FindResourceAddress(dependencyHcl, resourceType); len(target) > 0 {
				ref := target + "." + propertyName
				return ref
			} else {
				log.Printf("[WARN] dependency not found, resource type: %s", resourceType)
			}
		}
	}
	return fmt.Sprintf(`"%s"`, utils.GetParentIdFromId(r.ExampleId))
}

func (r Resource) DependencyHcl(existDeps []types.Dependency, deps []types.Dependency) string {
	dependencyHcl := ""
	for index, mapping := range r.PropertyDependencyMappings {
		for _, dep := range existDeps {
			if utils.IsValueMatchPattern(mapping.Value, dep.Pattern) {
				log.Printf("[INFO] found existing dependency: %s", dep.Address)
				r.PropertyDependencyMappings[index].Reference = dep.Address + "." + dep.ReferredProperty
				break
			}
		}
		if len(r.PropertyDependencyMappings[index].Reference) != 0 {
			continue
		}
		for _, dep := range deps {
			if utils.IsValueMatchPattern(mapping.Value, dep.Pattern) {
				log.Printf("[INFO] found dependency: %s", dep.ResourceType)
				r.PropertyDependencyMappings[index].Reference = dep.ResourceType + "." + dep.ReferredProperty
				dependencyHcl = hcl.Combine(dependencyHcl, hcl.RenameLabel(dep.ExampleConfiguration))
				break // take the first match
			}
		}
	}
	return dependencyHcl
}

// GetKeyValueMappings returns a list of key and value of input
func GetKeyValueMappings(parameters interface{}, path string) []PropertyDependencyMapping {
	if parameters == nil {
		return []PropertyDependencyMapping{}
	}
	results := make([]PropertyDependencyMapping, 0)
	switch param := parameters.(type) {
	case map[string]interface{}:
		for key, value := range param {
			results = append(results, GetKeyValueMappings(value, path+"."+key)...)
		}
	case []interface{}:
		for index, value := range param {
			results = append(results, GetKeyValueMappings(value, path+"."+strconv.Itoa(index))...)
		}
	case string:
		results = append(results, PropertyDependencyMapping{
			ValuePath: path,
			Value:     param,
		})
	default:

	}
	return results
}
