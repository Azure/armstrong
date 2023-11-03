# Fix Error Example - Microsoft.Purview/accounts/kafkaConfiguration - 409 Require Account Disable Managed Event Hub

## Error Message

```shell
Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest2258/providers/Microsoft.Purview/accounts/acctest2258/kafkaConfigurations/acctest2258\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest2258/providers/Microsoft.Purview/accounts/acctest2258/kafkaConfigurations/acctest2258
--------------------------------------------------------------------------------
RESPONSE 409: 409 Conflict
ERROR CODE: 32002
--------------------------------------------------------------------------------
{
  "error": {
    "code": "32002",
    "message": "Not able to create or update kafka configurations when managed eventhub enabled. Please disable before managing kafka configurations.",
    "target": null,
    "details": []
  }
}
--------------------------------------------------------------------------------


  with azapi_resource.kafkaConfiguration,
  on main.tf line 84, in resource "azapi_resource" "kafkaConfiguration":
  84: resource "azapi_resource" "kafkaConfiguration" {
```

## Config that triggers the error

```hcl
resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    identity = {
      type                   = "SystemAssigned"
      userAssignedIdentities = null
    }
    properties = {
      publicNetworkAccess = "Enabled"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
```

## Solution

In purview's doc, it says that the managed event hub will be disabled when the api-version is 2021-12-01. 

Changing the type from `Microsoft.Purview/accounts@2021-07-01` to `Microsoft.Purview/accounts@2021-12-01` will fix the error.

Unfortunatly, `managedEventHubState` could not be updated after the account is created. So we need to delete the account from the Azure Portal and run the test command again.

Refs: [Microsoft.Purview/accounts - 400 ManagedEventHubStateConflict](Microsoft.Purview_accounts_managedEventHubStateConflict.md).


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
  default = "acctest2258"
}

variable "location" {
  type    = string
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01" // changed from 2021-07-01
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    identity = {
      type                   = "SystemAssigned"
      userAssignedIdentities = null
    }
    properties = {
      publicNetworkAccess = "Enabled"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "userAssignedIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    identity = {
      type                   = "None"
      userAssignedIdentities = null
    }
    properties = {
      disableLocalAuth     = false
      isAutoInflateEnabled = false
      publicNetworkAccess  = "Enabled"
      zoneRedundant        = false
    }
    sku = {
      capacity = 1
      name     = "Standard"
      tier     = "Standard"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

// OperationId: KafkaConfigurations_CreateOrUpdate, KafkaConfigurations_Get, KafkaConfigurations_Delete
// PUT GET DELETE /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/kafkaConfigurations/{kafkaConfigurationName}
resource "azapi_resource" "kafkaConfiguration" {
  type      = "Microsoft.Purview/accounts/kafkaConfigurations@2021-12-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      consumerGroup = "consumerGroup"
      credentials = {
        identityId = azapi_resource.userAssignedIdentity.id
        type       = "UserAssigned"
      }
      eventHubPartitionId = "partitionId"
      eventHubResourceId  = azapi_resource.namespace.id
      eventHubType        = "Notification"
      eventStreamingState = "Enabled"
      eventStreamingType  = "Azure"
    }
  })
  schema_validation_enabled = false
}

// OperationId: KafkaConfigurations_ListByAccount
// GET /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/kafkaConfigurations
data "azapi_resource_list" "listKafkaConfigurationsByAccount" {
  type       = "Microsoft.Purview/accounts/kafkaConfigurations@2021-12-01"
  parent_id  = azapi_resource.account.id
  depends_on = [azapi_resource.kafkaConfiguration]
}

```