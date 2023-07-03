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
	name         string
	apiVersion   string
	apiPath      string
	rawRequest   []string
	resourceType string
}

func TestCoverageResourceGroup(t *testing.T) {
	tc := testCase{
		name:         "ResourceGroup",
		resourceType: "Microsoft.Resources/resourceGroups@2022-09-01",
		apiVersion:   "2022-09-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rgName",
		rawRequest: []string{
			`{"location": "westeurope"}`,
		},
	}

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 1 {
		t.Errorf("expected CoveredCount 1, got %d", model.CoveredCount)
	}

	if model.TotalCount != 3 {
		t.Errorf("expected TotalCount 3, got %d", model.TotalCount)
	}

	if model.Properties == nil {
		t.Errorf("expected properties, got none")
	}

	if v, ok := (*model.Properties)["location"]; !ok || v == nil {
		t.Errorf("expected location, got none")
	}

	if !(*model.Properties)["location"].IsAnyCovered {
		t.Errorf("expected location IsAnyCovered true, got false")
	}
}

func TestCoverageKeyVault(t *testing.T) {
	tc := testCase{
		name:         "KeyVault",
		resourceType: "Microsoft.KeyVault/vaults@2023-02-01",
		apiVersion:   "2023-02-01",
		apiPath:      "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/sample-resource-group/providers/Microsoft.KeyVault/vaults/sample-vault",
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 13 {
		t.Errorf("expected CoveredCount 13, got %d", model.CoveredCount)
	}

	if model.TotalCount != 28 {
		t.Errorf("expected TotalCount 28, got %d", model.TotalCount)
	}
}

func TestCoverageStorageAccount(t *testing.T) {
	tc := testCase{
		name:         "StorageAccount",
		resourceType: "Microsoft.Storage/storageAccounts@2022-09-01",
		apiVersion:   "2022-09-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/res9101/providers/Microsoft.Storage/storageAccounts/sto4445",
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
	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 24 {
		t.Errorf("expected CoveredCount 24, got %d", model.CoveredCount)
	}

	if model.TotalCount != 69 {
		t.Errorf("expected TotalCount 69, got %d", model.TotalCount)
	}
}

func TestCoverageVM(t *testing.T) {
	tc := testCase{
		name:         "VM",
		resourceType: "Microsoft.Compute/virtualMachines@2023-03-01",
		apiVersion:   "2023-03-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Compute/virtualMachines/vm",
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 20 {
		t.Errorf("expected CoveredCount 20, got %d", model.CoveredCount)
	}

	if model.TotalCount != 155 {
		t.Errorf("expected TotalCount 155, got %d", model.TotalCount)
	}
}

func TestCoverageVNet(t *testing.T) {
	tc := testCase{
		name:         "VNet",
		resourceType: "Microsoft.Network/virtualNetworks@2023-02-01",
		apiVersion:   "2023-02-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/virtualNetwork",
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 4 {
		t.Errorf("expected CoveredCount 4, got %d", model.CoveredCount)
	}

	if model.TotalCount != 95 {
		t.Errorf("expect TotalCount 95, got %d", model.TotalCount)
	}
}

func TestCoverageDataCollectionRule(t *testing.T) {
	tc := testCase{
		name:         "DataCollectionRule",
		resourceType: "Microsoft.Insights/dataCollectionRules@2022-06-01",
		apiVersion:   "2022-06-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/test-resources/providers/Microsoft.Insights/dataCollectionRules/testDCR",
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 65 {
		t.Errorf("expected CoveredCount 65, got %d", model.CoveredCount)
	}

	if model.TotalCount != 65 {
		t.Errorf("expected TotalCount 65, got %d", model.TotalCount)
	}
}

func TestCoverageWebSite(t *testing.T) {
	tc := testCase{
		name:         "WebSites",
		resourceType: "Microsoft.Web/sites@2022-09-01",
		apiVersion:   "2022-09-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/testrg123/providers/Microsoft.Web/sites/sitef6141",
		rawRequest: []string{`{
    "kind": "app",
    "location": "East US",
    "properties": {
        "serverFarmId": "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/testrg123/providers/Microsoft.Web/serverfarms/DefaultAsp"
    }
}`,
		},
	}

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 3 {
		t.Errorf("expected CoveredCount 3, got %d", model.CoveredCount)
	}

	if model.TotalCount != 186 {
		t.Errorf("expected TotalCount 186, got %d", model.TotalCount)
	}
}
func TestCoverageAKS(t *testing.T) {
	tc := testCase{
		name:         "AKS",
		resourceType: "Microsoft.ContainerService/ManagedClusters@2023-05-02-preview",
		apiVersion:   "2023-05-02-preview",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourcegroups/rg1/providers/Microsoft.ContainerService/managedClusters/clustername1",
		rawRequest: []string{`{
    "location": "location1",
    "tags": {
        "tier": "production",
        "archv2": ""
    },
    "sku": {
        "name": "Basic",
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 33 {
		t.Errorf("expected TotalCount 33, got %d", model.CoveredCount)
	}

	if model.TotalCount != 242 {
		t.Errorf("expected TotalCount 242, got %d", model.TotalCount)
	}
}

func TestCoverageCosmosDB(t *testing.T) {
	tc := testCase{
		name:         "CosmosDB",
		resourceType: "Microsoft.DocumentDB/databaseAccounts@2023-04-15",
		apiVersion:   "2023-04-15",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.DocumentDB/databaseAccounts/testdb",
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %v", err)
	}

	if model.CoveredCount != 34 {
		t.Errorf("expected CoveredCount 34, got %d", model.CoveredCount)
	}

	if model.TotalCount != 65 {
		t.Errorf("expected TotalCount 65, got %d", model.TotalCount)
	}

	if model.Properties == nil {
		t.Errorf("expected properties, got none")
	}

	if v, ok := (*model.Properties)["identity"]; !ok || v == nil {
		t.Errorf("expected identity property, got none")
	}

	if (*model.Properties)["identity"].CoveredCount != 2 {
		t.Errorf("expected identity CoveredCount 2, got %d", (*model.Properties)["identity"].CoveredCount)
	}

	if (*model.Properties)["identity"].Properties == nil {
		t.Errorf("expected identity properties, got none")
	}

	if v, ok := (*(*model.Properties)["identity"].Properties)["type"]; !ok || v == nil {
		t.Errorf("expected identity type property, got none")
	}

	if !(*(*model.Properties)["identity"].Properties)["type"].IsAnyCovered {
		t.Errorf("expected identity type IsAnyCovered true, got false")
	}

	if !(*(*model.Properties)["identity"].Properties)["type"].IsFullyCovered {
		t.Errorf("expected identity type IsFullyCovered true, got false")
	}

	if (*(*model.Properties)["identity"].Properties)["type"].EnumCoveredCount != 1 {
		t.Errorf("expected identity type EnumCoveredCount 1, got %d", (*(*model.Properties)["identity"].Properties)["type"].EnumCoveredCount)
	}

	if (*(*model.Properties)["identity"].Properties)["type"].Enum == nil {
		t.Errorf("expected identity type Enum, got none")
	}

	if isSet, ok := (*(*(*model.Properties)["identity"].Properties)["type"].Enum)["SystemAssigned,UserAssigned"]; !ok || !isSet {
		t.Errorf("expected identity type Enum SystemAssigned,UserAssigned to be set")
	}

	if v, ok := (*model.Properties)["tags"]; !ok || v == nil {
		t.Errorf("expected identity tags, got none")
	}

	if !(*model.Properties)["tags"].HasAdditionalProperties {
		t.Errorf("expected tags HasAdditionalProperties true, got false")
	}

	if !(*model.Properties)["tags"].IsAnyCovered {
		t.Errorf("expected tags IsAnyCovered true, got false")
	}

	if (*model.Properties)["tags"].TotalCount != 1 {
		t.Errorf("expected tags TotalCount 1, got %d", (*model.Properties)["tags"].TotalCount)
	}

	if (*model.Properties)["tags"].CoveredCount != 1 {
		t.Errorf("expected tags CoveredCount 1, got %d", (*model.Properties)["tags"].CoveredCount)
	}

	if v, ok := (*model.Properties)["properties"]; !ok || v == nil {
		t.Errorf("expected identity properties, got none")
	}

	if !(*model.Properties)["properties"].IsAnyCovered {
		t.Errorf("expected properties IsAnyCovered true, got false")
	}

	if (*model.Properties)["properties"].Properties == nil {
		t.Errorf("expected properties properties, got none")
	}

	if v, ok := (*(*model.Properties)["properties"].Properties)["locations"]; !ok || v == nil {
		t.Errorf("expected properties locations, got none")
	}

	if !(*(*model.Properties)["properties"].Properties)["locations"].IsAnyCovered {
		t.Errorf("expected locations IsAnyCovered true, got false")
	}

	if (*(*model.Properties)["properties"].Properties)["locations"].Item == nil {
		t.Errorf("expected locations Item, got none")
	}

	if (*(*model.Properties)["properties"].Properties)["locations"].Item.Properties == nil {
		t.Errorf("expected locations Item properties, got none")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["locationName"]; !ok || v == nil {
		t.Errorf("expected locations Item locationName, got none")
	}

	if !(*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["locationName"].IsAnyCovered {
		t.Errorf("expected locationName IsAnyCovered true, got false")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"]; !ok || v == nil {
		t.Errorf("expected locations Item isZoneRedundant, got none")
	}

	if !(*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"].IsAnyCovered {
		t.Errorf("expected isZoneRedundant IsAnyCovered true, got false")
	}

	if (*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"].Bool == nil {
		t.Errorf("expected isZoneRedundant Bool, got none")
	}

	if isSet, ok := (*(*(*(*model.Properties)["properties"].Properties)["locations"].Item.Properties)["isZoneRedundant"].Bool)["false"]; !ok || !isSet {
		t.Errorf("expected isZoneRedundant Bool false to be set")
	}

	if v, ok := (*(*model.Properties)["properties"].Properties)["ipRules"]; !ok || v == nil {
		t.Errorf("expected properties ipRules, got none")
	}

	if (*(*model.Properties)["properties"].Properties)["ipRules"].Item == nil {
		t.Errorf("expected ipRules Item, got none")
	}

	if (*(*model.Properties)["properties"].Properties)["ipRules"].Item.Properties == nil {
		t.Errorf("expected ipRules Item properties, got none")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Properties)["ipRules"].Item.Properties)["ipAddressOrRange"]; !ok || v == nil {
		t.Errorf("expected ipRules Item ipAddressOrRange, got none")
	}

}

func TestCoverageDataFactoryLinkedServices(t *testing.T) {
	tc := testCase{
		name:         "DataFactoryLinkedServices",
		resourceType: "Microsoft.DataFactory/factories/linkedServices@2018-06-01",
		apiVersion:   "2018-06-01",
		apiPath:      "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/rg1/providers/Microsoft.DataFactory/factories/factory1/linkedServices/linked",
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

	model, err := testCoverage(t, tc)
	if err != nil {
		t.Errorf("process coverage: %+v", err)
	}

	if model.CoveredCount != 3 {
		t.Errorf("expected TotalCount 3, got %d", model.CoveredCount)
	}

	if model.TotalCount != 2857 {
		t.Errorf("expected TotalCount 2857, got %d", model.TotalCount)
	}

	if model.Properties == nil {
		t.Errorf("expected properties, got none")
	}

	if (*model.Properties)["properties"].Properties == nil {
		t.Errorf("expected properties properties, got none")
	}

	if v, ok := (*(*model.Properties)["properties"].Properties)["type"]; !ok || v == nil {
		t.Errorf("expected properties type property, got none")
	}

	if !(*(*model.Properties)["properties"].Properties)["type"].IsAnyCovered {
		t.Errorf("expected properties type IsAnyCovered true, got false")
	}

	if (*model.Properties)["properties"].Discriminator == nil {
		t.Errorf("expected properties discriminator, got none")
	}

	if *(*model.Properties)["properties"].Discriminator != "type" {
		t.Errorf("expected properties discriminator 'type', got %s", *(*model.Properties)["properties"].Discriminator)
	}

	if (*model.Properties)["properties"].Variants == nil {
		t.Errorf("expected properties variants, got none")
	}

	if v, ok := (*(*model.Properties)["properties"].Variants)["AzureStorage"]; !ok || v == nil {
		t.Errorf("expected properties variant AzureStorage, got none")
	}

	if (*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties == nil {
		t.Errorf("expected properties variant AzureStorage properties, got none")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties)["type"]; !ok || v == nil {
		t.Errorf("expected properties variant AzureStorage type property, got none")
	}

	if !(*(*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties)["type"].IsAnyCovered {
		t.Errorf("expected properties variant AzureStorage type IsAnyCovered true, got false")
	}

	if v, ok := (*(*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties)["typeProperties"]; !ok || v == nil {
		t.Errorf("expected properties variant AzureStorage typeProperties property, got none")
	}

	if !(*(*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties)["typeProperties"].IsAnyCovered {
		t.Errorf("expected properties variant AzureStorage typeProperties IsAnyCovered true, got false")
	}

	if (*(*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties)["typeProperties"].Properties == nil {
		t.Errorf("expected properties variant AzureStorage typeProperties properties, got none")
	}

	if v, ok := (*(*(*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties)["typeProperties"].Properties)["connectionString"]; !ok || v == nil {
		t.Errorf("expected properties variant AzureStorage typeProperties connectionString property, got none")
	}

	if !(*(*(*(*model.Properties)["properties"].Variants)["AzureStorage"].Properties)["typeProperties"].Properties)["connectionString"].IsAnyCovered {
		t.Errorf("expected properties variant AzureStorage typeProperties connectionString IsAnyCovered true, got false")
	}
}

func testCoverage(t *testing.T, tc testCase) (*coverage.Model, error) {
	apiPath, modelName, modelSwaggerPath, err := coverage.GetModelInfoFromIndex(
		tc.apiPath,
		tc.apiVersion,
	)

	if err != nil {
		return nil, fmt.Errorf("get model info from index: %+v", err)
	}

	model, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		return nil, fmt.Errorf("expand model: %+v", err)
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
		Coverages: map[types.Resource]*coverage.Model{
			types.Resource{
				ApiPath: *apiPath,
				Type:    tc.resourceType,
				Address: "azapi_resource.test",
			}: model,
		},
	}

	storeCoverageReport(coverageReport, ".", fmt.Sprintf("test_coverage_report_%s.md", tc.name))

	return model, nil
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
