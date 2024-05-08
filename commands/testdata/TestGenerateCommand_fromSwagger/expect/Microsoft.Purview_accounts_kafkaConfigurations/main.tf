terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
  skip_provider_registration = true
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest5353"
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
  type      = "Microsoft.Purview/accounts@2021-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "namespace" {
  type      = "Microsoft.EventHub/namespaces@2022-01-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
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
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

// OperationId: KafkaConfigurations_CreateOrUpdate, KafkaConfigurations_Get, KafkaConfigurations_Delete
// PUT GET DELETE /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/kafkaConfigurations/{kafkaConfigurationName}
resource "azapi_resource" "kafkaConfiguration" {
  type      = "Microsoft.Purview/accounts/kafkaConfigurations@2021-12-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  body = {
    properties = {
      consumerGroup       = "consumerGroup"
      eventHubPartitionId = "partitionId"
      eventHubResourceId  = azapi_resource.namespace.id
      eventHubType        = "Notification"
      eventStreamingState = "Enabled"
      eventStreamingType  = "Azure"
    }
  }
  schema_validation_enabled = false
}

// OperationId: KafkaConfigurations_ListByAccount
// GET /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/kafkaConfigurations
data "azapi_resource_list" "listKafkaConfigurationsByAccount" {
  type       = "Microsoft.Purview/accounts/kafkaConfigurations@2021-12-01"
  parent_id  = azapi_resource.account.id
  depends_on = [azapi_resource.kafkaConfiguration]
}

