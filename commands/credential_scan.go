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

	"github.com/ms-henglu/armstrong/coverage"
	"github.com/ms-henglu/armstrong/hcl"
	"github.com/sirupsen/logrus"
)

type CredentialScanCommand struct {
	workingDir  string
	swaggerPath string
	verbose     bool
}

func (c *CredentialScanCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("test")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.StringVar(&c.swaggerPath, "swagger", "", "path to the swagger repo specification directory")
	fs.Usage = func() { logrus.Error(c.Help()) }
	return fs
}

func (c CredentialScanCommand) Help() string {
	helpText := `
Usage: armstrong credscan [-v] [-working-dir <path to directory containing Terraform configuration files>] [-swagger <path to the swagger repo specification directory>]
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
	if c.swaggerPath != "" {
		c.swaggerPath, err = filepath.Abs(c.swaggerPath)
		if err != nil {
			logrus.Errorf("swagger path %q is invalid: %+v", c.swaggerPath, err)
			return 1
		}

		if _, err := os.Stat(c.swaggerPath); os.IsNotExist(err) {
			logrus.Errorf("swagger path %q is invalid: path does not exist", c.swaggerPath)
			return 1
		}

		if !strings.HasSuffix(c.swaggerPath, "specification") {
			logrus.Errorf("swagger path %q is invalid: must point to \"specification\", e.g., /home/projects/azure-rest-api-specs/specification", c.swaggerPath)
			return 1
		}
	}

	tfFiles, err := hcl.FindTfFiles(wd)
	if err != nil {
		logrus.Errorf("failed to find tf files for %q: %+v", wd, err)
		return 1
	}
	if len(*tfFiles) == 0 {
		logrus.Warnf("no .tf file found in %q", wd)
	}
	logrus.Infof("find %v .tf files under %s", len(*tfFiles), wd)

	azapiResources := make([]hcl.AzapiResource, 0)
	vars := make(map[string]hcl.Variable, 0)
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
	}

	credScanErrors := make([]CredScanError, 0)

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
		if c.swaggerPath != "" {
			logrus.Infof("scan based on local swagger file: %s", c.swaggerPath)
			swaggerModel, err = coverage.GetModelInfoFromLocalIndex(mockedResourceId, apiVersion, c.swaggerPath)
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
			if !strings.HasPrefix(v, "$var.") {
				credScanErr := makeCredScanError(
					azapiResource,
					"must use variable for secret field",
					k,
				)
				credScanErrors = append(credScanErrors, credScanErr)
				logrus.Error(credScanErr)

				continue
			}

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

			if theVar.Default != "" {
				credScanErr := makeCredScanError(
					azapiResource,
					fmt.Sprintf("variable %q used in secret field but has a default value, please remove the default value", varName),
					k,
				)
				credScanErrors = append(credScanErrors, credScanErr)
				logrus.Error(credScanErr)
			}

			if !theVar.IsSensitive {
				credScanErr := makeCredScanError(
					azapiResource,
					fmt.Sprintf("variable %q used in secret field but is not marked as sensitive, please add \"sensitive: true\" for the variable", varName),
					k,
				)
				credScanErrors = append(credScanErrors, credScanErr)
				logrus.Error(credScanErr)
			}
		}
	}

	storeCredScanErrors(wd, credScanErrors)

	return 0
}

type CredScanError struct {
	FileName     string `json:"file_name"`
	ResourceName string `json:"resource_name"`
	ResourceType string `json:"resource_type"`
	PropertyName string `json:"property_name"`
	ErrorMessage string `json:"error_message"`
	LineNumber   int    `json:"line_number"`
}

func makeCredScanError(azapiResource hcl.AzapiResource, errMessage string, PropertyName string) CredScanError {
	result := CredScanError{
		FileName:     azapiResource.FileName,
		LineNumber:   azapiResource.LineNumber,
		ResourceName: fmt.Sprintf("azapi_resource.%s", azapiResource.Name),
		ResourceType: azapiResource.Type,
		ErrorMessage: errMessage,
	}

	if PropertyName != "" {
		result.PropertyName = PropertyName
	}

	return result
}

func (e CredScanError) Error() string {
	return fmt.Sprintf("%s:%d %s(%s) --%s: %s", e.FileName, e.LineNumber, e.ResourceName, e.ResourceType, e.PropertyName, e.ErrorMessage)
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
| File Name | Line Number | Resource Name | Resource Type | Property Name | Error Message |
| --- | --- | --- | --- | --- | --- |
`
	for _, r := range credScanErrors {
		credScanErrorsMarkdown += fmt.Sprintf("| %s | %d | %s | %s | %s | %s |\n", r.FileName, r.LineNumber, r.ResourceName, r.ResourceType, r.PropertyName, r.ErrorMessage)
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
