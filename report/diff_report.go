package report

import (
	_ "embed"
	"strings"

	"github.com/ms-henglu/armstrong/types"
	paltypes "github.com/ms-henglu/pal/types"
)

//go:embed diff_report.md
var diffReportTemplate string

func DiffMarkdownReport(report types.Diff, logs []paltypes.RequestTrace) string {
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

	content := diffReportTemplate
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

func RequestTracesContent(id string, logs []paltypes.RequestTrace) string {
	content := ""
	index := len(logs) - 1
	for ; index >= 0; index-- {
		if IsUrlMatchWithId(logs[index].Url, id) && logs[index].Method == "GET" {
			content = RequestTraceToString(logs[index])
			break
		}
	}
	for ; index >= 0; index-- {
		if IsUrlMatchWithId(logs[index].Url, id) && logs[index].Method == "PUT" {
			content = RequestTraceToString(logs[index]) + "\n\n\n" + content
			break
		}
	}
	return content
}
