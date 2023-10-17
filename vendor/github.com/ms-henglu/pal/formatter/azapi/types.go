package azapi

import (
	"fmt"

	"github.com/ms-henglu/pal/formatter/azapi/hcl"
)

type AzapiDefinition struct {
	Kind              string // resource or data
	ResourceName      string // azapi_resource, azapi_update_resource, azapi_resource_action
	Label             string // example: test
	AzureResourceType string // example: Microsoft.Network/virtualNetworks
	ApiVersion        string // example: 2020-06-01
	Body              interface{}
	Output            interface{}
	ResourceId        string
	Method            string
	AdditionalFields  map[string]Value
}

func (def AzapiDefinition) String() string {
	expressions := fmt.Sprintf(`  type = "%[1]s@%[2]s"`, def.AzureResourceType, def.ApiVersion)
	fields := []string{"resource_id", "parent_id", "name", "location", "action", "method"}
	for _, field := range fields {
		if value, ok := def.AdditionalFields[field]; ok {
			expressions += fmt.Sprintf(`
  %[1]s = %[2]s`, field, value)
		}
	}
	if def.Body != nil {
		bodyMap, ok := def.Body.(map[string]interface{})
		if ok {
			tagRaw, ok := bodyMap["tags"]
			if ok && tagRaw != nil {
				if tagMap, ok := tagRaw.(map[string]interface{}); ok && len(tagMap) == 0 {
					delete(bodyMap, "tags")
				}
			}
		}
		if len(bodyMap) > 0 {
			expressions += fmt.Sprintf(`
  body = jsonencode(%[1]s)`, hcl.MarshalIndent(bodyMap, "  ", "  "))
		}
	}
	return fmt.Sprintf(
		`%[1]s "%[2]s" "%[3]s" {
%[4]s
}
`, def.Kind, def.ResourceName, def.Label, expressions)
}

type Value interface {
	String() string
}

type RawValue struct {
	Raw string
}

func (v RawValue) String() string {
	return v.Raw
}

func NewRawValue(raw string) RawValue {
	return RawValue{
		Raw: raw,
	}
}

type ReferenceValue struct {
	Reference string
}

func (v ReferenceValue) String() string {
	return v.Reference
}

func NewReferenceValue(reference string) ReferenceValue {
	return ReferenceValue{
		Reference: reference,
	}
}

type LiteralValue struct {
	Literal string
}

func (v LiteralValue) String() string {
	return fmt.Sprintf(`"%s"`, v.Literal)
}

func NewLiteralValue(literal string) LiteralValue {
	return LiteralValue{
		Literal: literal,
	}
}

var _ Value = &RawValue{}
var _ Value = &ReferenceValue{}
var _ Value = &LiteralValue{}
