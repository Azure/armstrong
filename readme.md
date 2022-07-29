# Armstrong - A Terraform based ARM REST API testing Tool

## Introduction
The tool can simplify the process to test an ARM REST API. It can generate a terraform file containing dependencies and a
terraform file containing the testing resource which is based on the [azapi provider](https://github.com/Azure/terraform-provider-azapi).
It can also generate a markdown report when found API issues.

## Usage
```
Usage: armstrong [--version] [--help] <command> [<args>]

Available commands are:
auto        Run generate and test, if test passed, run cleanup
cleanup     Clean up dependencies and testing resource
generate    Generate testing files including terraform configuration for dependencies and testing resource.
setup       Update dependencies for tests
test        Update dependencies for tests and run tests
```

## How to use?
1. Install this tool: `go install github.com/ms-henglu/armstrong`, or download it from [releases](https://github.com/ms-henglu/armstrong/releases).
2. Generate terraform files and Test
    1. Generate testing files by running `armstrong generate -path {path to a swagger 'Create' example}`.
        Here's an example:
        
        `armstrong generate -path .\2021-07-01\examples\Compute\createOrUpdate\ComputeInstanceMinimal.json`.
        
        Then `dependency.tf` and `testing.tf` will be generated. It also supports generate `body` with raw json format, by adding option `-raw`.
        You can append more test cases by executing this command with different example paths, this feature is enabled by default,
        to get a clean working directory, use `-overwrite` option to clean up the existing files.
    2. Run API tests by running `armstrong test`. This command will set up dependencies and test the ARM resource API.
    3. There's an `auto` command, it can generate testing files, then run the tests and remove all resources if test is passed. Example:
    
       `armstrong auto -path .\2021-07-01\examples\Compute\createOrUpdate\ComputeInstanceMinimal.json`

## Troubleshooting
1. Q: When use `test` commands, server side validation error happens.
   
   A: You may need to modify the dependency.tf or testing.tf to meet the server side's requirements. It happens when the testing resource requires running in specific region or some configurations of its dependencies. After modification, run `test` command to continue the test.

2. Q: Will dependencies be removed after testing is done?
   
    A: If using `test` command, resources won't be removed after testing, user must use `cleanup` command to remove these resources. If using `auto` command, it will run `cleanup` command if test passed,


## Features
- [ ] Generate multiple test cases from given resource type and swagger file
- [ ] Generate test cases containing all defined properties
- [ ] Support complicated dependency analysis, ex: key vault id, key vault cert id
- [x] Support appending more testcases.
- [x] Support `body` in both `jsonencode` format and raw json format.
- [x] Improve accuracy in mapping between resourceId and azurerm resource type
- [x] Improve accuracy in azurerm resource example configuration: example configuration must be valid
- [x] Hide terraform logs and generate a more friendly report
- [x] Generate Markdown report when there are API issues


