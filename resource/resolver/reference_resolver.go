package resolver

import (
	"github.com/gertd/go-pluralize"
	"github.com/ms-henglu/armstrong/dependency"
	"github.com/ms-henglu/armstrong/resource/types"
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
