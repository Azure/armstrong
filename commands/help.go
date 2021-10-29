package commands

import "log"

const helpMessage = `
Usage: azurerm-rest-api-testing-tool <subcommand> <args>

Commands:
  generate		Generate testing files including terraform configuration for dependency and testing resource.
				Arguments: <filepath of rest api to create arm resource example>

  test			Run tests, it will update dependencies if necessary
				Arguments: [fist testcase name] [second testcase] ignore it to run all tests

  setup			Update dependency for tests

  cleanup		Clean up dependency

  auto			Run generate and test
				Arguments: <filepath of rest api to create arm resource example>

  help			Show help message

`

func Help() {
	log.Printf("[INFO] %v\n", helpMessage)
}
