## Microsoft.Insights/dataCollectionRules@2022-06-01 - Error

### Description

I found an error when creating this resource:

```bash
 "Resource: (ResourceId \"/subscriptions/{subscription_id}/resourceGroups/acctest7383/providers/Microsoft.Insights/dataCollectionRules/acctest8189\" / Api Version \"2022-06-01\")": PUT https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest7383/providers/Microsoft.Insights/dataCollectionRules/acctest8189
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: InvalidPayload
--------------------------------------------------------------------------------
{
  "error": {
    "code": "InvalidPayload",
    "message": "Data collection rule is invalid",
    "details": [
      {
        "code": "InvalidOutputTable",
        "message": "Table for output stream 'Microsoft-WindowsEvent' is not available for destination 'centralWorkspace'.",
        "target": "properties.dataFlows[0]"
      }
    ]
  }
}
--------------------------------------------------------------------------
```

### Details

1. ARM Fully-Qualified Resource Type
```
Microsoft.Insights/dataCollectionRules
```

2. API Version
```
2022-06-01
```

3. Swagger issue type
```
Other
```

4. OperationId
```
TODO
```

5. Swagger GitHub permalink
```
TODO, 
e.g., https://github.com/Azure/azure-rest-api-specs/blob/60723d13309c8f8060d020a7f3dd9d6e380f0bbd
/specification/compute/resource-manager/Microsoft.Compute/stable/2020-06-01/compute.json#L9065-L9101
```

6. Error code
```
TODO
```

7. Request traces
```
GET https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest7383/providers/Microsoft.Insights/dataCollectionRules/acctest8189?api-version=2022-06-01
   Accept: application/json
   Authorization: REDACTED
   User-Agent: HashiCorp Terraform/1.3.7 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/v1.4.0 pid-222c6c49-1b0a-5959-a213-6608f9eb8820
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
   RESPONSE Status: 404 Not Found
   Cache-Control: no-cache
   Content-Length: 233
   Content-Type: application/json; charset=utf-8
   Date: Mon, 27 Mar 2023 05:36:43 GMT
   Expires: -1
   Pragma: no-cache
   Strict-Transport-Security: REDACTED
   X-Content-Type-Options: REDACTED
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Failure-Cause: REDACTED
   X-Ms-Request-Id: 9ac6d1e2-8774-4233-b60b-d3248971b796
   X-Ms-Routing-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{"error":{"code":"ResourceNotFound","message":"The Resource 'Microsoft.Insights/dataCollectionRules/acctest8189' under resource group 'acctest7383' was not found. For more details please go to https://aka.ms/ARMResourceNotFoundFix"}}
   --------------------------------------------------------------------------------


PUT https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest7383/providers/Microsoft.Insights/dataCollectionRules/acctest8189?api-version=2022-06-01
   Accept: application/json
   Authorization: REDACTED
   Content-Length: 1507
   Content-Type: application/json
   User-Agent: HashiCorp Terraform/1.3.7 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/v1.4.0 pid-222c6c49-1b0a-5959-a213-6608f9eb8820
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{"location":"westeurope","name":"acctest8189","properties":{"dataFlows":[{"destinations":["centralWorkspace"],"streams":["Microsoft-Perf","Microsoft-Syslog","Microsoft-WindowsEvent"]}],"dataSources":{"performanceCounters":[{"counterSpecifiers":["\\Processor(_Total)\\% Processor Time","\\Memory\\Committed Bytes","\\LogicalDisk(_Total)\\Free Megabytes","\\PhysicalDisk(_Total)\\Avg. Disk Queue Length"],"name":"cloudTeamCoreCounters","samplingFrequencyInSeconds":15,"streams":["Microsoft-Perf"]},{"counterSpecifiers":["\\Process(_Total)\\Thread Count"],"name":"appTeamExtraCounters","samplingFrequencyInSeconds":30,"streams":["Microsoft-Perf"]}],"syslog":[{"facilityNames":["cron"],"logLevels":["Debug","Critical","Emergency"],"name":"cronSyslog","streams":["Microsoft-Syslog"]},{"facilityNames":["syslog"],"logLevels":["Alert","Critical","Emergency"],"name":"syslogBase","streams":["Microsoft-Syslog"]}],"windowsEventLogs":[{"name":"cloudSecurityTeamEvents","streams":["Microsoft-WindowsEvent"],"xPathQueries":["Security![System[(Level = 1 or Level = 2 or Level = 3)]]"]},{"name":"appTeam1AppEvents","streams":["Microsoft-WindowsEvent"],"xPathQueries":["System![System[(Level = 1 or Level = 2 or Level = 3)]]","Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]"]}]},"destinations":{"logAnalytics":[{"name":"centralWorkspace","workspaceResourceId":"/subscriptions/{subscription_id}/resourceGroups/acctest7383/providers/Microsoft.OperationalInsights/workspaces/acctest9962"}]}}}
   --------------------------------------------------------------------------------

RESPONSE Status: 400 Bad Request
   Api-Supported-Versions: REDACTED
   Cache-Control: no-cache
   Content-Length: 270
   Content-Type: application/json
   Date: Mon, 27 Mar 2023 05:36:49 GMT
   Expires: -1
   Pragma: no-cache
   Request-Context: REDACTED
   Server: Microsoft-HTTPAPI/2.0
   Strict-Transport-Security: REDACTED
   X-Content-Type-Options: REDACTED
   X-Ms-Client-Request-Id: 24229245-2c39-4c45-9bf0-3562844839d7
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests: REDACTED
   X-Ms-Request-Id: 1b4daa73-584a-4bf0-a5d3-e4f3a09c0a2e
   X-Ms-Routing-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{"error":{"code":"InvalidPayload","message":"Data collection rule is invalid","details":[{"code":"InvalidOutputTable","message":"Table for output stream 'Microsoft-WindowsEvent' is not available for destination 'centralWorkspace'.","target":"properties.dataFlows[0]"}]}}
   --------------------------------------------------------------------------------



```

### Links
1. [Semantic and Model Violations Reference](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/Semantic-and-Model-Violations-Reference.md)
2. [S360 action item generator for Swagger issues](https://aka.ms/swaggers360)