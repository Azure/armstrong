package resolver

import (
	"strings"

	"github.com/ms-henglu/armstrong/dependency"
	"github.com/sirupsen/logrus"
)

var _ ReferenceResolver = &AzurermDependencyResolver{}

type AzurermDependencyResolver struct {
	Dependencies []dependency.Dependency
}

func (r AzurermDependencyResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, dep := range r.Dependencies {
		if strings.EqualFold(dep.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				HclToAdd: dep.ExampleConfiguration,
			}, nil
		}
	}
	return nil, nil
}

func NewAzurermDependencyResolver() AzurermDependencyResolver {
	deps := dependency.LoadAzurermDependencies()
	logrus.Debugf("loaded %d azurerm dependencies", len(deps))
	return AzurermDependencyResolver{
		Dependencies: deps,
	}
}
