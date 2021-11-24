
resource "azurerm-restapi_resource" "test" {
	resource_id = "${azurerm_batch_account.test.id}/applications/acctest2793"
	type = "Microsoft.Batch/batchAccounts/applications@2021-06-01"
 	body = <<BODY
{
    "properties": {
        "allowUpdates": false,
        "displayName": "myAppName"
    }
}
BODY
}
