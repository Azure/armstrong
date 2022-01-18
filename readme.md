# Terraform based API testing Tool

## Introduction
The tool can simplify the process to test a ARM rest API. It can generate a terraform file containing dependencies and a terraform file containing the testing resource which is based on the [generic azurerm provider](https://github.com/Azure/terraform-provider-azurerm-restapi).

## Usage
```
Usage: azurerm-restapi-to-azurerm [--version] [--help] <command> [<args>]

Available commands are:
auto        Run generate and test, if test passed, run cleanup
cleanup     Clean up dependencies and testing resource
generate    Generate testing files including terraform configuration for dependencies and testing resource.
setup       Update dependencies for tests
test        Update dependencies for tests and run tests
```

## How to use?
1. Requisites
    1. Download and setup Terraform.
    2. Setup the generic terraform provider by following [this document](https://github.com/Azure/terraform-provider-azurerm-restapi/blob/main/README.md).
2. Install this tool: `go install github.com/ms-henglu/azurerm-restapi-testing-tool`
3. Generate terraform files and Test
    1.  Generate testing files by running `azurerm-restapi-testing-tool generate -path path_to_swagger_example`.
        Here's an example:
        
        `azurerm-restapi-testing-tool generate -path C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\machinelearningservices\resource-manager\Microsoft.MachineLearningServices\stable\2021-07-01\examples\Compute\createOrUpdate\ComputeInstanceMinimal.json`.
        
        Then `dependency.tf` and `testing.tf` will be generated.
    2. Run API tests by running `azurerm-restapi-testing-tool test`. This command will set up dependencies and test the ARM resource API.
    3. There's an `auto` command, it can generate testing files, then run the tests and remove all resources if test is passed. Example:
    
       `azurerm-restapi-testing-tool auto -path C:\Users\henglu\go\src\github.com\Azure\azure-rest-api-specs\specification\machinelearningservices\resource-manager\Microsoft.MachineLearningServices\stable\2021-07-01\examples\Compute\createOrUpdate\ComputeInstanceMinimal.json`

## Troubleshooting
1. Q: When use `test` commands, server side validation error happens.
   
   A: You may need to modify the dependency.tf or testing.tf to meet the server side's requirements. It happens when the testing resource requires running in specific region or some configurations of its dependencies. After modification, run `test` command to continue the test.

2. Q: Will dependencies be removed after testing is done?
   
    A: If using `test` command, resources won't be removed after testing, user must use `cleanup` command to remove these resources. If using `auto` command, it will run `cleanup` command if test passed,


## Features
- [ ] Generate multiple test cases from given resource type and swagger file
- [ ] Generate test cases containing all defined properties
- [ ] Support complicated dependency analysis, ex: key vault id, key vault cert id
- [x] Improve accuracy in mapping between resourceId and azurerm resource type
- [x] Improve accuracy in azurerm resource example configuration: example configuration must be valid
- [x] Hide terraform logs and generate a more friendly report


