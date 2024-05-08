package tf_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/azure/armstrong/tf"
	"github.com/azure/armstrong/types"
	tfjson "github.com/hashicorp/terraform-json"
)

func Test_GetChanges(t *testing.T) {
	var testcases = []struct {
		Input  *tfjson.Plan
		Expect []tf.Action
	}{
		{
			Input:  nil,
			Expect: []tf.Action{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionCreate},
						},
					},
				},
			},
			Expect: []tf.Action{tf.ActionCreate},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionRead},
						},
					},
				},
			},
			Expect: []tf.Action{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionNoop},
						},
					},
				},
			},
			Expect: []tf.Action{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionDelete},
						},
					},
				},
			},
			Expect: []tf.Action{tf.ActionDelete},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
						},
					},
				},
			},
			Expect: []tf.Action{tf.ActionUpdate},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionDelete, tfjson.ActionCreate},
						},
					},
				},
			},
			Expect: []tf.Action{tf.ActionReplace},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Change: &tfjson.Change{
							Actions: []tfjson.Action{},
						},
					},
				},
			},
			Expect: []tf.Action{},
		},
	}

	for _, testcase := range testcases {
		actual := tf.GetChanges(testcase.Input)
		if len(actual) != len(testcase.Expect) {
			t.Errorf("Expect %d changes, but got %d", len(testcase.Expect), len(actual))
			continue
		}
		for i := 0; i < len(actual); i++ {
			if actual[i] != testcase.Expect[i] {
				t.Errorf("Expect change %s, but got %s", testcase.Expect[i], actual[i])
			}
		}
	}
}

func Test_NewDiffReport(t *testing.T) {
	var testcases = []struct {
		Input  *tfjson.Plan
		Expect types.DiffReport
	}{
		{
			Input:  nil,
			Expect: types.DiffReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{},
				},
			},
			Expect: types.DiffReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azurerm_resource_group.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
							Before:  "before",
							After:   "after",
						},
					},
				},
			},
			Expect: types.DiffReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionNoop},
							Before:  "foo",
							After:   "foo",
						},
					},
				},
			},
			Expect: types.DiffReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
							Before:  "foo",
							After:   "foo",
						},
					},
				},
			},
			Expect: types.DiffReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
							Before: map[string]interface{}{
								"type": "Microsoft.AppPlatform/Spring",
								"body": map[string]interface{}{
									"foo": "bar",
								},
							},
							After: map[string]interface{}{
								"type": "Microsoft.AppPlatform/Spring",
								"body": map[string]interface{}{
									"foo": "bar",
								},
							},
						},
					},
				},
			},
			Expect: types.DiffReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
							Before: map[string]interface{}{
								"id":   "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220",
								"type": "Microsoft.AppPlatform/Spring",
								"body": `{"foo": "bar"}`,
							},
							After: map[string]interface{}{
								"id":   "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220",
								"type": "Microsoft.AppPlatform/Spring",
								"body": `{"foo": "after"}`,
							},
						},
					},
				},
			},
			Expect: types.DiffReport{
				Diffs: []types.Diff{
					{
						Id:      "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220",
						Type:    "Microsoft.AppPlatform/Spring",
						Address: "azapi_resource.test",
						Change: types.Change{
							Before: `{"foo": "bar"}`,
							After:  `{"foo": "after"}`,
						},
					},
				},
			},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
							Before: map[string]interface{}{
								"id":   "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220",
								"type": "Microsoft.AppPlatform/Spring",
								"body": map[string]interface{}{
									"foo": "bar",
								},
							},
							After: map[string]interface{}{
								"id":   "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220",
								"type": "Microsoft.AppPlatform/Spring",
								"body": map[string]interface{}{
									"foo": "after",
								},
							},
						},
					},
				},
			},
			Expect: types.DiffReport{
				Diffs: []types.Diff{
					{
						Id:      "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220",
						Type:    "Microsoft.AppPlatform/Spring",
						Address: "azapi_resource.test",
						Change: types.Change{
							Before: `{"foo":"bar"}`,
							After:  `{"foo":"after"}`,
						},
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		actual := tf.NewDiffReport(testcase.Input, nil)
		if len(actual.Diffs) != len(testcase.Expect.Diffs) {
			t.Errorf("Expect diff %v, but got %v", testcase.Expect.Diffs, actual.Diffs)
			continue
		}
		for i := 0; i < len(actual.Diffs); i++ {
			if actual.Diffs[i].Id != testcase.Expect.Diffs[i].Id {
				t.Errorf("Expect diff id %s, but got %s", testcase.Expect.Diffs[i].Id, actual.Diffs[i].Id)
			}
			if actual.Diffs[i].Type != testcase.Expect.Diffs[i].Type {
				t.Errorf("Expect diff type %s, but got %s", testcase.Expect.Diffs[i].Type, actual.Diffs[i].Type)
			}
			if actual.Diffs[i].Address != testcase.Expect.Diffs[i].Address {
				t.Errorf("Expect diff address %s, but got %s", testcase.Expect.Diffs[i].Address, actual.Diffs[i].Address)
			}
			if actual.Diffs[i].Change.Before != testcase.Expect.Diffs[i].Change.Before {
				t.Errorf("Expect diff before %s, but got %s", testcase.Expect.Diffs[i].Change.Before, actual.Diffs[i].Change.Before)
			}
			if actual.Diffs[i].Change.After != testcase.Expect.Diffs[i].Change.After {
				t.Errorf("Expect diff after %s, but got %s", testcase.Expect.Diffs[i].Change.After, actual.Diffs[i].Change.After)
			}
		}
	}
}

func Test_NewPassReportFromState(t *testing.T) {
	var testcases = []struct {
		Input  *tfjson.State
		Expect types.PassReport
	}{
		{
			Input:  nil,
			Expect: types.PassReport{},
		},
		{
			Input:  &tfjson.State{},
			Expect: types.PassReport{},
		},
		{
			Input: &tfjson.State{
				Values: &tfjson.StateValues{
					RootModule: &tfjson.StateModule{
						Resources: []*tfjson.StateResource{
							{
								Address: "azapi_resource.test",
								AttributeValues: map[string]interface{}{
									"type": "Microsoft.AppPlatform/Spring",
								},
							},
							{
								Address: "azapi_resource.test2",
								AttributeValues: map[string]interface{}{
									"type": "Microsoft.AppPlatform/Spring",
								},
							},
							{
								Address:         "azurerm_resource_group.test",
								AttributeValues: map[string]interface{}{},
							},
						},
					},
				},
			},
			Expect: types.PassReport{
				Resources: []types.Resource{
					{
						Type:    "Microsoft.AppPlatform/Spring",
						Address: "azapi_resource.test",
					},
					{
						Type:    "Microsoft.AppPlatform/Spring",
						Address: "azapi_resource.test2",
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		actual := tf.NewPassReportFromState(testcase.Input)
		if len(actual.Resources) != len(testcase.Expect.Resources) {
			t.Errorf("Expect %d resources, but got %d", len(testcase.Expect.Resources), len(actual.Resources))
			continue
		}
		for i := 0; i < len(actual.Resources); i++ {
			if actual.Resources[i].Address != testcase.Expect.Resources[i].Address {
				t.Errorf("Expect resource address %s, but got %s", testcase.Expect.Resources[i].Address, actual.Resources[i].Address)
			}
			if actual.Resources[i].Type != testcase.Expect.Resources[i].Type {
				t.Errorf("Expect resource type %s, but got %s", testcase.Expect.Resources[i].Type, actual.Resources[i].Type)
			}
		}
	}
}

func Test_NewPassReport(t *testing.T) {
	var testcases = []struct {
		Input  *tfjson.Plan
		Expect types.PassReport
	}{
		{
			Input:  nil,
			Expect: types.PassReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{},
				},
			},
			Expect: types.PassReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azurerm_resource_group.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
							Before:  "before",
							After:   "after",
						},
					},
				},
			},
			Expect: types.PassReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionNoop},
							Before:  "foo",
							After:   "foo",
						},
					},
				},
			},
			Expect: types.PassReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
					},
				},
			},
			Expect: types.PassReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionUpdate},
							Before:  "foo",
							After:   "foo",
						},
					},
				},
			},
			Expect: types.PassReport{},
		},
		{
			Input: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "azapi_resource.test",
						Change: &tfjson.Change{
							Actions: []tfjson.Action{tfjson.ActionNoop},
							Before: map[string]interface{}{
								"type": "Microsoft.AppPlatform/Spring",
							},
							After: map[string]interface{}{
								"type": "Microsoft.AppPlatform/Spring",
							},
						},
					},
				},
			},
			Expect: types.PassReport{
				Resources: []types.Resource{
					{
						Type:    "Microsoft.AppPlatform/Spring",
						Address: "azapi_resource.test",
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		actual := tf.NewPassReport(testcase.Input)
		if len(actual.Resources) != len(testcase.Expect.Resources) {
			t.Errorf("Expect %d resources, but got %d", len(testcase.Expect.Resources), len(actual.Resources))
			continue
		}
		for i := 0; i < len(actual.Resources); i++ {
			if actual.Resources[i].Address != testcase.Expect.Resources[i].Address {
				t.Errorf("Expect resource address %s, but got %s", testcase.Expect.Resources[i].Address, actual.Resources[i].Address)
			}
			if actual.Resources[i].Type != testcase.Expect.Resources[i].Type {
				t.Errorf("Expect resource type %s, but got %s", testcase.Expect.Resources[i].Type, actual.Resources[i].Type)
			}
		}
	}
}

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
		{
			Input: fmt.Errorf(`
Error: Failed to create/update resource

  with azapi_resource.dataCollectionRule,
  on testing.tf line 2, in resource "azapi_resource" "dataCollectionRule":
   2: resource "azapi_resource" "dataCollectionRule" {

creating/updating Resource: (ResourceId
"/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001"
/ Api Version "2022-06-01"): PUT
https://management.azure.com/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: InvalidPayload
--------------------------------------------------------------------------------
{
  "error": {
    "code": "InvalidPayload",
    "message": "Data collection rule is invalid",
    "details": [
      {
        "code": "InvalidProperty",
        "message": "XPath query is invalid. Query: 'Security!', Error: Query expression cannot be empty.",
        "target": "Properties.DataSources.WindowsEventLogs[0].XPathQueries[0]"
      }
    ]
  }
}
--------------------------------------------------------------------------------`),
			Expect: []types.Error{
				{
					Type:  "Microsoft.Insights/dataCollectionRules@2022-06-01",
					Label: "dataCollectionRule",
					Id:    "/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001",
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

func Test_NewCleanupErrorReport(t *testing.T) {
	var testcases = []struct {
		Input  error
		Expect types.ErrorReport
	}{
		{
			Input:  nil,
			Expect: types.ErrorReport{},
		},
		{
			Input: fmt.Errorf(`Error: Failed to delete resource

deleting Resource: (ResourceId
"/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001"
/ Api Version "2022-06-01"): DELETE
https://management.azure.com/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001
--------------------------------------------------------------------------------
RESPONSE 409: 409 Conflict
ERROR CODE: ScopeLocked
--------------------------------------------------------------------------------
{
  "error": {
    "code": "ScopeLocked",
    "message": "The scope '/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001' cannot perform delete operation because following scope(s) are locked: '/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001'. Please remove the lock and try again."
  }
}
--------------------------------------------------------------------------------`),
			Expect: types.ErrorReport{
				Errors: []types.Error{
					{
						Type: "Microsoft.Insights/dataCollectionRules@2022-06-01",
						Id:   "/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001",
					},
				},
			},
		},
		{
			Input: fmt.Errorf(`Error: deleting Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001\" / Api Version \"2022-06-01\"): 

DELETE https://management.azure.com/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001
--------------------------------------------------------------------------------
RESPONSE 409: 409 Conflict
ERROR CODE: ScopeLocked
--------------------------------------------------------------------------------
{
  "error": {
    "code": "ScopeLocked",
    "message": "The scope '/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001' cannot perform delete operation because following scope(s) are locked: '/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001'. Please remove the lock and try again."
  }
}
--------------------------------------------------------------------------------`),
			Expect: types.ErrorReport{
				Errors: []types.Error{
					{
						Type: "Microsoft.Insights/dataCollectionRules@2022-06-01",
						Id:   "/subscriptions/******/resourceGroups/acctest0001/providers/Microsoft.Insights/dataCollectionRules/acctest0001",
					},
				},
			},
		},
	}

	for _, testcase := range testcases {
		actual := tf.NewCleanupErrorReport(testcase.Input, nil)
		if len(actual.Errors) != len(testcase.Expect.Errors) {
			t.Errorf("Expect %d errors, but got %d", len(testcase.Expect.Errors), len(actual.Errors))
			continue
		}
		for i := 0; i < len(actual.Errors); i++ {
			if actual.Errors[i].Type != testcase.Expect.Errors[i].Type {
				t.Errorf("Expect error type %s, but got %s", testcase.Expect.Errors[i].Type, actual.Errors[i].Type)
			}
			if actual.Errors[i].Id != testcase.Expect.Errors[i].Id {
				t.Errorf("Expect error id %s, but got %s", testcase.Expect.Errors[i].Id, actual.Errors[i].Id)
			}
		}
	}
}

func Test_NewIdAddressFromState(t *testing.T) {
	testcases := []struct {
		Input  *tfjson.State
		Expect map[string]string
	}{
		{
			Input:  nil,
			Expect: map[string]string{},
		},
		{
			Input: &tfjson.State{
				Values: &tfjson.StateValues{
					RootModule: &tfjson.StateModule{
						Resources: []*tfjson.StateResource{
							{
								Address: "azapi_resource.test",
								AttributeValues: map[string]interface{}{
									"id": "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220",
								},
							},
							{
								Address: "azapi_resource.test2",
								AttributeValues: map[string]interface{}{
									"id": "/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1222",
								},
							},
						},
					},
				},
			},
			Expect: map[string]string{
				"/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1220": "azapi_resource.test",
				"/subscriptions/*****/resourceGroups/acctest1220/providers/Microsoft.AppPlatform/Spring/acctest1222": "azapi_resource.test2",
			},
		},
	}

	for _, testcase := range testcases {
		actual := tf.NewIdAddressFromState(testcase.Input)
		if !reflect.DeepEqual(actual, testcase.Expect) {
			t.Errorf("Expect %v, but got %v", testcase.Expect, actual)
		}
	}
}

func Test_DeepCopy(t *testing.T) {
	testcases := []interface{}{
		nil,
		"foo",
		1,
		1.1,
		true,
		map[string]interface{}{
			"foo": "bar",
		},
		[]interface{}{
			"foo",
		},
	}

	for _, testcase := range testcases {
		actual := tf.DeepCopy(testcase)
		if !reflect.DeepEqual(actual, testcase) {
			t.Errorf("Expect different pointer, but got the same")
		}
	}
}
