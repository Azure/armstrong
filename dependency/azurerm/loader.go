package loader

import "github.com/ms-henglu/armstrong/types"

type DependencyLoader interface {
	Load() ([]types.Dependency, error)
}
