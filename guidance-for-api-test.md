# Guidance for API Test

## Prerequisites

1. [Install](https://github.com/Azure/oav#how-to-install-the-tool) the latest oav
2. [Install](https://github.com/ms-henglu/armstrong#install) the latest armstrong
3. Prepare the swagger definitions and examples that you need to test

## Step-By-Step

Please follow the steps below to complete API Test. We use [2022-08-08/account.json](https://github.com/Azure/azure-rest-api-specs/blob/main/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/account.json) as an example.

### 1. Create a new Folder

Run the following command to create a new empty folder to save your test results.

```
mkdir {test-dir}
```

### 2. Generate Test Code

Run [armstrong generate](https://github.com/ms-henglu/armstrong#generate---generate-testing-files) to generate the test code. Here is an example:

```
cd {test-dir}

armstrong generate -swagger {swagger-repo}/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/account.json
```

Then the test code will be generated in `{test-dir}/Microsoft.Automation_automationAccounts/main.tf`. The API body definitions are based on the content of `x-ms-examples` defined in [2022-08-08/account.json](https://github.com/Azure/azure-rest-api-specs/blob/main/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/account.json).

### 3. Update Test Code

Please physically update the content of the API body and the order of API operations if needed.

### 4. Validate Test Code

Run [armstrong validate](https://github.com/ms-henglu/armstrong#validate---validate-the-changes) to validate the test code and generates a speculative execution plan. Please fix any reported errors by physically updating the test code. Here is an example:

```
armstrong validate -working-dir ./Microsoft.Automation_automationAccounts
```

### 5. Test API operations

Run [armstrong test](https://github.com/ms-henglu/armstrong#test---run-tests) to test API operations. Here is an example:

```
armstrong test -working-dir ./Microsoft.Automation_automationAccounts -swagger {swagger-repo}/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/account.json -destroy-after-test true
```

### 6. Generate Summary Report

### 7. Submit Test Code and Summary Report with Swagger Pull Request