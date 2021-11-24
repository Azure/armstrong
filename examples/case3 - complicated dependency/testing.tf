
resource "azurerm-restapi_resource" "test" {
	resource_id = "${azurerm_databricks_workspace.test.id}/virtualNetworkPeerings/acctest6988"
	type = "Microsoft.Databricks/workspaces/virtualNetworkPeerings@2021-04-01-preview"
	body = <<BODY
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
