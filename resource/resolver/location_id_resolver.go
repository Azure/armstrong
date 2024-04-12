package resolver

import (
	"github.com/azure/armstrong/dependency"
	"github.com/azure/armstrong/resource/types"
	"github.com/azure/armstrong/utils"
)

var _ ReferenceResolver = &LocationIDResolver{}

type LocationIDResolver struct {
}

func (r LocationIDResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	if utils.LastSegment(pattern.AzureResourceType) == "locations" {
		azapiDef := types.AzapiDefinition{
			Id:                pattern.Placeholder,
			Kind:              types.KindDataSource,
			ResourceName:      "azapi_resource_id",
			Label:             "location",
			BodyFormat:        types.BodyFormatHcl,
			AzureResourceType: pattern.AzureResourceType,
			ApiVersion:        "2023-12-12",
			AdditionalFields: map[string]types.Value{
				"parent_id": types.NewStringLiteralValue(utils.ParentIdOfResourceId(pattern.Placeholder)),
				"name":      types.NewReferenceValue("var.location"),
			},
		}
		return &ResolvedResult{
			AzapiDefinitionToAdd: &azapiDef,
		}, nil
	}
	return nil, nil
}

func NewLocationIDResolver() LocationIDResolver {
	return LocationIDResolver{}
}
