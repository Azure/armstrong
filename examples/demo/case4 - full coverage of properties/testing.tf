resource "azapi_resource" "sqr" {
  type      = "Microsoft.Insights/scheduledQueryRules@2021-08-01"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  name      = "test-sqr"
  body = jsonencode({
    kind = "LogAlert"
    properties = {
      criteria = {
        allOf = [
          {
            query               = <<-QUERY
 requests
 | summarize CountByCountry=count() by client_CountryOrRegion
QUERY
            timeAggregation     = "Count"
            operator            = "GreaterThan"
            threshold           = 5.0
            resourceIdColumn    = "client_CountryOrRegion"
            dimensions = [
              {
                "name" : "client_CountryOrRegion",
                "operator" : "Include",
                "values" : [
                  "*"
                ]
              }
            ]
            failingPeriods = {
              numberOfEvaluationPeriods = 1
              minFailingPeriodsToAlert  = 1
            }
          }
        ]
      }

      enabled                = false
      evaluationFrequency    = "PT5M"
      windowSize             = "PT5M"
      overrideQueryTimeRange = "PT10M"
      muteActionsDuration    = "PT10M",
      targetResourceTypes = [
        "microsoft.insights/components"
      ],
      actions = {
        actionGroups = [azurerm_monitor_action_group.test.id]
        customProperties = {
          email_subject          = "Email Header"
          custom_webhook_payload = ""
        }
      }
      scopes                                = [azurerm_application_insights.test.id]
      autoMitigate                          = false
      checkWorkspaceAlertsStorageConfigured = false
      skipQueryValidation                   = false
      description                           = "sqr description"
      displayName                           = "1"
      severity                              = 0
    }
  })
}

resource "azapi_resource" "sqr2" {
  type      = "Microsoft.Insights/scheduledQueryRules@2021-08-01"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  name      = "test-sqr2"
  body = jsonencode({
    kind = "LogAlert"
    properties = {
      criteria = {
        allOf = [
          {
            query               = <<-QUERY
 requests
 | summarize CountByCountry=count() by client_CountryOrRegion
QUERY
            timeAggregation     = "Maximum"
            threshold               = 17.5
            operator                = "LessThan"
            resourceIdColumn    = "client_CountryOrRegion"
            metricMeasureColumn = "CountByCountry"
            failingPeriods = {
              numberOfEvaluationPeriods = 1
              minFailingPeriodsToAlert  = 1
            }
          }
        ]
      }
      evaluationFrequency    = "PT5M"
      windowSize             = "PT5M"
      scopes                                = [azurerm_application_insights.test.id]
      severity                              = 0
    }
  })
}

resource "azapi_resource" "sqr3" {
  type      = "Microsoft.Insights/scheduledQueryRules@2021-08-01"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  name      = "test-sqr3"
  body = jsonencode({
    kind = "LogToMetric"
    tags = {
      ENV = "Test"
    }
    properties = {
      criteria = {
        allOf = [
          {
            timeAggregation     = "Average"
            operator                = "LessThan"
            metricName    = "Average_% Idle Time"
          }
        ]
      }
      scopes                                = [azurerm_log_analytics_workspace.test.id]
    }
  })
}
