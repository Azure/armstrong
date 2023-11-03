# Fix Error Example - Microsoft.Purview/accounts - 400 InvalidRequestContent

## Error Message

```shell
ERRO[0109] error running terraform apply: exit status 1

Error: performing action  of "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906\" / Api Version \"2021-12-01\")": PATCH https://management.azure.com/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: InvalidRequestContent
--------------------------------------------------------------------------------
{
  "error": {
    "code": "InvalidRequestContent",
    "message": "The request content was invalid and could not be deserialized: 'Could not find member 'cloudConnectors' on object of type 'ResourceDefinition'. Path 'cloudConnectors', line 1, position 20.'."
  }
}
--------------------------------------------------------------------------------


  with azapi_resource_action.patch_account,
  on main.tf line 51, in resource "azapi_resource_action" "patch_account":
  51: resource "azapi_resource_action" "patch_account" {

```

## Config that triggers the error

```hcl
resource "azapi_resource_action" "patch_account" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = ""
  method      = "PATCH"
  body = jsonencode({
    cloudConnectors = {
    }
    managedResourceGroupName            = "aaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    managedResourcesPublicNetworkAccess = "Disabled"
    publicNetworkAccess                 = "NotSpecified"
  })
}
```


## Solution

In above config, the `cloudConnectors`, `managedResourceGroupName`, `managedResourcesPublicNetworkAccess` and `publicNetworkAccess` should be under the `properties` bag, not the root of the body. The following config will fix the error.

```hcl
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
```


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
    # ========== fix: add the following configs ==========
    properties = {
      cloudConnectors = {
      }
      managedResourceGroupName            = "aaaaaaaaaaaaaaaaaaaaaaaaaaaa"
      managedResourcesPublicNetworkAccess = "Disabled"
      publicNetworkAccess                 = "NotSpecified"
    }
    # ========== fix ==========
    # ========== fix: remove the following configs  ==========
    cloudConnectors = {
    }
    managedResourceGroupName            = "aaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    managedResourcesPublicNetworkAccess = "Disabled"
    publicNetworkAccess                 = "NotSpecified"
    # ========== fix ==========
  })
}

// OperationId: Accounts_AddRootCollectionAdmin
// POST /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/addRootCollectionAdmin
resource "azapi_resource_action" "addRootCollectionAdmin" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "addRootCollectionAdmin"
  method      = "POST"
  body = jsonencode({
    objectId = "7e8de0e7-2bfc-4e1f-9659-2a5785e4356f"
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