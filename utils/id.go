package utils

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/sirupsen/logrus"
)

func IsResourceId(input string) bool {
	id := strings.Trim(input, "/")
	if len(strings.Split(id, "/"))%2 == 1 {
		return false
	}
	_, err := arm.ParseResourceID(input)
	return err == nil
}

func ActionName(input string) string {
	if !IsAction(input) {
		return ""
	}
	return LastSegment(input)
}

func LastSegment(input string) string {
	id := strings.Trim(input, "/")
	components := strings.Split(id, "/")
	if len(components) == 0 {
		return ""
	}
	return components[len(components)-1]
}

func IsAction(input string) bool {
	id := strings.Trim(input, "/")
	components := strings.Split(id, "/")
	return len(components)%2 == 1
}

func ResourceIdOfAction(input string) string {
	id := strings.Trim(input, "/")
	components := strings.Split(id, "/")
	return fmt.Sprintf("/%s", strings.Join(components[:len(components)-1], "/"))
}

func ScopeOfListAction(input string) string {
	id := fmt.Sprintf("%s/{placeholder}", input)
	return ParentIdOfResourceId(id)
}

func ParentIdOfResourceId(input string) string {
	resourceId, err := arm.ParseResourceID(input)
	if err == nil && resourceId.Parent != nil {
		if resourceId.Parent.ResourceType.String() == arm.TenantResourceType.String() {
			return "/"
		}
		return resourceId.Parent.String()
	}
	if !strings.Contains(ResourceTypeOfResourceId(input), "/") {
		id := strings.Trim(input, "/")
		components := strings.Split(id, "/")
		if len(components) <= 2 {
			logrus.Warnf("no parent_id found for resource id: %s", input)
			return ""
		}
		return fmt.Sprintf("/%s", strings.Join(components[:len(components)-2], "/"))
	}
	return ""
}

func ResourceTypeOfResourceId(input string) string {
	if input == "/" {
		return arm.TenantResourceType.String()
	}
	id := input
	if IsAction(id) {
		id = ResourceIdOfAction(id)
	}
	if resourceType, err := arm.ParseResourceType(id); err == nil {
		if resourceType.Type != arm.ProviderResourceType.Type {
			return resourceType.String()
		}
	}

	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return ""
	}

	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")

	resourceType := ""
	provider := ""
	for current := 0; current < len(components)-1; current += 2 {
		key := components[current]
		value := components[current+1]

		// Check key/value for empty strings.
		if key == "" || value == "" {
			return ""
		}

		if key == "providers" {
			provider = value
			resourceType = provider
		} else if len(provider) > 0 {
			resourceType += "/" + key
		}
	}
	return resourceType
}
