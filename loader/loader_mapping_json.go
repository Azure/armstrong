package loader

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/types"
)

type MappingJsonDependencyLoader struct {
	MappingJsonFilepath string
}

func (m MappingJsonDependencyLoader) Load() ([]types.Dependency, error) {
	var mappings []types.Mapping
	data, err := ioutil.ReadFile(m.MappingJsonFilepath)
	if err != nil {
		return []types.Dependency{}, err
	}
	err = json.Unmarshal(data, &mappings)
	if err != nil {
		return []types.Dependency{}, err
	}
	deps := make([]types.Dependency, 0)
	for _, mapping := range mappings {
		deps = append(deps, types.Dependency{
			Pattern:              mapping.IdPattern,
			ExampleConfiguration: mapping.ExampleConfiguration,
			ResourceType:         mapping.ResourceType,
			ReferredProperty:     "id",
		})
	}

	return deps, nil
}

var _ DependencyLoader = MappingJsonDependencyLoader{}
