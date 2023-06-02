package coverage_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"sort"
	"testing"

	"github.com/ms-henglu/armstrong/coverage"
	"github.com/ms-henglu/armstrong/report"
	"github.com/ms-henglu/armstrong/types"
)

func TestCoverage(t *testing.T) {
	//resourceId := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1/providers/Microsoft.Insights/dataCollectionRules/rule1"
	//swaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/monitor/resource-manager/Microsoft.Insights/stable/2022-06-01/dataCollectionRules_API.json"
	resourceId := "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ex-resources/providers/Microsoft.Media/mediaServices/mediatest/transforms/transform1"
	swaggerPath := "https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/specification/mediaservices/resource-manager/Microsoft.Media/Encoding/stable/2022-07-01/Encoding.json"

	apiPath, modelName, modelSwaggerPath, err := coverage.PathPatternFromId(resourceId, swaggerPath)
	if err != nil {
		t.Error(err)
	}

	apiPath = apiPath
	//fmt.Println(*apiPath, *modelName, *modelSwaggerPath)

	model, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	rawRequestJson := `
{
  "properties": {
	"description": "Example Transform to illustrate create and update.",
	"outputs": [
	  {
		"preset": {
		  "@odata.type": "#Microsoft.Media.BuiltInStandardEncoderPreset",
		  "presetName": "AdaptiveStreaming"
		}
	  }
	]
  }
}
`

	body := map[string]interface{}{}
	err = json.Unmarshal([]byte(rawRequestJson), &body)
	if err != nil {
		t.Error(err)
	}

	coverage.MarkCovered(body, model)
	coverage.ComputeCoverage(model)

	out, err := json.MarshalIndent(*model, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("model", string(out))

	var covered, uncovered []string
	coverage.SplitCovered(model, &covered, &uncovered)

	sort.Strings(covered)
	out, err = json.MarshalIndent(covered, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("coveredList", string(out))

	sort.Strings(uncovered)
	out, err = json.MarshalIndent(uncovered, "", "\t")
	if err != nil {
		t.Error(err)
	}

	fmt.Println("uncoveredList", string(out))

	fmt.Printf("covered:%v, uncoverd: %v, total: %v\n", len(covered), len(uncovered), len(covered)+len(uncovered))

}

func storeCoverageReport(coverageReport types.CoverageReport, reportDir string, reportName string) {
	if len(coverageReport.Coverages) != 0 {
		err := os.WriteFile(path.Join(reportDir, reportName), []byte(report.CoverageMarkdownReport(coverageReport)), 0644)
		if err != nil {
			log.Printf("[WARN] failed to save passed markdown report to %s: %+v", reportName, err)
		} else {
			log.Printf("[INFO] markdown report saved to %s", reportName)
		}
	}
}
