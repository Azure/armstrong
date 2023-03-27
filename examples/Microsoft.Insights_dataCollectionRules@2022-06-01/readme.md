# Example

## Introduction

This is an example of testing the API which manages `Microsoft.Insights/dataCollectionRules@2022-06-01` resource.

This example demonstrates how to test the API which has a implicit dependency. The dependency is not defined in the swagger file, but it's required by the API.

## Step-by-step guide

### 1. Run generate command with create example

The following command will generate the terraform files containing the testing resources and its dependencies.

```bash
armstrong generate --path ./Swagger_Create_Example.json
```

### 2. Run test command

The following command will run the test.

```bash
armstrong test
```

Then you'll see the following error:

```bash
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: InvalidPayload
--------------------------------------------------------------------------------
{
  "error": {
    "code": "InvalidPayload",
    "message": "Data collection rule is invalid",
    "details": [
      {
        "code": "InvalidOutputTable",
        "message": "Table for output stream 'Microsoft-WindowsEvent' is not available for destination 'centralWorkspace'.",
        "target": "properties.dataFlows[0]"
      }
    ]
  }
}
```

This error is caused by the implicit dependency. The `Microsoft-WindowsEvent` table is not available for the `centralWorkspace` destination. So we need to create the `Microsoft-WindowsEvent` table first.

### 3. Fix the configuration

You can manage this dependency by azurerm provider or azapi provider, here's an example:

```hcl
// azurerm example
resource "azurerm_log_analytics_solution" "test" {
  solution_name         = "WindowsEventForwarding"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  workspace_resource_id = azurerm_log_analytics_workspace.test.id
  workspace_name        = azurerm_log_analytics_workspace.test.name
  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/WindowsEventForwarding"
  }
}

// azapi example
resource "azapi_resource" "solution" {
  type = "Microsoft.OperationsManagement/solutions@2015-11-01-preview"
  parent_id = azurerm_resource_group.test.id
  name = "WindowsEventForwarding(${azurerm_log_analytics_workspace.test.name})"
  location = azurerm_resource_group.test.location
  body = jsonencode({
    plan = {
      name = "WindowsEventForwarding(${azurerm_log_analytics_workspace.test.name})"
      publisher = "Microsoft"
      product = "OMSGallery/WindowsEventForwarding"
      promotionCode = ""
    }
    properties = {
      workspaceResourceId = azurerm_log_analytics_workspace.test.id
    }
  })
}

```

### 4. Run test command again

The following command will run the test again.

```bash
armstrong test
```

Then you'll see the passed results.