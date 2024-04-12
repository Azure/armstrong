package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/azure/armstrong/coverage"
	"github.com/azure/armstrong/hcl"
	"github.com/sirupsen/logrus"
)

type CredentialScanCommand struct {
	workingDir      string
	swaggerRepoPath string
	verbose         bool
}

func (c *CredentialScanCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("test")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.StringVar(&c.swaggerRepoPath, "swagger-repo", "", "path to the swagger repo specification directory")
	fs.Usage = func() { logrus.Error(c.Help()) }
	return fs
}

func (c CredentialScanCommand) Help() string {
	helpText := `
Usage: armstrong credscan [-v] [-working-dir <path to directory containing Terraform configuration files>] [-swagger-repo <path to the swagger repo specification directory>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c CredentialScanCommand) Synopsis() string {
	return "Scan the credential in given Terraform configuration"
}

func (c CredentialScanCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		logrus.Errorf("Error parsing command-line flags: %s", err)
		return 1
	}
	if c.verbose {
		log.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Infof("verbose mode enabled")
	}
	return c.Execute()
}

func (c CredentialScanCommand) Execute() int {
	wd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("failed to get working directory: %+v", err)
		return 1
	}
	if c.workingDir != "" {
		wd, err = filepath.Abs(c.workingDir)
		if err != nil {
			logrus.Errorf("working directory is invalid: %+v", err)
			return 1
		}
	}
	if c.swaggerRepoPath != "" {
		c.swaggerRepoPath, err = filepath.Abs(c.swaggerRepoPath)
		if err != nil {
			logrus.Errorf("swagger repo path %q is invalid: %+v", c.swaggerRepoPath, err)
			return 1
		}

		if _, err := os.Stat(c.swaggerRepoPath); os.IsNotExist(err) {
			logrus.Errorf("swagger repo path %q is invalid: path does not exist", c.swaggerRepoPath)
			return 1
		}

		c.swaggerRepoPath = strings.TrimSuffix(c.swaggerRepoPath, "/")

		if !strings.HasSuffix(c.swaggerRepoPath, "specification") {
			logrus.Errorf("swagger repo path %q is invalid: must point to \"specification\", e.g., /home/projects/azure-rest-api-specs/specification", c.swaggerRepoPath)
			return 1
		}

		c.swaggerRepoPath += "/"
	}

	tfFiles, err := hcl.FindTfFiles(wd)
	if err != nil {
		logrus.Errorf("failed to find tf files for %q: %+v", wd, err)
		return 1
	}
	if len(*tfFiles) == 0 {
		logrus.Warnf("no tf file found in %q", wd)
	}
	logrus.Infof("find %v tf file(s) under %s", len(*tfFiles), wd)

	azapiResources := make([]hcl.AzapiResource, 0)
	vars := make(map[string]hcl.Variable, 0)
	azureProviders := make([]hcl.AzureProvider, 0)
	for _, tfFile := range *tfFiles {
		f, err := hcl.ParseHclFile(tfFile)
		if err != nil {
			logrus.Errorf("failed to parse hcl file %q: %+v", tfFile, err)
			return 1
		}

		azapiResourceInFile, err := hcl.ParseAzapiResource(*f)
		if err != nil {
			logrus.Errorf("failed to parse azapi resource for %q: %+v", tfFile, err)
			return 1
		}
		azapiResources = append(azapiResources, *azapiResourceInFile...)

		varsInFile, err := hcl.ParseVariables(*f)
		if err != nil {
			logrus.Errorf("failed to parse variables for %q: %+v", tfFile, err)
			return 1
		}

		for k, v := range *varsInFile {
			vars[k] = v
		}

		azureProvidersInFile, err := hcl.ParseAzureProvider(*f)
		if err != nil {
			logrus.Errorf("failed to parse azure provider for %q: %+v", tfFile, err)
			return 1
		}
		azureProviders = append(azureProviders, *azureProvidersInFile...)

	}

	credScanErrors := make([]CredScanError, 0)

	for _, azureProvider := range azureProviders {
		if v := azureProvider.SubscriptionId; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "subscription_id", v, vars)...)
		}

		if v := azureProvider.TenantId; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "tenant_id", v, vars)...)
		}

		if v := azureProvider.AuxiliaryTenantIds; len(v) > 0 {
			for i, tenant_id := range v {
				credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, fmt.Sprintf("auxiliary_tenant_ids[%v]", i), tenant_id, vars)...)
			}
		}

		if v := azureProvider.AuxiliaryTenantIdsString; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "auxiliary_tenant_ids", v, vars)...)
		}

		if v := azureProvider.ClientId; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "client_id", v, vars)...)
		}

		if v := azureProvider.ClientCertificate; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "client_certificate", v, vars)...)
		}

		if v := azureProvider.ClientCertificatePassword; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "client_certificate_password", v, vars)...)
		}

		if v := azureProvider.ClientSecret; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "client_secret", v, vars)...)
		}

		if v := azureProvider.OidcRequestToken; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "oidc_request_token", v, vars)...)
		}

		if v := azureProvider.OidcToken; v != "" {
			credScanErrors = append(credScanErrors, checkAzureProviderSecret(azureProvider, "oidc_token", v, vars)...)
		}

	}

	for _, azapiResource := range azapiResources {
		logrus.Infof("scaning azapi_resource.%s(%s)", azapiResource.Name, azapiResource.Type)

		if azapiResource.Body == "" {
			continue
		}
		var body interface{}
		err = json.Unmarshal([]byte(azapiResource.Body), &body)
		if err != nil {
			credScanErr := makeCredScanError(
				azapiResource,
				fmt.Sprintf("failed to unmarshal body: %+v", err),
				"",
			)
			credScanErrors = append(credScanErrors, credScanErr)
			logrus.Error(credScanErr)
			continue
		}

		mockedResourceId, apiVersion := coverage.MockResourceIDFromType(azapiResource.Type)
		logrus.Infof("azapi_resource.%s(%s): mocked possible resource ID: %s, API version: %s", azapiResource.Name, azapiResource.Type, mockedResourceId, apiVersion)

		var swaggerModel *coverage.SwaggerModel
		if c.swaggerRepoPath != "" {
			logrus.Infof("scan based on local swagger repo: %s", c.swaggerRepoPath)
			swaggerModel, err = coverage.GetModelInfoFromLocalIndex(mockedResourceId, apiVersion, c.swaggerRepoPath)
			if err != nil {
				credScanErr := makeCredScanError(
					azapiResource,
					fmt.Sprintf("fail to find swagger model from local swagger with possible resource ID(%s) API version(%s): %+v", mockedResourceId, apiVersion, err),
					"",
				)
				credScanErrors = append(credScanErrors, credScanErr)
				logrus.Error(credScanErr)

				continue
			}
		} else {
			swaggerModel, err = coverage.GetModelInfoFromIndex(mockedResourceId, apiVersion)
			if err != nil {
				credScanErr := makeCredScanError(
					azapiResource,
					fmt.Sprintf("fail to find swagger model with possible resource ID(%s) API version(%s): %+v", mockedResourceId, apiVersion, err),
					"",
				)
				credScanErrors = append(credScanErrors, credScanErr)
				logrus.Error(credScanErr)

				continue
			}
		}

		if swaggerModel == nil {
			credScanErr := makeCredScanError(
				azapiResource,
				fmt.Sprintf("unable to find swagger model with possible resource ID(%s) API version(%s)", mockedResourceId, apiVersion),
				"",
			)
			credScanErrors = append(credScanErrors, credScanErr)
			logrus.Error(credScanErr)

			continue
		}

		logrus.Infof("find swagger model for azapi_resource.%s(%s): %+v", azapiResource.Name, azapiResource.Type, *swaggerModel)

		model, err := coverage.Expand(swaggerModel.ModelName, swaggerModel.SwaggerPath)
		if err != nil {
			credScanErr := makeCredScanError(
				azapiResource,
				fmt.Sprintf("failed to expand model: %+v", err),
				"",
			)
			credScanErrors = append(credScanErrors, credScanErr)
			logrus.Error(credScanErr)

			continue
		}

		secrets := make(map[string]string)
		model.CredScan(body, secrets)

		logrus.Infof("find secrets for azapi_resource.%s(%s): %+v", azapiResource.Name, azapiResource.Type, secrets)

		for k, v := range secrets {
			if !strings.HasPrefix(v, "$") || strings.HasPrefix(v, "$local.") {
				credScanErr := makeCredScanError(
					azapiResource,
					"cannot use plain text or 'local' for secret, use 'variable' instead",
					k,
				)
				credScanErrors = append(credScanErrors, credScanErr)
				logrus.Error(credScanErr)

				continue
			}

			if strings.HasPrefix(v, "$var.") {
				varName := strings.TrimPrefix(v, "$var.")
				varName = strings.Split(varName, ".")[0]
				theVar, ok := vars[varName]
				if !ok {
					credScanErr := makeCredScanError(
						azapiResource,
						fmt.Sprintf("variable %q was not found", varName),
						k,
					)
					credScanErrors = append(credScanErrors, credScanErr)
					logrus.Error(credScanErr)

					continue
				}

				if theVar.HasDefault {
					credScanErr := makeCredScanError(
						azapiResource,
						fmt.Sprintf("variable %q (%v:%v) used in secret field but has a default value, please remove the default value", varName, theVar.FileName, theVar.LineNumber),
						k,
					)
					credScanErrors = append(credScanErrors, credScanErr)
					logrus.Error(credScanErr)
				}

				if !theVar.IsSensitive {
					credScanErr := makeCredScanError(
						azapiResource,
						fmt.Sprintf("variable %q (%v:%v) used in secret field but is not marked as sensitive, please add \"sensitive=true\" for the variable", varName, theVar.FileName, theVar.LineNumber),
						k,
					)
					credScanErrors = append(credScanErrors, credScanErr)
					logrus.Error(credScanErr)
				}
			}
		}
	}

	storeCredScanErrors(wd, credScanErrors)

	return 0
}

type CredScanError struct {
	FileName     string `json:"file_name"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	PropertyName string `json:"property_name"`
	ErrorMessage string `json:"error_message"`
	LineNumber   int    `json:"line_number"`
}

func makeCredScanError(azapiResource hcl.AzapiResource, errMessage, propertyName string) CredScanError {
	result := CredScanError{
		FileName:     azapiResource.FileName,
		LineNumber:   azapiResource.LineNumber,
		Name:         fmt.Sprintf("azapi_resource.%s", azapiResource.Name),
		Type:         azapiResource.Type,
		ErrorMessage: errMessage,
	}

	if propertyName != "" {
		result.PropertyName = propertyName
	}

	return result
}

func makeCredScanErrorForProvider(azureProvider hcl.AzureProvider, errMessage, propertyName string) CredScanError {
	result := CredScanError{
		FileName:     azureProvider.FileName,
		LineNumber:   azureProvider.LineNumber,
		Name:         azureProvider.Name(),
		Type:         "provider",
		ErrorMessage: errMessage,
	}

	if propertyName != "" {
		result.PropertyName = propertyName
	}

	return result
}

func (e CredScanError) Error() string {
	return fmt.Sprintf("%s:%d %s(%s) --%s: %s", e.FileName, e.LineNumber, e.Name, e.Type, e.PropertyName, e.ErrorMessage)
}

func storeCredScanErrors(wd string, credScanErrors []CredScanError) {
	reportDir := fmt.Sprintf("armstrong_credscan_%s", time.Now().Format(time.Stamp))
	reportDir = strings.ReplaceAll(reportDir, ":", "")
	reportDir = strings.ReplaceAll(reportDir, " ", "_")
	reportDir = path.Join(wd, reportDir)

	err := os.Mkdir(reportDir, 0755)
	if err != nil {
		logrus.Fatalf("error creating report dir %s: %+v", reportDir, err)
	}

	markdownFileName := "errors.md"
	credScanErrorsMarkdown := `
| File Name | Line Number | Name | Type | Property Name | Error Message |
| --- | --- | --- | --- | --- | --- |
`
	for _, r := range credScanErrors {
		credScanErrorsMarkdown += fmt.Sprintf("| %s | %d | %s | %s | %s | %s |\n", r.FileName, r.LineNumber, r.Name, r.Type, r.PropertyName, r.ErrorMessage)
	}

	err = os.WriteFile(path.Join(reportDir, markdownFileName), []byte(credScanErrorsMarkdown), 0644)
	if err != nil {
		logrus.Errorf("failed to save markdown report to %s: %+v", markdownFileName, err)
	} else {
		logrus.Infof("markdown report saved to %s", markdownFileName)
	}

	jsonFileName := "errors.json"
	jsonContent, err := json.MarshalIndent(credScanErrors, "", "  ")
	if err != nil {
		logrus.Errorf("failed to marshal json content %+v: %+v", credScanErrors, err)
	}

	err = os.WriteFile(path.Join(reportDir, jsonFileName), jsonContent, 0644)
	if err != nil {
		logrus.Errorf("failed to save json report to %s: %+v", jsonFileName, err)
	} else {
		logrus.Infof("json report saved to %s", jsonFileName)
	}
}

func checkAzureProviderSecret(azureProvider hcl.AzureProvider, propertyName, propertyValue string, vars map[string]hcl.Variable) []CredScanError {
	credScanErrors := make([]CredScanError, 0)

	if !strings.HasPrefix(propertyValue, "$") || strings.HasPrefix(propertyValue, "$local.") {
		credScanErr := makeCredScanErrorForProvider(
			azureProvider,
			"cannot use plain text or 'local' for secret, use 'variable' instead",
			propertyName,
		)
		credScanErrors = append(credScanErrors, credScanErr)
		logrus.Error(credScanErr)

		return credScanErrors
	}

	if strings.HasPrefix(propertyValue, "$var.") {
		varName := strings.TrimPrefix(propertyValue, "$var.")
		varName = strings.Split(varName, ".")[0]
		theVar, ok := vars[varName]
		if !ok {
			credScanErr := makeCredScanErrorForProvider(
				azureProvider,
				fmt.Sprintf("variable %q was not found", varName),
				propertyName,
			)
			credScanErrors = append(credScanErrors, credScanErr)
			logrus.Error(credScanErr)

			return credScanErrors
		}

		if theVar.HasDefault {
			credScanErr := makeCredScanErrorForProvider(
				azureProvider,
				fmt.Sprintf("variable %q (%v:%v) used in secret field but has a default value, please remove the default value", varName, theVar.FileName, theVar.LineNumber),
				propertyName,
			)
			credScanErrors = append(credScanErrors, credScanErr)
			logrus.Error(credScanErr)
		}

		if !theVar.IsSensitive {
			credScanErr := makeCredScanErrorForProvider(
				azureProvider,
				fmt.Sprintf("variable %q (%v:%v) used in secret field but is not marked as sensitive, please add \"sensitive=true\" for the variable", varName, theVar.FileName, theVar.LineNumber),
				propertyName,
			)
			credScanErrors = append(credScanErrors, credScanErr)
			logrus.Error(credScanErr)
		}
	}

	return credScanErrors
}
