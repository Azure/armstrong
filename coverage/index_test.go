package coverage_test

import (
	"testing"

	"github.com/ms-henglu/armstrong/coverage"
)

func TestGetModelInfoFromIndex(t *testing.T) {
	apiVersion := "2022-06-01"
	apiPath, modelName, modelSwaggerPath, err := coverage.GetModelInfoFromIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
	)
	if err != nil {
		t.Errorf("get model info from index error: %+v", err)
	}

	expectedApiPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/dataCollectionRules/{dataCollectionRuleName}"
	if *apiPath != expectedApiPath {
		t.Errorf("expected apiPath %s, got %s", expectedApiPath, *apiPath)
	}

	expectedModelName := "DataCollectionRuleResource"
	if *modelName != expectedModelName {
		t.Errorf("expected modelName %s, got %s", expectedModelName, *modelName)
	}

	expectedModelSwaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/monitor/resource-manager/Microsoft.Insights/stable/2022-06-01/dataCollectionRules_API.json"
	if *modelSwaggerPath != expectedModelSwaggerPath {
		t.Errorf("expected modelSwaggerPath %s, got %s", expectedModelSwaggerPath, *modelSwaggerPath)
	}
}
