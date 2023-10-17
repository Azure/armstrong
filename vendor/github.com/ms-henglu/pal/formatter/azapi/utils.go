package azapi

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
)

func parseResourceTypeApiVersion(input string) (string, string) {
	idUrl, err := url.Parse(input)
	if err != nil {
		return "", ""
	}
	apiVersion := idUrl.Query().Get("api-version")
	resourceType := GetResourceType(idUrl.Path)
	return resourceType, apiVersion
}

func GetResourceType(id string) string {
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
		value := ""
		value = components[current+1]

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
	if resourceType == "" {
		resourceId, err := arm.ParseResourceID(id)
		if err != nil {
			return ""
		}
		return resourceId.ResourceType.String()
	}
	return resourceType
}

// newLabel returns an unique label for a resource type
func newUniqueLabel(prefix string, label string, labels *map[string]bool) string {
	check := fmt.Sprintf("%s.%s", prefix, label)
	_, ok := (*labels)[check]
	if !ok {
		(*labels)[check] = true
		return label
	}
	for i := 2; i <= 100; i++ {
		newLabel := fmt.Sprintf("%s%d", label, i)
		check = fmt.Sprintf("%s.%s", prefix, newLabel)
		_, ok := (*labels)[check]
		if !ok {
			(*labels)[check] = true
			return newLabel
		}
	}
	return label
}

func defaultLabel(resourceType string) string {
	parts := strings.Split(resourceType, "/")
	label := "test"
	if len(parts) != 0 {
		label = parts[len(parts)-1]
		label = pluralizeClient.Singular(label)
	}
	return label
}

func GetName(id string) string {
	resourceId, err := arm.ParseResourceID(id)
	if err != nil {
		return ""
	}
	return resourceId.Name
}

func GetParentId(id string) string {
	resourceId, err := arm.ParseResourceID(id)
	if err != nil {
		return ""
	}
	if resourceId.Parent.ResourceType.String() == arm.TenantResourceType.String() {
		return "/"
	}
	return resourceId.Parent.String()
}

func GetId(input string) string {
	idUrl, err := url.Parse(input)
	if err != nil {
		return ""
	}
	return idUrl.Path
}

func IsResourceAction(resourceId string) bool {
	return len(strings.Split(resourceId, "/"))%2 == 0
}
