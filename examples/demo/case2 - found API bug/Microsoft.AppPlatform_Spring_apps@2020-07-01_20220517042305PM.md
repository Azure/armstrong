## Microsoft.AppPlatform/Spring/apps@2020-07-01 - ROUNDTRIP_INCONSISTENT_PROPERTY && ROUNDTRIP_MISSING_PROPERTY

### Description

I found differences between PUT request body and GET response:
```json
{
    "location": "westeurope",
    "properties": {
        "activeDeploymentName": "mydeployment1" is not returned from response,
        "fqdn": Got "acctest5486.azuremicroservices.io" in response, expect "myapp.mydomain.com",
        "httpsOnly": false,
        "persistentDisk": {
            "mountPath": Got "/persistent" in response, expect "/mypersistentdisk",
            "sizeInGB": Got 0 in response, expect 2
        },
        "public": Got false in response, expect true,
        "temporaryDisk": {
            "mountPath": Got "/tmp" in response, expect "/mytemporarydisk",
            "sizeInGB": Got 5 in response, expect 2
        }
    }
}
```


### Reference

1. ARM Fully-Qualified Resource Type
```
Microsoft.AppPlatform/Spring/apps
```

2. API Version
```
2020-07-01
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
ROUNDTRIP_INCONSISTENT_PROPERTY
ROUNDTRIP_MISSING_PROPERTY
```

7. Request traces
```
PUT https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest5486/providers/Microsoft.AppPlatform/Spring/acctest5486/apps/acctest7675?api-version=2020-07-01
   Accept: application/json
   Authorization: REDACTED
   Content-Length: 265
   Content-Type: application/json
   User-Agent: HashiCorp Terraform/1.1.4 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/dev
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
{"location":"westeurope","properties":{"activeDeploymentName":"mydeployment1","fqdn":"myapp.mydomain.com","httpsOnly":false,"persistentDisk":{"mountPath":"/mypersistentdisk","sizeInGB":2},"public":true,"temporaryDisk":{"mountPath":"/mytemporarydisk","sizeInGB":2}}}
   --------------------------------------------------------------------------------

RESPONSE Status: 201 Created
   Azure-Asyncoperation: REDACTED
   Cache-Control: no-cache
   Content-Length: 463
   Content-Type: application/json; charset=utf-8
   Date: Tue, 17 May 2022 08:22:00 GMT
   Expires: -1
   Location: REDACTED
   Pragma: no-cache
   Request-Context: REDACTED
   Server: nginx/1.17.7
   Strict-Transport-Security: REDACTED
   X-Content-Type-Options: REDACTED
   X-Ms-Client-Request-Id: ce346f11-60d1-4deb-8548-6be594650502
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests: REDACTED
   X-Ms-Request-Id: 816ccc42-8215-477f-aad1-a24e43778fda
   X-Ms-Routing-Request-Id: REDACTED
   X-Rp-Server-Mvid: REDACTED
   --------------------------------------------------------------------------------
{"properties":{"public":true,"provisioningState":"Creating","httpsOnly":false,"temporaryDisk":{"sizeInGB":2,"mountPath":"/mytemporarydisk"},"persistentDisk":{"sizeInGB":2,"mountPath":"/mypersistentdisk"}},"type":"Microsoft.AppPlatform/Spring/apps","identity":null,"location":"westeurope","id":"/subscriptions/{subscription_id}/resourceGroups/acctest5486/providers/Microsoft.AppPlatform/Spring/acctest5486/apps/acctest7675","name":"acctest7675"}
   --------------------------------------------------------------------------------


GET https://management.azure.com/subscriptions/{subscription_id}/resourceGroups/acctest5486/providers/Microsoft.AppPlatform/Spring/acctest5486/apps/acctest7675?api-version=2020-07-01
   Accept: application/json
   Authorization: REDACTED
   User-Agent: HashiCorp Terraform/1.1.4 (+https://www.terraform.io) Terraform Plugin SDK/2.8.0 terraform-provider-azapi/dev
   X-Ms-Correlation-Request-Id: REDACTED
   --------------------------------------------------------------------------------
   RESPONSE Status: 200 OK
   Cache-Control: no-cache
   Content-Type: application/json; charset=utf-8
   Date: Tue, 17 May 2022 08:23:03 GMT
   Expires: -1
   Pragma: no-cache
   Request-Context: REDACTED
   Server: nginx/1.17.7
   Strict-Transport-Security: REDACTED
   Vary: REDACTED
   X-Content-Type-Options: REDACTED
   X-Ms-Client-Request-Id: d5ad3c8a-344d-43f6-b404-9da4311410b3
   X-Ms-Correlation-Request-Id: REDACTED
   X-Ms-Ratelimit-Remaining-Subscription-Resource-Requests: REDACTED
   X-Ms-Request-Id: 97567c5c-a233-4512-b9ff-ff5f70bbe797
   X-Ms-Routing-Request-Id: REDACTED
   X-Rp-Server-Mvid: REDACTED
   --------------------------------------------------------------------------------
{"properties":{"public":false,"provisioningState":"Succeeded","fqdn":"acctest5486.azuremicroservices.io","httpsOnly":false,"createdTime":"2022-05-17T08:22:06.672Z","temporaryDisk":{"sizeInGB":5,"mountPath":"/tmp"},"persistentDisk":{"sizeInGB":0,"mountPath":"/persistent"}},"type":"Microsoft.AppPlatform/Spring/apps","identity":null,"location":"westeurope","id":"/subscriptions/{subscription_id}/resourceGroups/acctest5486/providers/Microsoft.AppPlatform/Spring/acctest5486/apps/acctest7675","name":"acctest7675"}
   --------------------------------------------------------------------------------

```

### Links
1. [Semantic and Model Violations Reference](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/Semantic-and-Model-Violations-Reference.md)
2. [S360 action item generator for Swagger issues](https://aka.ms/swaggers360)
