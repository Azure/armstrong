package report

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

const (
	TestReportDirName     = "ApiTestReport"
	TestResultDirName     = "ApiTestResult"
	TraceLogDirName       = "ApiTestTraces"
	ApiTestReportFileName = "ApiTestReport"
	ApiTestConfigFileName = "ApiTestConfig.json"
)

type ApiTestReport struct {
	ApiVersion      string           `json:"apiVersion"`
	CoverageResults []CoverageResult `json:"coverageResultsForRendering"`
}

type CoverageResult struct {
	GeneralErrorsInnerList  []GeneralErrorsInner `json:"generalErrorsInnerList"`
	UnCoveredOperationsList []UnCoveredOperation `json:"unCoveredOperationsList"`
}

type GeneralErrorsInner struct {
	ErrorsForRendering []ErrorForRendering `json:"errorsForRendering"`
	OperationInfo      OperationInfo       `json:"operationInfo"`
}

type OperationInfo struct {
	OperationId string `json:"operationId"`
}

type ErrorForRendering struct {
	Code                   string `json:"code"`
	Message                string `json:"message"`
	Link                   string `json:"link"`
	SchemaPathWithPosition string `json:"schemaPathWithPosition"`
}

type UnCoveredOperation struct {
	OperationId string `json:"operationId"`
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

func StoreApiTestReport(wd string, swaggerPath string) error {
	wd = filepath.ToSlash(wd)
	swaggerPath = filepath.ToSlash(swaggerPath)
	testResultPath := path.Join(wd, TestResultDirName)
	traceLogPath := path.Join(testResultPath, TraceLogDirName)
	if err := os.RemoveAll(testResultPath); err != nil {
		return fmt.Errorf("[ERROR] error removing trace log dir %s: %+v", testResultPath, err)
	}

	if err := os.MkdirAll(traceLogPath, 0755); err != nil {
		return fmt.Errorf("[ERROR] error creating trace log dir %s: %+v", testResultPath, err)
	}

	cmd := exec.Command("pal", "-i", path.Join(wd, "log.txt"), "-m", "oav", "-o", traceLogPath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run  `pal` command with %+v", err)
	}

	reportFileName := fmt.Sprintf("%s.html", ApiTestReportFileName)
	cmd = exec.Command("oav", "validate-traffic", traceLogPath, swaggerPath, "--report", path.Join(testResultPath, reportFileName))
	if err = cmd.Run(); err != nil {
		logrus.Infof("oav validates errors from traffic: %+v", err)
	}

	if _, err = os.Stat(path.Join(testResultPath, reportFileName)); os.IsNotExist(err) {
		return fmt.Errorf("failed to generate test report")
	}

	return nil
}

func GenerateApiTestReports(wd string, swaggerPath string) error {
	wd = filepath.ToSlash(wd)
	testReportPath := path.Join(wd, TestReportDirName)
	traceLogPath := path.Join(testReportPath, TraceLogDirName)
	apiTestReportFilePath := path.Join(testReportPath, fmt.Sprintf("%s.json", ApiTestReportFileName))

	if err := mergeApiTestTraceFiles(wd, traceLogPath); err != nil {
		return fmt.Errorf("[ERROR] failed to merge trace files: %+v", err)
	}

	swaggerFilePaths, err := getApiTestSwaggerPaths(swaggerPath)
	if err != nil {
		return fmt.Errorf("[ERROR] failed to get swagger paths: %+v", err)
	}

	result, err := readApiTestHistoryReport(swaggerPath, apiTestReportFilePath)
	if err != nil {
		return err
	}

	if err = generateApiTestJsonReport(result, swaggerFilePaths, testReportPath); err != nil {
		return fmt.Errorf("[ERROR] failed to generate oav reports: %+v", err)
	}

	if err = generateApiTestMarkdownReport(result, testReportPath, path.Join(wd, ApiTestConfigFileName)); err != nil {
		return fmt.Errorf("[ERROR] failed to generate markdown report: %+v", err)
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

func generateApiTestMarkdownReport(result map[string]*ApiTestReport, testReportPath string, apiTestConfigFilePath string) error {
	var config ApiTestConfig

	if _, err := os.Stat(apiTestConfigFilePath); err == nil {
		contentBytes, err := os.ReadFile(apiTestConfigFilePath)
		if err != nil {
			logrus.Errorf("error when opening file(%s): %+v", apiTestConfigFilePath, err)
		}

		err = json.Unmarshal(contentBytes, &config)
		if err != nil {
			logrus.Errorf("error during Unmarshal() for file(%s): %+v", apiTestConfigFilePath, err)
		}
	} else {
		logrus.Infof("no config file found")
	}

	mdTitle := "## API TEST ERROR REPORT<br>\n|Rule|Message|\n|---|---|"
	mdTable := make([]string, 0)

	for rk, rv := range result {
		if rv == nil {
			if !isSuppressedInApiTest(config.SuppressionList, "SWAGGER_NOT_TEST", rk, "") {
				mdTable = append(mdTable, fmt.Sprintf("|[SWAGGER_NOT_TEST](about:blank)|**message**: No operations in swagger is test.<br>**location**: %s", rk[strings.Index(rk, "/specification/"):]))
			}

			continue
		}

		for _, coverageResult := range rv.CoverageResults {
			for _, operation := range coverageResult.UnCoveredOperationsList {
				if !isSuppressedInApiTest(config.SuppressionList, "OPERATION_NOT_TEST", rk, operation.OperationId) {
					mdTable = append(mdTable, fmt.Sprintf("|[OPERATION_NOT_TEST](about:blank)|**message**: **%s** opeartion is not test.<br>**opeartion**: %s<br>**location**: %s", operation.OperationId, operation.OperationId, rk[strings.Index(rk, "/specification/"):]))
				}

			}

			for _, item := range coverageResult.GeneralErrorsInnerList {
				for _, errItem := range item.ErrorsForRendering {
					location := rk[strings.Index(rk, "/specification/"):]
					normalizedPath := filepath.ToSlash(errItem.SchemaPathWithPosition)
					if subIndex := strings.Index(normalizedPath, "/specification/"); subIndex != -1 {
						location = normalizedPath[subIndex:]
					}

					if !isSuppressedInApiTest(config.SuppressionList, errItem.Code, rk, item.OperationInfo.OperationId) {
						mdTable = append(mdTable, fmt.Sprintf("|[%s](%s)|**message**: %s.<br>**opeartion**: %s<br>**location**: %s", errItem.Code, errItem.Link, errItem.Message, item.OperationInfo.OperationId, location))
					}
				}
			}
		}
	}

	sort.Strings(mdTable)

	mdReportFilePath := path.Join(testReportPath, fmt.Sprintf("%s.md", ApiTestReportFileName))
	if err := os.WriteFile(mdReportFilePath, []byte(mdTitle+"\n"+strings.Join(mdTable, "\n")), 0644); err != nil {
		return fmt.Errorf("error when writing file(%s): %+v", mdReportFilePath, err)
	}

	return nil
}

func generateApiTestJsonReport(result map[string]*ApiTestReport, swaggerFilePaths []string, testReportPath string) error {
	for idx, filePath := range swaggerFilePaths {
		logrus.Infof("generating oav report for %d: %s", idx, filePath)
		if report, err := retrieveOavReport(filePath, testReportPath); err != nil {
			return err
		} else {
			result[filePath] = report
		}
	}

	jsonStr, err := json.Marshal(result)
	if err != nil {
		return err
	}

	jsonReportFilePath := path.Join(testReportPath, fmt.Sprintf("%s.json", ApiTestReportFileName))
	if err = os.WriteFile(jsonReportFilePath, jsonStr, 0644); err != nil {
		return fmt.Errorf("error when writing file(%s): %+v", jsonReportFilePath, err)
	}

	return nil
}

func readApiTestHistoryReport(swaggerPath string, apiTestReportFilePath string) (map[string]*ApiTestReport, error) {
	swaggerPathInfo, err := os.Stat(swaggerPath)
	if err != nil {
		return nil, err
	}

	if !swaggerPathInfo.IsDir() {
		swaggerPath = path.Dir(swaggerPath)
	}

	swaggerFiles, err := getApiTestSwaggerPathsFromDirectory(swaggerPath)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*ApiTestReport)
	for _, v := range swaggerFiles {
		result[v] = nil
	}

	if _, err = os.Stat(apiTestReportFilePath); os.IsNotExist(err) {
		logrus.Infof("no history report found")
		return result, nil
	}

	contentBytes, err := os.ReadFile(apiTestReportFilePath)
	if err != nil {
		return nil, fmt.Errorf("error when opening file(%s): %+v", apiTestReportFilePath, err)
	}

	var payload map[string]*ApiTestReport
	err = json.Unmarshal(contentBytes, &payload)
	if err != nil {
		return nil, fmt.Errorf("error during Unmarshal() for file(%s): %+v", apiTestReportFilePath, err)
	}

	for k := range result {
		if v, exists := payload[k]; exists {
			result[k] = v
		}
	}

	return result, nil
}

func retrieveOavReport(filePath string, testReportPath string) (*ApiTestReport, error) {
	traceLogPath := path.Join(testReportPath, TraceLogDirName)
	prefix := strings.Split(path.Base(filePath), ".")[0] + "_report"
	htmlReportFilePath := path.Join(testReportPath, fmt.Sprintf("%s.html", prefix))
	jsonReportFilePath := path.Join(testReportPath, fmt.Sprintf("%s.json", prefix))

	if err := os.RemoveAll(htmlReportFilePath); err != nil {
		return nil, fmt.Errorf("[ERROR] error removing test report file %s: %+v", htmlReportFilePath, err)
	}

	if err := os.RemoveAll(jsonReportFilePath); err != nil {
		return nil, fmt.Errorf("[ERROR] error removing test report file %s: %+v", jsonReportFilePath, err)
	}

	cmd := exec.Command("oav", "validate-traffic", traceLogPath, filePath, "--report", htmlReportFilePath)
	if err := cmd.Run(); err != nil {
		logrus.Infof("oav validates errors from traffic: %+v", err)
	}

	if _, err := os.Stat(path.Join(testReportPath, jsonReportFilePath)); os.IsNotExist(err) {
		return nil, fmt.Errorf("[ERROR] failed to generate report for: %s", testReportPath)
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

	if err = os.RemoveAll(jsonReportFilePath); err != nil {
		return nil, fmt.Errorf("[ERROR] error removing test report file %s: %+v", jsonReportFilePath, err)
	}

	if payload != nil && payload.ApiVersion != "unknown" {
		return payload, nil
	}

	return nil, nil
}

func getApiTestSwaggerPaths(swaggerPath string) ([]string, error) {
	logrus.Infof("loading swagger spec: %s...", swaggerPath)
	file, err := os.Stat(swaggerPath)
	if err != nil {
		return nil, fmt.Errorf("loading swagger spec: %+v", err)
	}

	apiPathsAll := make([]string, 0)
	if file.IsDir() {
		if apiPathsAll, err = getApiTestSwaggerPathsFromDirectory(swaggerPath); err != nil {
			return nil, err
		}

	} else {
		logrus.Infof("parsing swagger spec: %s...", swaggerPath)
		apiPathsAll = append(apiPathsAll, swaggerPath)
	}

	logrus.Infof("found %d api paths", len(apiPathsAll))
	return apiPathsAll, nil
}

func getApiTestSwaggerPathsFromDirectory(swaggerPath string) ([]string, error) {
	logrus.Infof("swagger spec is a directory")
	logrus.Infof("loading swagger spec directory: %s...", swaggerPath)
	files, err := os.ReadDir(swaggerPath)
	if err != nil {
		return nil, fmt.Errorf("reading swagger spec directory: %+v", err)
	}

	apiPathsAll := make([]string, 0)
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") || file.IsDir() {
			continue
		}

		filePath := path.Join(swaggerPath, file.Name())
		apiPathsAll = append(apiPathsAll, filePath)
	}

	return apiPathsAll, nil
}

func mergeApiTestTraceFiles(wd string, traceLogPath string) error {
	if err := os.RemoveAll(traceLogPath); err != nil {
		return fmt.Errorf("[ERROR] error removing test trace dir %s: %+v", traceLogPath, err)
	}

	if err := os.MkdirAll(traceLogPath, 0755); err != nil {
		return fmt.Errorf("[ERROR] error creating test report dir %s: %+v", traceLogPath, err)
	}

	destIndex := 1
	err := filepath.WalkDir(wd, func(walkPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			traceDir := filepath.Join(walkPath, TestResultDirName, TraceLogDirName)
			if _, err = os.Stat(traceDir); !os.IsNotExist(err) {
				destIndex, err = copyApiTestTraceFiles(traceDir, traceLogPath, destIndex)
				if err != nil {
					return fmt.Errorf("failed to copy trace files: %+v", err)
				}
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("impossible to walk directories: %s", err)
	}

	return nil
}

func copyApiTestTraceFiles(src string, dest string, destIndex int) (int, error) {
	err := filepath.WalkDir(src, func(walkPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}

		srcFile, err := os.Open(walkPath)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destPath := path.Join(dest, fmt.Sprintf("trace-%d.json", destIndex))
		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}

		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		destIndex++
		return nil
	})

	if err != nil {
		return destIndex, err
	}

	return destIndex, nil
}
