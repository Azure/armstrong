package coverage_test

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/ms-henglu/armstrong/coverage"
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

	expand, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	out, err := json.MarshalIndent(expand, "", "\t")
	if err != nil {
		t.Error(err)
	}
	out = out
	//fmt.Println("expand", string(out))

	lookupTable := map[string]bool{}
	discriminatorTable := map[string]string{}
	coverage.Flatten(*expand, "", lookupTable, discriminatorTable)

	out, err = json.MarshalIndent(lookupTable, "", "\t")
	if err != nil {
		t.Error(err)
	}
	//fmt.Println("lookupTable", string(out))

	out, err = json.MarshalIndent(discriminatorTable, "", "\t")
	if err != nil {
		t.Error(err)
	}
	//fmt.Println("discriminatorTable", string(out))

	//	rawRequestJson := `
	//{
	//  "location": "eastus",
	//  "properties": {
	//	"dataSources": {
	//	  "performanceCounters": [
	//		{
	//		  "name": "cloudTeamCoreCounters",
	//		  "streams": [
	//			"Microsoft-Perf"
	//		  ],
	//		  "samplingFrequencyInSeconds": 15,
	//		  "counterSpecifiers": [
	//			"\\Processor(_Total)\\% Processor Time",
	//			"\\Memory\\Committed Bytes",
	//			"\\LogicalDisk(_Total)\\Free Megabytes",
	//			"\\PhysicalDisk(_Total)\\Avg. Disk Queue Length"
	//		  ]
	//		},
	//		{
	//		  "name": "appTeamExtraCounters",
	//		  "streams": [
	//			"Microsoft-Perf"
	//		  ],
	//		  "samplingFrequencyInSeconds": 30,
	//		  "counterSpecifiers": [
	//			"\\Process(_Total)\\Thread Count"
	//		  ]
	//		}
	//	  ],
	//	  "windowsEventLogs": [
	//		{
	//		  "name": "cloudSecurityTeamEvents",
	//		  "streams": [
	//			"Microsoft-WindowsEvent"
	//		  ],
	//		  "xPathQueries": [
	//			"Security!"
	//		  ]
	//		},
	//		{
	//		  "name": "appTeam1AppEvents",
	//		  "streams": [
	//			"Microsoft-WindowsEvent"
	//		  ],
	//		  "xPathQueries": [
	//			"System![System[(Level = 1 or Level = 2 or Level = 3)]]",
	//			"Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]"
	//		  ]
	//		}
	//	  ],
	//	  "syslog": [
	//		{
	//		  "name": "cronSyslog",
	//		  "streams": [
	//			"Microsoft-Syslog"
	//		  ],
	//		  "facilityNames": [
	//			"cron"
	//		  ],
	//		  "logLevels": [
	//			"Debug",
	//			"Critical",
	//			"Emergency"
	//		  ]
	//		},
	//		{
	//		  "name": "syslogBase",
	//		  "streams": [
	//			"Microsoft-Syslog"
	//		  ],
	//		  "facilityNames": [
	//			"syslog"
	//		  ],
	//		  "logLevels": [
	//			"Alert",
	//			"Critical",
	//			"Emergency"
	//		  ]
	//		}
	//	  ]
	//	},
	//	"destinations": {
	//	  "logAnalytics": [
	//		{
	//		  "workspaceResourceId": "/subscriptions/703362b3-f278-4e4b-9179-c76eaf41ffc2/resourceGroups/myResourceGroup/providers/Microsoft.OperationalInsights/workspaces/centralTeamWorkspace",
	//		  "name": "centralWorkspace"
	//		}
	//	  ]
	//	},
	//	"dataFlows": [
	//	  {
	//		"streams": [
	//		  "Microsoft-Perf",
	//		  "Microsoft-Syslog",
	//		  "Microsoft-WindowsEvent"
	//		],
	//		"destinations": [
	//		  "centralWorkspace"
	//		]
	//	  }
	//	]
	//  }
	//}
	//`

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

	coverage.MarkCovered(body, "", lookupTable, discriminatorTable)

	out, err = json.MarshalIndent(lookupTable, "", "\t")
	if err != nil {
		t.Error(err)
	}
	//fmt.Println("coveredTable", string(out))

	var covered, uncovered []string
	for k, v := range lookupTable {
		if v {
			covered = append(covered, k)
		} else {
			uncovered = append(uncovered, k)
		}

	}

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

	fmt.Printf("covered:%v, uncoverd: %v, total: %v\n", len(covered), len(uncovered), len(lookupTable))
}
