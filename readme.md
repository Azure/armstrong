# Armstrong - A Terraform based ARM REST API testing tool

## Introduction
The tool can simplify the process to test an ARM REST API. It can generate a terraform file containing dependencies and a
terraform file containing the testing resource which is based on the [azapi provider](https://github.com/Azure/terraform-provider-azapi).
It can also generate a markdown report when found API issues.

## Guidance

For ARM review, please refer to [guidance for API test](./docs/guidance-for-api-test.md).

## Install

Install this tool: `go install github.com/ms-henglu/armstrong`, or download it from [releases](https://github.com/ms-henglu/armstrong/releases).

## Usage
```
Usage: armstrong [--version] [--help] <command> [<args>]

Available commands are:
cleanup     Clean up dependencies and testing resource
generate    Generate testing files including terraform configuration for dependencies and testing resource.
test        Update dependencies for tests and run tests
validate    Generates a speculative execution plan, showing what actions Terraform would take to apply the current configuration
```

## Commands

### generate - Generate testing files

Armstrong supports generate testcases from different kinds of input.

Supported options:
1. `-working-dir`: Specify the working directory which stores the output config, default is current directory.
2. `-raw`: Generate `body` with raw json format, default is false.
3. `-v`: Enable verbose mode, default is false.

Supported inputs:
1. Generate testcase from swagger 'Create' example:
```shell
armstrong generate -path {path to a swagger 'Create' example}
```
In this mode, it supports `-overwrite` option to clean up the existing files, default is false, which means append more testcases.
And it supports `-type` option to specify the resource type, allowed values: `resource`(supports CRUD) and `data`(read-only), default is `resource`.

2. Generate multiple testcases from swagger spec, it supports both path to the swagger file and the directory containing swagger files.
```shell
armstrong generate -swagger {path/dir to swagger spec}
```

3. Generate multiple testcases from an autorest configuration file and its tag:
```shell
armstrong generate -readme {path to autorest configuration file} -tag {tag name}
```

### validate - Validate the changes

This command generates a speculative execution plan, showing what actions Terraform would take to apply the current configuration.

```shell
armstrong validate
```

Supported options:
1. `-working-dir`: Specify the working directory which stores the output config, default is current directory.
2. `-v`: Enable verbose mode, default is false.

### test - Run tests

This command will set up dependencies and test the ARM resource API.

```shell
armstrong test
```

Supported options:
1. `-working-dir`: Specify the working directory which stores the output config, default is current directory.
2. `-v`: Enable verbose mode, default is false.
3. `-destroy-after-test`: Destroy the testing resource after test, default is false.
4. `-swagger`: Specify the swagger file path or directory path.

Armstrong also output different kinds of reports:
1. `all_passed_report.md`: A markdown report which contains all passed testcases. It will be generated when all testcases passed.
It also contains the `coverage report` which shows the tested properties and the total properties.
2. `partial_passed_report.md`: A markdown report which contains all passed testcases. It will be generated when there are failed testcases.
It also contains the `coverage report` which shows the tested properties and the total properties.
3. `api error report`: A markdown report which contains one API error when creating the testing resource. It will be generated when there are API issues.
It also contains other details like http traces to help debugging.
4. `api issue report`: A markdown report which contains one API issue when testing the resource. It will be generated when there are API issues.
It also contains other details like http traces to help debugging.
5. `swagger accuracy report`: A html report which contains the swagger accuracy analysis result. It will be generated when `-swagger` option is specified and `oav` is installed.

**Notice:**
1. How to install `oav`, please refer to [oav](https://github.com/Azure/oav).
2. The `coverage report` is generated based on the [public swagger repo](https://github.com/Azure/azure-rest-api-specs) by default, but it can be changed to the local swagger specs by specifying `-swagger` option.

### cleanup - Clean up dependencies and testing resource

```shell
armstrong cleanup
```

Supported options:
1. `-working-dir`: Specify the working directory which stores the output config, default is current directory.
2. `-v`: Enable verbose mode, default is false.

Armstrong also output different kinds of reports:
1. `cleanup_all_passed_report`: A markdown report which contains all passed testcases. It will be generated when all testcases passed.
2. `cleanup_partial_passed_report`: A markdown report which contains all passed testcases. It will be generated when there are failed testcases.
3. `cleanup_api error report`: A markdown report which contains one API error when deleting the testing resource. It will be generated when there are API issues.

### report - Generate a summary report

**Notice:** The `oav` must be installed, please refer to [oav](https://github.com/Azure/oav).

After multiple testcases are generated from swagger spec and tested, `swagger accuracy report` are generated in each testcase directory.
This command will generate a summary report which contains all `swagger accuracy report` from each testcase directory.

```shell
armstrong report -swagger {path/dir to swagger spec}
```

Supported options:
1. `-working-dir`: Specify the working directory which stores the output config, default is current directory.
2. `-swagger`: Specify the swagger file path or directory path.

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

## Troubleshooting
1. Q: When use `test` commands, server side validation error happens.
   
   A: You may need to modify the dependency.tf or testing.tf to meet the server side's requirements. It happens when the testing resource requires running in specific region or some configurations of its dependencies. After modification, run `test` command to continue the test.

2. Q: Will dependencies be removed after testing is done?
   
    A: If using `test` command, resources won't be removed after testing, user must use `cleanup` command to remove these resources.

## TODO
- [ ] Support extension scoped resource, ex: `{resourceUri}/providers/Microsoft.Advisor/recommendations/{recommendationId}`
- [ ] Support the API path whose key segment is a variable, ex: `.../providers/Microsoft.Web/sites/{name}/host/default/{keyType}/{keyName}`
- [ ] Generate `body` placeholder with all configurable properties when the example is invalid or empty.
- [ ] Support the API paths that don't follow the ARM guidelines: give warning message or handle them in a special way.
