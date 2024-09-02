package report

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/azure/armstrong/coverage"
	"github.com/azure/armstrong/types"
)

//go:embed pass_report.md
var passedReportTemplate string

func PassedMarkdownReport(passReport types.PassReport, coverageReport coverage.CoverageReport) string {
	resourceTypes := make([]string, 0)
	for _, resource := range passReport.Resources {
		resourceTypes = append(resourceTypes, fmt.Sprintf("%s (%s)", resource.Type, resource.Address))
	}

	content := passedReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", strings.Join(resourceTypes, "\n"))

	content += coverageReport.MarkdownContent()

	return content
}
