package resolver

import (
	"strings"

	"github.com/ms-henglu/armstrong/dependency"
	"github.com/ms-henglu/armstrong/resource/types"
	"github.com/ms-henglu/armstrong/utils"
)

var _ ReferenceResolver = &AzapiResourceIdResolver{}

type AzapiResourceIdResolver struct {
}

func (r AzapiResourceIdResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	if !strings.Contains(pattern.AzureResourceType, "/") {
		return nil, nil
	}
	return &ResolvedResult{
		AzapiDefinitionToAdd: &types.AzapiDefinition{
			Id:                pattern.Placeholder,
			Kind:              "data",
			ResourceName:      "azapi_resource_id",
			BodyFormat:        types.BodyFormatHcl,
			Label:             pluralizeClient.Singular(utils.LastSegment(pattern.AzureResourceType)),
			AzureResourceType: pattern.AzureResourceType,
			ApiVersion:        "2023-12-12",
			AdditionalFields: map[string]types.Value{
				"parent_id": types.NewStringLiteralValue(utils.ParentIdOfResourceId(pattern.Placeholder)),
				"name":      types.NewStringLiteralValue(utils.LastSegment(pattern.Placeholder)),
			},
		},
	}, nil
}

func NewAzapiResourceIdResolver() AzapiResourceIdResolver {
	return AzapiResourceIdResolver{}
}
