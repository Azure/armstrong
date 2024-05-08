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
  default = "acctest4496"
}

variable "location" {
  type    = string
  default = "westeurope"
}

data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2020-06-01"
  response_export_values = ["*"]
}

data "azapi_resource_id" "location" {
  type      = "Microsoft.Purview/locations@2023-12-12"
  parent_id = data.azapi_resource.subscription.id
  name      = var.location
}

// OperationId: Features_SubscriptionGet
// POST /subscriptions/{subscriptionId}/providers/Microsoft.Purview/locations/{locations}/listFeatures
resource "azapi_resource_action" "listFeatures" {
  type        = "Microsoft.Purview/locations@2021-12-01"
  resource_id = data.azapi_resource_id.location.id
  action      = "listFeatures"
  method      = "POST"
  body = {
    features = [
      "Feature1",
      "Feature2",
      "FeatureThatDoesntExist",
    ]
  }
}

// OperationId: Usages_Get
// GET /subscriptions/{subscriptionId}/providers/Microsoft.Purview/locations/{location}/usages
data "azapi_resource_action" "usages" {
  type        = "Microsoft.Purview/locations@2021-12-01"
  resource_id = data.azapi_resource_id.location.id
  action      = "usages"
  method      = "GET"
}

