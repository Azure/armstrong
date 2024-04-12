package resolver

import (
	"fmt"
	"strings"

	"github.com/azure/armstrong/dependency"
	"github.com/azure/armstrong/resource/types"
	"github.com/azure/armstrong/utils"
)

var _ ReferenceResolver = &ProviderIDResolver{}

type ProviderIDResolver struct {
}

func (r ProviderIDResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	if !strings.Contains(pattern.AzureResourceType, "/") && pattern.Scope != dependency.ScopeTenant {
		azapiDef := types.AzapiDefinition{
			Id:                pattern.Placeholder,
			Kind:              types.KindDataSource,
			ResourceName:      "azapi_resource_id",
			Label:             fmt.Sprintf("%sScopeProvider", pattern.Scope),
			AzureResourceType: "Microsoft.Resources/providers",
			ApiVersion:        "2020-06-01",
			BodyFormat:        types.BodyFormatHcl,
			AdditionalFields: map[string]types.Value{
				"parent_id": types.NewStringLiteralValue(utils.ParentIdOfResourceId(pattern.Placeholder)),
				"name":      types.NewStringLiteralValue(pattern.AzureResourceType),
			},
		}
		return &ResolvedResult{
			AzapiDefinitionToAdd: &azapiDef,
		}, nil
	}
	return nil, nil
}

func NewProviderIDResolver() ProviderIDResolver {
	return ProviderIDResolver{}
}
