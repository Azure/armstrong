package utils_test

import (
	"testing"

	"github.com/azure/armstrong/utils"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func Test_TypeValue(t *testing.T) {
	testcases := []struct {
		Input    string
		Expected []string
	}{
		{
			Input: `
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "Spring" {
  type      = "Microsoft.AppPlatform/Spring@2023-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      zoneRedundant = false
    }
    sku = {
      name = "E0"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "app" {
  type      = "Microsoft.AppPlatform/Spring/apps@2023-05-01-preview"
  parent_id = azapi_resource.Spring.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    identity = {
      type = "None"
    }
    properties = {
      customPersistentDisks = [
      ]
      enableEndToEndTLS = false
      public            = false
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "domain" {
  type      = "Microsoft.AppPlatform/Spring/apps/domains@2022-04-01"
  parent_id = azapi_resource.app.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      certName   = "mycert"
      thumbprint = "934367bf1c97033f877db0f15cb1b586957d3133"
    }
  })
  schema_validation_enabled = false
}

resource "azapi_resource_action" "patch_domain" {
  type        = "Microsoft.AppPlatform/Spring/apps/domains@2022-04-01"
  resource_id = azapi_resource.domain.id
  action      = ""
  method      = "PATCH"
  body = jsonencode({
    properties = {
      certName   = "mycert"
      thumbprint = "934367bf1c97033f877db0f15cb1b586957d3133"
    }
  })
}

data "azapi_resource_list" "listDomainsByApp" {
  type       = "Microsoft.AppPlatform/Spring/apps/domains@2022-04-01"
  parent_id  = azapi_resource.app.id
  depends_on = [azapi_resource.domain]
}

`,
			Expected: []string{
				"Microsoft.Resources/resourceGroups@2020-06-01",
				"Microsoft.AppPlatform/Spring@2023-05-01-preview",
				"Microsoft.AppPlatform/Spring/apps@2023-05-01-preview",
				"Microsoft.AppPlatform/Spring/apps/domains@2022-04-01",
				"Microsoft.AppPlatform/Spring/apps/domains@2022-04-01",
				"Microsoft.AppPlatform/Spring/apps/domains@2022-04-01",
			},
		},
	}

	for _, tc := range testcases {
		file, diags := hclwrite.ParseConfig([]byte(tc.Input), "", hcl.InitialPos)
		if diags.HasErrors() {
			t.Error(diags)
		}
		actuals := make([]string, 0)
		for _, block := range file.Body().Blocks() {
			if block.Type() != "data" && block.Type() != "resource" {
				continue
			}
			actual := utils.TypeValue(block)
			actuals = append(actuals, actual)
		}
		if len(actuals) != len(tc.Expected) {
			t.Errorf("expected %d, got %d", len(tc.Expected), len(actuals))
			continue
		}
		for i, actual := range actuals {
			if actual != tc.Expected[i] {
				t.Errorf("expected %s, got %s", tc.Expected[i], actual)
			}
		}
	}
}
