# Fix Error Example - Microsoft.Purview/accounts/kafkaConfiguration - 400 EventHubResourceId Invalid

## Error Message

```shell
ERRO[0130] error running terraform apply: exit status 1

Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: 1002
--------------------------------------------------------------------------------
{
  "error": {
    "code": "1002",
    "message": "eventHubResourceId:/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.EventHub/namespaces/acctest2248 is not a valid ARM resource ID",
    "target": "eventHubResourceId",
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
```

## Solution

In the example used to create the terraform config, it sets `eventHubResourceId` to the namespace's id placeholder. But the `eventHubResourceId` should be the event hub's id.

1. Add a new resource to create the event hub.

```hcl
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
```

2. Update the `eventHubResourceId` to use the event hub's id.

```hcl
resource "azapi_resource" "kafkaConfiguration" {
  type      = "Microsoft.Purview/accounts/kafkaConfigurations@2021-12-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      consumerGroup = "consumerGroup"
      credentials = {
        type = "SystemAssigned"
      }
      eventHubPartitionId = "partitionId"
      eventHubResourceId  = azapi_resource.eventhub.id // changed from azapi_resource.namespace.id
      eventHubType        = "Notification"
      eventStreamingState = "Enabled"
      eventStreamingType  = "Azure"
    }
  })
  schema_validation_enabled = false
}
```

3. It's recommended to fix the placeholder in the JSON examples.

```
before:
  "eventHubResourceId": "/subscriptions/225be6fe-ec1c-4d51-a368-f69348d2e6c5/resourceGroups/testRG/providers/Microsoft.EventHub/namespaces/eventHubNameSpaceName"

after:
  "eventHubResourceId": "/subscriptions/225be6fe-ec1c-4d51-a368-f69348d2e6c5/resourceGroups/testRG/providers/Microsoft.EventHub/namespaces/eventHubNameSpaceName/eventhubs/eventhubName"
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

# ========= add the following resources =========
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
# ========= fix =========

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
      eventHubResourceId  = azapi_resource.eventhub.id // changed from azapi_resource.namespace.id
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