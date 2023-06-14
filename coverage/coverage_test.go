package coverage_test

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"testing"

	"github.com/ms-henglu/armstrong/coverage"
	"github.com/ms-henglu/armstrong/report"
	"github.com/ms-henglu/armstrong/types"
)

func TestCoverage(t *testing.T) {
	apiVersion := "2022-06-01"
	apiPath, modelName, modelSwaggerPath, commitId, err := coverage.GetModelInfoFromIndex(
		"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		apiVersion,
	)

	if err != nil {
		t.Errorf("get model info from index error: %+v", err)
	}

	t.Logf("commit id: %s", *commitId)

	model, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	rawRequestJson := `
	{
	 "location": "eastus",
	 "properties": {
		"dataSources": {
		  "performanceCounters": [
			{
			  "name": "cloudTeamCoreCounters",
			  "streams": [
				"Microsoft-Perf"
			  ],
			  "samplingFrequencyInSeconds": 15,
			  "counterSpecifiers": [
				"\\Processor(_Total)\\% Processor Time",
				"\\Memory\\Committed Bytes",
				"\\LogicalDisk(_Total)\\Free Megabytes",
				"\\PhysicalDisk(_Total)\\Avg. Disk Queue Length"
			  ]
			},
			{
			  "name": "appTeamExtraCounters",
			  "streams": [
				"Microsoft-Perf"
			  ],
			  "samplingFrequencyInSeconds": 30,
			  "counterSpecifiers": [
				"\\Process(_Total)\\Thread Count"
			  ]
			}
		  ],
		  "windowsEventLogs": [
			{
			  "name": "cloudSecurityTeamEvents",
			  "streams": [
				"Microsoft-WindowsEvent"
			  ],
			  "xPathQueries": [
				"Security!"
			  ]
			},
			{
			  "name": "appTeam1AppEvents",
			  "streams": [
				"Microsoft-WindowsEvent"
			  ],
			  "xPathQueries": [
				"System![System[(Level = 1 or Level = 2 or Level = 3)]]",
				"Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]"
			  ]
			}
		  ],
		  "syslog": [
			{
			  "name": "cronSyslog",
			  "streams": [
				"Microsoft-Syslog"
			  ],
			  "facilityNames": [
				"cron"
			  ],
			  "logLevels": [
				"Debug",
				"Critical",
				"Emergency"
			  ]
			},
			{
			  "name": "syslogBase",
			  "streams": [
				"Microsoft-Syslog"
			  ],
			  "facilityNames": [
				"syslog"
			  ],
			  "logLevels": [
				"Alert",
				"Critical",
				"Emergency"
			  ]
			}
		  ]
		},
		"destinations": {
		  "logAnalytics": [
			{
			  "workspaceResourceId": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/myResourceGroup/providers/Microsoft.OperationalInsights/workspaces/centralTeamWorkspace",
			  "name": "centralWorkspace"
			}
		  ]
		},
		"dataFlows": [
		  {
			"streams": [
			  "Microsoft-Perf",
			  "Microsoft-Syslog",
			  "Microsoft-WindowsEvent"
			],
			"destinations": [
			  "centralWorkspace"
			]
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

	model.MarkCovered(body)
	model.CountCoverage()

	coverageReport := types.CoverageReport{
		CommitId: *commitId,
		Coverages: map[string]*coverage.Model{
			fmt.Sprintf("%s?api-version=%s", *apiPath, apiVersion): model,
		},
	}

	storeCoverageReport(coverageReport, ".", "test_coverage_report.md")
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
