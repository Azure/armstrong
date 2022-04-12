
resource "azapi_resource" "test" {
	name      = "acctest6988"
	parent_id = azurerm_databricks_workspace.test.id
	type      = "Microsoft.Databricks/workspaces/virtualNetworkPeerings@2021-04-01-preview"
	body      = <<BODY
{
    "properties": {
        "allowForwardedTraffic": false,
        "allowGatewayTransit": false,
        "allowVirtualNetworkAccess": true,
        "remoteVirtualNetwork": {
            "id": "${azurerm_virtual_network.test.id}"
        },
        "useRemoteGateways": false
    }
}
BODY
}
