package resolver

import (
	"strings"

	"github.com/azure/armstrong/dependency"
	"github.com/azure/armstrong/resource/types"
	"github.com/azure/armstrong/utils"
)

var _ ReferenceResolver = &AzapiResourcePlaceholderResolver{}

type AzapiResourcePlaceholderResolver struct {
}

func (a AzapiResourcePlaceholderResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	if !strings.Contains(pattern.AzureResourceType, "/") {
		return nil, nil
	}
	return &ResolvedResult{
		AzapiDefinitionToAdd: &types.AzapiDefinition{
			Id:                pattern.Placeholder,
			Kind:              "resource",
			ResourceName:      "azapi_resource",
			Label:             pluralizeClient.Singular(utils.LastSegment(pattern.AzureResourceType)),
			AzureResourceType: pattern.AzureResourceType,
			BodyFormat:        types.BodyFormatHcl,
			ApiVersion:        "TODO",
			AdditionalFields: map[string]types.Value{
				"parent_id":                 types.NewStringLiteralValue(utils.ParentIdOfResourceId(pattern.Placeholder)),
				"name":                      types.NewStringLiteralValue(utils.LastSegment(pattern.Placeholder)),
				"schema_validation_enabled": types.NewRawValue("false"),
			},
			Body: map[string]interface{}{
				"properties": map[string]interface{}{
					"TODO": "TODO",
				},
			},
		},
	}, nil
}

func NewAzapiResourcePlaceholderResolver() AzapiResourcePlaceholderResolver {
	return AzapiResourcePlaceholderResolver{}
}
