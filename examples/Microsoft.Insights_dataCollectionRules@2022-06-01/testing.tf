
resource "azapi_resource" "dataCollectionRule" {
  type      = "Microsoft.Insights/dataCollectionRules@2022-06-01"
  name      = "acctest818"
  parent_id = azurerm_resource_group.test.id

  body = jsonencode({
    location = azurerm_resource_group.test.location
    properties = {
      dataFlows = [
        {
          destinations = [
            "centralWorkspace",
          ]
          streams = [
            "Microsoft-Perf",
            "Microsoft-Syslog",
            "Microsoft-WindowsEvent",
          ]
        },
      ]
      dataSources = {
        performanceCounters = [
          {
            counterSpecifiers = [
              "\\Processor(_Total)\\% Processor Time",
              "\\Memory\\Committed Bytes",
              "\\LogicalDisk(_Total)\\Free Megabytes",
              "\\PhysicalDisk(_Total)\\Avg. Disk Queue Length",
            ]
            name                       = "cloudTeamCoreCounters"
            samplingFrequencyInSeconds = 15
            streams = [
              "Microsoft-Perf",
            ]
          },
          {
            counterSpecifiers = [
              "\\Process(_Total)\\Thread Count",
            ]
            name                       = "appTeamExtraCounters"
            samplingFrequencyInSeconds = 30
            streams = [
              "Microsoft-Perf",
            ]
          },
        ]
        syslog = [
          {
            facilityNames = [
              "cron",
            ]
            logLevels = [
              "Debug",
              "Critical",
              "Emergency",
            ]
            name = "cronSyslog"
            streams = [
              "Microsoft-Syslog",
            ]
          },
          {
            facilityNames = [
              "syslog",
            ]
            logLevels = [
              "Alert",
              "Critical",
              "Emergency",
            ]
            name = "syslogBase"
            streams = [
              "Microsoft-Syslog",
            ]
          },
        ]
        windowsEventLogs = [
          {
            name = "cloudSecurityTeamEvents"
            streams = [
              "Microsoft-WindowsEvent",
            ]
            xPathQueries = [
              "Security![System[(Level = 1 or Level = 2 or Level = 3)]]",
            ]
          },
          {
            name = "appTeam1AppEvents"
            streams = [
              "Microsoft-WindowsEvent",
            ]
            xPathQueries = [
              "System![System[(Level = 1 or Level = 2 or Level = 3)]]",
              "Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]",
            ]
          },
        ]
      }
      destinations = {
        logAnalytics = [
          {
            name                = "centralWorkspace"
            workspaceResourceId = azurerm_log_analytics_workspace.test.id
          },
        ]
      }
    }
  })
  
  depends_on = [
    azurerm_log_analytics_solution.test
  ]

  schema_validation_enabled = false
  ignore_missing_property   = false
}
