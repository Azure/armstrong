package resource

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ms-henglu/armstrong/resource/utils"
	"github.com/ms-henglu/armstrong/types"
)

type Base interface {
	Hcl(useRawJsonPayload bool) string
	RequiredDependencies(existingDependencies []types.Dependency, dependencies []types.Dependency) []types.Dependency
	UpdatePropertyDependencyMappingsReference(dependencies []types.Dependency, references []Reference)
	GenerateLabel(references []Reference) string
}

type PropertyDependencyMapping struct {
	IsKey        bool
	ValuePath    string
	LiteralValue string
	Reference    *Reference
}

func FindParentReference(propertyDependencyMappings []PropertyDependencyMapping) string {
	var parentMapping PropertyDependencyMapping
	for _, mapping := range propertyDependencyMappings {
		if mapping.ValuePath == "parent" {
			parentMapping = mapping
			break
		}
	}
	if parentMapping.Reference.IsKnown() {
		return parentMapping.Reference.String()
	} else {
		log.Printf("[WARN] reference is unkown, reference: %v", parentMapping.Reference)
	}

	return ""
}

// requiredDependencies returns the dependencies that are required by the resource
func requiredDependencies(propertyDependencyMappings []PropertyDependencyMapping, existDeps []types.Dependency, deps []types.Dependency) []types.Dependency {
	out := make([]types.Dependency, 0)
	for _, mapping := range propertyDependencyMappings {
		found := false
		for _, dep := range existDeps {
			if utils.IsValueMatchPattern(mapping.LiteralValue, dep.Pattern) {
				log.Printf("[INFO] found existing dependency: %s", dep.Address)
				found = true
				break
			}
		}
		if found {
			continue
		}
		for _, dep := range deps {
			if utils.IsValueMatchPattern(mapping.LiteralValue, dep.Pattern) {
				log.Printf("[INFO] found dependency: %s", dep.ResourceType)
				out = append(out, dep)
				break // take the first match
			}
		}
	}
	return out
}

// updatePropertyDependencyMappingsReference updates the reference of property dependency mappings
func updatePropertyDependencyMappingsReference(propertyDependencyMappings []PropertyDependencyMapping, deps []types.Dependency, refs []Reference) []PropertyDependencyMapping {
	refMap := make(map[string]Reference)
	for _, ref := range refs {
		refMap[ref.ResourceType] = ref
	}

	for index, mapping := range propertyDependencyMappings {
		for _, dep := range deps {
			if utils.IsValueMatchPattern(mapping.LiteralValue, dep.Pattern) {
				if dep.Address != "" {
					propertyDependencyMappings[index].Reference = NewReferenceFromAddress(fmt.Sprintf("%s.%s", dep.Address, dep.ReferredProperty))
				} else {
					propertyDependencyMappings[index].Reference = &Reference{
						Type:         "resource",
						ResourceType: dep.ResourceType,
						PropertyName: dep.ReferredProperty,
					}
					if target, ok := refMap[dep.ResourceType]; ok {
						propertyDependencyMappings[index].Reference.Label = target.Label
					} else {
						log.Printf("[WARN] dependency not found, resource type: %s", dep.ResourceType)
					}
				}
				break
			}
		}
	}
	return propertyDependencyMappings
}

// GetKeyValueMappings returns a list of key and value of input
func GetKeyValueMappings(parameters interface{}, path string) []PropertyDependencyMapping {
	if parameters == nil {
		return []PropertyDependencyMapping{}
	}
	results := make([]PropertyDependencyMapping, 0)
	switch param := parameters.(type) {
	case map[string]interface{}:
		for key, value := range param {
			results = append(results, GetKeyValueMappings(value, path+"."+key)...)
			results = append(results, PropertyDependencyMapping{
				ValuePath:    path + "." + key,
				LiteralValue: key,
				IsKey:        true,
			})
		}
	case []interface{}:
		for index, value := range param {
			results = append(results, GetKeyValueMappings(value, path+"."+strconv.Itoa(index))...)
		}
	case string:
		results = append(results, PropertyDependencyMapping{
			ValuePath:    path,
			LiteralValue: param,
			IsKey:        false,
		})
	default:

	}
	return results
}

// newLabel returns an unique label for a resource type
func newLabel(resourceType string, refs []Reference) string {
	idMap := make(map[string]bool)
	for _, ref := range refs {
		if ref.ResourceType == "azapi_resource" {
			idMap[ref.Label] = true
		}
	}
	parts := strings.Split(resourceType, "/")
	label := "test"
	if len(parts) != 0 {
		label = parts[len(parts)-1]
		label = pluralizeClient.Singular(label)
	}
	_, ok := idMap[label]
	if !ok {
		return label
	}
	for i := 2; i <= 100; i++ {
		newLabel := fmt.Sprintf("%s%d", label, i)
		_, ok = idMap[newLabel]
		if !ok {
			return newLabel
		}
	}
	return label
}
