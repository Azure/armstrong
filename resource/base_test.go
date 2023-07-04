package resource_test

import (
	"encoding/json"
	"testing"

	"github.com/ms-henglu/armstrong/resource"
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
	if len(outputs) != 15 {
		t.Fatalf("expect %d mappings, but got %d", 15, len(outputs))
	}
}
