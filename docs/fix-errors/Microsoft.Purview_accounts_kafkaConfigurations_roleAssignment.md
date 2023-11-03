# Fix Error Example - Microsoft.Purview/accounts/kafkaConfiguration - 409 No Permission
## Error Message

```shell
Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248
│ --------------------------------------------------------------------------------
│ RESPONSE 409: 409 Conflict
│ ERROR CODE: 19000
│ --------------------------------------------------------------------------------
│ {
│   "error": {
│     "code": "19000",
│     "message": "The caller does not have EventHubDataSender permission on resource: [/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.EventHub/namespaces/acctest2248/eventhubs/acctest2248].",
│     "target": null,
│     "details": []
│   }
│ }
│ --------------------------------------------------------------------------------
│ 
│ 
│   with azapi_resource.kafkaConfiguration,
│   on main.tf line 102, in resource "azapi_resource" "kafkaConfiguration":
│  102: resource "azapi_resource" "kafkaConfiguration" {
│ 
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
        type       = "SystemAssigned"
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

It requires the identity used to create the kafka configuration to have the access permission to the eventhub. Adding role assignment requires a collection of API operations, but `azurerm_role_assignment` resource simplify the process, we'll use this resource to add the role assignment.

Here're the steps to fix it.

1. Add `azurerm` provider definition to the config file to enable the provider.

```hcl
provider "azurerm" {
  features {}
}
```

2. Add the following resources to the config file to add the role assignments.

```hcl
resource "azurerm_role_assignment" "roleAssignment" {
  scope                = azapi_resource.namespace.id
  role_definition_name = "Owner"
  principal_id         = azapi_resource.account.identity[0].principal_id
}


resource "azurerm_role_assignment" "roleAssignment2" {
  scope                = azapi_resource.eventhub.id
  role_definition_name = "Owner"
  principal_id         = azapi_resource.account.identity[0].principal_id
}

resource "azurerm_role_assignment" "roleAssignment3" {
  scope                = azapi_resource.eventhub.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azapi_resource.account.identity[0].principal_id
}
```

3. Add the `depends_on` to add implicit dependency between the resources.

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
      eventHubResourceId  = azapi_resource.eventhub.id
      eventHubType        = "Notification"
      eventStreamingState = "Enabled"
      eventStreamingType  = "Azure"
    }
  })
  schema_validation_enabled = false

  // add the following configs:
  depends_on = [    
    azurerm_role_assignment.roleAssignment,
    azurerm_role_assignment.roleAssignment2,
    azurerm_role_assignment.roleAssignment3
  ]
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

// add the following provider definition
provider "azurerm" {
  features {
    
  }
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
  identity {
    type = "SystemAssigned"
  }
  body = jsonencode({
    properties = {
      publicNetworkAccess = "Enabled"
    }
  })
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

# =================== add the following resources ===================
resource "azurerm_role_assignment" "roleAssignment" {
  scope                = azapi_resource.namespace.id
  role_definition_name = "Owner"
  principal_id         = azapi_resource.account.identity[0].principal_id
}


resource "azurerm_role_assignment" "roleAssignment2" {
  scope                = azapi_resource.eventhub.id
  role_definition_name = "Owner"
  principal_id         = azapi_resource.account.identity[0].principal_id
}

resource "azurerm_role_assignment" "roleAssignment3" {
  scope                = azapi_resource.eventhub.id
  role_definition_name = "Azure Event Hubs Data Sender"
  principal_id         = azapi_resource.account.identity[0].principal_id
}
# =================== fix ===================

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
        type       = "SystemAssigned"
      }
      eventHubPartitionId = "partitionId"
      eventHubResourceId  = azapi_resource.eventhub.id
      eventHubType        = "Notification"
      eventStreamingState = "Enabled"
      eventStreamingType  = "Azure"
    }
  })
  schema_validation_enabled = false
  // add the following configs:
  depends_on = [    
    azurerm_role_assignment.roleAssignment,
    azurerm_role_assignment.roleAssignment2,
    azurerm_role_assignment.roleAssignment3
  ]
}

// OperationId: KafkaConfigurations_ListByAccount
// GET /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/kafkaConfigurations
data "azapi_resource_list" "listKafkaConfigurationsByAccount" {
  type       = "Microsoft.Purview/accounts/kafkaConfigurations@2021-12-01"
  parent_id  = azapi_resource.account.id
  depends_on = [azapi_resource.kafkaConfiguration]
}
```