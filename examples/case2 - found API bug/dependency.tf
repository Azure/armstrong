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
  name     = "acctest1731"
  location = "West Europe"
}
resource "azurerm_application_insights" "test" {
  name                = "acctest1731"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}
resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest1731"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "S0"

  config_server_git_setting {
    uri          = "https://github.com/Azure-Samples/piggymetrics"
    label        = "config"
    search_paths = ["dir1", "dir2"]
  }

  trace {
    connection_string = azurerm_application_insights.test.connection_string
    sample_rate       = 10.0
  }

  tags = {
    Env = "staging"
  }
}
