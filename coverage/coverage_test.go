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

type testCase struct {
	name       string
	apiVersion string
	apiPath    string
	rawRequest []string
}

func TestCoverageResourceGroup(t *testing.T) {
	tc := testCase{
		name:       "ResourceGroup",
		apiVersion: "2022-09-01",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/rgName",
		rawRequest: []string{
			`{"location": "westeurope"}`,
		},
	}

	testCoverage(t, tc)
}

func TestCoverageKeyVault(t *testing.T) {
	tc := testCase{
		name:       "KeyVault",
		apiVersion: "2023-02-01",
		apiPath:    "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sample-resource-group/providers/Microsoft.KeyVault/vaults/sample-vault",
		rawRequest: []string{`{
    "location": "westus",
    "properties": {
        "tenantId": "00000000-0000-0000-0000-000000000000",
        "sku": {
            "family": "A",
            "name": "standard"
        },
        "accessPolicies": [
            {
                "tenantId": "00000000-0000-0000-0000-000000000000",
                "objectId": "00000000-0000-0000-0000-000000000000",
                "permissions": {
                    "keys": [
                        "encrypt",
                        "decrypt",
                        "wrapKey",
                        "unwrapKey",
                        "sign",
                        "verify",
                        "get",
                        "list",
                        "create",
                        "update",
                        "import",
                        "delete",
                        "backup",
                        "restore",
                        "recover",
                        "purge"
                    ],
                    "secrets": [
                        "get",
                        "list",
                        "set",
                        "delete",
                        "backup",
                        "restore",
                        "recover",
                        "purge"
                    ],
                    "certificates": [
                        "get",
                        "list",
                        "delete",
                        "create",
                        "import",
                        "update",
                        "managecontacts",
                        "getissuers",
                        "listissuers",
                        "setissuers",
                        "deleteissuers",
                        "manageissuers",
                        "recover",
                        "purge"
                    ]
                }
            }
        ],
        "enabledForDeployment": true,
        "enabledForDiskEncryption": true,
        "enabledForTemplateDeployment": true,
        "publicNetworkAccess": "Enabled"
    }
}`,
		},
	}

	testCoverage(t, tc)
}

func TestCoverageStorageAccount(t *testing.T) {
	tc := testCase{
		name:       "StorageAccount",
		apiVersion: "2022-09-01",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/res9101/providers/Microsoft.Storage/storageAccounts/sto4445",
		rawRequest: []string{`{
    "sku": {
        "name": "Standard_GRS"
    },
    "kind": "Storage",
    "location": "eastus",
    "extendedLocation": {
        "type": "EdgeZone",
        "name": "losangeles001"
    },
    "properties": {
        "keyPolicy": {
            "keyExpirationPeriodInDays": 20
        },
        "sasPolicy": {
            "sasExpirationPeriod": "1.15:59:59",
            "expirationAction": "Log"
        },
        "isHnsEnabled": true,
        "isSftpEnabled": true,
        "allowBlobPublicAccess": false,
        "defaultToOAuthAuthentication": false,
        "minimumTlsVersion": "TLS1_2",
        "allowSharedKeyAccess": true,
        "routingPreference": {
            "routingChoice": "MicrosoftRouting",
            "publishMicrosoftEndpoints": true,
            "publishInternetEndpoints": true
        },
        "encryption": {
            "services": {
                "file": {
                    "keyType": "Account",
                    "enabled": true
                },
                "blob": {
                    "keyType": "Account",
                    "enabled": true
                }
            },
            "requireInfrastructureEncryption": false,
            "keySource": "Microsoft.Storage"
        }
    },
    "tags": {
        "key1": "value1",
        "key2": "value2"
    }
}`,
		},
	}

	testCoverage(t, tc)
}

func TestCoverageVM(t *testing.T) {
	tc := testCase{
		name:       "VM",
		apiVersion: "2023-03-01",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm",
		rawRequest: []string{`{
    "location": "westus",
    "properties": {
        "hardwareProfile": {
            "vmSize": "Standard_D4_v3",
            "vmSizeProperties": {
                "vCPUsAvailable": 1,
                "vCPUsPerCore": 1
            }
        },
        "storageProfile": {
            "imageReference": {
                "sku": "2016-Datacenter",
                "publisher": "MicrosoftWindowsServer",
                "version": "latest",
                "offer": "WindowsServer"
            },
            "osDisk": {
                "caching": "ReadWrite",
                "managedDisk": {
                    "storageAccountType": "Standard_LRS"
                },
                "name": "myVMosdisk",
                "createOption": "FromImage"
            }
        },
        "networkProfile": {
            "networkInterfaces": [
                {
                    "id": "/subscriptions/{subscription-id}/resourceGroups/myResourceGroup/providers/Microsoft.Network/networkInterfaces/{existing-nic-name}",
                    "properties": {
                        "primary": true
                    }
                }
            ]
        },
        "osProfile": {
            "adminUsername": "{your-username}",
            "computerName": "myVM",
            "adminPassword": "{your-password}"
        },
        "diagnosticsProfile": {
            "bootDiagnostics": {
                "storageUri": "http://{existing-storage-account-name}.blob.core.windows.net",
                "enabled": true
            }
        },
        "userData": "U29tZSBDdXN0b20gRGF0YQ=="
    }
}`,
		},
	}

	testCoverage(t, tc)
}

func TestCoverageVNet(t *testing.T) {
	tc := testCase{
		name:       "VNet",
		apiVersion: "2023-02-01",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/virtualNetwork",
		rawRequest: []string{`{
    "properties": {
        "addressSpace": {
            "addressPrefixes": [
                "10.0.0.0/16"
            ]
        },
        "subnets": [
            {
                "name": "test-1",
                "properties": {
                    "addressPrefix": "10.0.0.0/24"
                }
            }
        ]
    },
    "location": "eastus"
}`,
		},
	}
	testCoverage(t, tc)
}

func TestCoverageDataCollectionRule(t *testing.T) {
	tc := testCase{
		name:       "DataCollectionRule",
		apiVersion: "2022-06-01",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		rawRequest: []string{`{
    "location": "westeurope",
    "identity": {
        "type": "UserAssigned",
        "userAssignedIdentities": {
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
                    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.EventHub/namespaces/TestEventHubNamespace/eventhubs/wtacceptanceTestEventHub",
                    "name": "testev"
                }
            ],
            "eventHubsDirect": [
                {
                    "eventHubResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.EventHub/namespaces/TestEventHubNamespace/eventhubs/wtacceptanceTestEventHub",
                    "name": "testevd"
                }
            ],
            "storageTablesDirect": [
                {
                    "storageAccountResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Storage/storageAccounts/eaplstoraccount",
                    "tableName": "mysampletable",
                    "name": "testtb"
                }
            ],
            "storageBlobsDirect": [
                {
                    "storageAccountResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Storage/storageAccounts/eaplstoraccount",
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
                    "storageAccountResourceId": "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/ex-resources/providers/Microsoft.Storage/storageAccounts/eaplstoraccount",
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
                "builtInTransform": "Syslog-CRON",
                "transformKql": "source | project TimeGenerated = Time, Computer, Message = AdditionalContext",
                "outputStream": "Microsoft-Syslog"
            }
        ]
    }
}`,
			`{
    "location": "westeurope",
    "kind": "Windows",
    "properties": {
        "description": "test",
        "dataSources": {
            "extensions": [
                {
                    "streams": [
                        "Microsoft-WindowsEvent",
                        "Microsoft-ServiceMap"
                    ],
                    "inputDataSources": [
                        "test-datasource-wineventlog"
                    ],
                    "extensionName": "test-extension-name",
                    "extensionSettings": {
                        "b": "hello"
                    },
                    "name": "test-datasource-extension"
                }
            ]
        },
        "destinations": {
            "azureMonitorMetrics": {
                "name": "testDes1"
            }
        },
        "dataFlows": [
            {
                "streams": [
                    "Microsoft-InsightsMetrics"
                ],
                "destinations": [
                    "testDes1"
                ]
            }
        ]
    },
    "tags": {
        "env": "test"
    }
}`,
		},
	}

	testCoverage(t, tc)
}

func TestCoverageMediaTransform(t *testing.T) {
	tc := testCase{
		name:       "MediaTransform",
		apiVersion: "2022-07-01",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/ex-resources/providers/Microsoft.Media/mediaServices/mediatest/transforms/transform1",
		rawRequest: []string{`{	
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
`,
		},
	}

	testCoverage(t, tc)
}

func TestCoverageCosmosDB(t *testing.T) {
	tc := testCase{
		name:       "CosmosDB",
		apiVersion: "2023-04-15",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/testdb",
		rawRequest: []string{`{
    "location": "westus",
    "tags": {},
    "kind": "MongoDB",
    "identity": {
        "type": "SystemAssigned,UserAssigned",
        "userAssignedIdentities": {
            "/subscriptions/fa5fc227-a624-475e-b696-cdd604c735bc/resourceGroups/eu2cgroup/providers/Microsoft.ManagedIdentity/userAssignedIdentities/id1": {}
        }
    },
    "properties": {
        "databaseAccountOfferType": "Standard",
        "ipRules": [
            {
                "ipAddressOrRange": "23.43.230.120"
            },
            {
                "ipAddressOrRange": "110.12.240.0/12"
            }
        ],
        "isVirtualNetworkFilterEnabled": true,
        "virtualNetworkRules": [
            {
                "id": "/subscriptions/subId/resourceGroups/rg/providers/Microsoft.Network/virtualNetworks/vnet1/subnets/subnet1",
                "ignoreMissingVNetServiceEndpoint": false
            }
        ],
        "publicNetworkAccess": "Enabled",
        "locations": [
            {
                "failoverPriority": 0,
                "locationName": "southcentralus",
                "isZoneRedundant": false
            },
            {
                "failoverPriority": 1,
                "locationName": "eastus",
                "isZoneRedundant": false
            }
        ],
        "consistencyPolicy": {
            "defaultConsistencyLevel": "BoundedStaleness",
            "maxIntervalInSeconds": 10,
            "maxStalenessPrefix": 200
        },
        "keyVaultKeyUri": "https://myKeyVault.vault.azure.net",
        "defaultIdentity": "FirstPartyIdentity",
        "enableFreeTier": false,
        "apiProperties": {
            "serverVersion": "3.2"
        },
        "enableAnalyticalStorage": true,
        "analyticalStorageConfiguration": {
            "schemaType": "WellDefined"
        },
        "createMode": "Default",
        "backupPolicy": {
            "type": "Periodic",
            "periodicModeProperties": {
                "backupIntervalInMinutes": 240,
                "backupRetentionIntervalInHours": 8,
                "backupStorageRedundancy": "Geo"
            }
        },
        "cors": [
            {
                "allowedOrigins": "https://test"
            }
        ],
        "networkAclBypass": "AzureServices",
        "networkAclBypassResourceIds": [
            "/subscriptions/subId/resourcegroups/rgName/providers/Microsoft.Synapse/workspaces/workspaceName"
        ],
        "capacity": {
            "totalThroughputLimit": 2000
        },
        "minimalTlsVersion": "Tls12"
    }
}`,
		},
	}

	testCoverage(t, tc)
}

func TestCoverageDataFactoryLinkedServices(t *testing.T) {
	tc := testCase{
		name:       "DataFactoryLinkedServices",
		apiVersion: "2018-06-01",
		apiPath:    "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.DataFactory/factories/factory1/linkedservices/linkedservice1",
		rawRequest: []string{`{
    "properties": {
        "type": "AzureStorage",
        "typeProperties": {
            "connectionString": {
                "type": "SecureString",
                "value": "DefaultEndpointsProtocol=https;AccountName=examplestorageaccount;AccountKey=<storage key>"
            }
        }
    }
}`,
		},
	}

	testCoverage(t, tc)
}

func testCoverage(t *testing.T, tc testCase) {
	apiPath, modelName, modelSwaggerPath, err := coverage.GetModelInfoFromIndex(
		tc.apiPath,
		tc.apiVersion,
	)

	if err != nil {
		t.Errorf("get model info from index error: %+v", err)
	}

	model, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		t.Error(err)
	}

	for _, rq := range tc.rawRequest {
		request := map[string]interface{}{}
		err = json.Unmarshal([]byte(rq), &request)
		if err != nil {
			t.Error(err)
		}

		model.MarkCovered(request)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Error(err)
	}

	t.Logf("expanded model %s", string(out))

	model.CountCoverage()

	coverageReport := types.CoverageReport{
		Coverages: map[string]*coverage.Model{
			fmt.Sprintf("%s?api-version=%s", *apiPath, tc.apiVersion): model,
		},
	}

	storeCoverageReport(coverageReport, ".", fmt.Sprintf("test_coverage_report_%s.md", tc.name))
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
