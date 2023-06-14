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

func CoverageMarkdownReport2(report types.CoverageReport) string {
	content := coverageReportTemplate

	coverages := []string{}
	count := 0
	for k, v := range report.Coverages {
		count++
		var covered, uncovered []string
		v.SplitCovered(&covered, &uncovered)

		sort.Strings(covered)
		sort.Strings(uncovered)

		coverages = append(coverages, fmt.Sprintf("%v. %s\ncovered:%v total:%v\n\ncovered properties:\n- %s\n\nuncovered properties:\n\n- %s\n",
			count, k, len(covered), len(covered)+len(uncovered), strings.Join(covered, "\n- "), strings.Join(uncovered, "\n- ")))
	}
	content = strings.ReplaceAll(content, "${coverage}", strings.Join(coverages, "\n"))
	return content
}

func CoverageMarkdownReport(report types.CoverageReport) string {
	content := coverageReportTemplate

	var coverages []string
	count := 0
	for k, v := range report.Coverages {
		count++

		reportDetail := generateReport(v)
		sort.Strings(reportDetail)

		coverages = append(coverages, fmt.Sprintf(`<blockquote><details open><summary>%s</summary><blockquote>

<details open><summary><span %v>body(%v/%v)</span></summary><blockquote>

%v

</blockquote></details>

</blockquote></details>
</blockquote>`, k, getStyle(v.IsFullyCovered), v.CoveredCount, v.TotalCount, strings.Join(reportDetail, "\n\n")))

	}
	sort.Strings(coverages)
	content = strings.ReplaceAll(content, "${specs-commit-id}", report.CommitId)
	content = strings.ReplaceAll(content, "${coverage}", strings.Join(coverages, "\n"))
	return content
}

func generateReport(model *coverage.Model) []string {
	out := make([]string, 0)

	if model.Enum != nil {
		for k, isCovered := range *model.Enum {
			out = append(out, getEnumerableReport(k, isCovered))
		}
	}

	if model.Bool != nil {
		for k, isCovered := range *model.Bool {
			out = append(out, getEnumerableReport(fmt.Sprintf("%v", k), isCovered))
		}
	}

	if model.Item != nil {
		return generateReport(model.Item)
	}

	if model.Variants != nil {
		for k, v := range *model.Variants {
			out = append(out, getChildReport(k, v))
		}
	}

	if model.Properties != nil {
		for k, v := range *model.Properties {
			if v.IsReadOnly {
				continue
			}

			out = append(out, getChildReport(k, v))
		}
	}

	return out
}

func getEnumerableReport(name string, isCovered bool) string {
	return fmt.Sprintf("- <span %v>value=%v</span>", getStyle(isCovered), name)
}

func getChildReport(name string, model *coverage.Model) string {
	childReport := generateReport(model)
	var color, report string

	color = getStyle(model.IsFullyCovered)

	if model.TotalCount == 1 {
		// leaf property
		report = fmt.Sprintf("- <span %v>%v</span>", color, name)
	} else if len(childReport) == 0 {
		// leaf property with enum or bool type
		report = fmt.Sprintf("- <span %v>%v(%v/%v)</span>", color, name, model.CoveredCount, model.TotalCount)
	} else {
		// non-leaf property
		sort.Strings(childReport)
		report = fmt.Sprintf(`<details><summary><span %v>%v(%v/%v)</span></summary><blockquote>

%v

</blockquote></details>`, color, name, model.CoveredCount, model.TotalCount, strings.Join(childReport, "\n\n"))
	}

	return report
}

func getStyle(isFullyCovered bool) string {
	if isFullyCovered {
		return ""
	}
	return "style=\"color:red\""
}
