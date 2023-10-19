package resource_test

import (
	"regexp"
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
