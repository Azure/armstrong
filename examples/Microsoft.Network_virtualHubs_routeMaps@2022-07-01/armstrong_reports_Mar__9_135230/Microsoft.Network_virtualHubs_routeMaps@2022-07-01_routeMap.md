## Microsoft.Network/virtualHubs/routeMaps@2022-07-01 - ROUNDTRIP_MISSING_PROPERTY

### Description

I found differences between PUT request body and GET response:

- .properties.associatedInboundConnections: expect 1 in length, but got 0

```json
{
    "properties": {
        "associatedInboundConnections": [
            "/subscriptions/{subscription_id}/resourceGroups/acctest6981/providers/Microsoft.Network/expressRouteGateways/acctest7416/expressRouteConnections/acctest2735" is not returned from response
        ],
        "associatedOutboundConnections": [],
        "rules": [
            {
                "actions": [
                    {
                        "parameters": [
                            {
                                "asPath": [
                                    "22334"
                                ],
                                "community": [],
                                "routePrefix": []
                            }
                        ],
                        "type": "Add"
                    }
                ],
                "matchCriteria": [
                    {
                        "asPath": [],
                        "community": [],
                        "matchCondition": "Contains",
                        "routePrefix": [
                            "10.0.0.0/8"
                        ]
                    }
                ],
                "name": "rule1",
                "nextStepIfMatched": "Continue"
            }
        ]
    }
}
```

### Details

1. ARM Fully-Qualified Resource Type
```
Microsoft.Network/virtualHubs/routeMaps
```

2. API Version
```
2022-07-01
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
PUT https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest6981/providers/Microsoft.Network/virtualHubs/acctest3106/routeMaps/acctest836?api-version=2022-07-01
   Accept: application/json
   Authorization: REDACTED
   Content-Length: 539
   Content-Type: application/json
   User-Agent: HashiCorp Terraform/1.3.7 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/v1.4.0 pid-222c6c49-1b0a-5959-a213-6608f9eb8820
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{"name":"acctest836","properties":{"associatedInboundConnections":["/subscriptions/{subscription_id}/resourceGroups/acctest6981/providers/Microsoft.Network/expressRouteGateways/acctest7416/expressRouteConnections/acctest2735"],"associatedOutboundConnections":[],"rules":[{"actions":[{"parameters":[{"asPath":["22334"],"community":[],"routePrefix":[]}],"type":"Add"}],"matchCriteria":[{"asPath":[],"community":[],"matchCondition":"Contains","routePrefix":["10.0.0.0/8"]}],"name":"rule1","nextStepIfMatched":"Continue"}]}}
   --------------------------------------------------------------------------------

RESPONSE Status: 201 Created
   Azure-Asyncnotification: REDACTED
   Azure-Asyncoperation: REDACTED
   Cache-Control: no-cache
   Content-Length: 1107
   Content-Type: application/json; charset=utf-8
   Date: Thu, 09 Mar 2023 05:51:38 GMT
   Expires: -1
   Pragma: no-cache
   Retry-After: 10
   Server: Microsoft-HTTPAPI/2.0
   Strict-Transport-Security: REDACTED
   X-Content-Type-Options: REDACTED
   X-Ms-Arm-Service-Request-Id: REDACTED
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Ratelimit-Remaining-Subscription-Writes: REDACTED
   X-Ms-Request-Id: 955e6202-ad27-4ecf-a9d8-600080e6921f
   X-Ms-Routing-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{
  "name": "acctest836",
  "id": "/subscriptions/{subscription_id}/resourceGroups/acctest6981/providers/Microsoft.Network/virtualHubs/acctest3106/routeMaps/acctest836",
  "etag": "W/\"e02461c8-ef68-42b2-94e2-1804c1cf4a7f\"",
  "properties": {
    "provisioningState": "Updating",
    "rules": [
      {
        "name": "rule1",
        "matchCriteria": [
          {
            "matchCondition": "Contains",
            "routePrefix": [
              "10.0.0.0/8"
            ],
            "community": [],
            "asPath": []
          }
        ],
        "actions": [
          {
            "type": "Add",
            "parameters": [
              {
                "routePrefix": [],
                "community": [],
                "asPath": [
                  "22334"
                ]
              }
            ]
          }
        ],
        "nextStepIfMatched": "Continue"
      }
    ],
    "associatedInboundConnections": [],
    "associatedOutboundConnections": []
  },
  "type": "Microsoft.Network/virtualHubs/routeMaps"
}
   --------------------------------------------------------------------------------


GET https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest6981/providers/Microsoft.Network/virtualHubs/acctest3106/routeMaps/acctest836?api-version=2022-07-01
   Accept: application/json
   Authorization: REDACTED
   User-Agent: HashiCorp Terraform/1.3.7 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/v1.4.0 pid-222c6c49-1b0a-5959-a213-6608f9eb8820
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
   RESPONSE Status: 200 OK
   Cache-Control: no-cache
   Content-Type: application/json; charset=utf-8
   Date: Thu, 09 Mar 2023 05:52:29 GMT
   Etag: W/"d022c88f-2bba-4efd-92dd-d02f27c9b041"
   Expires: -1
   Pragma: no-cache
   Server: Microsoft-HTTPAPI/2.0
   Strict-Transport-Security: REDACTED
   Vary: REDACTED
   X-Content-Type-Options: REDACTED
   X-Ms-Arm-Service-Request-Id: REDACTED
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Ratelimit-Remaining-Subscription-Reads: REDACTED
   X-Ms-Request-Id: 4559eadc-bc4c-4c75-9511-1ae0e8a7b03f
   X-Ms-Routing-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{
  "name": "acctest836",
  "id": "/subscriptions/{subscription_id}/resourceGroups/acctest6981/providers/Microsoft.Network/virtualHubs/acctest3106/routeMaps/acctest836",
  "etag": "W/\"d022c88f-2bba-4efd-92dd-d02f27c9b041\"",
  "properties": {
    "provisioningState": "Succeeded",
    "rules": [
      {
        "name": "rule1",
        "matchCriteria": [
          {
            "matchCondition": "Contains",
            "routePrefix": [
              "10.0.0.0/8"
            ],
            "community": [],
            "asPath": []
          }
        ],
        "actions": [
          {
            "type": "Add",
            "parameters": [
              {
                "routePrefix": [],
                "community": [],
                "asPath": [
                  "22334"
                ]
              }
            ]
          }
        ],
        "nextStepIfMatched": "Continue"
      }
    ],
    "associatedInboundConnections": [],
    "associatedOutboundConnections": []
  },
  "type": "Microsoft.Network/virtualHubs/routeMaps"
}
   --------------------------------------------------------------------------------
```

### Links
1. [Semantic and Model Violations Reference](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/Semantic-and-Model-Violations-Reference.md)
2. [S360 action item generator for Swagger issues](https://aka.ms/swaggers360)
