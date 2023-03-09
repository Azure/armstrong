resource "azurerm_resource_group" "test" {
  name     = "acctest6981"
  location = "West Europe"
}

resource "azurerm_virtual_wan" "test" {
  name                = "acctest1829"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_virtual_hub" "test" {
  name                = "acctest7117"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_wan_id      = azurerm_virtual_wan.test.id
  address_prefix      = "10.0.1.0/24"
}

resource "azurerm_express_route_gateway" "test" {
  name                = "acctest7416"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  virtual_hub_id      = azurerm_virtual_hub.test.id
  scale_units         = 1
}

resource "azurerm_express_route_port" "test" {
  name                = "acctest2357"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  peering_location    = "Equinix-Seattle-SE2"
  bandwidth_in_gbps   = 10
  encapsulation       = "Dot1Q"
}

resource "azurerm_express_route_circuit" "test" {
  name                  = "acctest7855"
  location              = azurerm_resource_group.test.location
  resource_group_name   = azurerm_resource_group.test.name
  express_route_port_id = azurerm_express_route_port.test.id
  bandwidth_in_gbps     = 5

  sku {
    tier   = "Premium"
    family = "MeteredData"
  }
}

resource "azurerm_express_route_circuit_peering" "test" {
  peering_type                  = "AzurePrivatePeering"
  express_route_circuit_name    = azurerm_express_route_circuit.test.name
  resource_group_name           = azurerm_resource_group.test.name
  shared_key                    = "ItsASecret"
  peer_asn                      = 100
  primary_peer_address_prefix   = "192.168.1.0/30"
  secondary_peer_address_prefix = "192.168.2.0/30"
  vlan_id                       = 100
}

resource "azurerm_express_route_connection" "test" {
  name                             = "acctest2735"
  express_route_gateway_id         = azurerm_express_route_gateway.test.id
  express_route_circuit_peering_id = azurerm_express_route_circuit_peering.test.id
}

resource "azurerm_virtual_network" "test" {
  name                = "acctest5784"
  address_space       = ["10.0.0.0/16"]
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  tags = {
    environment = "Production"
  }
}

resource "azurerm_subnet" "test" {
  name                 = "RouteServerSubnet"
  virtual_network_name = azurerm_virtual_network.test.name
  resource_group_name  = azurerm_resource_group.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurerm_public_ip" "test" {
  name                = "acctest9073"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  allocation_method   = "Static"
  sku                 = "Standard"
}

resource "azurerm_route_server" "test" {
  name                             = "acctest3106"
  resource_group_name              = azurerm_resource_group.test.name
  location                         = azurerm_resource_group.test.location
  sku                              = "Standard"
  public_ip_address_id             = azurerm_public_ip.test.id
  subnet_id                        = azurerm_subnet.test.id
  branch_to_branch_traffic_enabled = true
}

