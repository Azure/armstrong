## v0.16.1
BUG FIXES:
- Fix a bug the azapi examples are not correctly loaded.

## v0.16.0

ENHANCEMENTS:
- Add more azapi examples for the dependency resolver.
- Update mapping of azurerm v4.20.0

BUG FIXES:
- Fix mockResourceIDFromType for special resource types.

## v0.15.1

BUG FIXES:
- Fix a bug that root level discriminator is not handled correctly.

## v0.15.0

ENHANCEMENTS:
- Add operation properties coverage report for API test scenarios
- Hide readonly array item in coverage report
- Update mapping of azurerm v4.0.1

BUG FIXES:
- Fix a bug that undocumented swagger properties caused inaccurate coverage report
- Remove duplicated items in the error report.

## v0.14.0
FEATURES:
- Support generating and testing with `azapi` dynamic property feature.

ENHANCEMENTS:
- Update mapping of azurerm v3.100.0

BUG FIXES:
- Fix the bug that the default resource name is not randomly generated.

## v0.13.0
FEATURES:
- Support `credscan` command to scan the credentials in the testing configuration files.

ENHANCEMENTS:
- Update mapping of azurerm v3.95.0

## v0.12.1
BUG FIXES:
- Fix the bug that resource types with same name but different casing are not handled correctly.
- Fix the bug that coverage reports are generated even if there are no valid test cases.
- Fix the bug that other dependency resolvers are not called when error occurs in the previous dependency resolver.
- Fix the bug that the generated resources are not in the correct order.
- Fix the bug that `HEAD` method and other methods combination couldn't be handled correctly.
- Fix the bug that error reports are not generated for other azapi resources.
- Fix the bug that error reports are not generated when error happens during the `terraform plan`.

## v0.12.0
FEATURES:
- Generate multiple test cases from one or multiple swagger spec files.
- Support using verified azapi examples as automatically generated dependencies.
- Support `azapi_resource_id` data source as automatically generated dependencies.
- Support generating coverage report from local swagger specs.
- Support swagger accuracy report.

BUG FIXES:
- Fix the panic when generating from the swagger example which doesn't have `api-version` field.

## v0.11.0
FEATURES:
- Support coverage report
- Support cleanup report

ENHANCEMENTS:

BUG FIXES:

## v0.10.0

FEATURES:
- Support generate data source from Swagger GET example.
- Support generate dependency automatically for identity.

ENHANCEMENTS:

BUG FIXES:
- Fix error message in reports.

## v0.9.0

ENHANCEMENTS:
- Update mapping of azurerm v3.50.0

## v0.8.1

BUG FIXES:
- Fix invalid characters in folder name on Windows

## v0.8.0

FEATURES:
- Support error report

ENHANCEMENTS:
- Update mapping of azurerm v3.41.0

## v0.7.0

FEATURES:
- Generate a passed report or partially passed report

ENHANCEMENTS:
- Update mapping of azurerm v3.30.0

## v0.6.0

BUG FIXES:
- Dependency detection failed when working-dir is specified

## v0.5.0

ENHANCEMENTS:
- Update mapping of azurerm v3.24.0

BUG FIXES:
- Wrap keys which start with numbers.

## v0.4.0

ENHANCEMENTS:
- Update mapping of azurerm v3.22.0

## v0.3.0

FEATURES:
- Supports validate command to preview the resource changes.
- Supports -working-dir option to specify the working directory.

## v0.2.1

FEATURES:
- Generated document improvement: now differences are highlighted.

## v0.2.0

FEATURES:
- Support install Terraform automatically.
- Support `-raw` option, which allows user to use raw json payload. The default payload will use jsonencode function.
- Support `-overwrite` option, which allows user to overwrite existing configurations. The default behavior is appending test cases on the existing configurations.
- Support markdown report generation: The `test` command can generate markdown report when it found bugs for each test cases.

## v0.1.0
Initial release.