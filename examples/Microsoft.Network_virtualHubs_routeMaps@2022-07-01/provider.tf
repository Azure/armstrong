
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
  
  // if you want to authenticate
  // using a service principal
  // add the following 4 fields

  /*
  subscription_id = "******"
  client_id       = "******"
  client_secret   = "******"
  tenant_id       = "******"
  */
}

provider "azapi" {
  skip_provider_registration = false

  // if you want to authenticate
  // using a service principal
  // add the following 4 fields

  /*
  subscription_id = "******"
  client_id       = "******"
  client_secret   = "******"
  tenant_id       = "******"
  */
}
