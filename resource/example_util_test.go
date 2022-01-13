package resource_test

import (
	"encoding/json"
	"testing"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/resource"
)

const inputJson = `
{
      "location": "eastus",
      "properties": {
        "computeType": "ComputeInstance",
        "properties": {
          "vmSize": "STANDARD_NC6",
          "subnet": "test-subnet-resource-id",
          "applicationSharingPolicy": "Personal",
          "sshSettings": {
            "sshPublicAccess": "Disabled"
          }
        }
      }
    }`

func Test_GetKeyValueMappings(t *testing.T) {
	var parameter interface{}
	_ = json.Unmarshal([]byte(inputJson), &parameter)
	outputs := resource.GetKeyValueMappings(parameter, "")
	if len(outputs) != 6 {
		t.Fatalf("expect %d mappings, but got %d", 6, len(outputs))
	}
}

func Test_GetUpdatedBody_Removes(t *testing.T) {
	var parameter interface{}
	_ = json.Unmarshal([]byte(inputJson), &parameter)

	removes := []string{".location"}
	output := resource.GetUpdatedBody(parameter, nil, removes, "").(map[string]interface{})
	if output["location"] != nil {
		t.Fatalf("expect nil but got %v", output["location"])
	}
}

func Test_GetUpdatedBody_Replacements(t *testing.T) {
	var parameter interface{}
	_ = json.Unmarshal([]byte(inputJson), &parameter)

	replacements := make(map[string]string)
	subnetId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1"
	replacements[".properties.properties.subnet"] = subnetId
	output := resource.GetUpdatedBody(parameter, replacements, nil, "").(map[string]interface{})
	output = output["properties"].(map[string]interface{})
	output = output["properties"].(map[string]interface{})
	if output["subnet"].(string) != subnetId {
		t.Fatalf("expect %s but got %v", subnetId, output["subnet"])
	}
}

func Test_GetParentIdFromId(t *testing.T) {
	testcases := []struct {
		Input  string
		Expect string
	}{
		{
			Input:  "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123/computes/compute123",
			Expect: "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123",
		},
		{
			Input:  "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123/providers/Microsoft.MachineLearningServices/workspaces/workspaces123",
			Expect: "/subscriptions/34adfa4f-cedf-4dc0-ba29-b6d1a69ab345/resourceGroups/testrg123",
		},
	}

	for _, testcase := range testcases {
		if output := resource.GetParentIdFromId(testcase.Input); output != testcase.Expect {
			t.Fatalf("expect %v but got %v", testcase.Expect, output)
		}
	}
}
