
resource "azurerm-restapi_resource" "test" {
	resource_id = "${azurerm_spring_cloud_service.test.id}/apps/acctest7600"
	type = "Microsoft.AppPlatform/Spring/apps@2020-07-01"
 	body = <<BODY
{
    "location": "westeurope",
    "properties": {
        "activeDeploymentName": "mydeployment1",
        "fqdn": "myapp.mydomain.com",
        "httpsOnly": false,
        "persistentDisk": {
            "mountPath": "/mypersistentdisk",
            "sizeInGB": 2
        },
        "public": true,
        "temporaryDisk": {
            "mountPath": "/mytemporarydisk",
            "sizeInGB": 2
        }
    }
}
BODY
}
