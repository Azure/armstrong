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
  description = "The administrator login for the MariaDB server"
}

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the MariaDB server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.DBforMariaDB/servers@2018-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = var.administrator_login
      administratorLoginPassword = var.administrator_login_password
      createMode                 = "Default"
      minimalTlsVersion          = "TLS1_2"
      publicNetworkAccess        = "Enabled"
      sslEnforcement             = "Enabled"
      storageProfile = {
        backupRetentionDays = 7
        storageAutogrow     = "Enabled"
        storageMB           = 51200
      }
      version = "10.2"
    }
    sku = {
      capacity = 2
      family   = "Gen5"
      name     = "B_Gen5_2"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.DBforMariaDB/servers/databases@2018-06-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  body = {
    properties = {
      charset   = "utf8"
      collation = "utf8_general_ci"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

