package resource

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pluralize "github.com/gertd/go-pluralize"
	"github.com/ms-henglu/armstrong/hcl"
	"github.com/ms-henglu/armstrong/resource/utils"
	"github.com/ms-henglu/armstrong/types"
)

var _ Base = &Resource{}

var pluralizeClient = pluralize.NewClient()

type Resource struct {
	ApiVersion                 string
	ExampleId                  string
	ExampleBody                interface{}
	PropertyDependencyMappings []PropertyDependencyMapping
	Label                      string
}

func NewResourceFromExample(filepath string) (*Resource, error) {
	exampleData, err := os.ReadFile(filepath)
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

				if apiVer, ok := parametersMap["api-version"].(string); ok {
					apiVersion = apiVer
				}

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
						ValuePath:    "parent",
						LiteralValue: utils.GetParentIdFromId(exampleId),
						IsKey:        false,
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
		Label:                      "resource",
	}, nil
}

func (r *Resource) RequiredDependencies(existingDependencies []types.Dependency, dependencies []types.Dependency) []types.Dependency {
	return requiredDependencies(r.PropertyDependencyMappings, existingDependencies, dependencies)
}

func (r *Resource) UpdatePropertyDependencyMappingsReference(dependencies []types.Dependency, references []Reference) {
	r.PropertyDependencyMappings = updatePropertyDependencyMappingsReference(r.PropertyDependencyMappings, dependencies, references)
}

func (r *Resource) GenerateLabel(references []Reference) string {
	resourceType := utils.GetResourceType(r.ExampleId)
	r.Label = newLabel(resourceType, references)
	return r.Label
}

func (r *Resource) Hcl(useRawJsonPayload bool) string {
	body := ""
	if useRawJsonPayload {
		jsonBody, _ := json.MarshalIndent(r.GetBody(), "", "    ")
		body = fmt.Sprintf(`<<BODY
%s
BODY`, jsonBody)
	} else {
		hclBody := hcl.MarshalIndent(r.GetBody(), "", "  ")
		body = fmt.Sprintf(`jsonencode(%s)`, hclBody)
	}
	resourceType := utils.GetResourceType(r.ExampleId)
	parentId := FindParentReference(r.PropertyDependencyMappings)
	if parentId == "" {
		parentId = fmt.Sprintf(`"%s"`, utils.GetParentIdFromId(r.ExampleId))
	}
	return fmt.Sprintf(`
resource "azapi_resource" "%[6]s" {
	type = "%[3]s@%[4]s"
	parent_id = %[2]s
    name = "%[1]s"

 	body = %[5]s

    schema_validation_enabled = false
    ignore_missing_property = false
}
`, hcl.RandomName(), parentId, resourceType, r.ApiVersion, body, r.Label)
}

func (r *Resource) GetBody() interface{} {
	replacements := make(map[string]string)
	for _, mapping := range r.PropertyDependencyMappings {
		if mapping.ValuePath != "parent" && mapping.Reference != nil {
			if mapping.Reference.IsKnown() {
				valuePath := mapping.ValuePath
				if mapping.IsKey {
					valuePath = fmt.Sprintf("key:%s", mapping.ValuePath)
				}
				replacements[valuePath] = fmt.Sprintf(`${%s}`, mapping.Reference)
			} else {
				log.Printf("[WARN] reference is unkown, reference: %v", mapping.Reference)
			}
		}
	}
	removes := []string{".name"}
	return utils.GetUpdatedBody(r.ExampleBody, replacements, removes, "")
}
