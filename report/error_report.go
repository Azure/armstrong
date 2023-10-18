package report

import (
	_ "embed"
	paltypes "github.com/ms-henglu/pal/types"
	"strings"

	"github.com/ms-henglu/armstrong/types"
)

//go:embed error_report.md
var errorReportTemplate string

func ErrorMarkdownReport(report types.Error, logs []paltypes.RequestTrace) string {
	parts := strings.Split(report.Type, "@")
	resourceType := ""
	apiVersion := ""
	if len(parts) == 2 {
		resourceType = parts[0]
		apiVersion = parts[1]
	}
	requestTraces := AllRequestTracesContent(report.Id, logs)
	content := errorReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", resourceType)
	content = strings.ReplaceAll(content, "${api_version}", apiVersion)
	content = strings.ReplaceAll(content, "${request_traces}", requestTraces)
	content = strings.ReplaceAll(content, "${error_message}", report.Message)
	return content
}

func AllRequestTracesContent(id string, logs []paltypes.RequestTrace) string {
	content := ""
	for i := len(logs) - 1; i >= 0; i-- {
		if IsUrlMatchWithId(logs[i].Url, id) {
			content += RequestTraceToString(logs[i]) + "\n\n\n"
		}
	}

	return content
}
