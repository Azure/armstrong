package utils_test

import (
	"encoding/json"
	"github.com/azure/armstrong/utils"
	"testing"
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

func Test_UpdatedBody(t *testing.T) {
	var parameter interface{}
	_ = json.Unmarshal([]byte(inputJson), &parameter)

	replacements := make(map[string]string)
	subnetId := "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1"
	replacements[".properties.properties.subnet"] = subnetId
	output := utils.UpdatedBody(parameter, replacements, "").(map[string]interface{})
	output = output["properties"].(map[string]interface{})
	output = output["properties"].(map[string]interface{})
	if output["subnet"].(string) != subnetId {
		t.Fatalf("expect %s but got %v", subnetId, output["subnet"])
	}
}
