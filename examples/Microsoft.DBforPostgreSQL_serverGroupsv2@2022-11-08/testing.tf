resource "azapi_resource" "serverGroupsv2" {
  type      = "Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08"
  name      = "acctest9810"
  parent_id = azurerm_resource_group.test.id

  body = jsonencode({
    location = "westus"
    properties = {
      administratorLoginPassword      = "MSDFL8923@#$"
      citusVersion                    = "11.1"
      coordinatorEnablePublicIpAccess = true
      coordinatorServerEdition        = "GeneralPurpose"
      coordinatorStorageQuotaInMb     = 524288
      coordinatorVCores               = 4
      enableHa                        = true
      enableShardsOnCoordinator       = false
      nodeCount                       = 3
      nodeEnablePublicIpAccess        = false
      nodeServerEdition               = "MemoryOptimized"
      nodeStorageQuotaInMb            = 524288
      nodeVCores                      = 8
      postgresqlVersion               = "15"
      preferredPrimaryZone            = "1"
    }
    tags = {
    }
  })

  schema_validation_enabled = false
  ignore_missing_property   = true
}

resource "azapi_resource" "serverGroupsv22" {
  type      = "Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08"
  name      = "acctest1681"
  parent_id = azurerm_resource_group.test.id

  body = jsonencode({
    location = "westus"
    properties = {
      pointInTimeUTC   = "2023-03-23T08:30:37.467Z"
      sourceLocation   = "westus"
      sourceResourceId = azapi_resource.serverGroupsv2.id
    }
  })

  schema_validation_enabled = false
  ignore_missing_property   = false
}

