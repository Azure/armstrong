package helper_test

import (
	"testing"

	"github.com/ms-henglu/armstrong/helper"
)

func Test_GetResourceFromHcl(t *testing.T) {
	config :=
		`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "henglu112-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "henglu0113"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "henglu0113"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "henglu0113sw"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}
data "azurerm_client_config" "current" {}

resource "azurerm_synapse_role_assignment" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  role_name            = "Synapse SQL Administrator"
  principal_id         = data.azurerm_client_config.current.object_id

}
`
	testcases := []struct {
		Input  string
		Expect string
	}{
		{
			Input:  "azurerm_synapse_role_assignment",
			Expect: "azurerm_synapse_role_assignment.example",
		},
		{
			Input:  "azurerm_resource_group",
			Expect: "azurerm_resource_group.example",
		},
		{
			Input:  "azurerm_storage_data_lake_gen2_filesystem",
			Expect: "azurerm_storage_data_lake_gen2_filesystem.example",
		},
		{
			Input:  "azurerm_client_config",
			Expect: "azurerm_client_config.current",
		},
	}

	for _, testcase := range testcases {
		if output := helper.GetResourceFromHcl(config, testcase.Input); output != testcase.Expect {
			t.Fatalf("expect %v but got %v", testcase.Expect, output)
		}
	}
}

func Test_GetCombinedHcl(t *testing.T) {
	old :=
		`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "henglu112-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "henglu0113"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}
`
	new := `
resource "azurerm_storage_account" "example" {
  name                     = "henglu0113"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "henglu0113"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "henglu0113sw"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}
data "azurerm_client_config" "current" {}

resource "azurerm_synapse_role_assignment" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  role_name            = "Synapse SQL Administrator"
  principal_id         = data.azurerm_client_config.current.object_id

}
`
	output := helper.GetCombinedHcl(old, new)
	expect := `provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "henglu112-resources"
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  name                     = "henglu0113"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  name               = "henglu0113"
  storage_account_id = azurerm_storage_account.example.id
}

resource "azurerm_synapse_workspace" "example" {
  name                                 = "henglu0113sw"
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}

data "azurerm_client_config" "current" {}

resource "azurerm_synapse_role_assignment" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  role_name            = "Synapse SQL Administrator"
  principal_id         = data.azurerm_client_config.current.object_id

}

`
	if expect != output {
		t.Fatalf("expect %s but got %s", expect, output)
	}
}

func Test_GetRenamedHcl(t *testing.T) {
	config :=
		`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  location = "West Europe"
}

resource "azurerm_storage_account" "example" {
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}

resource "azurerm_storage_data_lake_gen2_filesystem" "example" {
  storage_account_id = azurerm_storage_account.example.id
}
resource "azurerm_storage_data_lake_gen2_filesystem" "example1" {
  storage_account_id = azurerm_storage_account.example.id
}
resource "azurerm_synapse_workspace" "example" {
  resource_group_name                  = azurerm_resource_group.example.name
  location                             = azurerm_resource_group.example.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.example.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}
data "azurerm_client_config" "current" {}

resource "azurerm_synapse_role_assignment" "example" {
  synapse_workspace_id = azurerm_synapse_workspace.example.id
  role_name            = "Synapse SQL Administrator"
  principal_id         = data.azurerm_client_config.current.object_id

}
`
	output := helper.GetRenamedHcl(config)

	expect := `resource "azurerm_resource_group" "test" {
  location = "West Europe"
}
resource "azurerm_storage_account" "test" {
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
  account_kind             = "StorageV2"
  is_hns_enabled           = "true"
}
resource "azurerm_storage_data_lake_gen2_filesystem" "test" {
  storage_account_id = azurerm_storage_account.test.id
}
resource "azurerm_storage_data_lake_gen2_filesystem" "test2" {
  storage_account_id = azurerm_storage_account.test.id
}
resource "azurerm_synapse_workspace" "test" {
  resource_group_name                  = azurerm_resource_group.test.name
  location                             = azurerm_resource_group.test.location
  storage_data_lake_gen2_filesystem_id = azurerm_storage_data_lake_gen2_filesystem.test.id
  sql_administrator_login              = "sqladminuser"
  sql_administrator_login_password     = "H@Sh1CoR3!"
}
data "azurerm_client_config" "test" {}
resource "azurerm_synapse_role_assignment" "test" {
  synapse_workspace_id = azurerm_synapse_workspace.test.id
  role_name            = "Synapse SQL Administrator"
  principal_id         = data.azurerm_client_config.test.object_id

}
`
	if expect != output {
		t.Fatalf("expect %s but got %s", expect, output)
	}
}
