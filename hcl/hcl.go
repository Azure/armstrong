package hcl

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/ms-henglu/armstrong/types"
	"github.com/zclconf/go-cty/cty"
)

// RenameLabel is used to rename resource name to make name unique
func RenameLabel(input string) string {
	f, parseDiags := hclwrite.ParseConfig([]byte(input), "temp.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Println(parseDiags.Error())
		return ""
	}
	resourceCountMap := make(map[string]int)
	resourceNameMap := make(map[string]string)
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
			block.Body().SetAttributeValue("name", cty.StringVal(RandomName()))
		}
		resHcl.Body().AppendBlock(block)
	}

	labelRenamedHcl := string(resHcl.BuildTokens(nil).Bytes())
	for key, value := range resourceNameMap {
		labelRenamedHcl = strings.ReplaceAll(labelRenamedHcl, key, value)
	}
	return string(hclwrite.Format([]byte(labelRenamedHcl)))
}

// Combine is used to merge hcl and avoid create duplicate resources
func Combine(old, new string) string {
	oldHcl, parseDiags := hclwrite.ParseConfig([]byte(old), "old.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Println(parseDiags.Error())
		return ""
	}
	newHcl, parseDiags := hclwrite.ParseConfig([]byte(new), "new.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Println(parseDiags.Error())
		return ""
	}
	resHcl, _ := hclwrite.ParseConfig([]byte(""), "res.tf", hcl.InitialPos)

	blocks := make(map[string]hclwrite.Block)
	for _, block := range oldHcl.Body().Blocks() {
		labels := block.Labels()
		resourceName := strings.Join(labels, ".")
		blocks[resourceName] = *block
		resHcl.Body().AppendBlock(block)
		resHcl.Body().AppendNewline()
	}
	for _, block := range newHcl.Body().Blocks() {
		labels := block.Labels()
		resourceName := labels[0] + "." + labels[1]
		if _, ok := blocks[resourceName]; ok {
			// TODO: check whether exist and block are equal, if not, they must both exist??
		} else {
			resHcl.Body().AppendBlock(block)
			resHcl.Body().AppendNewline()
		}
	}

	return string(hclwrite.Format([]byte(resHcl.BuildTokens(nil).Bytes())))
}

// FindResourceAddress returns first resource address which is resourceType from config
func FindResourceAddress(config, resourceType string) string {
	f, parseDiags := hclwrite.ParseConfig([]byte(config), "old.tf", hcl.InitialPos)
	if parseDiags != nil && parseDiags.HasErrors() {
		log.Println(parseDiags.Error())
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

func RandomName() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("acctest%d", rand.Intn(10000))
}

func LoadExistingDependencies() []types.Dependency {
	dir, _ := os.Getwd()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("[WARN] reading dir %s: %+v", dir, err)
		return nil
	}
	existDeps := make([]types.Dependency, 0)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".tf") {
			continue
		}
		src, err := ioutil.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			log.Printf("[WARN] reading file %s: %+v", file.Name(), err)
			continue
		}
		f, diag := hclwrite.ParseConfig(src, file.Name(), hcl.InitialPos)
		if diag.HasErrors() {
			log.Printf("[WARN] parsing file %s: %+v", file.Name(), diag.Error())
			continue
		}
		if f == nil || f.Body() == nil {
			continue
		}
		for _, block := range f.Body().Blocks() {
			labels := block.Labels()
			if len(labels) >= 2 {
				pattern := ""
				if labels[0] == "azapi_resource" {
					pattern = GetAzApiResourceIdPattern(block)
				}
				existDeps = append(existDeps, types.Dependency{
					Pattern:          pattern,
					ResourceType:     labels[0],
					ReferredProperty: "id",
					Address:          strings.Join(labels, "."),
				})
			}
		}
	}
	return existDeps
}

func GetAzApiResourceIdPattern(block *hclwrite.Block) string {
	if block == nil || block.Body() == nil {
		return ""
	}
	attribute := block.Body().GetAttribute("type")
	if attribute == nil || attribute.Expr() == nil {
		return ""
	}
	typeValue := string(attribute.Expr().BuildTokens(nil).Bytes())
	typeValue = strings.Trim(typeValue, ` "`)
	typeValue = typeValue[0:strings.Index(typeValue, "@")]
	return fmt.Sprintf("/subscriptions/resourceGroups/providers/%s", typeValue)
}

const ProviderHcl = `
terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
  skip_provider_registration = false
}
`
