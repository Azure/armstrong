package report

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/azure/armstrong/types"
)

//go:embed cleanup_passed_report.md
var cleanupReportTemplate string

// for now, we don't display bool and enum detail in coverage detail

func CleanupMarkdownReport(passReport types.PassReport) string {
	resourceTypes := make([]string, 0)
	for _, resource := range passReport.Resources {
		resourceTypes = append(resourceTypes, fmt.Sprintf("%s (%s)", resource.Type, resource.Address))
	}

	content := cleanupReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", strings.Join(resourceTypes, "\n"))

	return content
}
