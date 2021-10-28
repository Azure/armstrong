package loader

import "github.com/ms-henglu/azurerm-rest-api-testing-tool/types"

type HardcodeDependencyLoader struct {
}

func (h HardcodeDependencyLoader) Load() ([]types.Dependency, error) {
	return []types.Dependency{}, nil
}

var _ DependencyLoader = HardcodeDependencyLoader{}
