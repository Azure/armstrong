# Example

## Introduction
This is an exmaple of testing the API which manages Machine Learning Compute/Compute Instance resource.

## Commands and logs
Command
`azurerm-rest-api-testing-tool auto C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\machinelearningservices\resource-manager\Microsoft.MachineLearningServices\stable\2021-07-01\examples\Compute\createOrUpdate\ComputeInstanceMinimal.json`

Logs

1. Genreate testing files, [dependency.tf](https://github.com/ms-henglu/azurerm-rest-api-testing-tool/blob/master/examples/dependency.tf) and [testing.tf](https://github.com/ms-henglu/azurerm-rest-api-testing-tool/blob/master/examples/testing.tf).
```
2021/11/01 13:27:33 [INFO] command: auto, args: [C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\machinelearningservices\resource-manager\Microsoft.MachineLearningServices\stable\2021-07-01\examples\Compute\createOrUpdate\ComputeInstanceMinimal.json]
2021/11/01 13:27:33 [INFO] loading dependencies
2021/11/01 13:27:33 [INFO] generating testing files
2021/11/01 13:27:33 [INFO] found dependency: azurerm_machine_learning_workspace
2021/11/01 13:27:33 [INFO] dependency.tf generated
2021/11/01 13:27:33 [INFO] testing.tf generated
```

2. Deploy depdencies and test ARM API.
```
2021/11/01 13:27:33 [INFO] prepare working directory
2021/11/01 13:27:33 [INFO] skip running init command because .terraform folder exist
[INFO] running Terraform command: C:\Apps\terraform.exe plan -no-color -input=false -detailed-exitcode -lock-timeout=0s -out=tfplan -lock=true -parallelism=10 -refresh=true

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # azurerm_application_insights.test will be created
  + resource "azurerm_application_insights" "test" {
        ignore details.....
    }

  # azurerm_key_vault.test will be created
  + resource "azurerm_key_vault" "test" {
        ignore details.....
    }

  # azurerm_machine_learning_workspace.test will be created
  + resource "azurerm_machine_learning_workspace" "test" {
        ignore details.....
    }

  # azurerm_resource_group.test will be created
  + resource "azurerm_resource_group" "test" {
      + id       = (known after apply)
      + location = "westeurope"
      + name     = "acctest1022"
    }

  # azurerm_storage_account.test will be created
  + resource "azurerm_storage_account" "test" {
        ignore details.....
    }

  # azurermg_resource.test will be created
  + resource "azurermg_resource" "test" {
        ignore details....
    }

Plan: 6 to add, 0 to change, 0 to destroy.
[INFO] running Terraform command: C:\Apps\terraform.exe version -json
{
  "terraform_version": "0.14.7",
  "terraform_revision": "",
  "provider_selections": {},
  "terraform_outdated": true
}
[INFO] running Terraform command: C:\Apps\terraform.exe show -json -no-color tfplan
...
ignore plan json content
...
2021/11/01 13:27:55 [INFO] Running plan completed, found 6 changes
[INFO] running Terraform command: C:\Apps\terraform.exe apply -no-color -auto-approve -input=false -lock=true -parallelism=10 -refresh=true

Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - hashicorp/azurerm in C:\Users\henglu\go\bin
 - ms-henglu/azurermg in C:\Users\henglu\go\bin

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.

azurerm_resource_group.test: Creating...
azurerm_resource_group.test: Creation complete after 3s [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022]
azurerm_application_insights.test: Creating...
azurerm_key_vault.test: Creating...
azurerm_storage_account.test: Creating...
azurerm_key_vault.test: Still creating... [10s elapsed]
azurerm_application_insights.test: Still creating... [10s elapsed]
azurerm_storage_account.test: Still creating... [10s elapsed]
azurerm_storage_account.test: Still creating... [20s elapsed]
azurerm_key_vault.test: Still creating... [20s elapsed]
azurerm_application_insights.test: Still creating... [20s elapsed]
azurerm_application_insights.test: Creation complete after 27s [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/microsoft.insights/components/acctest1022]
azurerm_storage_account.test: Still creating... [30s elapsed]
azurerm_key_vault.test: Still creating... [30s elapsed]
azurerm_storage_account.test: Creation complete after 33s [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.Storage/storageAccounts/acctest1022]
azurerm_key_vault.test: Still creating... [40s elapsed]
azurerm_key_vault.test: Still creating... [2m20s elapsed]
azurerm_key_vault.test: Creation complete after 2m29s [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.KeyVault/vaults/acctest1022]
azurerm_machine_learning_workspace.test: Creating...
azurerm_machine_learning_workspace.test: Still creating... [10s elapsed]
azurerm_machine_learning_workspace.test: Still creating... [20s elapsed]
azurerm_machine_learning_workspace.test: Still creating... [30s elapsed]
azurerm_machine_learning_workspace.test: Creation complete after 38s [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.MachineLearningServices/workspaces/acctest1022]
azurermg_resource.test: Creating...
azurermg_resource.test: Still creating... [10s elapsed]
azurermg_resource.test: Still creating... [4m41s elapsed]
azurermg_resource.test: Creation complete after 4m41s [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.MachineLearningServices/workspaces/acctest1022/computes/acctest1275?api-version=2021-07-01]

Apply complete! Resources: 6 added, 0 changed, 0 destroyed.
[INFO] running Terraform command: C:\Apps\terraform.exe plan -no-color -input=false -detailed-exitcode -lock-timeout=0s -out=tfplan -lock=true -parallelism=10 -refresh=true
azurerm_resource_group.test: Refreshing state... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022]
azurerm_key_vault.test: Refreshing state... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.KeyVault/vaults/acctest1022]
azurerm_application_insights.test: Refreshing state... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/microsoft.insights/components/acctest1022]
azurerm_storage_account.test: Refreshing state... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.Storage/storageAccounts/acctest1022]
azurerm_machine_learning_workspace.test: Refreshing state... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.MachineLearningServices/workspaces/acctest1022]
azurermg_resource.test: Refreshing state... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.MachineLearningServices/workspaces/acctest1022/computes/acctest1275?api-version=2021-07-01]

No changes. Infrastructure is up-to-date.

This means that Terraform did not detect any differences between your
configuration and real physical resources that exist. As a result, no
actions need to be performed.
2021/11/01 13:37:01 [INFO] Test passed! 
```
3. Clean up dependencies and testing resources
```
2021/11/01 13:37:01 [INFO] prepare working directory
2021/11/01 13:37:01 [INFO] skip running init command because .terraform folder exist
[INFO] running Terraform command: C:\Apps\terraform.exe destroy -no-color -auto-approve -input=false -lock-timeout=0s -lock=true -parallelism=10 -refresh=true

Warning: Provider development overrides are in effect

The following provider development overrides are set in the CLI configuration:
 - hashicorp/azurerm in C:\Users\henglu\go\bin
 - ms-henglu/azurermg in C:\Users\henglu\go\bin

The behavior may therefore not match any released version of the provider and
applying changes may cause the state to become incompatible with published
releases.

azurermg_resource.test: Destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.MachineLearningServices/workspaces/acctest1022/computes/acctest1275?api-version=2021-07-01]
azurermg_resource.test: Still destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-...tes/acctest1275?api-version=2021-07-01, 10s elapsed]
azurermg_resource.test: Still destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-...tes/acctest1275?api-version=2021-07-01, 2m20s elapsed]
azurermg_resource.test: Destruction complete after 2m36s
azurerm_machine_learning_workspace.test: Destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.MachineLearningServices/workspaces/acctest1022]
azurerm_machine_learning_workspace.test: Still destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-...earningServices/workspaces/acctest1022, 10s elapsed]
azurerm_machine_learning_workspace.test: Still destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-...earningServices/workspaces/acctest1022, 50s elapsed]
azurerm_machine_learning_workspace.test: Destruction complete after 1m4s
azurerm_application_insights.test: Destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/microsoft.insights/components/acctest1022]
azurerm_key_vault.test: Destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.KeyVault/vaults/acctest1022]
azurerm_storage_account.test: Destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022/providers/Microsoft.Storage/storageAccounts/acctest1022]
azurerm_application_insights.test: Destruction complete after 9s
azurerm_storage_account.test: Destruction complete after 9s
azurerm_key_vault.test: Still destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-.../Microsoft.KeyVault/vaults/acctest1022, 10s elapsed]
azurerm_resource_group.test: Destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022]
azurerm_resource_group.test: Still destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022, 10s elapsed]
azurerm_resource_group.test: Still destroying... [id=/subscriptions/67a9759d-d099-4aa8-8675-e6cfd669c3f4/resourceGroups/acctest1022, 2m20s elapsed]
azurerm_resource_group.test: Destruction complete after 2m23s
```
