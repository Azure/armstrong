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

variable "vm_password" {
  type        = string
  description = "The password for the DevTest Lab virtual machine"
  sensitive   = true
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "lab" {
  type      = "Microsoft.DevTestLab/labs@2018-09-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      labStorageType = "Premium"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
}

data "azapi_resource_id" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2023-04-01"
  parent_id = data.azapi_resource_id.virtualNetwork.id
  name      = "${var.resource_name}Subnet"
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.DevTestLab/labs/virtualNetworks@2018-09-15"
  parent_id = azapi_resource.lab.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      subnetOverrides = [
        {
          labSubnetName                = data.azapi_resource_id.subnet.name
          resourceId                   = data.azapi_resource_id.subnet.id
          useInVmCreationPermission    = "Allow"
          usePublicIpAddressPermission = "Allow"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.DevTestLab/labs/virtualMachines@2018-09-15"
  parent_id = azapi_resource.lab.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      allowClaim              = true
      disallowPublicIpAddress = false
      galleryImageReference = {
        offer     = "WindowsServer"
        osType    = "Windows"
        publisher = "MicrosoftWindowsServer"
        sku       = "2012-Datacenter"
        version   = "latest"
      }
      isAuthenticationWithSshKey = false
      labSubnetName              = data.azapi_resource_id.subnet.name
      labVirtualNetworkId        = azapi_resource.virtualNetwork.id
      networkInterface = {
      }
      notes       = ""
      osType      = "Windows"
      password    = var.vm_password
      size        = "Standard_F2"
      storageType = "Standard"
      userName    = "acct5stU5er"
    }
  }
  ignore_casing             = true
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

