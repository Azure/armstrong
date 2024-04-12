package resolver

import (
	"os"
	"path"
	"strings"

	"github.com/azure/armstrong/dependency"
	"github.com/azure/armstrong/resource/types"
	"github.com/azure/armstrong/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
)

var _ ReferenceResolver = &ExistingDependencyResolver{}

type ExistingDependencyResolver struct {
	ExistingDependencies []dependency.Dependency
}

func (r ExistingDependencyResolver) Resolve(pattern dependency.Pattern) (*ResolvedResult, error) {
	for _, dep := range r.ExistingDependencies {
		if strings.EqualFold(dep.AzureResourceType, pattern.AzureResourceType) {
			return &ResolvedResult{
				Reference: &types.Reference{
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
			logrus.Debugf("empty file %s", file.Name())
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
						logrus.Warnf("invalid type value %s (labels: %v, filename: %s)", typeValue, labels, file.Name())
						continue
					}
					logrus.Debugf("found existing azapi dependency: %s (labels: %v, filename: %s)", parts[0], labels, file.Name())
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
					logrus.Debugf("found existing azurerm dependency: %s (labels: %v, filename: %s)", azureResourceType, labels, file.Name())
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
	logrus.Infof("found %d existing dependencies", len(existDeps))
	return ExistingDependencyResolver{
		ExistingDependencies: existDeps,
	}
}
