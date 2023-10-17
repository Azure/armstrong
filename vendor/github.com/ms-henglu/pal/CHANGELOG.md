## v0.4.0
BREAKING CHANGES:
- Only support logs from terraform-provider-azapi v1.10.0 or above.

ENHANCEMENTS:
- The redundant query parameters are removed in the `markdown` format.
- Remove the `/providers` API from parsed logs, because its response couldn't be parsed.

BUG FIXES:
- Fix the issue that some resources may not be outputted to `azapi` format.

## v0.3.0

FEATURES:
- Support parsing terraform logs to `azapi` format.

BUG FIXES:
- Fix the issue that the parsed URL paths are not normalized.
- Fix the issue that the request body from azurerm provider may not be parsed correctly, when the request body is pretty printed JSON.

## v0.2.0

FEATURES:
- Support parsing terraform logs to `oav` traffic format.
- Support `-version` option to show the version information.
- Support `-help` option to show the help information.
- Support `-o` option to specify the output directory.
- Support `-i` option to specify the input file path.
- Support `-m` option to specify the output format.

BUG FIXES:
- Fix the issue that response headers may contain duplicated values.
- Fix the issue that logs from released `azurerm` provider may not be parsed correctly.

## v0.1.0

FEATURES:

- Support parsing terraform logs to markdown format.
- Support `azurerm`, `azuread` and `azapi` providers.