package utils

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TypeValue(block *hclwrite.Block) string {
	typeAttribute := block.Body().GetAttribute("type")
	return AttributeValue(typeAttribute)
}

func AttributeValue(attribute *hclwrite.Attribute) string {
	if attribute == nil {
		return ""
	}
	value := string(attribute.Expr().BuildTokens(nil).Bytes())
	value = strings.Trim(value, ` "`)
	return value
}
