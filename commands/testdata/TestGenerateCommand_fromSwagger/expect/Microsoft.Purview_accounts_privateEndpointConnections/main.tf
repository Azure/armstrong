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
  default = "acctest2506"
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

// OperationId: PrivateEndpointConnections_CreateOrUpdate, PrivateEndpointConnections_Get, PrivateEndpointConnections_Delete
// PUT GET DELETE /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/privateEndpointConnections/{privateEndpointConnectionName}
resource "azapi_resource" "privateEndpointConnection" {
  type      = "Microsoft.Purview/accounts/privateEndpointConnections@2021-12-01"
  parent_id = azapi_resource.account.id
  name      = var.resource_name
  body = {
    properties = {
      privateLinkServiceConnectionState = {
        description = "Approved by johndoe@company.com"
        status      = "Approved"
      }
    }
  }
  schema_validation_enabled = false
}

// OperationId: PrivateEndpointConnections_ListByAccount
// GET /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/privateEndpointConnections
data "azapi_resource_list" "listPrivateEndpointConnectionsByAccount" {
  type       = "Microsoft.Purview/accounts/privateEndpointConnections@2021-12-01"
  parent_id  = azapi_resource.account.id
  depends_on = [azapi_resource.privateEndpointConnection]
}

