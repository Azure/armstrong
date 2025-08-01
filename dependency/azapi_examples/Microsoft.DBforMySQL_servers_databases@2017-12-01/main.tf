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
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

variable "administrator_login" {
  type        = string
  description = "The administrator login for the MySQL server"
}

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the MySQL server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.DBforMySQL/servers@2017-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = var.administrator_login
      administratorLoginPassword = var.administrator_login_password
      createMode                 = "Default"
      infrastructureEncryption   = "Disabled"
      minimalTlsVersion          = "TLS1_1"
      publicNetworkAccess        = "Enabled"
      sslEnforcement             = "Enabled"
      storageProfile = {
        storageAutogrow = "Enabled"
        storageMB       = 51200
      }
      version = "5.7"
    }
    sku = {
      capacity = 2
      family   = "Gen5"
      name     = "GP_Gen5_2"
      tier     = "GeneralPurpose"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.DBforMySQL/servers/databases@2017-12-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = {
    properties = {
      charset   = "utf8"
      collation = "utf8_unicode_ci"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

