## Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08 - ROUNDTRIP_MISSING_PROPERTY

### Description

I found differences between PUT request body and GET response:



```json
{
    "location": "westus",
    "properties": {
        "pointInTimeUTC": "2023-03-23T08:30:37.467Z" is not returned from response,
        "sourceLocation": "westus" is not returned from response,
        "sourceResourceId": "/subscriptions/{subscription_id}/resourceGroups/acctest4304/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/acctest9810" is not returned from response
    }
}
```

### Details

1. ARM Fully-Qualified Resource Type
```
Microsoft.DBforPostgreSQL/serverGroupsv2
```

2. API Version
```
2022-11-08
```

3. Swagger issue type
```
Swagger Correctness
```

4. OperationId
```
TODO
e.g., VirtualMachines_Get
```

5. Swagger GitHub permalink
```
TODO, 
e.g., https://github.com/Azure/azure-rest-api-specs/blob/60723d13309c8f8060d020a7f3dd9d6e380f0bbd
/specification/compute/resource-manager/Microsoft.Compute/stable/2020-06-01/compute.json#L9065-L9101
```

6. Error code
```
ROUNDTRIP_MISSING_PROPERTY
```

7. Request traces
```
PUT https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest4304/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/acctest1681?api-version=2022-11-08
   Accept: application/json
   Authorization: REDACTED
   Content-Length: 290
   Content-Type: application/json
   User-Agent: HashiCorp Terraform/1.3.7 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/v1.4.0 pid-222c6c49-1b0a-5959-a213-6608f9eb8820
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{"location":"westus","name":"acctest1681","properties":{"pointInTimeUTC":"2023-03-23T08:30:37.467Z","sourceLocation":"westus","sourceResourceId":"/subscriptions/{subscription_id}/resourceGroups/acctest4304/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/acctest9810"}}
   --------------------------------------------------------------------------------

RESPONSE Status: 201 Created
   Azure-Asyncoperation: REDACTED
   Cache-Control: no-cache
   Content-Length: 2002
   Content-Security-Policy: REDACTED
   Content-Type: application/json; charset=utf-8
   Date: Thu, 23 Mar 2023 08:36:48 GMT
   Expires: -1
   Location: REDACTED
   Pragma: no-cache
   Referrer-Policy: REDACTED
   Request-Id: ba2246ad-d896-41fd-98c6-821bcdb13326
   Strict-Transport-Security: REDACTED
   Vary: REDACTED
   X-Content-Type-Options: REDACTED
   X-Download-Options: REDACTED
   X-Frame-Options: REDACTED
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Ratelimit-Remaining-Subscription-Writes: REDACTED
   X-Ms-Request-Id: 52046e79-98fd-45af-85cb-eb2724502836
   X-Ms-Routing-Request-Id: REDACTED
   X-Permitted-Cross-Domain-Policies: REDACTED
   X-Xss-Protection: REDACTED
   --------------------------------------------------------------------------------
{
  "id": "/subscriptions/{subscription_id}/resourceGroups/acctest4304/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/acctest1681",
  "name": "acctest1681",
  "type": "Microsoft.DBforPostgreSQL/serverGroupsv2",
  "tags": {
},
  "location": "westus",
  "systemData": {
  "createdAt": "2023-03-23T08:36:47.1114738Z",
  "createdBy": "0ddf9866-48e9-4e1d-a880-17a2ea0c9ec6",
  "createdByType": "Application",
  "lastModifiedAt": "2023-03-23T08:36:47.1114738Z",
  "lastModifiedBy": "0ddf9866-48e9-4e1d-a880-17a2ea0c9ec6",
  "lastModifiedByType": "Application"
},
  "properties": {
    "provisioningState": "Provisioning",
    "state": "Provisioning",
    "administratorLogin": "citus",
    "postgresqlVersion": "15",
    "citusVersion": "11.1",
    "maintenanceWindow": {
      "customWindow": "Disabled",
      "dayOfWeek": 0,
      "startHour": 0,
      "startMinute": 0
    },
    "preferredPrimaryZone": null,
    "enableShardsOnCoordinator": false,
    "earliestRestoreTime": null,
    "sourceResourceId": null,
    "enableHa": false,
    "coordinatorServerEdition": "GeneralPurpose",
    "nodeServerEdition": "MemoryOptimized",
    "coordinatorStorageQuotaInMb": 524288,
    "nodeStorageQuotaInMb": 524288,
    "coordinatorVCores": 4,
    "nodeVCores": 8,
    "coordinatorEnablePublicIpAccess": true,
    "nodeEnablePublicIpAccess": false,
    "nodeCount": 3,
    "serverNames": [
      {
        "name": "acctest1681-c",
        "fullyQualifiedDomainName": "c.acctest1681.postgres.database.azure.com"
      },
      {
        "name": "acctest1681-w0",
        "fullyQualifiedDomainName": "w0.acctest1681.postgres.database.azure.com"
      },
      {
        "name": "acctest1681-w1",
        "fullyQualifiedDomainName": "w1.acctest1681.postgres.database.azure.com"
      },
      {
        "name": "acctest1681-w2",
        "fullyQualifiedDomainName": "w2.acctest1681.postgres.database.azure.com"
      }
    ],
    "readReplicas": [],
    "privateEndpointConnections": []
  }
}
   --------------------------------------------------------------------------------


GET https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest4304/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/acctest1681?api-version=2022-11-08
   Accept: application/json
   Authorization: REDACTED
   User-Agent: HashiCorp Terraform/1.3.7 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/v1.4.0 pid-222c6c49-1b0a-5959-a213-6608f9eb8820
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
   RESPONSE Status: 200 OK
   Cache-Control: no-cache
   Content-Security-Policy: REDACTED
   Content-Type: application/json; charset=utf-8
   Date: Thu, 23 Mar 2023 08:48:24 GMT
   Expires: -1
   Pragma: no-cache
   Referrer-Policy: REDACTED
   Request-Id: 32907678-21d5-44a3-89d7-1b2454e4bb07
   Strict-Transport-Security: REDACTED
   Vary: REDACTED
   X-Content-Type-Options: REDACTED
   X-Download-Options: REDACTED
   X-Frame-Options: REDACTED
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Ratelimit-Remaining-Subscription-Reads: REDACTED
   X-Ms-Request-Id: 623c5e9f-8059-433a-af52-b6839d410fe6
   X-Ms-Routing-Request-Id: REDACTED
   X-Permitted-Cross-Domain-Policies: REDACTED
   X-Xss-Protection: REDACTED
   --------------------------------------------------------------------------------
{
  "id": "/subscriptions/{subscription_id}/resourceGroups/acctest4304/providers/Microsoft.DBforPostgreSQL/serverGroupsv2/acctest1681",
  "name": "acctest1681",
  "type": "Microsoft.DBforPostgreSQL/serverGroupsv2",
  "tags": {
},
  "location": "westus",
  "systemData": {
  "createdAt": "2023-03-23T08:36:47.1114738Z",
  "createdBy": "0ddf9866-48e9-4e1d-a880-17a2ea0c9ec6",
  "createdByType": "Application",
  "lastModifiedAt": "2023-03-23T08:36:47.1114738Z",
  "lastModifiedBy": "0ddf9866-48e9-4e1d-a880-17a2ea0c9ec6",
  "lastModifiedByType": "Application"
},
  "properties": {
    "provisioningState": "Succeeded",
    "state": "Ready",
    "administratorLogin": "citus",
    "postgresqlVersion": "15",
    "citusVersion": "11.1",
    "maintenanceWindow": {
      "customWindow": "Disabled",
      "dayOfWeek": 0,
      "startHour": 0,
      "startMinute": 0
    },
    "preferredPrimaryZone": null,
    "enableShardsOnCoordinator": false,
    "earliestRestoreTime": null,
    "sourceResourceId": null,
    "enableHa": false,
    "coordinatorServerEdition": "GeneralPurpose",
    "nodeServerEdition": "MemoryOptimized",
    "coordinatorStorageQuotaInMb": 524288,
    "nodeStorageQuotaInMb": 524288,
    "coordinatorVCores": 4,
    "nodeVCores": 8,
    "coordinatorEnablePublicIpAccess": true,
    "nodeEnablePublicIpAccess": false,
    "nodeCount": 3,
    "serverNames": [
      {
        "name": "acctest1681-c",
        "fullyQualifiedDomainName": "c.acctest1681.postgres.database.azure.com"
      },
      {
        "name": "acctest1681-w0",
        "fullyQualifiedDomainName": "w0.acctest1681.postgres.database.azure.com"
      },
      {
        "name": "acctest1681-w1",
        "fullyQualifiedDomainName": "w1.acctest1681.postgres.database.azure.com"
      },
      {
        "name": "acctest1681-w2",
        "fullyQualifiedDomainName": "w2.acctest1681.postgres.database.azure.com"
      }
    ],
    "readReplicas": [],
    "privateEndpointConnections": []
  }
}
   --------------------------------------------------------------------------------
```

### Links
1. [Semantic and Model Violations Reference](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/Semantic-and-Model-Violations-Reference.md)
2. [S360 action item generator for Swagger issues](https://aka.ms/swaggers360)