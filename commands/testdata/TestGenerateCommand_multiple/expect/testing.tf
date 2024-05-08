
data "azapi_resource" "sqlVulnerabilityAssessment" {
  type      = "Microsoft.Sql/servers/databases/sqlVulnerabilityAssessments@2022-02-01-preview"
  parent_id = azapi_resource.database.id
  name      = "default"
}


resource "azapi_resource" "baseline" {
  type      = "Microsoft.Sql/servers/databases/sqlVulnerabilityAssessments/baselines@2022-02-01-preview"
  parent_id = data.azapi_resource.sqlVulnerabilityAssessment.id
  name      = var.resource_name
  body = {
    properties = {
      latestScan = false
      results = {
        VA2063 = [
          [
            "AllowAll",
            "0.0.0.0",
            "255.255.255.255",
          ],
        ]
        VA2065 = [
          [
            "AllowAll",
            "0.0.0.0",
            "255.255.255.255",
          ],
        ]
      }
    }
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}
