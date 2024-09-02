package report

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/azure/armstrong/utils"
	"github.com/sirupsen/logrus"
)

const (
	TestReportDirName     = "ArmstrongReport"
	TraceLogDirName       = "traces"
	ApiTestReportFileName = "SwaggerAccuracyReport"
	ApiTestConfigFileName = "ApiTestConfig.json"
)

type ApiTestReport struct {
	CoveredSpecFiles        []string              `json:"coveredSpecFiles"`
	UnCoveredOperationsList []UnCoveredOperations `json:"unCoveredOperationsList"`
	Errors                  []ErrorItem           `json:"errors"`
}

type UnCoveredOperations struct {
	Spec         string   `json:"spec"`
	OperationIds []string `json:"operationIds"`
}

type ErrorItem struct {
	Spec                   string `json:"spec"`
	ErrorCode              string `json:"errorCode"`
	ErrorLink              string `json:"errorLink"`
	ErrorMessage           string `json:"errorMessage"`
	OperationId            string `json:"operationId"`
	SchemaPathWithPosition string `json:"schemaPathWithPosition"`
}

type ApiTestConfig struct {
	SuppressionList []Suppression `json:"suppressionList"`
}

type Suppression struct {
	Code      string `json:"rule"`
	File      string `json:"file"`
	Operation string `json:"operation"`
	Reason    string `json:"reason"`
}

func OavValidateTraffic(traceDir string, swaggerPath string, outputDir string) (*ApiTestReport, error) {
	htmlReportFilePath := path.Join(outputDir, fmt.Sprintf("%s.html", ApiTestReportFileName))
	jsonReportFilePath := path.Join(outputDir, fmt.Sprintf("%s.json", ApiTestReportFileName))

	logrus.Debugf("oav validate-traffic %s %s --report %s --jsonReport %s", traceDir, swaggerPath, htmlReportFilePath, jsonReportFilePath)
	cmd := exec.Command("oav", "validate-traffic", traceDir, swaggerPath, "--report", htmlReportFilePath, "--jsonReport", jsonReportFilePath)
	if err := cmd.Run(); err != nil {
		logrus.Warnf("oav validates-traffic: %+v", err)
	}

	contentBytes, err := os.ReadFile(jsonReportFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when opening file(%s): %+v", jsonReportFilePath, err)
	}

	var payload *ApiTestReport
	err = json.Unmarshal(contentBytes, &payload)
	if err != nil {
		return nil, fmt.Errorf("error during Unmarshal() for file(%s): %+v", jsonReportFilePath, err)
	}

	if payload == nil {
		return nil, fmt.Errorf("oav report is empty")
	}

	// remove duplicated error items
	errorMap := make(map[string]ErrorItem)
	for _, errItem := range payload.Errors {
		errorMap[fmt.Sprintf("%s-%s-%s-%s", errItem.ErrorCode, errItem.ErrorMessage, errItem.OperationId, errItem.SchemaPathWithPosition)] = errItem
	}

	errors := make([]ErrorItem, 0)
	for _, v := range errorMap {
		errors = append(errors, v)
	}

	payload.Errors = errors

	return payload, nil
}

func GenerateApiTestReports(wd string, swaggerPath string) error {
	testReportPath := path.Join(wd, TestReportDirName)
	traceLogPath := path.Join(testReportPath, TraceLogDirName)
	swaggerPath, _ = filepath.Abs(swaggerPath)

	logrus.Infof("copying trace files to %s...", traceLogPath)
	if err := mergeApiTestTraceFiles(wd, traceLogPath); err != nil {
		return fmt.Errorf("[ERROR] failed to merge trace files: %+v", err)
	}

	logrus.Infof("validating traces...")
	report, err := OavValidateTraffic(traceLogPath, swaggerPath, testReportPath)
	if err != nil {
		return fmt.Errorf("[ERROR] failed to retrieve oav report: %+v", err)
	}

	logrus.Infof("generating markdown report...")
	if err = generateApiTestMarkdownReport(*report, swaggerPath, testReportPath, path.Join(wd, ApiTestConfigFileName)); err != nil {
		return fmt.Errorf("[ERROR] failed to generate markdown report: %+v", err)
	}

	return nil
}

func mergeApiTestTraceFiles(wd string, traceLogPath string) error {
	if err := os.RemoveAll(traceLogPath); err != nil {
		return fmt.Errorf("error removing test trace dir %s: %+v", traceLogPath, err)
	}

	if err := os.MkdirAll(traceLogPath, 0755); err != nil {
		return fmt.Errorf("error creating test report dir %s: %+v", traceLogPath, err)
	}

	dirs, err := os.ReadDir(wd)
	if err != nil {
		return fmt.Errorf("failed to read working directory: %+v", err)
	}
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		traceDir := filepath.Join(wd, d.Name(), TraceLogDirName)
		if utils.Exists(traceDir) {
			err := utils.CopyWithOptions(traceDir, traceLogPath, fmt.Sprintf("%s-", d.Name()))
			if err != nil {
				return fmt.Errorf("failed to copy trace files: %+v", err)
			}
		}
	}

	return nil
}

func isSuppressedInApiTest(suppressionList []Suppression, rule string, filePath string, operation string) bool {
	segments := strings.Split(filepath.ToSlash(filePath), "/")
	file := segments[len(segments)-1]

	for _, suppression := range suppressionList {
		if strings.EqualFold(suppression.Code, rule) && strings.EqualFold(suppression.File, file) && (strings.EqualFold(suppression.Code, "SWAGGER_NOT_TEST") || strings.EqualFold(suppression.Operation, operation)) {
			return true
		}
	}

	return false
}

func generateApiTestMarkdownReport(result ApiTestReport, swaggerPath string, testReportPath string, apiTestConfigFilePath string) error {
	var config ApiTestConfig

	if utils.Exists(apiTestConfigFilePath) {
		contentBytes, err := os.ReadFile(apiTestConfigFilePath)
		if err != nil {
			logrus.Errorf("error when opening file(%s): %+v", apiTestConfigFilePath, err)
		}

		err = json.Unmarshal(contentBytes, &config)
		if err != nil {
			logrus.Errorf("error during Unmarshal() for file(%s): %+v", apiTestConfigFilePath, err)
		}
	} else {
		logrus.Debugf("no config file found")
	}

	mdTitle := "## API TEST ERROR REPORT<br>\n|Rule|Message|\n|---|---|"
	mdTable := make([]string, 0)

	testedMap := make(map[string]bool)
	for _, v := range result.CoveredSpecFiles {
		v = strings.ReplaceAll(v, "\\", "/")
		testedMap[v] = true
	}

	swaggerFiles, err := utils.ListFiles(swaggerPath, ".json", 1)
	if err != nil {
		return err
	}

	for _, v := range swaggerFiles {
		v = strings.ReplaceAll(v, "\\", "/")
		if _, exists := testedMap[v]; !exists {
			if !isSuppressedInApiTest(config.SuppressionList, "SWAGGER_NOT_TEST", v, "") {
				mdTable = append(mdTable, fmt.Sprintf("|[SWAGGER_NOT_TEST](about:blank)|**message**: No operations in swagger is test.<br>**location**: %s", v[strings.Index(v, "/specification/"):]))
			}
		}
	}

	for _, operationsItem := range result.UnCoveredOperationsList {
		for _, id := range operationsItem.OperationIds {
			if !isSuppressedInApiTest(config.SuppressionList, "OPERATION_NOT_TEST", operationsItem.Spec, id) {
				mdTable = append(mdTable, fmt.Sprintf("|[OPERATION_NOT_TEST](about:blank)|**message**: **%s** opeartion is not test.<br>**opeartion**: %s<br>**location**: %s", id, id, operationsItem.Spec[strings.Index(operationsItem.Spec, "/specification/"):]))
			}
		}
	}

	for _, errItem := range result.Errors {
		location := errItem.Spec[strings.Index(errItem.Spec, "/specification/"):]
		normalizedPath := filepath.ToSlash(errItem.SchemaPathWithPosition)
		if subIndex := strings.Index(normalizedPath, "/specification/"); subIndex != -1 {
			location = normalizedPath[subIndex:]
		}

		if !isSuppressedInApiTest(config.SuppressionList, errItem.ErrorCode, errItem.Spec, errItem.OperationId) {
			mdTable = append(mdTable, fmt.Sprintf("|[%s](%s)|**message**: %s.<br>**opeartion**: %s<br>**location**: %s", errItem.ErrorCode, errItem.ErrorLink, errItem.ErrorMessage, errItem.OperationId, location))
		}
	}

	sort.Strings(mdTable)

	mdReportFilePath := path.Join(testReportPath, fmt.Sprintf("%s.md", ApiTestReportFileName))
	if err := os.WriteFile(mdReportFilePath, []byte(mdTitle+"\n"+strings.Join(mdTable, "\n")), 0644); err != nil {
		return fmt.Errorf("error when writing file(%s): %+v", mdReportFilePath, err)
	}
	logrus.Infof("markdown report saved to %s", mdReportFilePath)

	return nil
}
