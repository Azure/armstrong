package report

import (
	_ "embed"
	"strings"

	"github.com/azure/armstrong/types"
	paltypes "github.com/ms-henglu/pal/types"
)

//go:embed cleanup_error_report.md
var cleanupErrorReportTemplate string

func CleanupErrorMarkdownReport(report types.Error, logs []paltypes.RequestTrace) string {
	parts := strings.Split(report.Type, "@")
	resourceType := ""
	apiVersion := ""
	if len(parts) == 2 {
		resourceType = parts[0]
		apiVersion = parts[1]
	}
	requestTraces := CleanupAllRequestTracesContent(report.Id, logs)
	content := cleanupErrorReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", resourceType)
	content = strings.ReplaceAll(content, "${api_version}", apiVersion)
	content = strings.ReplaceAll(content, "${request_traces}", requestTraces)
	content = strings.ReplaceAll(content, "${error_message}", report.Message)
	return content
}

func CleanupAllRequestTracesContent(id string, logs []paltypes.RequestTrace) string {
	content := ""
	index := len(logs) - 1
	for ; index >= 0; index-- {
		if IsUrlMatchWithId(logs[index].Url, id) && logs[index].Method == "DELETE" {
			content = RequestTraceToString(logs[index])
			break
		}
	}
	for ; index >= 0; index-- {
		if IsUrlMatchWithId(logs[index].Url, id) && logs[index].Method == "GET" {
			content = RequestTraceToString(logs[index]) + "\n\n\n" + content
			break
		}
	}
	return content
}
