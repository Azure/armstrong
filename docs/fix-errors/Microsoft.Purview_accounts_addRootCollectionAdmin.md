# Fix Error Example - Microsoft.Purview/accounts - 400 addRootCollectionAdmin

## Error Message

```shell
ERRO[0041] error running terraform apply: exit status 1

Error: performing action addRootCollectionAdmin of "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906\" / Api Version \"2021-12-01\")": POST https://management.azure.com/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906/addRootCollectionAdmin
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: 1002
--------------------------------------------------------------------------------
{
  "error": {
    "code": "1002",
    "message": "The payload is invalid. Error: Failed to find ODataType for objectId:7e8de0e7-2bfc-4e1f-9659-2a5785e4356f tenantId:72f988bf-86f1-41af-91ab-2d7cd011db47..",
    "target": null,
    "details": []
  }
}
--------------------------------------------------------------------------------


  with azapi_resource_action.addRootCollectionAdmin,
  on main.tf line 69, in resource "azapi_resource_action" "addRootCollectionAdmin":
  69: resource "azapi_resource_action" "addRootCollectionAdmin" {
```

## Config that triggers the error

```hcl
resource "azapi_resource_action" "addRootCollectionAdmin" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "addRootCollectionAdmin"
  method      = "POST"
  body = jsonencode({
    objectId = "7e8de0e7-2bfc-4e1f-9659-2a5785e4356f"
  })
}
```

## Solution

Armstrong will automatically generate the necessary dependencies for the test case, but there're cases that Armstrong can't generate the dependencies automatically. For example, this API operation depends on the identity's principalId, but Armstrong can't identify the identity's principalId automatically. In this case, we need to add the dependency manually.

1. Adding the following config to the `main.tf` to add a user identity resource.

```hcl
resource "azapi_resource" "user" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
```

2. Update the config to reference the principalId of the user identity resource.

```hcl
resource "azapi_resource_action" "addRootCollectionAdmin" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "addRootCollectionAdmin"
  method      = "POST"
  body = jsonencode({
    objectId = jsondecode(azapi_resource.user.output).properties.principalId
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

  # ========== fix: add the following block ==========
  identity {
    type = "SystemAssigned"
    identity_ids = []
  }
  # ========== fix ==========

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
    cloudConnectors = {
    }
    managedResourceGroupName            = "aaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    managedResourcesPublicNetworkAccess = "Disabled"
    publicNetworkAccess                 = "NotSpecified"
  })
}

# ========== fix: add the following block ==========
resource "azapi_resource" "user" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
# ========== fix ==========

// OperationId: Accounts_AddRootCollectionAdmin
// POST /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/addRootCollectionAdmin
resource "azapi_resource_action" "addRootCollectionAdmin" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "addRootCollectionAdmin"
  method      = "POST"
  body = jsonencode({
    # ========== fix: add the following line ==========
    objectId = jsondecode(azapi_resource.user.output).properties.principalId
    # ========== fix ==========
    # ========== fix: remove the following line ==========
    objectId = "7e8de0e7-2bfc-4e1f-9659-2a5785e4356f"
    # ========== fix ==========
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