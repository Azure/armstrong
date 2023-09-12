package resolver

import (
	"strings"
	
	"github.com/ms-henglu/armstrong/dependency"
	"github.com/ms-henglu/armstrong/resource/types"
)

var _ ReferenceResolver = &AzapiDefinitionResolver{}

type AzapiDefinitionResolver struct {
	AzapiDefinitions []types.AzapiDefinition
}

func (r AzapiDefinitionResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, def := range r.AzapiDefinitions {
		if def.AdditionalFields["action"] != nil {
			continue
		}
		if def.ResourceName == "azapi_resource_list" {
			continue
		}
		if strings.EqualFold(def.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				AzapiDefinitionToAdd: &def,
			}, nil
		}
	}
	return nil, nil
}

func NewAzapiDefinitionResolver(azapiDefinitions []types.AzapiDefinition) AzapiDefinitionResolver {
	return AzapiDefinitionResolver{
		AzapiDefinitions: azapiDefinitions,
	}
}
