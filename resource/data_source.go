package resource

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ms-henglu/armstrong/hcl"
	"github.com/ms-henglu/armstrong/resource/utils"
	"github.com/ms-henglu/armstrong/types"
)

var _ Base = &DataSource{}

type DataSource struct {
	ApiVersion                 string
	ExampleId                  string
	PropertyDependencyMappings []PropertyDependencyMapping
	Label                      string
}

func (r *DataSource) GenerateLabel(references []Reference) string {
	resourceType := utils.GetResourceType(r.ExampleId)
	r.Label = newLabel(resourceType, references)
	return r.Label
}

func NewDataSourceFromExample(filepath string) (*DataSource, error) {
	exampleData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var example interface{}
	err = json.Unmarshal(exampleData, &example)
	if err != nil {
		return nil, err
	}

	exampleId := ""
	apiVersion := ""
	mappings := make([]PropertyDependencyMapping, 0)
	if exampleMap, ok := example.(map[string]interface{}); ok {
		if exampleMap["parameters"] != nil {
			if parametersMap, ok := exampleMap["parameters"].(map[string]interface{}); ok {
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
						ValuePath:    "parent",
						LiteralValue: utils.GetParentIdFromId(exampleId),
						IsKey:        false,
					})
				}
			}
		}
	}

	return &DataSource{
		ApiVersion:                 apiVersion,
		ExampleId:                  exampleId,
		PropertyDependencyMappings: mappings,
		Label:                      "",
	}, nil
}

func (r *DataSource) RequiredDependencies(existingDependencies []types.Dependency, dependencies []types.Dependency) []types.Dependency {
	return requiredDependencies(r.PropertyDependencyMappings, existingDependencies, dependencies)
}

func (r *DataSource) UpdatePropertyDependencyMappingsReference(dependencies []types.Dependency, references []Reference) {
	r.PropertyDependencyMappings = updatePropertyDependencyMappingsReference(r.PropertyDependencyMappings, dependencies, references)
}

func (r *DataSource) Hcl(_ bool) string {
	resourceType := utils.GetResourceType(r.ExampleId)
	name := utils.GetName(r.ExampleId)
	if name != "default" {
		name = hcl.RandomName()
	}
	parentId := FindParentReference(r.PropertyDependencyMappings)
	if parentId == "" {
		parentId = fmt.Sprintf(`"%s"`, utils.GetParentIdFromId(r.ExampleId))
	}
	return fmt.Sprintf(`
data "azapi_resource" "%[5]s" {
	type = "%[3]s@%[4]s"
	parent_id = %[2]s
    name = "%[1]s"
}
`, name, parentId, resourceType, r.ApiVersion, r.Label)
}
