terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "acctest0706"
  location = "West Europe"
}

resource "azurerm_monitor_action_group" "test" {
  name                = "test-action-group"
  resource_group_name = azurerm_resource_group.test.name
  short_name          = "action"
}


resource "azurerm_application_insights" "test" {
  name                = "appinsights"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  application_type    = "web"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "wangta-law-01"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}