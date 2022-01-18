# Example

## Introduction
This is an example of testing the API which manages Microsoft.Batch/batchAccounts/applications@2021-06-01 resource

## Commands and logs
Command
`azurerm-restapi-testing-tool auto -path C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\batch\resource-manager\Microsoft.Batch\stable\2021-06-01\examples\ApplicationCreate.json`

Logs

1. Genreate testing files, [dependency.tf](https://github.com/ms-henglu/azurerm-restapi-testing-tool/blob/master/examples/case1%20-%20test%20passed/dependency.tf) and [testing.tf](https://github.com/ms-henglu/azurerm-restapi-testing-tool/blob/master/examples/case1%20-%20test%20passed/testing.tf).
```
2021/11/24 13:39:42 [INFO] command: auto, args: [C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\batch\resource-manager\Microsoft.Batch\stable\2021-06-01\examples\ApplicationCreate.json]
2021/11/24 13:39:42 [INFO] ----------- generate dependency and test resource ---------
2021/11/24 13:39:42 [INFO] loading dependencies
2021/11/24 13:39:42 [INFO] generating testing files
2021/11/24 13:39:42 [INFO] found dependency: azurerm_batch_account
2021/11/24 13:39:42 [INFO] dependency.tf generated
2021/11/24 13:39:42 [INFO] testing.tf generated
```

2. Deploy depdencies and test ARM API.
```
2021/11/24 13:39:42 [INFO] ----------- run tests ---------
2021/11/24 13:39:42 [INFO] prepare working directory
2021/11/24 13:39:42 [INFO] skip running init command because .terraform folder exist
2021/11/24 13:39:42 [INFO] running plan command to check changes...
2021/11/24 13:39:58 [INFO] found 4 changes in total, create: 4, replace: 0, update: 0, delete: 0
2021/11/24 13:39:58 [INFO] running apply command to provision test resource...
2021/11/24 13:41:27 [INFO] test resource has been provisioned
2021/11/24 13:41:27 [INFO] running plan command to verify test resource...
2021/11/24 13:41:50 [INFO] Test passed!
```
3. Clean up dependencies and testing resources
```
2021/11/24 13:41:50 [INFO] ----------- cleanup resources ---------
2021/11/24 13:41:50 [INFO] prepare working directory
2021/11/24 13:41:50 [INFO] skip running init command because .terraform folder exist
2021/11/24 13:41:50 [INFO] running destroy command to cleanup resources...
2021/11/24 13:43:51 [INFO] all resources have been deleted
2021/11/24 13:43:51 [INFO] Test passed!
```
