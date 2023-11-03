# Fix Error Example - Microsoft.Purview/accounts/kafkaConfiguration - 400 EventHubResourceId Invalid

## Error Message

```shell
ERRO[0037] error running terraform apply: exit status 1

Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248
--------------------------------------------------------------------------------
RESPONSE 500: 500 Internal Server Error
ERROR CODE: 500
--------------------------------------------------------------------------------
{
  "error": {
    "code": "500",
    "message": "Unknown error",
    "target": null,
    "details": null
  }
}
--------------------------------------------------------------------------------


  with azapi_resource.kafkaConfiguration,
  on main.tf line 99, in resource "azapi_resource" "kafkaConfiguration":
  99: resource "azapi_resource" "kafkaConfiguration" {
```

## Config that triggers the error

```hcl
resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
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
      eventHubResourceId  = azapi_resource.eventhub.id
      eventHubType        = "Notification"
      eventStreamingState = "Enabled"
      eventStreamingType  = "Azure"
    }
  })
  schema_validation_enabled = false
}
```

## Solution

The error from the server side is not clear, after some investigation, the root cause is that the user assigned idenity used to created the kafka configuration is not associated with the purview account.

There're two ways to fix it.

1. Use the system assigned identity to create the kafka configuration.

```hcl
resource "azapi_resource" "kafkaConfiguration" {
  type      = "Microsoft.Purview/accounts/kafkaConfigurations@2021-12-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      consumerGroup = "consumerGroup"
      credentials = {
        // identityId = azapi_resource.userAssignedIdentity.id, comment out
        type       = "SystemAssigned" // changed from UserAssigned
      }
      eventHubPartitionId = "partitionId"
      eventHubResourceId  = azapi_resource.eventhub.id
      eventHubType        = "Notification"
      eventStreamingState = "Enabled"
      eventStreamingType  = "Azure"
    }
  })
  schema_validation_enabled = false
}
```

2. Associated the user assigned identity with the purview account.

```hcl
resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    identity = {
      type                   = "SystemAssigned, UserAssigned"
      userAssignedIdentities = {
        (azapi_resource.userAssignedIdentity.id) = {
        }
      }
    }
    properties = {
      publicNetworkAccess = "Enabled"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
```

Another way to manage the identity:
```hcl
resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = jsonencode({
    properties = {
      publicNetworkAccess = "Enabled"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
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
  default = "acctest2248"
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
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  # =================== fix: add the following block ===================  
  identity {
    type = "SystemAssigned, UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  # =================== fix ===================
  body = jsonencode({
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

resource "azapi_resource" "eventhub" {
  type      = "Microsoft.EventHub/namespaces/eventhubs@2023-01-01-preview"
  parent_id = azapi_resource.namespace.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      messageRetentionInDays = 1
      partitionCount         = 2
      status                 = "Active"
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
      eventHubResourceId  = azapi_resource.eventhub.id
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