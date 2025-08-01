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

variable "sql_admin_username" {
  type        = string
  description = "The administrator username for the SQL server credential"
}

variable "sql_admin_password" {
  type        = string
  description = "The administrator password for the SQL server credential"
  sensitive   = true
}

variable "administrator_login_password" {
  type        = string
  description = "The administrator login password for the SQL server"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2021-02-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin            = "4dministr4t0r"
      administratorLoginPassword    = var.administrator_login_password
      minimalTlsVersion             = "1.2"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      version                       = "12.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Sql/servers/databases@2021-02-01-preview"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoPauseDelay                   = 0
      collation                        = "SQL_Latin1_General_CP1_CI_AS"
      createMode                       = "Default"
      elasticPoolId                    = ""
      highAvailabilityReplicaCount     = 0
      isLedgerOn                       = false
      minCapacity                      = 0
      readScale                        = "Disabled"
      requestedBackupStorageRedundancy = "Geo"
      zoneRedundant                    = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "jobAgent" {
  type      = "Microsoft.Sql/servers/jobAgents@2020-11-01-preview"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      databaseId = azapi_resource.database.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "credential" {
  type      = "Microsoft.Sql/servers/jobAgents/credentials@2020-11-01-preview"
  parent_id = azapi_resource.jobAgent.id
  name      = var.resource_name
  body = {
    properties = {
      password = var.sql_admin_password
      username = var.sql_admin_username
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

