package dependency

import (
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/ms-henglu/armstrong/utils"
)

type Pattern struct {
	AzureResourceType string
	Scope             Scope
	Placeholder       string
}

type Scope string

const (
	ScopeTenant        Scope = "tenant"
	ScopeSubscription  Scope = "subscription"
	ScopeResourceGroup Scope = "resource_group"
	ScopeResource      Scope = "resource"
)

func (p Pattern) String() string {
	return strings.ToLower(string(p.Scope) + ":" + p.AzureResourceType)
}

func (p Pattern) IsMatch(input string) bool {
	resourceType := utils.ResourceTypeOfResourceId(input)
	if resourceType != p.AzureResourceType {
		return false
	}
	parentId := utils.ParentIdOfResourceId(input)
	scope := scopeOfResourceId(parentId)
	return scope == p.Scope
}

func NewPattern(idPlaceholder string) Pattern {
	resourceType := utils.ResourceTypeOfResourceId(idPlaceholder)
	parentId := utils.ParentIdOfResourceId(idPlaceholder)
	scope := scopeOfResourceId(parentId)
	return Pattern{
		AzureResourceType: resourceType,
		Scope:             scope,
		Placeholder:       idPlaceholder,
	}
}

func scopeOfResourceId(input string) Scope {
	resourceType := utils.ResourceTypeOfResourceId(input)
	switch resourceType {
	case arm.TenantResourceType.String():
		return ScopeTenant
	case arm.SubscriptionResourceType.String():
		return ScopeSubscription
	case arm.ResourceGroupResourceType.String():
		return ScopeResourceGroup
	default:
		return ScopeResource
	}
}
