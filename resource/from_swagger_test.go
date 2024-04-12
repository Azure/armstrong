package resource_test

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/azure/armstrong/resource"
	"github.com/azure/armstrong/resource/types"
	"github.com/azure/armstrong/swagger"
)

func Test_Format(t *testing.T) {
	wd, _ := os.Getwd()
	testcases := []struct {
		ApiPath  swagger.ApiPath
		Expected []types.AzapiDefinition
	}{
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/providers/Microsoft.Automation/automationAccounts",
				ResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeList,
				Methods:      []string{"GET"},
				ExampleMap: map[string]string{
					"GET": path.Clean(path.Join(wd, "testdata", "./examples/listAutomationAccountsBySubscription.json")),
				},
			},
			Expected: []types.AzapiDefinition{{
				Kind:              types.KindDataSource,
				ResourceName:      "azapi_resource_list",
				Label:             "listAutomationAccountsBySubscription",
				AzureResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:        "2022-08-08",
				AdditionalFields: map[string]types.Value{
					"parent_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}"),
				},
				Body: nil,
			}},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts",
				ResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeList,
				Methods:      []string{"GET"},
				ExampleMap: map[string]string{
					"GET": path.Clean(path.Join(wd, "testdata", "./examples/listAutomationAccountsByResourceGroup.json")),
				},
			},
			Expected: []types.AzapiDefinition{{
				Kind:              types.KindDataSource,
				ResourceName:      "azapi_resource_list",
				Label:             "listAutomationAccountsByResourceGroup",
				AzureResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:        "2022-08-08",
				AdditionalFields: map[string]types.Value{
					"parent_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}"),
				},
				Body: nil,
			}},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",
				ResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeResource,
				Methods:      []string{"DELETE", "GET", "PATCH", "PUT"},
				ExampleMap: map[string]string{
					"GET":    path.Clean(path.Join(wd, "testdata", "./examples/getAutomationAccount.json")),
					"PATCH":  path.Clean(path.Join(wd, "testdata", "./examples/updateAutomationAccount.json")),
					"PUT":    path.Clean(path.Join(wd, "testdata", "./examples/createOrUpdateAutomationAccount.json")),
					"DELETE": path.Clean(path.Join(wd, "testdata", "./examples/deleteAutomationAccount.json")),
				},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource",
					Label:             "automationAccount",
					AzureResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"parent_id":                 types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}"),
						"name":                      types.NewStringLiteralValue("{automationAccountName}"),
						"location":                  types.NewStringLiteralValue("East US 2"),
						"schema_validation_enabled": types.NewRawValue("false"),
					},
					Body: map[string]interface{}{
						"properties": map[string]interface{}{
							"sku": map[string]interface{}{
								"name": "Free",
							},
						},
					},
				},
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "patch_automationAccount",
					AzureResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("PATCH"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}"),
						"action":      types.NewStringLiteralValue(""),
					},
					Body: map[string]interface{}{
						"properties": map[string]interface{}{
							"sku": map[string]interface{}{
								"name": "Free",
							},
						},
					},
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/listKeys",
				ResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeResourceAction,
				Methods:      []string{"POST"},
				ExampleMap: map[string]string{
					"POST": path.Clean(path.Join(wd, "testdata", "./examples/listAutomationAccountKeys.json")),
				},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "listKeys",
					AzureResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("POST"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}"),
						"action":      types.NewStringLiteralValue("listKeys"),
					},
					Body: nil,
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/statistics",
				ResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeResourceAction,
				Methods:      []string{"GET"},
				ExampleMap: map[string]string{
					"GET": path.Clean(path.Join(wd, "testdata", "./examples/getStatisticsOfAutomationAccount.json")),
				},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindDataSource,
					ResourceName:      "azapi_resource_action",
					Label:             "statistics",
					AzureResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("GET"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}"),
						"action":      types.NewStringLiteralValue("statistics"),
					},
					Body: nil,
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/usages",
				ResourceType: "Microsoft.Automation/automationAccounts",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeResourceAction,
				Methods:      []string{"GET"},
				ExampleMap: map[string]string{
					"GET": path.Clean(path.Join(wd, "testdata", "./examples/getUsagesOfAutomationAccount.json")),
				},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindDataSource,
					ResourceName:      "azapi_resource_action",
					Label:             "usages",
					AzureResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("GET"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}"),
						"action":      types.NewStringLiteralValue("usages"),
					},
					Body: nil,
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}",
				ResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeList,
				Methods:      []string{"DELETE", "GET", "PATCH", "PUT"},
				ExampleMap: map[string]string{
					"GET":    path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/getSourceControl.json")),
					"PATCH":  path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/updateSourceControl_patch.json")),
					"PUT":    path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/createOrUpdateSourceControl.json")),
					"DELETE": path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/deleteSourceControl.json")),
				},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource",
					Label:             "sourceControl",
					AzureResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"parent_id":                 types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}"),
						"name":                      types.NewStringLiteralValue("{sourceControlName}"),
						"schema_validation_enabled": types.NewRawValue("false"),
					},
					Body: map[string]interface{}{
						"properties": map[string]interface{}{
							"repoUrl":        "https://sampleUser.visualstudio.com/myProject/_git/myRepository",
							"branch":         "master",
							"folderPath":     "/folderOne/folderTwo",
							"autoSync":       true,
							"publishRunbook": true,
							"sourceType":     "VsoGit",
							"securityToken": map[string]interface{}{
								"accessToken": "******",
								"tokenType":   "PersonalAccessToken",
							},
							"description": "my description",
						},
					},
				},
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "patch_sourceControl",
					AzureResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("PATCH"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}"),
						"action":      types.NewStringLiteralValue(""),
					},
					Body: map[string]interface{}{
						"properties": map[string]interface{}{
							"branch":         "master",
							"folderPath":     "/folderOne/folderTwo",
							"autoSync":       true,
							"publishRunbook": true,
							"securityToken": map[string]interface{}{
								"accessToken": "******",
								"tokenType":   "PersonalAccessToken",
							},
							"description": "my description",
						},
					},
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls",
				ResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeList,
				Methods:      []string{"GET"},
				ExampleMap: map[string]string{
					"GET": path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/getAllSourceControls.json")),
				},
			},
			Expected: []types.AzapiDefinition{{
				Kind:              types.KindDataSource,
				ResourceName:      "azapi_resource_list",
				Label:             "listSourceControlsByAutomationAccount",
				AzureResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
				ApiVersion:        "2022-08-08",
				AdditionalFields: map[string]types.Value{
					"parent_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}"),
				},
			}},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/providers/Microsoft.Automation/operations",
				ResourceType: "Microsoft.Automation/operations",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeList,
				Methods:      []string{"GET"},
				ExampleMap: map[string]string{
					"GET": path.Clean(path.Join(wd, "testdata", "./examples/listRestAPIOperations.json")),
				},
			},
			Expected: []types.AzapiDefinition{{
				Kind:              types.KindDataSource,
				ResourceName:      "azapi_resource_list",
				Label:             "listOperationsByTenant",
				AzureResourceType: "Microsoft.Automation/operations",
				ApiVersion:        "2022-08-08",
				AdditionalFields: map[string]types.Value{
					"parent_id": types.NewStringLiteralValue("/"),
				},
			}},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/convertGraphRunbookContent",
				ResourceType: "Microsoft.Automation",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeProviderAction,
				Methods:      []string{"POST"},
				ExampleMap: map[string]string{
					"POST": path.Clean(path.Join(wd, "testdata", "./examples/deserializeGraphRunbookContent.json")),
				},
			},
			Expected: []types.AzapiDefinition{{
				Kind:              types.KindResource,
				ResourceName:      "azapi_resource_action",
				Label:             "convertGraphRunbookContent",
				AzureResourceType: "Microsoft.Automation",
				ApiVersion:        "2022-08-08",
				AdditionalFields: map[string]types.Value{
					"method":      types.NewStringLiteralValue("POST"),
					"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation"),
					"action":      types.NewStringLiteralValue("convertGraphRunbookContent"),
				},
				Body: map[string]interface{}{
					"rawContent": map[string]interface{}{
						"schemaVersion":     "1.10",
						"runbookDefinition": "AAEAAADAQAAAAAAAAAMAgAAAGJPcmNoZXN0cmF0b3IuR3JhcGhSdW5ib29rLk1vZGVsLCBWZXJzaW9uPTcuMy4wLjAsIEN1bHR....",
						"runbookType":       "GraphPowerShell",
					},
				},
			}},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}",
				ResourceType: "Microsoft.ApiManagement/locations/deletedservices",
				ApiVersion:   "2022-09-01-preview",
				ApiType:      swagger.ApiTypeResource,
				Methods:      []string{"GET", "DELETE"},
				ExampleMap:   map[string]string{},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "get_deletedservice",
					AzureResourceType: "Microsoft.ApiManagement/locations/deletedservices",
					ApiVersion:        "2022-09-01-preview",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("GET"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}"),
					},
				},
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "delete_deletedservice",
					AzureResourceType: "Microsoft.ApiManagement/locations/deletedservices",
					ApiVersion:        "2022-09-01-preview",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("DELETE"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}"),
						"depends_on":  types.NewRawValue("[ azapi_resource_action.get_deletedservice ]"),
					},
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}",
				ResourceType: "Microsoft.ApiManagement/locations/deletedservices",
				ApiVersion:   "2022-09-01-preview",
				ApiType:      swagger.ApiTypeResource,
				Methods:      []string{"HEAD"},
				ExampleMap:   map[string]string{},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "head_deletedservice",
					AzureResourceType: "Microsoft.ApiManagement/locations/deletedservices",
					ApiVersion:        "2022-09-01-preview",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("HEAD"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}"),
						"action":      types.NewStringLiteralValue(""),
					},
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}",
				ResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeList,
				Methods:      []string{"GET", "PUT"},
				ExampleMap: map[string]string{
					"GET": path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/getSourceControl.json")),
					"PUT": path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/createOrUpdateSourceControl.json")),
				},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "put_sourceControl",
					AzureResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}"),
						"method":      types.NewStringLiteralValue("PUT"),
					},
					Body: map[string]interface{}{
						"properties": map[string]interface{}{
							"repoUrl":        "https://sampleUser.visualstudio.com/myProject/_git/myRepository",
							"branch":         "master",
							"folderPath":     "/folderOne/folderTwo",
							"autoSync":       true,
							"publishRunbook": true,
							"sourceType":     "VsoGit",
							"securityToken": map[string]interface{}{
								"accessToken": "******",
								"tokenType":   "PersonalAccessToken",
							},
							"description": "my description",
						},
					},
				},
				{
					Kind:              types.KindDataSource,
					ResourceName:      "azapi_resource",
					Label:             "sourceControl",
					AzureResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}"),
						"depends_on":  types.NewRawValue("[ azapi_resource_action.put_sourceControl ]"),
					},
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}",
				ResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
				ApiVersion:   "2022-08-08",
				ApiType:      swagger.ApiTypeList,
				Methods:      []string{"GET", "PATCH"},
				ExampleMap: map[string]string{
					"GET":   path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/getSourceControl.json")),
					"PATCH": path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/createOrUpdateSourceControl.json")),
				},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "patch_sourceControl",
					AzureResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}"),
						"method":      types.NewStringLiteralValue("PATCH"),
					},
					Body: map[string]interface{}{
						"properties": map[string]interface{}{
							"repoUrl":        "https://sampleUser.visualstudio.com/myProject/_git/myRepository",
							"branch":         "master",
							"folderPath":     "/folderOne/folderTwo",
							"autoSync":       true,
							"publishRunbook": true,
							"sourceType":     "VsoGit",
							"securityToken": map[string]interface{}{
								"accessToken": "******",
								"tokenType":   "PersonalAccessToken",
							},
							"description": "my description",
						},
					},
				},
				{
					Kind:              types.KindDataSource,
					ResourceName:      "azapi_resource",
					Label:             "sourceControl",
					AzureResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
					ApiVersion:        "2022-08-08",
					AdditionalFields: map[string]types.Value{
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls/{sourceControlName}"),
						"depends_on":  types.NewRawValue("[ azapi_resource_action.patch_sourceControl ]"),
					},
				},
			},
		},
		{
			ApiPath: swagger.ApiPath{
				Path:         "/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}",
				ResourceType: "Microsoft.ApiManagement/locations/deletedservices",
				ApiVersion:   "2022-09-01-preview",
				ApiType:      swagger.ApiTypeResource,
				Methods:      []string{"HEAD", "DELETE"},
				ExampleMap:   map[string]string{},
			},
			Expected: []types.AzapiDefinition{
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "head_deletedservice",
					AzureResourceType: "Microsoft.ApiManagement/locations/deletedservices",
					ApiVersion:        "2022-09-01-preview",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("HEAD"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}"),
						"action":      types.NewStringLiteralValue(""),
					},
				},
				{
					Kind:              types.KindResource,
					ResourceName:      "azapi_resource_action",
					Label:             "delete_deletedservice",
					AzureResourceType: "Microsoft.ApiManagement/locations/deletedservices",
					ApiVersion:        "2022-09-01-preview",
					AdditionalFields: map[string]types.Value{
						"method":      types.NewStringLiteralValue("DELETE"),
						"resource_id": types.NewStringLiteralValue("/subscriptions/{subscriptionId}/providers/Microsoft.ApiManagement/locations/{location}/deletedservices/{serviceName}"),
					},
				},
			},
		},
	}
	for _, testcase := range testcases {
		t.Logf("[DEBUG] Testing path: %v", testcase.ApiPath.Path)
		actuals := resource.NewAzapiDefinitionsFromSwagger(testcase.ApiPath)
		expecteds := testcase.Expected
		if len(actuals) != len(expecteds) {
			t.Errorf("expected %d definitions, got %d", len(expecteds), len(actuals))
		}
		for i := range actuals {
			actual := actuals[i]
			expected := expecteds[i]
			if actual.Kind != expected.Kind {
				t.Errorf("expected Kind %s, got %s", expected.Kind, actual.Kind)
			}
			if actual.ResourceName != expected.ResourceName {
				t.Errorf("expected ResourceName %s, got %s", expected.ResourceName, actual.ResourceName)
			}
			if actual.Label != expected.Label {
				t.Errorf("expected Label %s, got %s", expected.Label, actual.Label)
			}
			if actual.AzureResourceType != expected.AzureResourceType {
				t.Errorf("expected AzureResourceType %s, got %s", expected.AzureResourceType, actual.AzureResourceType)
			}
			if actual.ApiVersion != expected.ApiVersion {
				t.Errorf("expected ApiVersion %s, got %s", expected.ApiVersion, actual.ApiVersion)
			}
			if len(actual.AdditionalFields) != len(expected.AdditionalFields) {
				t.Errorf("expected %d AdditionalFields, got %d", len(expected.AdditionalFields), len(actual.AdditionalFields))
			}
			for key, value := range expected.AdditionalFields {
				actualValue, ok := actual.AdditionalFields[key]
				if !ok {
					t.Errorf("expected AdditionalFields %s, got %s", key, actualValue)
				}
				if actualValue.String() != value.String() {
					t.Errorf("expected AdditionalFields key %s: %s, got %s", key, value, actualValue)
				}
			}
			if !reflect.DeepEqual(actual.Body, expected.Body) {
				t.Errorf("expected Body %v, got %v", expected.Body, actual.Body)
			}
		}

	}
}
