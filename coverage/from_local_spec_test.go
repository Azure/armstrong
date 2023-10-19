package coverage

import (
	"os"
	"path"
	"testing"
)

func Test_isPathKeyMatchWithResourceId(t *testing.T) {
	testcases := []struct {
		PathKey    string
		ResourceId string
		Expected   bool
	}{
		{
			PathKey:    "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",
			ResourceId: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Automation/automationAccounts/test-automation-account",
			Expected:   true,
		},
		{
			PathKey:    "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",
			ResourceId: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Automation/automationAccounts/test-automation-account/modules/test-module",
			Expected:   false,
		},
		{
			PathKey:    "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/modules/{moduleName}",
			ResourceId: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Automation/automationAccounts/test-automation-account/modules/test-module",
			Expected:   true,
		},
		{
			PathKey:    "{resourceId}/providers/Microsoft.Security/deviceSecurityGroups/{deviceSecurityGroupName}",
			ResourceId: "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SampleRG/providers/Microsoft.Devices/iotHubs/sampleiothub/providers/Microsoft.Security/deviceSecurityGroups/samplesecuritygroup",
			Expected:   true,
		},
	}
	for _, testcase := range testcases {
		t.Logf("testcase: %+v", testcase)
		actual := isPathKeyMatchWithResourceId(testcase.PathKey, testcase.ResourceId)
		if actual != testcase.Expected {
			t.Fatalf("expected %v, got %v", testcase.Expected, actual)
		}
	}
}

func Test_GetModelInfoFromLocalDir(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory error: %+v", err)
	}
	swaggerPath := path.Join(wd, "testdata", "Microsoft.Automation", "stable", "2022-08-08")
	testcases := []struct {
		ResourceId string
		ApiVersion string
		Expected   *SwaggerModel
	}{
		{
			ResourceId: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Automation/automationAccounts/test-automation-account",
			ApiVersion: "2022-08-08",
			Expected: &SwaggerModel{
				ApiPath:     "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}",
				ModelName:   "AutomationAccountCreateOrUpdateParameters",
				SwaggerPath: path.Join(swaggerPath, "account.json"),
			},
		},
		{
			ResourceId: "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Automation/automationAccounts/test-automation-account/certificates/test-certificate",
			ApiVersion: "2022-08-08",
			Expected: &SwaggerModel{
				ApiPath:     "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Automation/automationAccounts/{automationAccountName}/certificates/{certificateName}",
				ModelName:   "CertificateCreateOrUpdateParameters",
				SwaggerPath: path.Join(swaggerPath, "certificate.json"),
			},
		},
	}

	for _, testcase := range testcases {
		t.Logf("testcase: %+v", testcase.ResourceId)
		actual, err := GetModelInfoFromLocalDir(testcase.ResourceId, testcase.ApiVersion, swaggerPath)
		if err != nil {
			t.Fatalf("get model info from local dir error: %+v", err)
		}
		if actual == nil {
			t.Fatalf("expected %+v, got nil", testcase.Expected)
		}
		if actual.ApiPath != testcase.Expected.ApiPath {
			t.Fatalf("expected apiPath %s, got %s", testcase.Expected.ApiPath, actual.ApiPath)
		}
		if actual.ModelName != testcase.Expected.ModelName {
			t.Fatalf("expected modelName %s, got %s", testcase.Expected.ModelName, actual.ModelName)
		}
		if actual.SwaggerPath != testcase.Expected.SwaggerPath {
			t.Fatalf("expected swaggerPath %s, got %s", testcase.Expected.SwaggerPath, actual.SwaggerPath)
		}
	}
}
