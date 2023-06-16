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
    "location": "westeurope",
    "identity": {
        "type": "UserAssigned",
        "UserAssignedIdentities": {
            "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.ManagedIdentity/userAssignedIdentities/acctest": {}
        }
    },
    "kind": "Linux",
    "properties": {
        "dataCollectionEndpointId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Insights/dataCollectionEndpoints/example-mdce",
        "dataSources": {
            "dataImports": {
                "eventHub": {
                    "name": "testdataimports",
                    "stream": "Custom-qweqwe_",
                    "consumerGroup": "$Default"
                }
            },
            "windowsFirewallLogs": [
                {
                    "streams": [
                        "Microsoft-Event"
                    ],
                    "name": "testfl"
                }
            ],
            "prometheusForwarder": [
                {
                    "labelIncludeFilter": {
                        "microsoft_metrics_include_label": "value"
                    },
                    "streams": [
                        "Microsoft-PrometheusMetrics"
                    ],
                    "name": "testpf"
                }
            ],
            "platformTelemetry": [
                {
                    "streams": [
                        "Microsoft.Cache/redis:Metrics-Group-All"
                    ],
                    "name": "testpt"
                }
            ],
            "iisLogs": [
                {
                    "streams": [
                        "Microsoft-W3CIISLog"
                    ],
                    "name": "testIIS",
                    "logDirectories": [
                        "C:\\Logs\\W3SVC1"
                    ]
                }
            ],
            "logFiles": [
                {
                    "streams": [
                        "Custom-MyTableRawData"
                    ],
                    "filePatterns": [
                        "C:\\JavaLogs\\*.log"
                    ],
                    "format": "text",
                    "settings": {
                        "text": {
                            "recordStartTimestampFormat": "ISO 8601"
                        }
                    },
                    "name": "myLogFileFormat-Windows"
                }
            ],
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
                        "Security![]"
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
                    "workspaceResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.OperationalInsights/workspaces/ewtacctest-01",
                    "name": "centralWorkspace2"
                }
            ],
            "eventHubs": [
                {
                    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.EventHub/namespaces/awtcceptanceTestEventHubNamespace/eventhubs/wtacceptanceTestEventHub",
                    "name": "testev"
                }
            ],
             "storageTablesDirect": [
                 {
                     "storageAccountResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Storage/storageAccounts/wtteaplstoraccount",
                     "tableName":"mysampletable",
                     "name": "testtb"
                 }
             ],
             "storageBlobsDirect": [
                 {
                     "storageAccountResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Storage/storageAccounts/wtteaplstoraccount",
                     "containerName": "vhds",
                     "name": "test1tb"
                 }
             ],
            "monitoringAccounts": [
                {
                    "accountResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Monitor/accounts/testaccount",
                    "name": "testmonitor"
                }
            ],
            "storageAccounts": [
                {
                    "storageAccountResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Storage/storageAccounts/wtteaplstoraccount",
                    "name": "teststg",
                    "containerName": "vhds"
                }
            ]
        },
        "streamDeclarations": {
            "Custom-MyTableRawData": {
                "columns": [
                    {
                        "name": "Time",
                        "type": "datetime"
                    },
                    {
                        "name": "Computer",
                        "type": "string"
                    },
                    {
                        "name": "AdditionalContext",
                        "type": "string"
                    }
                ]
            },
            "Custom-MyTableRawData2": {
                "columns": [
                    {
                        "name": "Time",
                        "type": "datetime"
                    },
                    {
                        "name": "Computer",
                        "type": "string"
                    },
                    {
                        "name": "AdditionalContext",
                        "type": "string"
                    }
                ]
            }
        },
        "dataFlows": [
            {
                "streams": [
                    "Microsoft-Syslog"
                ],
                "destinations": [
                    "centralWorkspace2"
                ]
            },
            {
                "streams": [
                    "Microsoft-Perf"
                ],
                "destinations": [
                    "centralWorkspace2"
                ]
            },
            {
                "streams": [
                    "Custom-MyTableRawData"
                ],
                "destinations": [
                    "centralWorkspace2"
                ],
                "transformKql": "source | project TimeGenerated = Time, Computer, Message = AdditionalContext",
                "outputStream": "Microsoft-Syslog"
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
