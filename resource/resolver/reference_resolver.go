package resolver

import (
	"github.com/azure/armstrong/dependency"
	"github.com/azure/armstrong/resource/types"
	"github.com/gertd/go-pluralize"
)

var pluralizeClient = pluralize.NewClient()

type ResolvedResult struct {
	Reference            *types.Reference
	HclToAdd             string
	AzapiDefinitionToAdd *types.AzapiDefinition
}

type ReferenceResolver interface {
	Resolve(pattern dependency.Pattern) (*ResolvedResult, error)
}
