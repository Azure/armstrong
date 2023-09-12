package resource

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"

	"github.com/ms-henglu/armstrong/dependency"
	"github.com/ms-henglu/armstrong/utils"
)

type ResolvedResult struct {
	Reference            *Reference
	HclToAdd             string
	AzapiDefinitionToAdd *AzapiDefinition
}

type ReferenceResolver interface {
	Resolve(pattern dependency.Pattern) (*ResolvedResult, error)
}

var _ ReferenceResolver = &KnownReferenceResolver{}

type KnownReferenceResolver struct {
	knownPatterns map[string]Reference
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

func NewKnownReferenceResolver(knownPatterns map[string]Reference) KnownReferenceResolver {
	return KnownReferenceResolver{
		knownPatterns: knownPatterns,
	}
}

var _ ReferenceResolver = &AzapiDependencyResolver{}

type AzapiDependencyResolver struct {
	dependencies []dependency.Dependency
}

func (r AzapiDependencyResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, dep := range r.dependencies {
		if strings.EqualFold(dep.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				HclToAdd: dep.ExampleConfiguration,
			}, nil
		}
	}
	return nil, nil
}

func NewAzapiDependencyResolver() AzapiDependencyResolver {
	logrus.Debugf("loading azapi dependencies...")
	azapiDeps, err := dependency.LoadAzapiDependencies()
	if err != nil {
		logrus.Fatalf("loading azapi dependencies: %+v", err)
		return AzapiDependencyResolver{}
	}
	return AzapiDependencyResolver{
		dependencies: azapiDeps,
	}
}

var _ ReferenceResolver = &AzapiDefinitionResolver{}

type AzapiDefinitionResolver struct {
	AzapiDefinitions []AzapiDefinition
}

func (r AzapiDefinitionResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, def := range r.AzapiDefinitions {
		if def.AdditionalFields["action"] != nil {
			continue
		}
		if def.ResourceName == "azapi_resource_list" {
			continue
		}
		if strings.EqualFold(def.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				AzapiDefinitionToAdd: &def,
			}, nil
		}
	}
	return nil, nil
}

func NewAzapiDefinitionResolver(azapiDefinitions []AzapiDefinition) AzapiDefinitionResolver {
	return AzapiDefinitionResolver{
		AzapiDefinitions: azapiDefinitions,
	}
}

var _ ReferenceResolver = &ProviderIDResolver{}

type ProviderIDResolver struct {
}

func (r ProviderIDResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	if !strings.Contains(pattern.AzureResourceType, "/") && pattern.Scope != dependency.ScopeTenant {
		azapiDef := AzapiDefinition{
			Id:                pattern.Placeholder,
			Kind:              KindDataSource,
			ResourceName:      "azapi_resource_id",
			Label:             fmt.Sprintf("%sScopeProvider", pattern.Scope),
			AzureResourceType: "Microsoft.Resources/providers",
			ApiVersion:        "2020-06-01",
			AdditionalFields: map[string]Value{
				"parent_id": NewStringLiteralValue(utils.ParentIdOfResourceId(pattern.Placeholder)),
				"name":      NewStringLiteralValue(pattern.AzureResourceType),
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

var _ ReferenceResolver = &LocationIDResolver{}

type LocationIDResolver struct {
}

func (r LocationIDResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	if utils.LastSegment(pattern.AzureResourceType) == "locations" {
		azapiDef := AzapiDefinition{
			Id:                pattern.Placeholder,
			Kind:              KindDataSource,
			ResourceName:      "azapi_resource_id",
			Label:             "location",
			AzureResourceType: pattern.AzureResourceType,
			ApiVersion:        "2023-12-12",
			AdditionalFields: map[string]Value{
				"parent_id": NewStringLiteralValue(utils.ParentIdOfResourceId(pattern.Placeholder)),
				"name":      NewReferenceValue("var.location"),
			},
		}
		return &ResolvedResult{
			AzapiDefinitionToAdd: &azapiDef,
		}, nil
	}
	return nil, nil
}

func NewLocationIDResolver() LocationIDResolver {
	return LocationIDResolver{}
}

var _ ReferenceResolver = &AzapiResourceIdResolver{}

type AzapiResourceIdResolver struct {
}

func (r AzapiResourceIdResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	if !strings.Contains(pattern.AzureResourceType, "/") {
		return nil, nil
	}
	return &ResolvedResult{
		AzapiDefinitionToAdd: &AzapiDefinition{
			Id:                pattern.Placeholder,
			Kind:              "data",
			ResourceName:      "azapi_resource_id",
			Label:             pluralizeClient.Singular(utils.LastSegment(pattern.AzureResourceType)),
			AzureResourceType: pattern.AzureResourceType,
			ApiVersion:        "2023-12-12",
			AdditionalFields: map[string]Value{
				"parent_id": NewStringLiteralValue(utils.ParentIdOfResourceId(pattern.Placeholder)),
				"name":      NewStringLiteralValue(utils.LastSegment(pattern.Placeholder)),
			},
		},
	}, nil
}

func NewAzapiResourceIdResolver() AzapiResourceIdResolver {
	return AzapiResourceIdResolver{}
}

var _ ReferenceResolver = &AzurermDependencyResolver{}

type AzurermDependencyResolver struct {
	Dependencies []dependency.Dependency
}

func (r AzurermDependencyResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, dep := range r.Dependencies {
		if strings.EqualFold(dep.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				HclToAdd: dep.ExampleConfiguration,
			}, nil
		}
	}
	return nil, nil
}

func NewAzurermDependencyResolver() AzurermDependencyResolver {
	return AzurermDependencyResolver{
		Dependencies: dependency.LoadAzurermDependencies(),
	}
}

var _ ReferenceResolver = &ExistingDependencyResolver{}

type ExistingDependencyResolver struct {
	ExistingDependencies []dependency.Dependency
}

func (r ExistingDependencyResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, dep := range r.ExistingDependencies {
		if strings.EqualFold(dep.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				Reference: &Reference{
					Label:    dep.ResourceLabel,
					Kind:     dep.ResourceKind,
					Name:     dep.ResourceName,
					Property: dep.ReferredProperty,
				},
			}, nil
		}
	}
	return nil, nil
}

func NewExistingDependencyResolver(workingDirectory string) ExistingDependencyResolver {
	azurermDeps := dependency.LoadAzurermDependencies()
	azurermTypeAzureResourceTypeMap := make(map[string]string)
	for _, dep := range azurermDeps {
		azurermTypeAzureResourceTypeMap[dep.ResourceName] = dep.AzureResourceType
	}
	files, err := os.ReadDir(workingDirectory)
	if err != nil {
		logrus.Warnf("reading dir %s: %+v", workingDirectory, err)
		return ExistingDependencyResolver{}
	}
	existDeps := make([]dependency.Dependency, 0)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}
		src, err := os.ReadFile(path.Join(workingDirectory, file.Name()))
		if err != nil {
			logrus.Warnf("reading file %s: %+v", file.Name(), err)
			continue
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if diag.HasErrors() {
			logrus.Warnf("parsing file %s: %+v", file.Name(), diag.Error())
			continue
		}
		if f == nil || f.Body() == nil {
			continue
		}
		for _, block := range f.Body().Blocks() {
			labels := block.Labels()
			if len(labels) >= 2 {
				switch {
				case strings.HasPrefix(labels[0], "azapi_"):
					typeValue := utils.AttributeValue(block.Body().GetAttribute("type"))
					parts := strings.Split(typeValue, "@")
					if len(parts) != 2 {
						continue
					}
					existDeps = append(existDeps, dependency.Dependency{
						AzureResourceType:    parts[0],
						ApiVersion:           parts[1],
						ExampleConfiguration: string(block.BuildTokens(nil).Bytes()),
						ResourceKind:         block.Type(),
						ReferredProperty:     "id",
						ResourceName:         labels[0],
						ResourceLabel:        labels[1],
					})
				case strings.HasPrefix(labels[0], "azurerm_"):
					azureResourceType := azurermTypeAzureResourceTypeMap[labels[0]]
					existDeps = append(existDeps, dependency.Dependency{
						AzureResourceType:    azureResourceType,
						ApiVersion:           "",
						ExampleConfiguration: string(block.BuildTokens(nil).Bytes()),
						ResourceKind:         block.Type(),
						ReferredProperty:     "id",
						ResourceName:         labels[0],
						ResourceLabel:        labels[1],
					})
				}
			}
		}
	}
	return ExistingDependencyResolver{
		ExistingDependencies: existDeps,
	}
}
