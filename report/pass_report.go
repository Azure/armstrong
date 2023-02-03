package report

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/ms-henglu/armstrong/types"
)

//go:embed pass_report.md
var passedReportTemplate string

func PassedMarkdownReport(report types.PassReport) string {
	resourceTypes := make([]string, 0)
	for _, resource := range report.Resources {
		resourceTypes = append(resourceTypes, fmt.Sprintf("%s (%s)", resource.Type, resource.Address))
	}

	content := passedReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", strings.Join(resourceTypes, "\n"))
	return content
}
