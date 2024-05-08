
resource "azapi_resource" "configurationStore" {
  type      = "Microsoft.AppConfiguration/configurationStores@2019-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    identity = {
      type = "SystemAssigned, UserAssigned"
      userAssignedIdentities = {
        (azapi_resource.userAssignedIdentity.id) = {
        }
      }
    }
    sku = {
      name = "Standard"
    }
    tags = {
      myTag = "myTagValue"
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}
