package loader

import "github.com/ms-henglu/azurerm-rest-api-testing-tool/types"

type DependencyLoader interface {
	Load() ([]types.Dependency, error)
}
