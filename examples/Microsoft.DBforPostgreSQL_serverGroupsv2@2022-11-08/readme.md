# Example

## Introduction

This is an example of testing the API which manages `Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08` resource.

This example demonstrates how to test the API which has a dependency on itself.

This resource can be restored from another `serverGroupsv2` resource. It uses the `sourceResourceId` property, which is a reference to the `id` property of the `Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08` resource.

But this resource is not supported by `azurerm` provider, so `armstrong` can't generate the dependency automatically.

## Step-by-step guide

### 1. Run generate command with create example

The following command will generate the terraform files containing the testing resources and its dependencies.

```bash
armstrong generate --path ./Swagger_Create_Example.json
```

Result:
```hcl
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
```

### 2. Run generate command with restore example

The restore example creates the resource from another existing resource. The following command will analyze the generated terraform files and find the dependency, then append the new testing resource.

```bash
armstrong generate --path ./Swagger_Restore_Example.json
```

Result:
```hcl
resource "azapi_resource" "serverGroupsv22" {
  type      = "Microsoft.DBforPostgreSQL/serverGroupsv2@2022-11-08"
  name      = "acctest1681"
  parent_id = azurerm_resource_group.test.id

  body = jsonencode({
    location = "westus"
    properties = {
      pointInTimeUTC   = "2023-03-23T08:30:37.467Z"
      sourceLocation   = "westus"
      sourceResourceId = azapi_resource.serverGroupsv2.id // it refers to the id of the resource created by the create example
    }
  })

  schema_validation_enabled = false
  ignore_missing_property   = false
}
```

### 3. Run test command

The following command will run the test and generate testing results.

```bash
armstrong test
```