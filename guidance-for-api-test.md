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

Then the test code will be generated in `{test-dir}/Microsoft.Automation_automationAccounts/main.tf`. The property values of API body are based on the content of `x-ms-examples` defined in [2022-08-08/account.json](https://github.com/Azure/azure-rest-api-specs/blob/main/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/account.json). The the folder structure will look like below:

```
├── {test-dir}
│   ├── Microsoft.Automation_automationAccounts
│   │   ├── main.tf
```

### 3. Update Test Code

Please physically update the content of the API body and the order of API operations in `main.tf` if needed.

### 4. Validate Test Code

Run [armstrong validate](https://github.com/ms-henglu/armstrong#validate---validate-the-changes) to validate the test code and generates a speculative execution plan. Here is an example:

```
armstrong validate -working-dir ./Microsoft.Automation_automationAccounts
```

If it runs successfully, you will get the log like below:

```
Plan: xx to add, 0 to change, 0 to destroy.
```

If it reports errors, please fix them by physically updating the test code and run this command again. 

### 5. Test API operations

Run [armstrong test](https://github.com/ms-henglu/armstrong#test---run-tests) to test API operations. Here is an example:

```
armstrong test -working-dir ./Microsoft.Automation_automationAccounts -swagger {swagger-repo}/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/account.json -destroy-after-test true
```

If it runs successfully, validated results will be generated in `SwaggerAccuracyReport.html`. Please open `SwaggerAccuracyReport.html` to fix errors in `Failed Operations`.

```
├── {test-dir}
│   ├── Microsoft.Automation_automationAccounts
│   │   ├── main.tf
│   │   ├── armstrong_reports_{month}_{day}_{random_number}
│   │   │   ├── traces
│   │   │   ├── SwaggerAccuracyReport.html
```

If it reports errors, please fix them by physically updating the test code and run this command again. 

### 6. Generate Summary Report

Run [armstrong report](https://github.com/ms-henglu/armstrong#report---generate-a-summary-report) to generate a summary report. Here is an example:

```
armstrong report -swagger {swagger-repo}/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/account.json
```

If it runs successfully, the summary reports will be generated in `{test-dir}/ArmstrongReport/SwaggerAccuracyReport.html` and `{test-dir}/ArmstrongReport/SwaggerAccuracyReport.md`. Please open `SwaggerAccuracyReport.html` to make sure there is no `Failed Operations` and `Untested Operations`.

```
├── {test-dir}
│   ├── Microsoft.Automation_automationAccounts
│   ├── ArmstrongReport
│   │   ├── SwaggerAccuracyReport.html
│   │   ├── SwaggerAccuracyReport.md
```

### 7. Submit Test Code and Summary Report with Swagger Pull Request

Create a new folder named `Armstrong` if it does not exist in the folder where the swagger json file resides in. Here is an example:

```
cd {swagger-repo}/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08
mkdir Armstrong
```

Copy all content under {test-dir} into the new created `Armstrong` folder(for example: {swagger-repo}/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08/Armstrong), the folder structure will look like below:

```
├── {swagger-repo}/specification/automation/resource-manager/Microsoft.Automation/stable/2022-08-08
│   ├── account.json
│   ├── examples
│   ├── Armstrong
│   │   ├── ArmstrongReport
│   │   ├── Microsoft.Automation_automationAccounts
```


Then submit these files with Swagger Pull Request.(Some unnecessary files will be automatically filtered by the .gitignore).

Copy all content in the generated `{test-dir}/ArmstrongReport/SwaggerAccuracyReport.md` and paste it in a new comment of the Swagger Pull Request for ARM reviewers to review.