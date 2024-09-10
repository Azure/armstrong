package coverage_test

import (
	"os"
	"strings"
	"testing"

	"github.com/azure/armstrong/coverage"
	"github.com/go-openapi/jsonreference"
	openapispec "github.com/go-openapi/spec"
)

const indexFilePath = "testdata/index.json"

func TestGetModelInfoFromIndex_DataCollectionRule(t *testing.T) {
	apiVersion := "2022-06-01"
	swaggerModel, err := coverage.GetModelInfoFromIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
		"PUT",
		"",
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

func TestGetModelInfoFromIndexWithCache_DataCollectionRule(t *testing.T) {
	apiVersion := "2022-06-01"
	swaggerModel, err := coverage.GetModelInfoFromIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
		"PUT",
		indexFilePath,
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
		"PUT",
		"",
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

func TestGetModelInfoFromIndexWithType_DataCollectionRule(t *testing.T) {
	azapiResourceType := "Microsoft.Insights/dataCollectionRules@2022-06-01"
	swaggerModel, err := coverage.GetModelInfoFromIndexWithType(azapiResourceType, "PUT", "")
	if err != nil {
		t.Fatalf("get model info from index error: %+v", err)
	}

	expectedApiPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if swaggerModel.ApiPath != expectedApiPath {
		t.Fatalf("expected apiPath %s, got %s", expectedApiPath, swaggerModel.ApiPath)
	}
}

func TestGetModelInfoFromIndexWithTypeWithCache_DataCollectionRule(t *testing.T) {
	azapiResourceType := "Microsoft.Insights/dataCollectionRules@2022-06-01"
	swaggerModel, err := coverage.GetModelInfoFromIndexWithType(azapiResourceType, "PUT", indexFilePath)
	if err != nil {
		t.Fatalf("get model info from index error: %+v", err)
	}

	expectedApiPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if swaggerModel.ApiPath != expectedApiPath {
		t.Fatalf("expected apiPath %s, got %s", expectedApiPath, swaggerModel.ApiPath)
	}

	_, err = coverage.GetModelInfoFromIndexWithType(azapiResourceType, "PUT", indexFilePath)
	if err != nil {
		t.Fatalf("get model info from index error: %+v", err)
	}
}

func TestGetModelInfoFromIndexWithType_DeviceSecurityGroups(t *testing.T) {
	azapiResourceType := "Microsoft.Security/deviceSecurityGroups@2019-08-01"
	swaggerModel, err := coverage.GetModelInfoFromIndexWithType(azapiResourceType, "PUT", "")
	if err != nil {
		t.Fatalf("get model info from index error: %+v", err)
	}

	expectedApiPath := "/{resourceId}/providers/Microsoft.Security/deviceSecurityGroups/{deviceSecurityGroupName}"
	if swaggerModel.ApiPath != expectedApiPath {
		t.Fatalf("expected apiPath %s, got %s", expectedApiPath, swaggerModel.ApiPath)
	}
}

func TestGetModelInfoFromLocalIndex_DataCollectionRule(t *testing.T) {
	azureRepoDir := os.Getenv("AZURE_REST_REPO_DIR")
	if azureRepoDir == "" {
		t.Skip("AZURE_REST_REPO_DIR is not set")
	}

	apiVersion := "2022-06-01"
	swaggerModel, err := coverage.GetModelInfoFromLocalIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
		"PUT",
		azureRepoDir,
		"",
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

	expectedModelSwaggerPathSuffix := "/monitor/resource-manager/Microsoft.Insights/stable/2022-06-01/dataCollectionRules_API.json"
	if !strings.HasSuffix(normarlizePath(swaggerModel.SwaggerPath), expectedModelSwaggerPathSuffix) {
		t.Fatalf("expected modelSwaggerPath has suffix %s, got %s", expectedModelSwaggerPathSuffix, swaggerModel.SwaggerPath)
	}
}

func TestGetModelInfoFromLocalIndexWithCache_DataCollectionRule(t *testing.T) {
	azureRepoDir := os.Getenv("AZURE_REST_REPO_DIR")
	if azureRepoDir == "" {
		t.Skip("AZURE_REST_REPO_DIR is not set")
	}

	apiVersion := "2022-06-01"
	swaggerModel, err := coverage.GetModelInfoFromLocalIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
		"PUT",
		azureRepoDir,
		indexFilePath,
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

	expectedModelSwaggerPathSuffix := "/monitor/resource-manager/Microsoft.Insights/stable/2022-06-01/dataCollectionRules_API.json"
	if !strings.HasSuffix(normarlizePath(swaggerModel.SwaggerPath), expectedModelSwaggerPathSuffix) {
		t.Fatalf("expected modelSwaggerPath has suffix %s, got %s", expectedModelSwaggerPathSuffix, swaggerModel.SwaggerPath)
	}

	_, err = coverage.GetModelInfoFromLocalIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
		"PUT",
		azureRepoDir,
		indexFilePath,
	)
	if err != nil {
		t.Fatalf("get model info from index error: %+v", err)
	}

}

func TestGetModelInfoFromIndexRef(t *testing.T) {
	if os.PathSeparator != '\\' {
		t.Skip("this test is only for windows with backslashes(\\) as the file path seperator")
	}

	azureRepoDir := os.Getenv("AZURE_REST_REPO_DIR")
	if azureRepoDir == "" {
		t.Skip("AZURE_REST_REPO_DIR is not set")
	}

	pathRef := jsonreference.MustCreateRef("monitor%5Cresource-manager%5CMicrosoft.Insights%5Cstable%5C2021-08-01%5CscheduledQueryRule_API.json#/paths/~1subscriptions~1%7BsubscriptionId%7D~1resourceGroups~1%7BresourceGroupName%7D~1providers~1Microsoft.Insights~1scheduledQueryRules~1%7BruleName%7D/put")
	_, err := coverage.GetModelInfoFromIndexRef(openapispec.Ref{Ref: pathRef}, azureRepoDir)
	if err != nil {
		t.Fatal(err)
	}
}
