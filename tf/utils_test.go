package tf_test

import (
	"fmt"
	"testing"

	"github.com/ms-henglu/armstrong/tf"
	"github.com/ms-henglu/armstrong/types"
)

func Test_NewErrorReport(t *testing.T) {
	var testcases = []struct {
		Input  error
		Expect []types.Error
	}{
		{
			Input:  nil,
			Expect: []types.Error{},
		},
		{
			Input: fmt.Errorf(`error running terraform apply: exit status 1

Error: performing action  of "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220/buildServices/default/agentPools/acctest1220\" / Api Version \"2023-11-01-preview\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220/buildServices/default/agentPools/acctest1220
--------------------------------------------------------------------------------
RESPONSE 404: 404 Not Found
ERROR CODE: NotFound
--------------------------------------------------------------------------------
{
  "error": {
    "code": "NotFound",
    "message": "build service default's agent pool acctest1220 does not exist",
    "target": "/subscriptions/******/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220/buildServices/default/agentPools/acctest1220",
    "details": null
  }
}
--------------------------------------------------------------------------------


  with azapi_resource_action.put_agentPool,
  on main.tf line 71, in resource "azapi_resource_action" "put_agentPool":
  71: resource "azapi_resource_action" "put_agentPool" {

 `),
			Expect: []types.Error{
				{
					Type:  "Microsoft.AppPlatform/Spring/buildServices/agentPools@2023-11-01-preview",
					Label: "put_agentPool",
					Id:    "/subscriptions/******/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220/buildServices/default/agentPools/acctest1220",
				},
			},
		},
		{
			Input: fmt.Errorf(`error running terraform apply: exit status 1

Error: performing action  of "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest8179/providers/Microsoft.AppPlatform/Spring/acctest8179/eurekaServers/default\" / Api Version \"2023-11-01-preview\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest8179/providers/Microsoft.AppPlatform/Spring/acctest8179/eurekaServers/default
--------------------------------------------------------------------------------
RESPONSE 404: 404 Not Found
ERROR CODE UNAVAILABLE
--------------------------------------------------------------------------------
Response contained no body
--------------------------------------------------------------------------------


  with azapi_resource_action.put_eurekaServer,
  on main.tf line 54, in resource "azapi_resource_action" "put_eurekaServer":
  54: resource "azapi_resource_action" "put_eurekaServer" {`),
			Expect: []types.Error{
				{
					Type:  "Microsoft.AppPlatform/Spring/eurekaServers@2023-11-01-preview",
					Label: "put_eurekaServer",
					Id:    "/subscriptions/******/resourceGroups/acctest8179/providers/Microsoft.AppPlatform/Spring/acctest8179/eurekaServers/default",
				},
			},
		},
	}

	for _, testcase := range testcases {
		actual := tf.NewErrorReport(testcase.Input, nil)
		if len(actual.Errors) != len(testcase.Expect) {
			t.Errorf("Expect %d errors, but got %d", len(testcase.Expect), len(actual.Errors))
			continue
		}
		for i := 0; i < len(actual.Errors); i++ {
			if actual.Errors[i].Type != testcase.Expect[i].Type {
				t.Errorf("Expect error type %s, but got %s", testcase.Expect[i].Type, actual.Errors[i].Type)
			}
			if actual.Errors[i].Id != testcase.Expect[i].Id {
				t.Errorf("Expect error id %s, but got %s", testcase.Expect[i].Id, actual.Errors[i].Id)
			}
			if actual.Errors[i].Label != testcase.Expect[i].Label {
				t.Errorf("Expect error label %s, but got %s", testcase.Expect[i].Label, actual.Errors[i].Label)
			}
		}
	}
}
