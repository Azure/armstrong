
resource "azurermg_resource" "test" {
	url = "${azurerm_machine_learning_workspace.test.id}/computes/acctest1275"
	api_version = "2021-07-01"
 	body = <<BODY
{
    "location": "westeurope",
    "properties": {
        "computeType": "ComputeInstance",
        "properties": {
            "vmSize": "STANDARD_NC6"
        }
    }
}
BODY
}
