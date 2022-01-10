
resource "azurerm-restapi_resource" "test" {
	name      = "acctest2793"
	parent_id = azurerm_batch_account.test.id
	type      = "Microsoft.Batch/batchAccounts/applications@2021-06-01"
 	body      = <<BODY
{
    "properties": {
        "allowUpdates": false,
        "displayName": "myAppName"
    }
}
BODY
}
