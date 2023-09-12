package resource_test

import (
	"testing"

	"github.com/ms-henglu/armstrong/resource"
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
`
	actual := context.String()
	if actual != expected {
		t.Fatalf("expected: %s, got: %s", expected, actual)
	}
}
