package swagger_test

import (
	"log"
	"os"
	"path"
	"testing"

	"github.com/ms-henglu/armstrong/swagger"
)

func Test_LoadSwagger(t *testing.T) {
	wd, _ := os.Getwd()
	testcases := []struct {
		Input       string
		ApiPaths    []swagger.ApiPath
		ExpectError bool
	}{
		{
			Input: path.Clean(path.Join(wd, "testdata", "account.json")),
			ApiPaths: []swagger.ApiPath{
				{
					Path:         "/subscriptions/{subscriptionId}/providers/Microsoft.Automation/automationAccounts",
					ResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeList,
					Methods:      []string{"GET"},
					ExampleMap: map[string]string{
						"GET": path.Clean(path.Join(wd, "testdata", "./examples/listAutomationAccountsBySubscription.json")),
					},
				},
				{
					Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts",
					ResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeList,
					Methods:      []string{"GET"},
					ExampleMap: map[string]string{
						"GET": path.Clean(path.Join(wd, "testdata", "./examples/listAutomationAccountsByResourceGroup.json")),
					},
				},
				{
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
				{
					Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/listKeys",
					ResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeResourceAction,
					Methods:      []string{"POST"},
					ExampleMap: map[string]string{
						"POST": path.Clean(path.Join(wd, "testdata", "./examples/listAutomationAccountKeys.json")),
					},
				},
				{
					Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/statistics",
					ResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeResourceAction,
					Methods:      []string{"GET"},
					ExampleMap: map[string]string{
						"GET": path.Clean(path.Join(wd, "testdata", "./examples/getStatisticsOfAutomationAccount.json")),
					},
				},
				{
					Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/usages",
					ResourceType: "Microsoft.Automation/automationAccounts",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeResourceAction,
					Methods:      []string{"GET"},
					ExampleMap: map[string]string{
						"GET": path.Clean(path.Join(wd, "testdata", "./examples/getUsagesOfAutomationAccount.json")),
					},
				},
			},
			ExpectError: false,
		},
		{
			Input: path.Clean(path.Join(wd, "testdata", "sourceControl.json")),
			ApiPaths: []swagger.ApiPath{
				{
					Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/sourceControls",
					ResourceType: "Microsoft.Automation/automationAccounts/sourceControls",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeList,
					Methods:      []string{"GET"},
					ExampleMap: map[string]string{
						"GET": path.Clean(path.Join(wd, "testdata", "./examples/sourceControl/getAllSourceControls.json")),
					},
				},
				{
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
			},
			ExpectError: false,
		},
		{
			Input: path.Clean(path.Join(wd, "testdata", "operations.json")),
			ApiPaths: []swagger.ApiPath{
				{
					Path:         "/providers/Microsoft.Automation/operations",
					ResourceType: "Microsoft.Automation/operations",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeList,
					Methods:      []string{"GET"},
					ExampleMap: map[string]string{
						"GET": path.Clean(path.Join(wd, "testdata", "./examples/listRestAPIOperations.json")),
					},
				},
				{
					Path:         "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/convertGraphRunbookContent",
					ResourceType: "Microsoft.Automation",
					ApiVersion:   "2022-08-08",
					ApiType:      swagger.ApiTypeProviderAction,
					Methods:      []string{"POST"},
					ExampleMap: map[string]string{
						"POST": path.Clean(path.Join(wd, "testdata", "./examples/deserializeGraphRunbookContent.json")),
					},
				},
			},
			ExpectError: false,
		},
	}

	for _, testcase := range testcases {
		log.Printf("[DEBUG] testcase: %s", testcase.Input)
		apiPaths, err := swagger.Load(testcase.Input)
		if err != nil && !testcase.ExpectError {
			t.Errorf("unexpected error: %+v", err)
		}
		if err == nil && testcase.ExpectError {
			t.Errorf("expected error but got nil")
		}
		if len(apiPaths) != len(testcase.ApiPaths) {
			t.Errorf("expected %d api paths but got %d", len(testcase.ApiPaths), len(apiPaths))
		}
		for i, apiPath := range apiPaths {
			log.Printf("[DEBUG] api path: %+v", testcase.ApiPaths[i].Path)
			if apiPath.Path != testcase.ApiPaths[i].Path {
				t.Errorf("expected api path %s but got %s", testcase.ApiPaths[i].Path, apiPath.Path)
			}
			if apiPath.ResourceType != testcase.ApiPaths[i].ResourceType {
				t.Errorf("expected resource type %s but got %s", testcase.ApiPaths[i].ResourceType, apiPath.ResourceType)
			}
			if apiPath.ApiVersion != testcase.ApiPaths[i].ApiVersion {
				t.Errorf("expected api version %s but got %s", testcase.ApiPaths[i].ApiVersion, apiPath.ApiVersion)
			}
			if len(apiPath.Methods) != len(testcase.ApiPaths[i].Methods) {
				t.Errorf("expected %d methods but got %d", len(testcase.ApiPaths[i].Methods), len(apiPath.Methods))
			}
			for j, method := range apiPath.Methods {
				if method != testcase.ApiPaths[i].Methods[j] {
					t.Errorf("expected method %s but got %s", testcase.ApiPaths[i].Methods[j], method)
				}
			}
			if len(apiPath.ExampleMap) != len(testcase.ApiPaths[i].ExampleMap) {
				t.Errorf("expected %d examples but got %d", len(testcase.ApiPaths[i].ExampleMap), len(apiPath.ExampleMap))
			}
			for method, example := range apiPath.ExampleMap {
				if example != testcase.ApiPaths[i].ExampleMap[method] {
					t.Errorf("expected example %s but got %s", testcase.ApiPaths[i].ExampleMap[method], example)
				}
			}
		}
	}
}
