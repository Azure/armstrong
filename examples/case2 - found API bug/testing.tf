
resource "azapi_resource" "test" {
	name      = "acctest7600"
	parent_id = azurerm_spring_cloud_service.test.id
	type      = "Microsoft.AppPlatform/Spring/apps@2020-07-01"
 	body      = <<BODY
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
