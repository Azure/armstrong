# Example

## Introduction
This is an example of testing the API which manages Microsoft.Databricks/workspaces/virtualNetworkPeerings@2021-04-01-preview resource

## Commands and logs
Command
`armstrong auto -path C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\databricks\resource-manager\Microsoft.Databricks\preview\2021-04-01-preview\examples\WorkspaceVirtualNetworkPeeringCreateOrUpdate.json`

Logs

1. Genreate testing files, [dependency.tf](https://github.com/azure/armstrong/blob/master/examples/demos/case3%20-%20complicated%20dependency/dependency.tf) and [testing.tf](https://github.com/azure/armstrong/blob/master/examples/demos/case3%20-%20complicated%20dependency/testing.tf).
```
2021/11/24 15:47:11 [INFO] ----------- generate dependency and test resource ---------
2021/11/24 15:47:11 [INFO] loading dependencies
2021/11/24 15:47:11 [INFO] generating testing files
2021/11/24 15:47:11 [INFO] found dependency: azurerm_virtual_network     
2021/11/24 15:47:11 [INFO] found dependency: azurerm_databricks_workspace
2021/11/24 15:47:11 [INFO] dependency.tf generated
2021/11/24 15:47:11 [INFO] testing.tf generated
```

2. Deploy depdencies and test ARM API.
```
2021/11/24 15:47:35 [INFO] ----------- run tests ---------
2021/11/24 15:47:35 [INFO] prepare working directory
2021/11/24 15:47:35 [INFO] skip running init command because .terraform folder exist
2021/11/24 15:47:35 [INFO] running plan command to check changes...
2021/11/24 15:47:50 [INFO] found 6 changes in total, create: 6, replace: 0, update: 0, delete: 0
2021/11/24 15:47:50 [INFO] running apply command to provision test resource...
2021/11/24 15:51:02 [INFO] test resource has been provisioned
2021/11/24 15:51:02 [INFO] running plan command to verify test resource...
2021/11/24 15:51:19 [INFO] Test passed!
```
3. Clean up dependencies and testing resources
```
2021/11/24 16:03:52 [INFO] ----------- cleanup resources ---------
2021/11/24 16:03:52 [INFO] prepare working directory
2021/11/24 16:03:52 [INFO] skip running init command because .terraform folder exist
2021/11/24 16:03:52 [INFO] running destroy command to cleanup resources...
2021/11/24 16:07:38 [INFO] all resources have been deleted
2021/11/24 16:07:38 [INFO] Test passed!
```
