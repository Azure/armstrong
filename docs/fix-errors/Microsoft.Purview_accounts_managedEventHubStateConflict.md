# Fix Error Example - Microsoft.Purview/accounts - 409 Conflict

## Error Message

```shell
ERRO[0015] error running terraform apply: exit status 1

Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906
--------------------------------------------------------------------------------
RESPONSE 409: 409 Conflict
ERROR CODE: 2008
--------------------------------------------------------------------------------
{
  "error": {
    "code": "2008",
    "message": "The managed event hub namespace cannot be re-enabled for account acctest5906 once it has been disabled.",
    "target": null,
    "details": []
  }
}
--------------------------------------------------------------------------------


  with azapi_resource.account,
  on main.tf line 31, in resource "azapi_resource" "account":
  31: resource "azapi_resource" "account" {
```

## Config that triggers the error

```hcl
resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type = "SystemAssigned"
    identity_ids = []
  }
  body = jsonencode({
    properties = {
      managedResourceGroupName            = "custom-rgname"
      managedResourcesPublicNetworkAccess = "Enabled"
      managedEventHubState = "Enabled" // newly added
    }
  })
  schema_validation_enabled = false
}
```

## Solution

Some APIs don't support to be updated. In this case, we must remove this resource from the portal and run the test command again.


## Terraform config

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest5906"
}

variable "location" {
  type    = string
  default = "West US 2"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

// OperationId: Accounts_CreateOrUpdate, Accounts_Get, Accounts_Delete
// PUT GET DELETE /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}
resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type = "SystemAssigned"
    identity_ids = []
  }
  body = jsonencode({
    properties = {
      managedResourceGroupName            = "custom-rgname"
      managedResourcesPublicNetworkAccess = "Enabled"
      managedEventHubState = "Enabled"
    }
  })
  schema_validation_enabled = false
}

// OperationId: Accounts_Update
// PATCH /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}
resource "azapi_resource_action" "patch_account" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = ""
  method      = "PATCH"
  body = jsonencode({
    properties = {
      cloudConnectors = {
      }
      managedResourceGroupName            = "aaaaaaaaaaaaaaaaaaaaaaaaaaaa"
      managedResourcesPublicNetworkAccess = "Disabled"
      publicNetworkAccess                 = "NotSpecified"
    }
  })
}

resource "azapi_resource" "user" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

// OperationId: Accounts_AddRootCollectionAdmin
// POST /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/addRootCollectionAdmin
resource "azapi_resource_action" "addRootCollectionAdmin" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "addRootCollectionAdmin"
  method      = "POST"
  body = jsonencode({
    objectId = jsondecode(azapi_resource.user.output).properties.principalId
  })
}

// OperationId: Features_AccountGet
// POST /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/listFeatures
resource "azapi_resource_action" "listFeatures" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "listFeatures"
  method      = "POST"
  body = jsonencode({
    features = [
      "Feature1",
      "Feature2",
      "FeatureThatDoesntExist",
    ]
  })
}

// OperationId: Accounts_ListKeys
// POST /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/listkeys
resource "azapi_resource_action" "listkeys" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "listkeys"
  method      = "POST"
}

data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2020-06-01"
  response_export_values = ["*"]
}

// OperationId: Accounts_ListBySubscription
// GET /subscriptions/{subscriptionId}/providers/Microsoft.Purview/accounts
data "azapi_resource_list" "listAccountsBySubscription" {
  type       = "Microsoft.Purview/accounts@2021-12-01"
  parent_id  = data.azapi_resource.subscription.id
  depends_on = [azapi_resource.account]
}

// OperationId: Accounts_ListByResourceGroup
// GET /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts
data "azapi_resource_list" "listAccountsByResourceGroup" {
  type       = "Microsoft.Purview/accounts@2021-12-01"
  parent_id  = azapi_resource.resourceGroup.id
  depends_on = [azapi_resource.account]
}

```