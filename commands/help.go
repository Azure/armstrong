package commands

import "log"

const helpMessage = `
Usage: azurerm-rest-api-testing-tool <subcommand> <args>

Commands:
  generate		Generate testing files including terraform configuration for dependency and testing resource.
				Arguments: <filepath of rest api to create arm resource example>

  test			Run tests, it will update dependencies if necessary
				Arguments: [fist testcase name] [second testcase] ignore it to run all tests
				Arguments: -v: show terraform logs

  setup			Update dependency for tests
				Arguments: -v: show terraform logs

  cleanup		Clean up dependency
				Arguments: -v: show terraform logs

  auto			Run generate and test
				Arguments: <filepath of rest api to create arm resource example>
				Arguments: -v: show terraform logs

  help			Show help message

`

func Help() {
	log.Printf("[INFO] %v\n", helpMessage)
}
