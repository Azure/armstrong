package coverage_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/azure/armstrong/coverage"
	"github.com/azure/armstrong/report"
	"github.com/azure/armstrong/types"
	"github.com/sirupsen/logrus"
)

type testCase struct {
	name                 string
	apiPath              string
	rawRequest           []string
	resourceType         string
	method               string
	expectedCoveredCount int
	expectedTotalCount   int
}

func normarlizePath(path string) string {
	return strings.ReplaceAll(path, string(os.PathSeparator), "/")
}

func TestCoverage_ResourceGroup(t *testing.T) {
	tc := testCase{
		name:                 "ResourceGroup",
		resourceType:         "Microsoft.Resources/resourceGroups@2022-09-01",
		method:               "PUT",
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rgName",
		expectedCoveredCount: 1,
		expectedTotalCount:   3,
		rawRequest: []string{
			`{
    "location": "westeurope"
}`,
		},
	}

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}

	if model.Properties == nil {
		t.Fatalf("expected properties, got none")
	}

	if v, ok := (*model.Properties)["location"]; !ok || v == nil {
		t.Fatalf("expected location, got none")
	}

	if !(*model.Properties)["location"].IsAnyCovered {
		t.Fatalf("expected location IsAnyCovered true, got false")
	}
}

func TestCoverage_HealthcareDicom(t *testing.T) {
	tc := testCase{
		name:                 "HealthcareApisDicom",
		resourceType:         "Microsoft.HealthcareApis/workspaces/dicomservices@2024-03-31",
		method:               "PUT",
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rgName/providers/Microsoft.HealthcareApis/workspaces/workspaceName/dicomservices/dicomServiceName",
		expectedCoveredCount: 12,
		expectedTotalCount:   14,
		rawRequest: []string{
			`{
    "identity": {
        "type": "UserAssigned",
        "userAssignedIdentities": {
            "${azurerm_user_assigned_identity.userAssignedIdentity.id}": {}
        }
    },
    "location": "westus",
    "properties": {
        "corsConfiguration": {
            "allowCredentials": false,
            "headers": [
                "*"
            ],
            "maxAge": 1440,
            "methods": [
                "DELETE",
                "GET",
                "OPTIONS",
                "PATCH",
                "POST",
                "PUT"
            ],
            "origins": [
                "*"
            ]
        },
        "enableDataPartitions": false,
        "encryption": {
            "customerManagedKeyEncryption": {
                "keyEncryptionKeyUrl": "https://${azapi_resource.vault.name}.vault.azure.net/keys/${azurerm_key_vault_key.key.name}/${azurerm_key_vault_key.key.version}"
            }
        },
        "storageConfiguration": {
            "fileSystemName": "fileSystemName",
            "storageResourceId": "${azapi_resource.storageAccount.id}"
        }
    }
}`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_MachineLearningServicesWorkspacesJobs(t *testing.T) {
	tc := testCase{
		name:                 "MachineLearningServicesWorkspacesJobs",
		resourceType:         "Microsoft.MachineLearningServices/workspaces/jobs@2024-04-01",
		method:               "PUT",
		expectedCoveredCount: 11,
		expectedTotalCount:   895,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.MachineLearningServices/workspaces/works1/jobs/job1",
		rawRequest: []string{
			`{
    "properties": {
        "description": "string",
        "tags": {
            "string": "string"
        },
        "properties": {
            "string": "string"
        },
        "displayName": "string",
        "experimentName": "string",
        "services": {
            "string": {
                "jobServiceType": "string",
                "port": 1,
                "endpoint": "string",
                "properties": {
                    "string": "string"
                }
            }
        },
        "computeId": "string",
        "jobType": "Pipeline",
        "settings": {},
        "inputs": {
            "string": {
                "description": "string",
                "jobInputType": "literal",
                "value": "string"
            }
        },
        "outputs": {
            "string": {
                "description": "string",
                "jobOutputType": "uri_file",
                "mode": "Upload",
                "uri": "string"
            }
        }
    }
}`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_MachineLearningServicesWorkspacesDataVersions(t *testing.T) {
	tc := testCase{
		name:                 "MachineLearningServicesWorkspacesDataVersions",
		resourceType:         "Microsoft.MachineLearningServices/workspaces/data/versions@2024-04-01",
		method:               "PUT",
		expectedCoveredCount: 8,
		expectedTotalCount:   29,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.MachineLearningServices/workspaces/works1/data/data1/versions/version1",
		rawRequest: []string{
			`{
    "properties": {
        "description": "string",
        "tags": {
            "string": "string"
        },
        "properties": {
            "string": "string"
        },
        "isArchived": false,
        "isAnonymous": false,
        "dataUri": "string",
        "dataType": "mltable",
        "referencedUris": [
            "string"
        ]
    }
}`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_DeviceSecurityGroup(t *testing.T) {
	tc := testCase{
		name:                 "DeviceSecurityGroup",
		resourceType:         "Microsoft.Security/deviceSecurityGroups@2019-08-01",
		method:               "PUT",
		expectedCoveredCount: 5,
		expectedTotalCount:   192,
		apiPath:              "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/SampleRG/providers/Microsoft.Devices/iotHubs/sampleiothub/providers/Microsoft.Security/deviceSecurityGroups/samplesecuritygroup",
		rawRequest: []string{
			`{
    "properties": {
        "timeWindowRules": [
            {
                "ruleType": "ActiveConnectionsNotInAllowedRange",
                "isEnabled": true,
                "minThreshold": 0,
                "maxThreshold": 30,
                "timeWindowSize": "PT05M"
            }
        ]
    }
}
`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}

}

func TestCoverage_SCVMM(t *testing.T) {
	tc := testCase{
		name:                 "SCVMM",
		resourceType:         "Microsoft.ScVmm/virtualMachineInstances@2023-10-07",
		method:               "PUT",
		expectedCoveredCount: 5,
		expectedTotalCount:   39,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/cli-test-rg-vmm/providers/Microsoft.HybridCompute/machines/test-vm-tf/providers/Microsoft.ScVmm/virtualMachineInstances/default",
		rawRequest: []string{
			`{
    "extendedLocation": {
        "name": "/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/syntheticsscvmmcan/providers/microsoft.extendedlocation/customlocations/terraform-test-cl",
        "type": "customLocation"
    },
    "properties": {
        "infrastructureProfile": {
            "cloudId": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/cli-test-rg-vmm/providers/Microsoft.SCVMM/Clouds/azcli-test-cloud",
            "templateId": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/cli-test-rg-vmm/providers/Microsoft.SCVMM/VirtualMachineTemplates/azcli-test-vm-template-win19",
            "vmmServerId": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/SyntheticsSCVMMCan/providers/microsoft.scvmm/vmmservers/tf-test-vmmserver"
        }
    }
}`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_DataMigrationServiceTasks(t *testing.T) {
	// Do we need to support cross file discriminator reference? Now seems only DataMigration has this. e.g., https://github.com/Azure/azure-rest-api-specs/blob/0ab5469dc0d75594f5747493dcfe8774e22d728f/specification/datamigration/resource-manager/Microsoft.DataMigration/stable/2021-06-30/definitions/ServiceTasks.json#L39
	tc := testCase{
		name:                 "DataMigrationServiceTasks",
		resourceType:         "Microsoft.DataMigration/services/serviceTasks@2021-06-30",
		method:               "PUT",
		expectedCoveredCount: 2,
		expectedTotalCount:   615,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/DmsSdkRg/providers/Microsoft.DataMigration/services/DmsSdkService/serviceTasks/DmsSdkTask",
		rawRequest: []string{
			`{
    "properties": {
        "taskType": "ConnectToSource.MySql",
        "input": {
          "sourceConnectionInfo": {
            "serverName": "mySqlService"
          }
        }
    }
}`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_DataMigrationTasks(t *testing.T) {
	tc := testCase{
		name:                 "DataMigrationTasks",
		resourceType:         "Microsoft.DataMigration/services/projects/tasks@2021-06-30",
		method:               "PUT",
		expectedCoveredCount: 8,
		expectedTotalCount:   615,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/DmsSdkRg/providers/Microsoft.DataMigration/services/DmsSdkService/projects/DmsSdkProject/tasks/DmsSdkTask",
		rawRequest: []string{
			`{
    "properties": {
        "taskType": "ConnectToTarget.SqlDb",
        "input": {
            "targetConnectionInfo": {
                "type": "SqlConnectionInfo",
                "dataSource": "ssma-test-server.database.windows.net",
                "authentication": "SqlAuthentication",
                "encryptConnection": true,
                "trustServerCertificate": true,
                "userName": "testuser",
                "password": "testpassword"
            }
        }
    }
}`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_KeyVault(t *testing.T) {
	tc := testCase{
		name:                 "KeyVault",
		resourceType:         "Microsoft.KeyVault/vaults@2023-02-01",
		method:               "PUT",
		expectedCoveredCount: 13,
		expectedTotalCount:   28,
		apiPath:              "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sample-resource-group/providers/Microsoft.KeyVault/vaults/sample-vault",
		rawRequest: []string{
			`{
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

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_StorageAccount(t *testing.T) {
	tc := testCase{
		name:                 "StorageAccount",
		resourceType:         "Microsoft.Storage/storageAccounts@2023-01-01",
		method:               "PUT",
		expectedCoveredCount: 24,
		expectedTotalCount:   69,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/res9101/providers/Microsoft.Storage/storageAccounts/sto4445",
		rawRequest: []string{
			`{
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
	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_VM(t *testing.T) {
	tc := testCase{
		name:                 "VM",
		resourceType:         "Microsoft.Compute/virtualMachines@2024-03-01",
		method:               "PUT",
		expectedCoveredCount: 20,
		expectedTotalCount:   166,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm",
		rawRequest: []string{
			`{
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

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_VNet(t *testing.T) {
	tc := testCase{
		name:                 "VNet",
		resourceType:         "Microsoft.Network/virtualNetworks@2024-01-01",
		method:               "PUT",
		expectedCoveredCount: 4,
		expectedTotalCount:   104,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/virtualNetwork",
		rawRequest: []string{
			`{
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

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_DataCollectionRule(t *testing.T) {
	tc := testCase{
		name:                 "DataCollectionRule",
		resourceType:         "Microsoft.Insights/dataCollectionRules@2022-06-01",
		method:               "PUT",
		expectedCoveredCount: 65,
		expectedTotalCount:   65,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
		rawRequest: []string{
			`{
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
                        "Microsoft-WindowsEvent"
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}

	expected := 65
	if model.CoveredCount != expected {
		t.Fatalf("expected CoveredCount %d, got %d", expected, model.CoveredCount)
	}
}

func TestCoverage_WebSite(t *testing.T) {
	tc := testCase{
		name:                 "WebSites",
		resourceType:         "Microsoft.Web/sites@2023-01-01",
		method:               "PUT",
		expectedCoveredCount: 3,
		expectedTotalCount:   198,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/testrg123/providers/Microsoft.Web/sites/sitef6141",
		rawRequest: []string{
			`{
    "kind": "app",
    "location": "East US",
    "properties": {
        "serverFarmId": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/testrg123/providers/Microsoft.Web/serverfarms/DefaultAsp"
    }
}`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_AKS(t *testing.T) {
	tc := testCase{
		name:                 "AKS",
		resourceType:         "Microsoft.ContainerService/ManagedClusters@2024-05-01",
		method:               "PUT",
		expectedCoveredCount: 33,
		expectedTotalCount:   234,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1",
		rawRequest: []string{`{
    "location": "location1",
    "tags": {
        "tier": "production",
        "archv2": ""
    },
    "sku": {
        "name": "Base",
        "tier": "Free"
    },
    "properties": {
        "kubernetesVersion": "",
        "dnsPrefix": "dnsprefix1",
        "agentPoolProfiles": [
            {
                "name": "nodepool1",
                "count": 3,
                "vmSize": "Standard_DS2_v2",
                "osType": "Linux",
                "osSKU": "AzureLinux",
                "type": "VirtualMachineScaleSets",
                "enableNodePublicIP": true,
                "mode": "System"
            }
        ],
        "linuxProfile": {
            "adminUsername": "azureuser",
            "ssh": {
                "publicKeys": [
                    {
                        "keyData": "keydata"
                    }
                ]
            }
        },
        "networkProfile": {
            "loadBalancerSku": "standard",
            "outboundType": "loadBalancer",
            "loadBalancerProfile": {
                "managedOutboundIPs": {
                    "count": 2
                }
            }
        },
        "autoScalerProfile": {
            "scan-interval": "20s",
            "scale-down-delay-after-add": "15m"
        },
        "windowsProfile": {
            "adminUsername": "azureuser",
            "adminPassword": "replacePassword1234$"
        },
        "servicePrincipalProfile": {
            "clientId": "clientid",
            "secret": "secret"
        },
        "addonProfiles": {},
        "enableRBAC": true,
        "diskEncryptionSetID": "/subscriptions/subid1/resourceGroups/rg1/providers/Microsoft.Compute/diskEncryptionSets/des",
        "enablePodSecurityPolicy": true,
        "httpProxyConfig": {
            "httpProxy": "http://myproxy.server.com:8080",
            "httpsProxy": "https://myproxy.server.com:8080",
            "noProxy": [
                "localhost",
                "127.0.0.1"
            ],
            "trustedCa": "Q29uZ3JhdHMhIFlvdSBoYXZlIGZvdW5kIGEgaGlkZGVuIG1lc3NhZ2U="
        }
    }
}`},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_CosmosDB(t *testing.T) {
	tc := testCase{
		name:                 "CosmosDB",
		resourceType:         "Microsoft.DocumentDB/databaseAccounts@2024-05-15",
		method:               "PUT",
		expectedCoveredCount: 33,
		expectedTotalCount:   67,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/testdb",
		rawRequest: []string{
			`{
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %v", err)
	}

	if model.Properties == nil {
		t.Fatalf("expected properties, got none")
	}

	if v, ok := (*model.Properties)["identity"]; !ok || v == nil {
		t.Fatalf("expected identity property, got none")
	}

	if (*model.Properties)["identity"].CoveredCount != 2 {
		t.Fatalf("expected identity CoveredCount 2, got %d", (*model.Properties)["identity"].CoveredCount)
	}

	if (*model.Properties)["identity"].Properties == nil {
		t.Fatalf("expected identity properties, got none")
	}

	if v, ok := (*(*model.Properties)["identity"].Properties)["type"]; !ok || v == nil {
		t.Fatalf("expected identity type property, got none")
	}

	if !(*(*model.Properties)["identity"].Properties)["type"].IsAnyCovered {
		t.Fatalf("expected identity type IsAnyCovered true, got false")
	}

	if !(*(*model.Properties)["identity"].Properties)["type"].IsFullyCovered {
		t.Fatalf("expected identity type IsFullyCovered true, got false")
	}

	if (*(*model.Properties)["identity"].Properties)["type"].EnumCoveredCount != 1 {
		t.Fatalf("expected identity type EnumCoveredCount 1, got %d", (*(*model.Properties)["identity"].Properties)["type"].EnumCoveredCount)
	}

	if (*(*model.Properties)["identity"].Properties)["type"].Enum == nil {
		t.Fatalf("expected identity type Enum, got none")
	}

	if isSet, ok := (*(*(*model.Properties)["identity"].Properties)["type"].Enum)["SystemAssigned,UserAssigned"]; !ok || !isSet {
		t.Fatalf("expected identity type Enum SystemAssigned,UserAssigned to be set")
	}

	if v, ok := (*model.Properties)["tags"]; !ok || v == nil {
		t.Fatalf("expected identity tags, got none")
	}

	if !(*model.Properties)["tags"].HasAdditionalProperties {
		t.Fatalf("expected tags HasAdditionalProperties true, got false")
	}

	if !(*model.Properties)["tags"].IsAnyCovered {
		t.Fatalf("expected tags IsAnyCovered true, got false")
	}

	if (*model.Properties)["tags"].TotalCount != 1 {
		t.Fatalf("expected tags TotalCount 1, got %d", (*model.Properties)["tags"].TotalCount)
	}

	if (*model.Properties)["tags"].CoveredCount != 1 {
		t.Fatalf("expected tags CoveredCount 1, got %d", (*model.Properties)["tags"].CoveredCount)
	}

	if v, ok := (*model.Properties)["properties"]; !ok || v == nil {
		t.Fatalf("expected identity properties, got none")
	}

	if !(*model.Properties)["properties"].IsAnyCovered {
		t.Fatalf("expected properties IsAnyCovered true, got false")
	}

	if (*model.Properties)["properties"].Properties == nil {
		t.Fatalf("expected properties properties, got none")
	}

	if v, ok := (*(*model.Properties)["properties"].Properties)["locations"]; !ok || v == nil {
		t.Fatalf("expected properties locations, got none")
	}

	if !(*(*model.Properties)["properties"].Properties)["locations"].IsAnyCovered {
		t.Fatalf("expected locations IsAnyCovered true, got false")
	}

	if (*(*model.Properties)["properties"].Properties)["locations"].Item == nil {
		t.Fatalf("expected locations Item, got none")
	}

	if (*(*model.Properties)["properties"].Properties)["locations"].Item.Properties == nil {
		t.Fatalf("expected locations Item properties, got none")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["locationName"]; !ok || v == nil {
		t.Fatalf("expected locations Item locationName, got none")
	}

	if !(*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["locationName"].IsAnyCovered {
		t.Fatalf("expected locationName IsAnyCovered true, got false")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"]; !ok || v == nil {
		t.Fatalf("expected locations Item isZoneRedundant, got none")
	}

	if !(*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"].IsAnyCovered {
		t.Fatalf("expected isZoneRedundant IsAnyCovered true, got false")
	}

	if (*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"].Bool == nil {
		t.Fatalf("expected isZoneRedundant Bool, got none")
	}

	if isSet, ok := (*(*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"].Bool)["false"]; !ok || !isSet {
		t.Fatalf("expected isZoneRedundant Bool false to be set")
	}

	if v, ok := (*(*model.Properties)["properties"].Properties)["ipRules"]; !ok || v == nil {
		t.Fatalf("expected properties ipRules, got none")
	}

	if (*(*model.Properties)["properties"].Properties)["ipRules"].Item == nil {
		t.Fatalf("expected ipRules Item, got none")
	}

	if (*(*model.Properties)["properties"].Properties)["ipRules"].Item.Properties == nil {
		t.Fatalf("expected ipRules Item properties, got none")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Properties)["ipRules"].Item.Properties)["ipAddressOrRange"]; !ok || v == nil {
		t.Fatalf("expected ipRules Item ipAddressOrRange, got none")
	}
}

func TestCoverage_DataFactoryPipelines(t *testing.T) {
	tc := testCase{
		name:                 "DataFactoryPipelines",
		resourceType:         "Microsoft.DataFactory/factories/pipelines@2018-06-01",
		method:               "PUT",
		expectedCoveredCount: 11,
		expectedTotalCount:   7239,
		apiPath:              "/subscriptions/12345678-1234-1234-1234-12345678abc/resourceGroups/exampleResourceGroup/providers/Microsoft.DataFactory/factories/exampleFactoryName/pipelines/examplePipeline",
		rawRequest: []string{
			`{
    "properties": {
        "activities": [
            {
                "type": "ForEach",
                "typeProperties": {
                    "isSequential": true,
                    "items": {
                        "value": "@pipeline().parameters.OutputBlobNameList",
                        "type": "Expression"
                    },
                    "activities": [
                        {
                            "type": "Copy",
                            "typeProperties": {
                                "source": {
                                    "type": "BlobSource"
                                },
                                "sink": {
                                    "type": "BlobSink"
                                },
                                "dataIntegrationUnits": 32
                            },
                            "inputs": [
                                {
                                    "referenceName": "exampleDataset",
                                    "parameters": {
                                        "MyFolderPath": "examplecontainer",
                                        "MyFileName": "examplecontainer.csv"
                                    },
                                    "type": "DatasetReference"
                                }
                            ],
                            "outputs": [
                                {
                                    "referenceName": "exampleDataset",
                                    "parameters": {
                                        "MyFolderPath": "examplecontainer",
                                        "MyFileName": {
                                            "value": "@item()",
                                            "type": "Expression"
                                        }
                                    },
                                    "type": "DatasetReference"
                                }
                            ],
                            "name": "ExampleCopyActivity"
                        }
                    ]
                },
                "name": "ExampleForeachActivity"
            }
        ],
        "parameters": {
            "OutputBlobNameList": {
                "type": "Array"
            },
            "JobId": {
                "type": "String"
            }
        },
        "variables": {
            "TestVariableArray": {
                "type": "Array"
            }
        },
        "runDimensions": {
            "JobId": {
                "value": "@pipeline().parameters.JobId",
                "type": "Expression"
            }
        },
        "policy": {
            "elapsedTimeMetric": {
                "duration": "0.00:10:00"
            }
        }
    }
}
`,
		},
	}

	_, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}
}

func TestCoverage_DataFactoryLinkedServices(t *testing.T) {
	tc := testCase{
		name:                 "DataFactoryLinkedServices",
		resourceType:         "Microsoft.DataFactory/factories/linkedServices@2018-06-01",
		method:               "PUT",
		expectedCoveredCount: 2,
		expectedTotalCount:   3450,
		apiPath:              "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.DataFactory/factories/factory1/linkedServices/linked",
		rawRequest: []string{
			`{
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}

	if model.Properties == nil {
		t.Fatalf("expected properties, got none")
	}

	if (*model.Properties)["properties"].Properties == nil {
		t.Fatalf("expected properties properties, got none")
	}

	if v, ok := (*(*model.Properties)["properties"].Properties)["type"]; !ok || v == nil {
		t.Fatalf("expected properties type property, got none")
	}

	if (*(*model.Properties)["properties"].Properties)["type"].IsAnyCovered {
		t.Fatalf("expected properties type IsAnyCovered false, got true")
	}

	if (*model.Properties)["properties"].Discriminator == nil {
		t.Fatalf("expected properties discriminator, got none")
	}

	if *(*model.Properties)["properties"].Discriminator != "type" {
		t.Fatalf("expected properties discriminator 'type', got %s", *(*model.Properties)["properties"].Discriminator)
	}

	if (*model.Properties)["properties"].Variants == nil {
		t.Fatalf("expected properties variants, got none")
	}

	if v, ok := (*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"]; !ok || v == nil {
		t.Fatalf("expected properties variant AzureStorageLinkedService, got none")
	}

	if v := (*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].VariantType; v == nil || *v != "AzureStorage" {
		t.Fatalf("expected properties variant AzureStorageLinkedService variant type AzureStorage")
	}

	if (*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties == nil {
		t.Fatalf("expected properties variant AzureStorageLinkedService properties, got none")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties)["type"]; !ok || v == nil {
		t.Fatalf("expected properties variant AzureStorageLinkedService type property, got none")
	}

	if !(*(*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties)["type"].IsAnyCovered {
		t.Fatalf("expected properties variant AzureStorageLinkedService type IsAnyCovered true, got false")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties)["typeProperties"]; !ok || v == nil {
		t.Fatalf("expected properties variant AzureStorageLinkedService typeProperties property, got none")
	}

	if !(*(*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties)["typeProperties"].IsAnyCovered {
		t.Fatalf("expected properties variant AzureStorageLinkedService typeProperties IsAnyCovered true, got false")
	}

	if (*(*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties)["typeProperties"].Properties == nil {
		t.Fatalf("expected properties variant AzureStorageLinkedService typeProperties properties, got none")
	}

	if v, ok := (*(*(*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties)["typeProperties"].Properties)["connectionString"]; !ok || v == nil {
		t.Fatalf("expected properties variant AzureStorageLinkedService typeProperties connectionString property, got none")
	}

	if !(*(*(*(*model.Properties)["properties"].Variants)["AzureStorageLinkedService"].Properties)["typeProperties"].Properties)["connectionString"].IsAnyCovered {
		t.Fatalf("expected properties variant AzureStorageLinkedService typeProperties connectionString IsAnyCovered true, got false")
	}
}

func testCoverage(t *testing.T, tc testCase) (*coverage.Model, error) {
	apiVersion := strings.Split(tc.resourceType, "@")[1]

	var swaggerModel *coverage.SwaggerModel
	var err error
	azureRepoDir := os.Getenv("AZURE_REST_REPO_DIR")
	if azureRepoDir != "" {
		if strings.HasSuffix(azureRepoDir, "specification") {
			azureRepoDir += "/"
		}
		if !strings.HasSuffix(normarlizePath(azureRepoDir), "specification/") {
			t.Fatalf("AZURE_REST_REPO_DIR must specify the specification folder, e.g., AZURE_REST_REPO_DIR=\"/home/test/go/src/github.com/azure/azure-rest-api-specs/specification/\"")
		}

		t.Logf("AZURE_REST_REPO_DIR: %s", azureRepoDir)

		swaggerModel, err = coverage.GetModelInfoFromLocalIndex(
			tc.apiPath,
			apiVersion,
			tc.method,
			azureRepoDir,
			"./testdata/index.json",
		)
	} else {
		swaggerModel, err = coverage.GetModelInfoFromIndex(
			tc.apiPath,
			apiVersion,
			tc.method,
			"./testdata/index.json",
		)
	}
	t.Logf("swaggerModel: %+v", swaggerModel)

	if err != nil {
		return nil, fmt.Errorf("get model info from index: %+v", err)
	}

	model, err := coverage.Expand(swaggerModel.ModelName, swaggerModel.SwaggerPath)
	if err != nil {
		return nil, fmt.Errorf("expand model: %+v", err)
	}

	// out, err := json.MarshalIndent(model, "", "\t")
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Logf("expand model %s", string(out))

	for _, rq := range tc.rawRequest {
		request := map[string]interface{}{}
		err = json.Unmarshal([]byte(rq), &request)
		if err != nil {
			t.Errorf("error unmarshal request json: %v", err)
		}

		model.MarkCovered(request)
	}

	model.CountCoverage()

	if model.RootCoveredCount != tc.expectedCoveredCount {
		t.Errorf("expected CoveredCount %d, got %d", tc.expectedCoveredCount, model.RootCoveredCount)
	}

	if model.RootTotalCount != tc.expectedTotalCount {
		t.Errorf("expected TotalCount %d, got %d", tc.expectedTotalCount, model.RootTotalCount)
	}

	// out, err = json.MarshalIndent(model, "", "\t")
	// if err != nil {
	// 	t.Error(err)
	// }
	// t.Logf("coverage model %s", string(out))

	coverageReport := coverage.CoverageReport{
		Coverages: map[string]*coverage.CoverageItem{
			swaggerModel.ApiPath: {
				ApiPath:     swaggerModel.ApiPath,
				DisplayName: tc.resourceType,
				Model:       model,
			},
		},
	}

	passReport := types.PassReport{
		Resources: []types.Resource{
			{
				Type:    tc.resourceType,
				Address: "azapi_resource.test",
			},
		},
	}

	storeCoverageReport(passReport, coverageReport, "./testdata/", fmt.Sprintf("test_coverage_report_%s.md", tc.name))

	return model, nil
}

func storeCoverageReport(passReport types.PassReport, coverageReport coverage.CoverageReport, reportDir string, reportName string) {
	if len(coverageReport.Coverages) != 0 {
		err := os.WriteFile(path.Join(reportDir, reportName), []byte(report.PassedMarkdownReport(passReport, coverageReport)), 0644)
		if err != nil {
			logrus.Warnf("failed to save passed markdown report to %s: %+v", reportName, err)
		} else {
			logrus.Infof("markdown report saved to %s", reportName)
		}
	}
}

func testCredScan(t *testing.T, tc testCase) (*map[string]string, error) {
	apiVersion := strings.Split(tc.resourceType, "@")[1]
	swaggerModel, err := coverage.GetModelInfoFromIndex(
		tc.apiPath,
		apiVersion,
		"PUT",
		"",
	)

	t.Logf("swaggerModel: %+v", swaggerModel)

	if err != nil {
		return nil, fmt.Errorf("get model info from index: %+v", err)
	}

	model, err := coverage.Expand(swaggerModel.ModelName, swaggerModel.SwaggerPath)
	if err != nil {
		return nil, fmt.Errorf("expand model: %+v", err)
	}

	out, err := json.MarshalIndent(model, "", "\t")
	if err != nil {
		t.Error(err)
	}
	t.Logf("expand model %s", string(out))

	secrets := map[string]string{}

	if len(tc.rawRequest) != 1 {
		return nil, fmt.Errorf("TestCredScan rawRequest should be of length 1")
	}

	rq := tc.rawRequest[0]
	request := map[string]interface{}{}
	err = json.Unmarshal([]byte(rq), &request)
	if err != nil {
		t.Error(err)
	}

	model.CredScan(request, secrets)

	return &secrets, nil
}

func TestCredScan(t *testing.T) {
	tc := testCase{
		name:         "VM",
		resourceType: "Microsoft.Compute/virtualMachines@2023-03-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm",
		rawRequest: []string{
			`{
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

	secrets, err := testCredScan(t, tc)
	if err != nil {
		t.Fatalf("process coverage: %+v", err)
	}

	t.Logf("secrets: %+v", secrets)

	if _, ok := (*secrets)["#.properties.osProfile.adminPassword"]; !ok {
		t.Fatalf("expected adminPassword to be secrets")
	}
}
