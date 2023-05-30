package report

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/ms-henglu/armstrong/coverage"
	"github.com/ms-henglu/armstrong/types"
)

//go:embed coverage_report.md
var coverageReportTemplate string

func CoverageMarkdownReport(report types.CoverageReport) string {

	content := coverageReportTemplate

	coverages := []string{}
	count := 0
	for k, v := range report.Coverages {
		count++
		var covered, uncovered []string
		coverage.SplitCovered(v, &covered, &uncovered)

		sort.Strings(covered)
		sort.Strings(uncovered)

		coverages = append(coverages, fmt.Sprintf("%v. %s\ncovered:%v total:%v\n\ncovered properties:\n- %s\n\nuncovered properties:\n\n- %s\n",
			count, k, len(covered), len(covered)+len(uncovered), strings.Join(covered, "\n- "), strings.Join(uncovered, "\n- ")))
	}
	content = strings.ReplaceAll(content, "${coverage}", strings.Join(coverages, "\n"))
	return content
}

func CoverageMarkdownReport2(report types.CoverageReport) string {

	content := coverageReportTemplate

	coverages := []string{}
	count := 0
	for k, v := range report.Coverages {
		count++
		var covered, uncovered []string
		coverage.ComputeCoverage(v)

		coverages = append(coverages, fmt.Sprintf("%v. %s\ncovered:%v total:%v\n\ncovered properties:\n- %s\n\nuncovered properties:\n\n- %s\n",
			count, k, len(covered), len(covered)+len(uncovered), strings.Join(covered, "\n- "), strings.Join(uncovered, "\n- ")))
	}
	content = strings.ReplaceAll(content, "${coverage}", strings.Join(coverages, "\n"))
	return content
}
