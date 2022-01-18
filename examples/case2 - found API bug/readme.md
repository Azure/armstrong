# Example

## Introduction
This is an example of testing the API which manages Microsoft.AppPlatform/Spring/apps@2020-07-01 resource

## Commands and logs
Command
`azurerm-restapi-testing-tool.exe auto -path C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\appplatform\resource-manager\Microsoft.AppPlatform\stable\2020-07-01\examples\Apps_CreateOrUpdate.json`

Logs

1. Genreate testing files, [dependency.tf](https://github.com/ms-henglu/azurerm-restapi-testing-tool/blob/master/examples/case2%20-%20found%20API%20bug/dependency.tf) and [testing.tf](https://github.com/ms-henglu/azurerm-restapi-testing-tool/blob/master/examples/case2%20-%20found%20API%20bug/testing.tf).
```
2021/11/24 13:42:17 [INFO] command: auto, args: [C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\appplatform\resource-manager\Microsoft.AppPlatform\stable\2020-07-01\examples\Apps_CreateOrUpdate.json]       
2021/11/24 13:42:17 [INFO] ----------- generate dependency and test resource ---------
2021/11/24 13:42:17 [INFO] loading dependencies
2021/11/24 13:42:17 [INFO] generating testing files
2021/11/24 13:42:17 [INFO] found dependency: azurerm_spring_cloud_service
2021/11/24 13:42:17 [INFO] dependency.tf generated
2021/11/24 13:42:17 [INFO] testing.tf generated
```

2. Deploy depdencies and test ARM API. If found bugs, detailed report will be printed in console. The report supports highlight in console.
```
2021/11/24 13:42:17 [INFO] ----------- run tests ---------
2021/11/24 13:42:17 [INFO] prepare working directory
2021/11/24 13:42:17 [INFO] skip running init command because .terraform folder exist
2021/11/24 13:42:17 [INFO] running plan command to check changes...
2021/11/24 13:42:32 [INFO] found 4 changes in total, create: 4, replace: 0, update: 0, delete: 0
2021/11/24 13:42:32 [INFO] running apply command to provision test resource...
2021/11/24 13:48:43 [INFO] test resource has been provisioned
2021/11/24 13:48:43 [INFO] running plan command to verify test resource...
2021/11/24 13:49:07 [INFO] found differences between response and configuration:
{
    "location": "westeurope",
    "properties": {
        "activeDeploymentName": "mydeployment1",
        "fqdn": "acctest1731.azuremicroservices.io" => "myapp.mydomain.com",
        "httpsOnly": false,
        "persistentDisk": {
            "mountPath": "/persistent" => "/mypersistentdisk",
            "sizeInGB": 0 => 2
        },
        "public": false => true,
        "temporaryDisk": {
            "mountPath": "/tmp" => "/mytemporarydisk",
            "sizeInGB": 5 => 2
        }
    }
}
2021/11/24 13:49:07 [INFO] report:
{
    "location": "westeurope",
    "properties": {
        "activeDeploymentName": "mydeployment1" is not returned from response,
        "fqdn":  Got "acctest1731.azuremicroservices.io" in response, expect "myapp.mydomain.com",
        "httpsOnly": false,
        "persistentDisk": {
            "mountPath":  Got "/persistent" in response, expect "/mypersistentdisk",
            "sizeInGB":  Got 0 in response, expect 2
        },
        "public":  Got false in response, expect true,
        "temporaryDisk": {
            "mountPath":  Got "/tmp" in response, expect "/mytemporarydisk",
            "sizeInGB":  Got 5 in response, expect 2
        }
    }
}
```
3. After modification dependency and test resource, remember to run `cleanup` command to remove all test resources.
```
2021/11/24 15:41:50 [INFO] ----------- cleanup resources ---------
2021/11/24 15:41:50 [INFO] prepare working directory
2021/11/24 15:41:50 [INFO] skip running init command because .terraform folder exist
2021/11/24 15:41:50 [INFO] running destroy command to cleanup resources...
2021/11/24 15:43:51 [INFO] all resources have been deleted
```
