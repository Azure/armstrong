package resource_test

import (
	"reflect"
	"testing"

	"github.com/ms-henglu/armstrong/resource"
	"github.com/ms-henglu/armstrong/resource/types"
)

func Test_NewAzapiDefinitionFromExample(t *testing.T) {
	testcases := []struct {
		Input       string
		Kind        string
		Want        types.AzapiDefinition
		ExpectError bool
	}{
		{
			Input: "testdata/examples/createOrUpdateAutomationAccount.json",
			Kind:  "resource",
			Want: types.AzapiDefinition{
				AzureResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:        "2022-08-08",
				BodyFormat:        types.BodyFormatHcl,
				Body: map[string]interface{}{
					"properties": map[string]interface{}{
						"sku": map[string]interface{}{
							"name": "Free",
						},
					},
				},
				AdditionalFields: map[string]types.Value{
					"parent_id":                 types.NewStringLiteralValue("/subscriptions/subid/resourceGroups/rg"),
					"name":                      types.NewStringLiteralValue("myAutomationAccount9"),
					"location":                  types.NewStringLiteralValue("East US 2"),
					"schema_validation_enabled": types.NewRawValue("false"),
					"ignore_casing":             types.NewRawValue("false"),
					"ignore_missing_property":   types.NewRawValue("false"),
				},
				ResourceName: "azapi_resource",
				Label:        "automationAccount",
				Kind:         types.KindResource,
			},
			ExpectError: false,
		},
		{
			Input: "testdata/examples/getAutomationAccount.json",
			Kind:  "data",
			Want: types.AzapiDefinition{
				AzureResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:        "2022-08-08",
				BodyFormat:        types.BodyFormatHcl,
				Body:              nil,
				AdditionalFields: map[string]types.Value{
					"parent_id": types.NewStringLiteralValue("/subscriptions/subid/resourceGroups/rg"),
					"name":      types.NewStringLiteralValue("myAutomationAccount9"),
				},
				ResourceName: "azapi_resource",
				Label:        "automationAccount",
				Kind:         types.KindDataSource,
			},
		},
	}

	for _, tc := range testcases {
		t.Logf("test case: %s", tc.Input)
		got, err := resource.NewAzapiDefinitionFromExample(tc.Input, tc.Kind)
		if tc.ExpectError != (err != nil) {
			t.Errorf("got error %v, want error %v", err, tc.ExpectError)
			continue
		}
		if got.AzureResourceType != tc.Want.AzureResourceType {
			t.Errorf("got %s, want %s", got.AzureResourceType, tc.Want.AzureResourceType)
		}
		if got.ApiVersion != tc.Want.ApiVersion {
			t.Errorf("got %s, want %s", got.ApiVersion, tc.Want.ApiVersion)
		}
		if got.BodyFormat != tc.Want.BodyFormat {
			t.Errorf("got %s, want %s", got.BodyFormat, tc.Want.BodyFormat)
		}
		if !reflect.DeepEqual(got.Body, tc.Want.Body) {
			t.Errorf("got %v, want %v", got.Body, tc.Want.Body)
		}
		if len(got.AdditionalFields) != len(tc.Want.AdditionalFields) {
			t.Errorf("got %v, want %v", got.AdditionalFields, tc.Want.AdditionalFields)
			continue
		}
		for key, value := range got.AdditionalFields {
			if tc.Want.AdditionalFields[key] != value {
				t.Errorf("got %v, want %v", got.AdditionalFields, tc.Want.AdditionalFields)
			}
		}
		if got.ResourceName != tc.Want.ResourceName {
			t.Errorf("got %s, want %s", got.ResourceName, tc.Want.ResourceName)
		}
		if got.Label != tc.Want.Label {
			t.Errorf("got %s, want %s", got.Label, tc.Want.Label)
		}
		if got.Kind != tc.Want.Kind {
			t.Errorf("got %s, want %s", got.Kind, tc.Want.Kind)
		}
	}
}
