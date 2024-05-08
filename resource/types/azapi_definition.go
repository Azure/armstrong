package types

import (
	"encoding/json"
	"fmt"
	"github.com/azure/armstrong/hcl"
)

type AzapiDefinition struct {
	Id                string
	Kind              Kind   // resource or data
	ResourceName      string // azapi_resource, azapi_update_resource, azapi_resource_action
	Label             string // example: test
	AzureResourceType string // example: Microsoft.Network/virtualNetworks
	ApiVersion        string // example: 2020-06-01
	Body              interface{}
	AdditionalFields  map[string]Value // fields like resource_id, parent_id, name, location, action, method
	BodyFormat        BodyFormat       // hcl or json
	LeadingComments   []string
}

type BodyFormat string

const (
	BodyFormatHcl  BodyFormat = "hcl"
	BodyFormatJson BodyFormat = "json"
)

type Kind string

const (
	KindDataSource Kind = "data"
	KindResource   Kind = "resource"
)

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
			if def.BodyFormat == BodyFormatJson {
				jsonBody, _ := json.MarshalIndent(bodyMap, "", "    ")
				expressions += fmt.Sprintf(`<<BODY
%s
BODY`, jsonBody)
			} else {
				expressions += fmt.Sprintf(`
  body = %[1]s`, hcl.MarshalIndent(bodyMap, "  ", "  "))
			}
		}
	}
	for _, field := range []string{"schema_validation_enabled", "ignore_casing", "ignore_missing_property", "depends_on"} {
		if value, ok := def.AdditionalFields[field]; ok {
			expressions += fmt.Sprintf(`
  %[1]s = %[2]s`, field, value)
		}
	}
	config := fmt.Sprintf(
		`%[1]s "%[2]s" "%[3]s" {
%[4]s
}
`, def.Kind, def.ResourceName, def.Label, expressions)
	if len(def.LeadingComments) > 0 {
		comment := ""
		for _, line := range def.LeadingComments {
			comment += fmt.Sprintf("// %s\n", line)
		}
		config = comment + config
	}
	return config
}

func (def AzapiDefinition) DeepCopy() AzapiDefinition {
	additionalFields := make(map[string]Value)
	for k, v := range def.AdditionalFields {
		additionalFields[k] = v.DeepCopy()
	}

	leadingComments := make([]string, len(def.LeadingComments))
	copy(leadingComments, def.LeadingComments)
	return AzapiDefinition{
		Id:                def.Id,
		Kind:              def.Kind,
		ResourceName:      def.ResourceName,
		Label:             def.Label,
		AzureResourceType: def.AzureResourceType,
		ApiVersion:        def.ApiVersion,
		Body:              def.Body,
		AdditionalFields:  additionalFields,
		LeadingComments:   leadingComments,
	}
}

func (def AzapiDefinition) Identifier() string {
	return fmt.Sprintf("%s-%s-%s", def.Kind, def.ResourceName, def.Id)
}
