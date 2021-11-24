package helper

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// GetRenamedHcl is used to rename resource name to make name unique
func GetRenamedHcl(input string) string {
	f, parseDiags := hclwrite.ParseConfig([]byte(input), "temp.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Printf(parseDiags.Error())
		return ""
	}
	resourceCountMap := make(map[string]int, 0)
	resourceNameMap := make(map[string]string, 0)
	resHcl, _ := hclwrite.ParseConfig([]byte(""), "res.tf", hcl.InitialPos)
	for _, block := range f.Body().Blocks() {
		// remove terraform and provider blocks
		if block.Type() == "terraform" || block.Type() == "provider" || block.Type() == "output" {
			continue
		}
		if block.Type() == "variable" {
			resHcl.Body().AppendBlock(block)
			continue
		}
		labels := block.Labels()
		count := resourceCountMap[labels[0]]
		resourceName := labels[0] + "." + labels[1]
		if count == 0 {
			resourceNameMap[resourceName] = labels[0] + "." + "test"
			labels[1] = "test"
		} else {
			resourceNameMap[resourceName] = labels[0] + ".test" + strconv.Itoa(count+1)
			labels[1] = "test" + strconv.Itoa(count+1)
		}
		resourceCountMap[labels[0]] = count + 1
		block.SetLabels(labels)
		if block.Body() != nil && block.Body().GetAttribute("name") != nil {
			rand.Seed(time.Now().UnixNano())
			block.Body().SetAttributeValue("name", cty.StringVal(GetRandomResourceName()))
		}
		resHcl.Body().AppendBlock(block)
	}

	labelRenamedHcl := string(resHcl.BuildTokens(nil).Bytes())
	for key, value := range resourceNameMap {
		labelRenamedHcl = strings.ReplaceAll(labelRenamedHcl, key, value)
	}
	return string(hclwrite.Format([]byte(labelRenamedHcl)))
}

// GetCombinedHcl is used to merge hcl and avoid create duplicate resources
func GetCombinedHcl(old, new string) string {
	oldHcl, parseDiags := hclwrite.ParseConfig([]byte(old), "old.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Printf(parseDiags.Error())
		return ""
	}
	newHcl, parseDiags := hclwrite.ParseConfig([]byte(new), "new.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Printf(parseDiags.Error())
		return ""
	}
	resHcl, _ := hclwrite.ParseConfig([]byte(""), "res.tf", hcl.InitialPos)

	blocks := make(map[string]hclwrite.Block, 0)
	for _, block := range oldHcl.Body().Blocks() {
		labels := block.Labels()
		resourceName := strings.Join(labels, ".")
		blocks[resourceName] = *block
		resHcl.Body().AppendBlock(block)
	}
	for _, block := range newHcl.Body().Blocks() {
		labels := block.Labels()
		resourceName := labels[0] + "." + labels[1]
		if _, ok := blocks[resourceName]; ok {
			// TODO: check whether exist and block are equal, if not, they must both exist??
		} else {
			resHcl.Body().AppendBlock(block)
		}
	}

	return string(hclwrite.Format([]byte(resHcl.BuildTokens(nil).Bytes())))
}

func GetResourceFromHcl(config, resourceType string) string {
	f, parseDiags := hclwrite.ParseConfig([]byte(config), "old.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Printf(parseDiags.Error())
		return ""
	}
	if f == nil || f.Body() == nil {
		return ""
	}
	for _, block := range f.Body().Blocks() {
		labels := block.Labels()
		if len(labels) >= 2 && labels[0] == resourceType {
			return labels[0] + "." + labels[1]
		}
	}
	return ""
}

func GetRandomResourceName() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("acctest%d", rand.Intn(10000))
}

const ProviderHcl = `
terraform {
  required_providers {
    azurerm-restapi = {
      source  = "Azure/azurerm-restapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azurerm-restapi" {
  schema_validation_enabled = false
}
`
