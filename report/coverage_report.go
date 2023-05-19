package report

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

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
		for k, isCovered := range v {
			if isCovered {
				covered = append(covered, k)
			} else {
				uncovered = append(uncovered, k)
			}
		}
		sort.Strings(covered)
		sort.Strings(uncovered)

		coverages = append(coverages, fmt.Sprintf("%v. %s\ncovered:%v total:%v\n\ncovered properties:\n```\n%s\n```\n\nuncovered properties:\n```\n%s\n```\n",
			count, k, len(covered), len(covered)+len(uncovered), strings.Join(covered, "\n"), strings.Join(uncovered, "\n")))
	}
	content = strings.ReplaceAll(content, "${coverage}", strings.Join(coverages, "\n"))
	return content
}
