# Guidance for API Test

## Video

This [video](https://microsoftapc-my.sharepoint.com/:v:/g/personal/henglu_microsoft_com/EfdL8LOuGD9OoKHEXmmlYg4BOpq2rMHDMdta_drq35QW8A?nav=eyJyZWZlcnJhbEluZm8iOnsicmVmZXJyYWxBcHAiOiJTdHJlYW1XZWJBcHAiLCJyZWZlcnJhbFZpZXciOiJTaGFyZURpYWxvZyIsInJlZmVycmFsQXBwUGxhdGZvcm0iOiJXZWIiLCJyZWZlcnJhbE1vZGUiOiJ2aWV3In19&e=RJw6hN) demonstrates how to use this tool. 

## Prerequisites

1. [Install](https://github.com/Azure/oav#how-to-install-the-tool) the latest oav
2. [Install](https://github.com/ms-henglu/armstrong#install) the latest armstrong
3. Prepare the swagger definitions and examples that you need to test

## Step-By-Step

Please follow the steps below to complete API Test. We use [Microsoft.Purview/stable/2021-12-01/purview.json](https://github.com/Azure/azure-rest-api-specs/blob/main/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/purview.json) as an example.

### 1. Create a new folder

In this step, we will create a new folder to save the test code and test results. 


1. Move to the folder where stores the swagger json file and examples.

```shell
cd /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01
```

2. Run the following command to create a new empty folder named `terraform`.

```shell
mkdir terraform
```

### 2. Generate test cases

In this step, we will generate test cases based on the swagger json file and examples. Each test case will be generated in a separate folder categorized by the resource type.
The test case is written in [HCL](https://developer.hashicorp.com/terraform/language/syntax/configuration) which is the configuration language for Terraform. It describes the desired state of the resources that you want to manage. Armstrong will use the test case to create, update, and delete resources in Azure.

1. Move to the folder created in the previous step.

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01

cd terraform
```

2. Run the following command to generate test cases.

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform

armstrong generate -swagger /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/purview.json
```

You could use the relative path to specify the swagger json file. For example:

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform

armstrong generate -swagger ../purview.json
```

It also supports to specify the directory where the swagger json file resides in. For example:

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform

armstrong generate -swagger ..
```

More details about the usage of `armstrong generate` can be found [here](https://github.com/ms-henglu/armstrong#generate---generate-testing-files).

Then the test cases will be generated in the folder. The test case folder name is in the format of `{provider}_{resourceType}`. For example, the test case folder name for `Microsoft.Purview/accounts` is `Microsoft.Purview_accounts`.
The folder structure will look like below:

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01
├── examples
├── purview.json
└── terraform
    ├── Microsoft.Purview
    │   └── main.tf
    ├── Microsoft.Purview_accounts
    │   └── main.tf
    ├── Microsoft.Purview_accounts_kafkaConfigurations
    │   └── main.tf
    ├── Microsoft.Purview_accounts_privateEndpointConnections
    │   └── main.tf
    ├── Microsoft.Purview_accounts_privateLinkResources
    │   └── main.tf
    ├── Microsoft.Purview_locations
    │   └── main.tf
    └── Microsoft.Purview_operations
        └── main.tf
```

In each test case, `main.tf` file contains the terraform configuration about how to manage the resources. For example, the content of `main.tf` for `Microsoft.Purview/accounts` is like below:

```hcl
// provider definition: azapi provider will be used
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

// provider configuration
provider "azapi" {
  skip_provider_registration = false
}

// variable definition
variable "resource_name" {
  type    = string
  default = "acctest5906"
}

variable "location" {
  type    = string
  default = "West US 2"
}

// The purview account depends on the resource group, the armstrong will generate the resource group's definition automatically
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

// OperationId: Accounts_CreateOrUpdate, Accounts_Get, Accounts_Delete
// PUT GET DELETE /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}
resource "azapi_resource" "account" {
  type      = "Microsoft.Purview/accounts@2021-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  // the following payload is generated based on the content of /specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/examples/Accounts_CreateOrUpdate.json, which is defined in x-ms-examples of this operation
  body = jsonencode({
    properties = {
      managedResourceGroupName            = "custom-rgname"
      managedResourcesPublicNetworkAccess = "Enabled"
    }
  })
  schema_validation_enabled = false
}

// OperationId: Accounts_Update
// PATCH /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}
resource "azapi_resource_action" "patch_account" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = ""
  method      = "PATCH"
  body = jsonencode({
    cloudConnectors = {
    }
    managedResourceGroupName            = "aaaaaaaaaaaaaaaaaaaaaaaaaaaa"
    managedResourcesPublicNetworkAccess = "Disabled"
    publicNetworkAccess                 = "NotSpecified"
  })
}

...

// The operation definition
// OperationId: Accounts_ListKeys
// POST /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/listkeys
resource "azapi_resource_action" "listkeys" {
  type        = "Microsoft.Purview/accounts@2021-12-01"
  resource_id = azapi_resource.account.id
  action      = "listkeys"
  method      = "POST"
}

...

// The list operation definition
// OperationId: Accounts_ListByResourceGroup
// GET /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts
data "azapi_resource_list" "listAccountsByResourceGroup" {
  type       = "Microsoft.Purview/accounts@2021-12-01"
  parent_id  = azapi_resource.resourceGroup.id

  // The following config specifies the implicit dependency for this resource, it will only be executed after the resource "azapi_resource.account" is created
  depends_on = [azapi_resource.account]
}
```

More details about the `azapi` provider can be found [here](https://registry.terraform.io/providers/Azure/azapi/latest/docs).

### 3. Test each test case

In this step, we will test each test case generated in the previous step.

1. Move to one of the test case folders.

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform

cd Microsoft.Purview_accounts
```

2. Authenticate with Azure.

The easiest way to authenticate with Azure is to use [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli). Run the following command to login with Azure CLI.

```shell
az login
```

Armstrong supports a number of different methods for authenticating with Azure. More details can be found [here](https://registry.terraform.io/providers/Azure/azapi/latest/docs#authenticating-to-azure).

3. View the test plan. (Optional)

The following command generates a speculative execution plan, showing what actions Terraform would take to apply the current configuration.

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform/Microsoft.Purview_accounts

armstrong validate
```

More details about the usage of `armstrong validate` can be found [here](https://github.com/ms-henglu/armstrong#validate---validate-the-changes).

4. Test the test case.

The following command will test the test case and record the API traffic during the process. Then it will compare the API traffic with the examples to validate the correctness of the swagger definitions.

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform/Microsoft.Purview_accounts

armstrong test -swagger ../.. --destroy-after-test
```

Above command uses a relative path to the directory where the swagger json file resides in.

More details about the usage of `armstrong test` can be found [here](https://github.com/ms-henglu/armstrong#test---run-tests).

5. Check the output.

The test results will be generated in the folder `armstrong_reports_{month}_{day}_{random_number}`.

If there's no error reported, the test case passes. The folder structure will look like below:

```shell
├── armstrong_reports_Nov__2_122604
│   ├── SwaggerAccuracyReport.html          <- The swagger accuracy report 
│   ├── SwaggerAccuracyReport.json          <- For Armstrong internal use only - The swagger accuracy report in JSON format
│   ├── all_passed_report.md                <- Armstrong report, it contains the coverage report of the tested fields
│   └── traces                              <- For debugging purpose - The API traffic during the test
        ├── trace-1.json ... trace-10.json
├── log.txt                                 <- For debugging purpose - The log of the test
├── main.tf
├── terraform.tfstate                       <- For debugging purpose - The state file of the test
├── tfplan                                  <- For debugging purpose - The plan file of the test
└── traces                                  <- For debugging purpose - The API traffic during the test
    ├── trace-1.json ... trace-10.json
```

Please check the swagger accurcy report `SwaggerAccuracyReport.html` and fix the errors if there's any. 

If there's an error reported, the test case fails. You could find an error report in the folder `armstrong_reports_{month}_{day}_{random_number}`. The folder structure will look like below: (unrelated files are omitted)

```shell
├── armstrong_reports_Nov__2_122604
│   ├── Microsoft.Purview_accounts@2021-12-01_account.md    <- Armstrong error report
├── main.tf
``````

Please check the error report and fix the errors. The error report contains the error message and http request/response details to help you debug the error.

About how to fix the errors, please refer to [Fix Errors](#fix-errors).

After the errors are fixed, please go to step 4 to test the test case again.


### 4. Generate summary report

In this step, we will generate a summary report for all test cases. 

1. Move to the folder where stores the test cases.

```shell
cd /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform
```

2. Run the following command to generate a summary report.

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform

armstrong report -swagger ..
```

Above command uses a relative path to the directory where the swagger json file resides in.

More details about the usage of `armstrong report` can be found [here](https://github.com/ms-henglu/armstrong#report---generate-a-summary-report).

If it runs successfully, the summary reports will be generated in `ArmstrongReport` folder. The folder structure will look like below:

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01/terraform

terraform
├── ArmstrongReport
│   ├── SwaggerAccuracyReport.html      <- The swagger accuracy report for all test cases
│   ├── SwaggerAccuracyReport.json      <- For Armstrong internal use only - The swagger accuracy report in JSON format
│   ├── SwaggerAccuracyReport.md        <- The swagger accuracy report for all test cases, markdown format
│   └── traces                          <- For debugging purpose - The API traffic during the test
│       ├── Microsoft.Purview-trace-1.json
│       ├── ...
│       └── Microsoft.Purview_operations-trace-10.json
├── Microsoft.Purview
├── Microsoft.Purview_accounts
├── Microsoft.Purview_accounts_kafkaConfigurations
├── Microsoft.Purview_accounts_privateEndpointConnections
├── Microsoft.Purview_accounts_privateLinkResources
├── Microsoft.Purview_locations
└── Microsoft.Purview_operations
```

Please check the swagger accuracy report `SwaggerAccuracyReport.html` and fix the errors if there's any and make sure there's no `Untested Operations`.


### 5. Submit results with swagger pull request

In this step, we will submit the test cases and summary report with swagger pull request.

1. For the test cases, please check in the test cases. (Some unnecessary files will be automatically filtered by the .gitignore).

The folder structure will look like below:

```shell
# working directory: /Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification/purview/resource-manager/Microsoft.Purview/stable/2021-12-01
├── examples
├── purview.json
└── terraform
    ├── Microsoft.Purview
    │   └── main.tf
    ├── Microsoft.Purview_accounts
    │   └── main.tf
    ├── Microsoft.Purview_accounts_kafkaConfigurations
    │   └── main.tf
    ├── Microsoft.Purview_accounts_privateEndpointConnections
    │   └── main.tf
    ├── Microsoft.Purview_accounts_privateLinkResources
    │   └── main.tf
    ├── Microsoft.Purview_locations
    │   └── main.tf
    └── Microsoft.Purview_operations
        └── main.tf
```

2. For the summary report, please copy all content in the generated `ArmstrongReport/SwaggerAccuracyReport.md` and paste it in a new comment of the swagger pull request for ARM reviewers to review.


## Fix Errors

In this section, we will introduce how to fix the errors reported by Armstrong.

Case 1. Bad Request. The error message is like below:

```shell
ERRO[0014] error running terraform apply: exit status 1

Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: 1001
--------------------------------------------------------------------------------
{
  "error": {
    "code": "1001",
    "message": "IdentityUrl is required",
    "target": "IdentityUrl",
    "details": []
  }
}
--------------------------------------------------------------------------------


  with azapi_resource.account,
  on main.tf line 31, in resource "azapi_resource" "account":
  31: resource "azapi_resource" "account" {
```

This kind of error is caused by the incorrect payload in the test case. Please check the error message and fix the payload in the test case. 

Please refer to the following examples for how to fix the payload.

1. example: [Microsoft.Purview/accounts - 400 IdentityUrl is required](./fix-errors/Microsoft.Purview_accounts_identityUrl_is_required.md)
2. example: [Microsoft.Purview/accounts - 400 InvalidRequestContent](./fix-errors/Microsoft.Purview_accounts_InvalidRequestContent.md)

Please notice that the payload is generated based on the example defined in `x-ms-examples` of the operation. It's recommended to update the example in the swagger definition if the example is incorrect.

Case 2. Bad Request, because the referenced resource is not found. The error message is like below:

```shell
ERRO[0041] error running terraform apply: exit status 1

Error: performing action addRootCollectionAdmin of "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906\" / Api Version \"2021-12-01\")": POST https://management.azure.com/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906/addRootCollectionAdmin
--------------------------------------------------------------------------------
RESPONSE 400: 400 Bad Request
ERROR CODE: 1002
--------------------------------------------------------------------------------
{
  "error": {
    "code": "1002",
    "message": "The payload is invalid. Error: Failed to find ODataType for objectId:7e8de0e7-2bfc-4e1f-9659-2a5785e4356f tenantId:72f988bf-86f1-41af-91ab-2d7cd011db47..",
    "target": null,
    "details": []
  }
}
--------------------------------------------------------------------------------


  with azapi_resource_action.addRootCollectionAdmin,
  on main.tf line 69, in resource "azapi_resource_action" "addRootCollectionAdmin":
  69: resource "azapi_resource_action" "addRootCollectionAdmin" {
```

Please refer to the following example for how to add the referenced resource in the test case.

1. example: [Microsoft.Purview/accounts - addRootCollectionAdmin](./fix-errors/Microsoft.Purview_accounts_addRootCollectionAdmin.md)


case 3. Bad Request, because the dependency resource does not meet the requirements. The error message is like below:

```shell
ERRO[0056] error running terraform apply: exit status 1

Error: performing action listkeys of "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906\" / Api Version \"2021-12-01\")": POST https://management.azure.com/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906/listkeys
--------------------------------------------------------------------------------
RESPONSE 404: 404 Not Found
ERROR CODE: 21000
--------------------------------------------------------------------------------
{
  "error": {
    "code": "21000",
    "message": "Managed Event Hub does not exist for account acctest5906.",
    "target": null,
    "details": []
  }
}
--------------------------------------------------------------------------------


  with azapi_resource_action.listkeys,
  on main.tf line 106, in resource "azapi_resource_action" "listkeys":
 106: resource "azapi_resource_action" "listkeys" {
```

Please refer to the following example for how to update existing resources to meet the requirements.

1. example: [Microsoft.Purview/accounts - listKeys](./fix-errors/Microsoft.Purview_accounts_listKeys.md)
2. example: [Microsoft.Purview/accounts/kafkaConfigurations - Require Account Disable Managed Event Hub](./fix-errors/Microsoft.Purview_accounts_kafkaConfigurations_requireAccountDisableManagedEventHub.md)

case 4. Conflict, because the resource couldn't be updated. The error message is like below:

```shell
ERRO[0015] error running terraform apply: exit status 1

Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest5906/providers/Microsoft.Purview/accounts/acctest5906
--------------------------------------------------------------------------------
RESPONSE 409: 409 Conflict
ERROR CODE: 2008
--------------------------------------------------------------------------------
{
  "error": {
    "code": "2008",
    "message": "The managed event hub namespace cannot be re-enabled for account acctest5906 once it has been disabled.",
    "target": null,
    "details": []
  }
}
--------------------------------------------------------------------------------


  with azapi_resource.account,
  on main.tf line 31, in resource "azapi_resource" "account":
  31: resource "azapi_resource" "account" {
```

Please refer to the following example for how to fix the error.

1. example: [Microsoft.Purview/accounts - managedEventHubState Conflict](./fix-errors/Microsoft.Purview_accounts_managedEventHubStateConflict.md)

case 5. Internal Server Error. The error message is like below:

```shell
ERRO[0037] error running terraform apply: exit status 1

Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248
--------------------------------------------------------------------------------
RESPONSE 500: 500 Internal Server Error
ERROR CODE: 500
--------------------------------------------------------------------------------
{
  "error": {
    "code": "500",
    "message": "Unknown error",
    "target": null,
    "details": null
  }
}
--------------------------------------------------------------------------------


  with azapi_resource.kafkaConfiguration,
  on main.tf line 99, in resource "azapi_resource" "kafkaConfiguration":
  99: resource "azapi_resource" "kafkaConfiguration" {
```

Some of the internal server errors are caused by the incorrect payload/depedency requirements in the test case.

Please refer to the following examples for how to fix the errors.

1. example: [Microsoft.Purview/accounts/kafkaConfigurations - identity Internal Error](./fix-errors/Microsoft.Purview_accounts_kafkaConfigurations_identityInternalError.md)

case 6. Does not have permission. The error message is like below:

```shell
 Error: creating/updating "Resource: (ResourceId \"/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248\" / Api Version \"2021-12-01\")": PUT https://management.azure.com/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.Purview/accounts/acctest2248/kafkaConfigurations/acctest2248
│ --------------------------------------------------------------------------------
│ RESPONSE 409: 409 Conflict
│ ERROR CODE: 19000
│ --------------------------------------------------------------------------------
│ {
│   "error": {
│     "code": "19000",
│     "message": "The caller does not have EventHubDataSender permission on resource: [/subscriptions/******/resourceGroups/acctest2248/providers/Microsoft.EventHub/namespaces/acctest2248/eventhubs/acctest2248].",
│     "target": null,
│     "details": []
│   }
│ }
│ --------------------------------------------------------------------------------
│ 
│ 
│   with azapi_resource.kafkaConfiguration,
│   on main.tf line 102, in resource "azapi_resource" "kafkaConfiguration":
│  102: resource "azapi_resource" "kafkaConfiguration" {
```

This kind of error is caused by lack of permission. Please refer to the following example for how to add role assignment and fix the error.

1. example: [Microsoft.Purview/accounts/kafkaConfigurations - 409 No Permission](./fix-errors/Microsoft.Purview_accounts_kafkaConfigurations_roleAssignment.md)


## Frequently Asked Questions

1. Q: In each test case, how to find the untested operations?

   A: We're working on improving it. Currently, you could compare the OperationIds in the test case with the swagger accuracy report, and see if there're any untested operations.
   
   Below example, the block tests `KafkaConfigurations_CreateOrUpdate`, `KafkaConfigurations_Get`, `KafkaConfigurations_Delete` operations.
```hcl
// OperationId: KafkaConfigurations_CreateOrUpdate, KafkaConfigurations_Get, KafkaConfigurations_Delete
// PUT GET DELETE /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/kafkaConfigurations/{kafkaConfigurationName}
resource "azapi_resource" "kafkaConfiguration" {
  ...
```

2. Q: After I generated the summary report, there're some untested operations, how to find which test case that these operations belong to?

   A: In each testcase, the operationId is added in the comment, like the following example:
   the block tests `KafkaConfigurations_CreateOrUpdate`, `KafkaConfigurations_Get`, `KafkaConfigurations_Delete` operations.

```hcl
// OperationId: KafkaConfigurations_CreateOrUpdate, KafkaConfigurations_Get, KafkaConfigurations_Delete
// PUT GET DELETE /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Purview/accounts/{accountName}/kafkaConfigurations/{kafkaConfigurationName}
resource "azapi_resource" "kafkaConfiguration" {
  ...
```

3. Q: I've fixed the API issue, how to reset the test case and test it again?

   A: You could delete the `traces` folder under the test case folder.


## Samples

Please refer to the [sample repo](https://github.com/ms-henglu/armstrong-demo) for more details.