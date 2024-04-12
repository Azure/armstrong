package dependency

import (
	"strings"
	"sync"

	"github.com/azure/armstrong/dependency/azurerm"
)

var azurermMutex = sync.Mutex{}

var azurermDeps = make([]Dependency, 0)

func LoadAzurermDependencies() []Dependency {
	azurermMutex.Lock()
	defer azurermMutex.Unlock()
	if len(azurermDeps) != 0 {
		return azurermDeps
	}

	mappingJsonLoader := azurerm.MappingJsonDependencyLoader{}
	hardcodeLoader := azurerm.HardcodeDependencyLoader{}

	depsMap := make(map[string]azurerm.Mapping)
	if temp, err := mappingJsonLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType] = dep
		}
	}
	if temp, err := hardcodeLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType] = dep
		}
	}

	deps := make([]azurerm.Mapping, 0)
	for _, dep := range depsMap {
		deps = append(deps, dep)
	}

	azurermDeps = make([]Dependency, 0)
	for _, dep := range deps {
		startStr := "providers/"
		index := strings.LastIndex(dep.IdPattern, startStr)
		if index == -1 {
			continue
		}
		azureResourceType := dep.IdPattern[index+len(startStr):]
		azurermDeps = append(azurermDeps, Dependency{
			AzureResourceType:    azureResourceType,
			ExampleConfiguration: dep.ExampleConfiguration,
			ResourceKind:         "resource",
			ReferredProperty:     "id",
			ResourceName:         dep.ResourceType,
			ResourceLabel:        "",
		})
	}
	return azurermDeps
}
