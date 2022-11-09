package report

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ms-henglu/armstrong/types"
)

//go:embed report.md
var reportTemplate string

//go:embed passed_report.md
var passedReportTemplate string

func MarkdownReport(report types.Report, logs []types.RequestTrace) string {
	parts := strings.Split(report.Type, "@")
	resourceType := ""
	apiVersion := ""
	if len(parts) == 2 {
		resourceType = parts[0]
		apiVersion = parts[1]
	}

	operationId := "TODO\ne.g., VirtualMachines_Get"

	diffDescription := DiffMessageDescription(report.Change)
	diffJson := DiffMessageMarkdown(report.Change)

	errCodes := make([]string, 0)
	if strings.Contains(diffJson, "in response, expect") {
		errCodes = append(errCodes, "ROUNDTRIP_INCONSISTENT_PROPERTY")
	}
	if strings.Contains(diffJson, "is not returned from response") {
		errCodes = append(errCodes, "ROUNDTRIP_MISSING_PROPERTY")
	}

	requestTraces := RequestTracesContent(report.Id, logs)

	content := reportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", resourceType)
	content = strings.ReplaceAll(content, "${api_version}", apiVersion)
	content = strings.ReplaceAll(content, "${operation_id}", operationId)
	content = strings.ReplaceAll(content, "${error_code_in_title}", strings.Join(errCodes, " && "))
	content = strings.ReplaceAll(content, "${error_code_in_block}", strings.Join(errCodes, "\n"))
	content = strings.ReplaceAll(content, "${request_traces}", requestTraces)
	content = strings.ReplaceAll(content, "${diff_description}", diffDescription)
	content = strings.ReplaceAll(content, "${diff_json}", diffJson)
	return content
}

func PassedMarkdownReport(reports []types.Report) string {
	resourceTypes := make([]string, 0)
	for _, report := range reports {
		resourceTypes = append(resourceTypes, fmt.Sprintf("%s (%s)", report.Type, report.Address))
	}

	content := passedReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", strings.Join(resourceTypes, "\n"))
	return content
}

func RequestTracesContent(id string, logs []types.RequestTrace) string {
	content := ""
	index := len(logs) - 1
	if log, i := findLastLog(logs, id, "GET", "REQUEST/RESPONSE", index); i != -1 {
		st := strings.Index(log.Content, "GET https")
		ed := strings.Index(log.Content, ": timestamp=")
		trimContent := log.Content
		if st < ed {
			trimContent = log.Content[st:ed]
		}
		content = trimContent
		index = i
	}
	if log, i := findLastLog(logs, id, "PUT", "REQUEST/RESPONSE", index); i != -1 {
		st := strings.Index(log.Content, "RESPONSE Status")
		ed := strings.Index(log.Content, ": timestamp=")
		trimContent := log.Content
		if st < ed {
			trimContent = log.Content[st:ed]
		}
		content = trimContent + "\n\n\n" + content
		index = i
	}
	if log, i := findLastLog(logs, id, "PUT", "OUTGOING REQUEST", index); i != -1 {
		st := strings.Index(log.Content, "PUT https")
		ed := strings.Index(log.Content, ": timestamp=")
		trimContent := log.Content
		if st < ed {
			trimContent = log.Content[st:ed]
		}
		content = trimContent + "\n\n" + content
	}
	return content
}

func findLastLog(logs []types.RequestTrace, id string, method string, substr string, index int) (types.RequestTrace, int) {
	for i := index; i >= 0; i-- {
		log := logs[i]
		if log.ID == id && log.HttpMethod == method && strings.Contains(log.Content, substr) {
			return log, i
		}
	}
	return types.RequestTrace{}, -1
}
