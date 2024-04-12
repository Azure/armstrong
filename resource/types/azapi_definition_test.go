package types_test

import (
	"testing"

	"github.com/azure/armstrong/resource/types"
)

func Test_AzapiDefinitionString(t *testing.T) {
	testcases := []struct {
		Input types.AzapiDefinition
		Want  string
	}{
		{
			Input: types.AzapiDefinition{
				Kind:              types.KindResource,
				ResourceName:      "azapi_resource",
				AzureResourceType: "Microsoft.Network/virtualNetworks",
				ApiVersion:        "2020-06-01",
				Label:             "test",
				Body: map[string]interface{}{
					"properties": map[string]interface{}{
						"addressSpace": map[string]interface{}{
							"addressPrefixes": []interface{}{
								"192.0.0.0/16",
							},
						},
					},
				},
				BodyFormat: types.BodyFormatHcl,
				AdditionalFields: map[string]types.Value{
					"parent_id": types.NewStringLiteralValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test"),
					"name":      types.NewStringLiteralValue("test"),
				},
			},
			Want: `resource "azapi_resource" "test" {
  type = "Microsoft.Network/virtualNetworks@2020-06-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test"
  name = "test"
  body = jsonencode({
    properties = {
      addressSpace = {
        addressPrefixes = [
          "192.0.0.0/16",
        ]
      }
    }
  })
}
`,
		},
	}

	for _, tc := range testcases {
		got := tc.Input.String()
		if got != tc.Want {
			t.Errorf("got %s, want %s", got, tc.Want)
		}
	}
}
