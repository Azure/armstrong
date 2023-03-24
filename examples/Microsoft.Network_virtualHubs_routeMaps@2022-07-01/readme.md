# Example

## Introduction

This is an example of testing the API which manages `Microsoft.Network/virtualHubs/routeMaps@2022-07-01` resource.

This example demonstrates how to fix the configuration and continue the test when encounters an error.

## Step-by-step guide

### 1. Generate testing files

The following command will generate the terraform files containing the testing resources and dependencies.

```bash
armstrong generate -path .\Swagger_Create_Example.json

2023/03/24 12:28:33 [INFO] ----------- generate dependency and test resource ---------
2023/03/24 12:28:33 [INFO] loading dependencies
2023/03/24 12:28:33 [INFO] generating testing files
2023/03/24 12:28:33 [INFO] found dependency: azurerm_express_route_connection
2023/03/24 12:28:33 [INFO] found dependency: azurerm_route_server
2023/03/24 12:28:33 [INFO] dependency.tf generated
2023/03/24 12:28:33 [INFO] testing.tf generated
```


### 2. Deploy depdencies and test ARM API

The following command will deploy the dependencies and test the ARM API.

```bash
armstrong test
```

In this step, you'll see the following error:
```
Error: creating/updating IP Configuration "ipConfig1" of Route Server "acctest8596" (Resource Group Name "acctest6251"): network.VirtualHubIPConfigurationClient#CreateOrUpdate: Failure sending request: StatusCode=400 -- Original Error: Code="InvalidRouteServerSubnetName" Message="Subnet name used to deploy Route Server should be RouteServerSubnet." Details=[]

  with azurerm_route_server.test,
  on dependency.tf line 93, in resource "azurerm_route_server" "test":
  93: resource "azurerm_route_server" "test" {
```

This error is caused by the `azurerm_route_server` resource. The `azurerm_route_server` resource requires the subnet name to be `RouteServerSubnet`. So we need to fix the configuration(change the subnet name to `RouteServerSubnet`) and continue the test.

```hcl
resource "azurerm_subnet" "test" {
  name                 = "RouteServerSubnet" // change the subnet name to RouteServerSubnet
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.1.0/24"]
}
```

Here's another error we'll see:
```
Error: waiting for creation of Express Route Connection: (Name "acctest4032" / Express Route Gateway Name "acctest4639" / Resource Group "acctest6251"): Code="InvalidParameter" Message="The creation of the virtual network gateway connection: acctest4032 failed because your circuit in Equinix-Seattle-SE2 cannot be connected to West Europe on a standard circuit. Virtual network gateway connections on a standard ExpressRoute circuit are only allowed within the same geopolitical region. Please upgrade to a premium SKU." Details=[]

  with azurerm_express_route_connection.test,
  on dependency.tf line 61, in resource "azurerm_express_route_connection" "test":
  61: resource "azurerm_express_route_connection" "test" {
```

To fix this error, we can upgrade the SKU of the Express Route Circuit to `Premium` and continue the test.

```hcl
resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest9466"
  ... 

  sku {
    tier   = "Premium" // change the SKU to Premium
    family = "MeteredData"
  }
}
```

And finally we'll see the generated testing reports.