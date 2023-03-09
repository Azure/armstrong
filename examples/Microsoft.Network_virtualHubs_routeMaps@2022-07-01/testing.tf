
resource "azapi_resource" "routeMap" {
  type      = "Microsoft.Network/virtualHubs/routeMaps@2022-07-01"
  name      = "acctest836"
  parent_id = azurerm_route_server.test.id

  body = jsonencode({
    properties = {
      associatedInboundConnections = [
        azurerm_express_route_connection.test.id,
      ]
      associatedOutboundConnections = [
      ]
      rules = [
        {
          actions = [
            {
              parameters = [
                {
                  asPath = [
                    "22334",
                  ]
                  community = [
                  ]
                  routePrefix = [
                  ]
                },
              ]
              type = "Add"
            },
          ]
          matchCriteria = [
            {
              asPath = [
              ]
              community = [
              ]
              matchCondition = "Contains"
              routePrefix = [
                "10.0.0.0/8",
              ]
            },
          ]
          name              = "rule1"
          nextStepIfMatched = "Continue"
        },
      ]
    }
  })

  schema_validation_enabled = false
  ignore_missing_property   = false
}
