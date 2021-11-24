terraform {
  required_providers {
    azurerm-restapi = {
      source = "Azure/azurerm-restapi"
    }
  }
}
provider "azurerm" {
  features {}
}
provider "azurerm-restapi" {
}
resource "azurerm_resource_group" "test" {
  name     = "acctest6746"
  location = "West Europe"
}
resource "azurerm_storage_account" "test" {
  name                     = "acctest6746"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}
resource "azurerm_batch_account" "test" {
  name                 = "acctest6746"
  resource_group_name  = azurerm_resource_group.test.name
  location             = azurerm_resource_group.test.location
  pool_allocation_mode = "BatchService"
  storage_account_id   = azurerm_storage_account.test.id

  tags = {
    env = "test"
  }
}
