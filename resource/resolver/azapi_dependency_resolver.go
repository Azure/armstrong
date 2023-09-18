package resolver

import (
	"strings"

	"github.com/ms-henglu/armstrong/dependency"
	"github.com/sirupsen/logrus"
)

var _ ReferenceResolver = &AzapiDependencyResolver{}

type AzapiDependencyResolver struct {
	dependencies []dependency.Dependency
}

func (r AzapiDependencyResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, dep := range r.dependencies {
		if strings.EqualFold(dep.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				HclToAdd: dep.ExampleConfiguration,
			}, nil
		}
	}
	return nil, nil
}

func NewAzapiDependencyResolver() AzapiDependencyResolver {
	azapiDeps, err := dependency.LoadAzapiDependencies()
	if err != nil {
		logrus.Fatalf("loading azapi dependencies: %+v", err)
		return AzapiDependencyResolver{}
	}
	logrus.Debugf("loaded %d azapi dependencies", len(azapiDeps))
	return AzapiDependencyResolver{
		dependencies: azapiDeps,
	}
}
