resource "azurerm_resource_group" "test" {
  name     = "acctest-rg"
  location = "East US"
}

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa37"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-37"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/spring/storages@2024-01-01-preview"
  name      = "acctest-ss-37"
  parent_id = azurerm_spring_cloud_service.test.id

  body = jsonencode({
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "StorageAccount"
    }
  })
}

resource "azapi_resource" "test_dynamic" {
  type      = "Microsoft.AppPlatform/spring/storages@2024-01-01-preview"
  name      = "acctest-ss-39"
  parent_id = azurerm_spring_cloud_service.test.id

  body = {
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "StorageAccount"
    }
  }
}
