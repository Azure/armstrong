package resolver

import (
	"github.com/ms-henglu/armstrong/dependency"
	"github.com/ms-henglu/armstrong/resource/types"
)

var _ ReferenceResolver = &KnownReferenceResolver{}

type KnownReferenceResolver struct {
	knownPatterns map[string]types.Reference
}

func (r KnownReferenceResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	ref, ok := r.knownPatterns[pattern.String()]
	if !ok || !ref.IsKnown() {
		return nil, nil
	}
	return &ResolvedResult{
		Reference: &ref,
	}, nil
}

func NewKnownReferenceResolver(knownPatterns map[string]types.Reference) KnownReferenceResolver {
	return KnownReferenceResolver{
		knownPatterns: knownPatterns,
	}
}
