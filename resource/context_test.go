package resource_test

import (
	"regexp"
	"testing"

	"github.com/azure/armstrong/resource"
)

func Test_NewContextInit(t *testing.T) {
	context := resource.NewContext(nil)
	expected := `terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    key_vault {
      purge_soft_delete_on_destroy       = false
      purge_soft_deleted_keys_on_destroy = false
    }
  }
  skip_provider_registration = true
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest\d+"
}

variable "location" {
  type    = string
  default = "westeurope"
}
`
	actual := context.String()
	r := regexp.MustCompile(expected)
	if !r.MatchString(actual) {
		t.Fatalf("expected: %s, got: %s", expected, actual)
	}
}
