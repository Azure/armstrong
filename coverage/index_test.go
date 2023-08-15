package coverage_test

import (
	"testing"

	"github.com/ms-henglu/armstrong/coverage"
)

func TestGetModelInfoFromIndex_DataCollectionRule(t *testing.T) {
	apiVersion := "2022-06-01"
	swaggerModel, err := coverage.GetModelInfoFromIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
	)
	if err != nil {
		t.Fatalf("get model info from index error: %+v", err)
	}

	expectedApiPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if swaggerModel.ApiPath != expectedApiPath {
		t.Fatalf("expected apiPath %s, got %s", expectedApiPath, swaggerModel.ApiPath)
	}

	expectedModelName := "DataCollectionRuleResource"
	if swaggerModel.ModelName != expectedModelName {
		t.Fatalf("expected modelName %s, got %s", expectedModelName, swaggerModel.ModelName)
	}

	expectedModelSwaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/monitor/resource-manager/Microsoft.Insights/stable/2022-06-01/dataCollectionRules_API.json"
	if swaggerModel.SwaggerPath != expectedModelSwaggerPath {
		t.Fatalf("expected modelSwaggerPath %s, got %s", expectedModelSwaggerPath, swaggerModel.SwaggerPath)
	}
}

func TestGetModelInfoFromIndex_DeviceSecurityGroups(t *testing.T) {
	apiVersion := "2019-08-01"
	swaggerModel, err := coverage.GetModelInfoFromIndex(
		"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SampleRG/providers/Microsoft.Devices/iotHubs/sampleiothub/providers/Microsoft.Security/deviceSecurityGroups/samplesecuritygroup",
		apiVersion,
	)
	if err != nil {
		t.Fatalf("get model info from index error: %+v", err)
	}

	expectedApiPath := "/{resourceId}/providers/Microsoft.Security/deviceSecurityGroups/{deviceSecurityGroupName}"
	if swaggerModel.ApiPath != expectedApiPath {
		t.Fatalf("expected apiPath %s, got %s", expectedApiPath, swaggerModel.ApiPath)
	}

	expectedModelName := "DeviceSecurityGroup"
	if swaggerModel.ModelName != expectedModelName {
		t.Fatalf("expected modelName %s, got %s", expectedModelName, swaggerModel.ModelName)
	}

	expectedModelSwaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/security/resource-manager/Microsoft.Security/stable/2019-08-01/deviceSecurityGroups.json"
	if swaggerModel.SwaggerPath != expectedModelSwaggerPath {
		t.Fatalf("expected modelSwaggerPath %s, got %s", expectedModelSwaggerPath, swaggerModel.SwaggerPath)
	}
}
