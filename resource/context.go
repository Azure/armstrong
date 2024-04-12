package resource

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/azure/armstrong/dependency"
	"github.com/azure/armstrong/resource/resolver"
	"github.com/azure/armstrong/resource/types"
	"github.com/azure/armstrong/utils"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
)

type Context struct {
	File               *hclwrite.File
	locationVarBlock   *hclwrite.Block
	terraformBlock     *hclwrite.Block
	KnownPatternMap    map[string]types.Reference
	ReferenceResolvers []resolver.ReferenceResolver
	azapiAddingMap     map[string]bool
}

const DefaultProviderConfig = `terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

`

func NewContext(referenceResolvers []resolver.ReferenceResolver) *Context {
	knownPatternMap := make(map[string]types.Reference)
	referenceResolvers = append([]resolver.ReferenceResolver{resolver.NewKnownReferenceResolver(knownPatternMap)}, referenceResolvers...)
	c := Context{
		KnownPatternMap:    knownPatternMap,
		ReferenceResolvers: referenceResolvers,
		azapiAddingMap:     make(map[string]bool),
	}
	err := c.InitFile(DefaultProviderConfig)
	if err != nil {
		logrus.Errorf("failed to init context: %v", err)
	}
	return &c
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func (c *Context) InitFile(content string) error {
	file, diags := hclwrite.ParseConfig([]byte(content), "", hcl.InitialPos)
	if diags.HasErrors() {
		logrus.Errorf("failed to parse input:\n%v", content)
		return diags
	}
	var locationVarBlock, terraformBlock *hclwrite.Block
	for _, block := range file.Body().Blocks() {
		switch block.Type() {
		case "variable":
			switch block.Labels()[0] {
			case "location":
				locationVarBlock = block
			case "resource_name":
				block.Body().SetAttributeValue("default", cty.StringVal(fmt.Sprintf("acctest%04d", r.Intn(10000))))
			}
		case "terraform":
			terraformBlock = block
		}
	}
	if terraformBlock == nil {
		logrus.Warnf("terraform block not found in the input.")
	}
	if locationVarBlock == nil {
		logrus.Warnf("location variable block not found in the input.")
	}
	c.File = file
	c.terraformBlock = terraformBlock
	c.locationVarBlock = locationVarBlock
	return nil
}

func (c *Context) AddAzapiDefinition(input types.AzapiDefinition) error {
	logrus.Debugf("adding azapi definition: \n%v", input)
	if c.azapiAddingMap[input.Identifier()] {
		return fmt.Errorf("azapi definition already added: %v", input.Identifier())
	}
	c.azapiAddingMap[input.Identifier()] = true
	defer func() {
		c.azapiAddingMap[input.Identifier()] = false
	}()
	def := input.DeepCopy()
	// find all id placeholders from def
	placeHolders := make([]types.PropertyDependencyMapping, 0)
	rootFields := []string{"parent_id", "resource_id"}
	for _, field := range rootFields {
		if value, ok := def.AdditionalFields[field]; ok {
			if literalValue, ok := value.(types.StringLiteralValue); ok {
				placeHolders = append(placeHolders, types.PropertyDependencyMapping{
					ValuePath:    field,
					LiteralValue: literalValue.Literal,
				})
			}
		}
	}
	if def.Body != nil {
		mappings := GetKeyValueMappings(def.Body, "")
		for _, mapping := range mappings {
			if utils.IsResourceId(mapping.LiteralValue) {
				placeHolders = append(placeHolders, mapping)
			}
		}
	}
	logrus.Debugf("found %d id placeholders", len(placeHolders))

	// find all dependencies that match the id placeholders
	for i, placeHolder := range placeHolders {
		logrus.Debugf("processing id placeholder: %s", placeHolder.LiteralValue)
		if utils.IsAction(placeHolder.LiteralValue) {
			logrus.Debugf("skip action: %s", placeHolder.LiteralValue)
			continue
		}

		pattern := dependency.NewPattern(placeHolder.LiteralValue)

		for _, resolver := range c.ReferenceResolvers {
			result, err := resolver.Resolve(pattern)
			if err != nil {
				return err
			}
			if result == nil {
				continue
			}
			logrus.Debugf("found dependency by resolver: %T", resolver)
			switch {
			case result.Reference.IsKnown():
				placeHolders[i].Reference = result.Reference
				logrus.Debugf("dependency resolved, ref: %v", result.Reference)
			case result.HclToAdd != "":
				logrus.Debugf("found dependency:\n %v", result.HclToAdd)
				ref, err := c.AddHcl(result.HclToAdd, true)
				if err != nil {
					logrus.Warnf("failed to add hcl as a dependency, will continue to try other dependency resolvers: %v", err)
					continue
				}
				c.KnownPatternMap[pattern.String()] = *ref
				placeHolders[i].Reference = ref
				logrus.Debugf("dependency resolved, ref: %v", ref)
			case result.AzapiDefinitionToAdd != nil:
				logrus.Debugf("found dependency:\n %v", result.AzapiDefinitionToAdd)
				err = c.AddAzapiDefinition(*result.AzapiDefinitionToAdd)
				if err != nil {
					logrus.Warnf("failed to add azapi definition as a dependency, will continue to try other dependency resolvers: %v", err)
					continue
				}
				ref := c.KnownPatternMap[pattern.String()]
				if !ref.IsKnown() {
					return fmt.Errorf("resource type address not found: %v after adding azapi definition to the context, azapi def: %v", pattern, result.AzapiDefinitionToAdd)
				}
				placeHolders[i].Reference = &ref
				logrus.Debugf("dependency resolved, ref: %v", ref)
			}
			break
		}
	}

	// replace the id placeholders with the dependency address
	logrus.Debugf("replacing id placeholders with dependency address...")
	for _, filed := range rootFields {
		for _, placeHolder := range placeHolders {
			if placeHolder.ValuePath == filed && placeHolder.Reference.IsKnown() {
				def.AdditionalFields[filed] = types.NewReferenceValue(placeHolder.Reference.String())
				break
			}
		}
	}
	if def.Body != nil {
		replacements := make(map[string]string)
		for _, placeHolder := range placeHolders {
			if !placeHolder.Reference.IsKnown() {
				continue
			}
			valuePath := placeHolder.ValuePath
			if placeHolder.IsKey {
				valuePath = fmt.Sprintf("key:%s", placeHolder.ValuePath)
			}
			replacements[valuePath] = fmt.Sprintf(`${%s}`, placeHolder.Reference)
		}
		def.Body = utils.UpdatedBody(def.Body, replacements, "")
	}

	// add extra dependencies
	if def.ResourceName == "azapi_resource_list" {
		logrus.Debugf("adding extra dependencies for azapi_resource_list...")
		var ref *types.Reference
		for pattern, r := range c.KnownPatternMap {
			if strings.HasSuffix(pattern, strings.ToLower(":"+def.AzureResourceType)) {
				ref = &r
				break
			}
		}
		if ref.IsKnown() {
			addr := fmt.Sprintf(`%s.%s`, ref.Name, ref.Label)
			if ref.Kind == "data" {
				addr = fmt.Sprintf(`data.%s.%s`, ref.Name, ref.Label)
			}
			def.AdditionalFields["depends_on"] = types.NewRawValue(fmt.Sprintf(`[%s]`, addr))
		} else {
			logrus.Debugf("no `depends_on` dependency found for azapi_resource_list: %v", def)
		}
	}

	ref, err := c.AddHcl(def.String(), false)
	if err != nil {
		return err
	}
	if def.AdditionalFields["action"] == nil && def.ResourceName != "azapi_resource_list" {
		pattern := dependency.NewPattern(def.Id)
		c.KnownPatternMap[pattern.String()] = *ref
		logrus.Debugf("adding known pattern: %s, ref: %s", pattern, *ref)
	}
	return nil
}

func (c *Context) AddHcl(input string, skipWhenDuplicate bool) (*types.Reference, error) {
	logrus.Debugf("adding hcl:\n%v, skipWhenDuplicate: %v", input, skipWhenDuplicate)
	inputFile, diags := hclwrite.ParseConfig([]byte(input), "", hcl.InitialPos)
	if diags.HasErrors() {
		logrus.Warnf("failed to parse input:\n%v", input)
		return nil, diags
	}

	labelsMap := make(map[string]map[string]*hclwrite.Block)
	for _, block := range c.File.Body().Blocks() {
		if block.Type() != "data" && block.Type() != "resource" {
			continue
		}
		labels := block.Labels()
		if len(labels) != 2 {
			return nil, fmt.Errorf("label is invalid: %v, input:\n%v", labels, input)
		}
		key := fmt.Sprintf("%s.%s", block.Type(), labels[0])
		if labelsMap[key] == nil {
			labelsMap[key] = make(map[string]*hclwrite.Block)
		}
		labelsMap[key][labels[1]] = block
	}

	// resource/data blocks with same resource name and labels will have conflicts,
	// merge them if their resource types are the same,
	// otherwise, rename the labels
	for _, block := range inputFile.Body().Blocks() {
		if block.Type() != "data" && block.Type() != "resource" {
			continue
		}
		labels := block.Labels()
		if len(labels) != 2 {
			return nil, fmt.Errorf("label is invalid: %v, input:\n%v", labels, input)
		}
		key := fmt.Sprintf("%s.%s", block.Type(), labels[0])
		// no conflict
		conflictBlock := labelsMap[key][labels[1]]
		if conflictBlock == nil {
			continue
		}
		logrus.Debugf("found conflict: %s %v", conflictBlock.Type(), conflictBlock.Labels())
		// if the resource types are the same, there is no conflict
		if utils.TypeValue(block) == utils.TypeValue(conflictBlock) && skipWhenDuplicate {
			continue
		}
		// conflict, rename the labels
		newLabel := labels[1]
		for i := 1; i < 100; i++ {
			newLabel = fmt.Sprintf("%s_%d", labels[1], i)
			if labelsMap[key][newLabel] == nil {
				break
			}
		}
		logrus.Debugf("renaming labels: %v -> %v", labels[1], newLabel)
		block.SetLabels([]string{labels[0], newLabel})
		input = string(inputFile.BuildTokens(nil).Bytes())

		// TODO: improve the following renaming labels logic
		oldAddressPrefix := fmt.Sprintf("%s.", strings.Join(labels, "."))
		if block.Type() == "data" {
			oldAddressPrefix = "data." + oldAddressPrefix
		}
		newAddressPrefix := fmt.Sprintf("%s.", strings.Join(block.Labels(), "."))
		if block.Type() == "data" {
			newAddressPrefix = "data." + newAddressPrefix
		}
		input = strings.ReplaceAll(input, oldAddressPrefix, newAddressPrefix)
	}

	inputFile, diags = hclwrite.ParseConfig([]byte(input), "", hcl.InitialPos)
	if diags.HasErrors() {
		logrus.Warnf("failed to parse input:\n%v", input)
		return nil, diags
	}

	// update the location and name fields
	for _, block := range inputFile.Body().Blocks() {
		if block.Type() != "data" && block.Type() != "resource" {
			continue
		}
		locationAttr := block.Body().GetAttribute("location")
		if locationAttr != nil {
			defaultLocation := utils.AttributeValue(c.locationVarBlock.Body().GetAttribute("default"))
			currentLocation := utils.AttributeValue(locationAttr)
			if currentLocation != defaultLocation && !strings.Contains(currentLocation, "var.") {
				c.locationVarBlock.Body().SetAttributeValue("default", cty.StringVal(currentLocation))
				block.Body().SetAttributeTraversal("location", hcl.Traversal{hcl.TraverseRoot{Name: "var"}, hcl.TraverseAttr{Name: "location"}})
			}
		}

		//TODO: replace location value in the body payload

		nameAttr := block.Body().GetAttribute("name")
		if nameAttr != nil {
			currentName := utils.AttributeValue(nameAttr)
			if isRandomName(currentName) {
				block.Body().SetAttributeTraversal("name", hcl.Traversal{hcl.TraverseRoot{Name: "var"}, hcl.TraverseAttr{Name: "resource_name"}})
			}
		}
	}

	varMap := make(map[string]bool)
	providerMap := make(map[string]bool)
	for _, block := range c.File.Body().Blocks() {
		switch block.Type() {
		case "provider":
			providerMap[strings.Join(block.Labels(), ".")] = true
		case "variable":
			varMap[strings.Join(block.Labels(), ".")] = true
		}
	}

	var lastBlock *hclwrite.Block
	for _, block := range inputFile.Body().Blocks() {
		switch block.Type() {
		case "terraform":
			newProvidersBlock := block.Body().FirstMatchingBlock("required_providers", []string{})
			if newProvidersBlock == nil {
				continue
			}
			oldProvidersBlock := c.terraformBlock.Body().FirstMatchingBlock("required_providers", []string{})
			if oldProvidersBlock != nil {
				for attrName, attr := range newProvidersBlock.Body().Attributes() {
					if oldProvidersBlock.Body().GetAttribute(attrName) == nil {
						oldProvidersBlock.Body().SetAttributeRaw(attrName, attr.Expr().BuildTokens(nil))
					}
				}
			} else {
				logrus.Errorf("required_providers block not found in the input.")
			}
			continue
		case "variable":
			label := strings.Join(block.Labels(), ".")
			if varMap[label] {
				continue
			}
			c.File.Body().AppendBlock(block)
		case "provider":
			label := strings.Join(block.Labels(), ".")
			if providerMap[label] {
				continue
			}
			c.File.Body().AppendBlock(block)
		case "output":
			continue
		case "locals":
			c.File.Body().AppendBlock(block)
			c.File.Body().AppendNewline()
		case "data", "resource":
			labels := block.Labels()
			if len(labels) != 2 {
				return nil, fmt.Errorf("label is invalid: %v, input:\n%v", labels, input)
			}
			key := fmt.Sprintf("%s.%s", block.Type(), labels[0])
			conflictBlock := labelsMap[key][labels[1]]
			if conflictBlock != nil {
				lastBlock = conflictBlock
				continue
			}
			c.File.Body().AppendBlock(block)
			c.File.Body().AppendNewline()
			lastBlock = block
		default:
			c.File.Body().AppendBlock(block)
			c.File.Body().AppendNewline()
		}
	}

	if lastBlock != nil {
		labels := lastBlock.Labels()
		if len(labels) != 2 {
			return nil, fmt.Errorf("label is invalid: %v, input:\n%v", labels, input)
		}
		return &types.Reference{
			Label:    labels[1],
			Kind:     lastBlock.Type(),
			Name:     labels[0],
			Property: "id",
		}, nil
	}
	return nil, fmt.Errorf("no resource or data block found in the input, input:\n%v", input)
}

func (c *Context) String() string {
	return string(hclwrite.Format(c.File.Bytes()))
}

// GetKeyValueMappings returns a list of key and value of input
func GetKeyValueMappings(parameters interface{}, path string) []types.PropertyDependencyMapping {
	if parameters == nil {
		return []types.PropertyDependencyMapping{}
	}
	results := make([]types.PropertyDependencyMapping, 0)
	switch param := parameters.(type) {
	case map[string]interface{}:
		for key, value := range param {
			results = append(results, GetKeyValueMappings(value, path+"."+key)...)
			results = append(results, types.PropertyDependencyMapping{
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
		results = append(results, types.PropertyDependencyMapping{
			ValuePath:    path,
			LiteralValue: param,
			IsKey:        false,
		})
	default:

	}
	return results
}

func isRandomName(input string) bool {
	if input == "default" {
		return false
	}
	if input == "current" {
		return false
	}
	if strings.Contains(input, "Microsoft.") {
		return false
	}
	if strings.Contains(input, "var.") {
		return false
	}
	return true
}
