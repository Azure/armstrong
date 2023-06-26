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

	fullyCoveredPath := make([]string, 0)
	partiallyCoveredPath := make([]string, 0)
	for k, v := range report.Coverages {
		if v.IsFullyCovered {
			fullyCoveredPath = append(fullyCoveredPath, k)
		} else {
			partiallyCoveredPath = append(partiallyCoveredPath, k)
		}
	}

	summary := ""
	if len(fullyCoveredPath) > 0 {
		summary += fmt.Sprintf("Congratulations! The following API paths are 100%% covered:\n\n- %s\n\n", strings.Join(fullyCoveredPath, "\n- "))
	}
	if len(partiallyCoveredPath) > 0 {
		summary += fmt.Sprintf("The following API paths are partially covered, please help add more test cases:\n\n- %s\n\n", strings.Join(partiallyCoveredPath, "\n- "))
	}

	content = strings.ReplaceAll(content, "${summary}", summary)

	var coverages []string
	count := 0
	for k, v := range report.Coverages {
		count++

		reportDetail := generateReport(v)
		sort.Strings(reportDetail)

		coverages = append(coverages, fmt.Sprintf(`### <!-- %[1]v -->
<details open>
<summary>%[1]v</summary>

[swagger](%[2]v)
<blockquote>
<details open>
<summary><span%[3]v>body(%[4]v/%[5]v)</span></summary>
<blockquote>

%[6]v

</blockquote>
</details>

</blockquote>
</details>

---
`, k, v.SourceFile, getStyle(v.IsFullyCovered), v.CoveredCount, v.TotalCount, strings.Join(reportDetail, "\n\n")))
	}

	sort.Strings(coverages)
	content = strings.ReplaceAll(content, "${coverage}", strings.Join(coverages, "\n"))
	return content
}

func generateReport(model *coverage.Model) []string {
	out := make([]string, 0)

	if model.Enum != nil {
		for k, isCovered := range *model.Enum {
			out = append(out, getEnumBoolReport(k, isCovered))
		}
	}

	if model.Bool != nil {
		for k, isCovered := range *model.Bool {
			out = append(out, getEnumBoolReport(fmt.Sprintf("%v", k), isCovered))
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

func getEnumBoolReport(name string, isCovered bool) string {
	return fmt.Sprintf("- <span %v>value=%v</span>", getStyle(isCovered), name)
}

func getCoverageCount(model *coverage.Model) string {
	if model.Bool != nil {
		return fmt.Sprintf("(bool=%v/%v)", model.BoolCoveredCount, 2)
	}
	if model.Enum != nil {
		return fmt.Sprintf("(enum=%v/%v)", model.EnumCoveredCount, model.EnumTotalCount)
	}
	return fmt.Sprintf("(%v/%v)", model.CoveredCount, model.TotalCount)
}

func getChildReport(name string, model *coverage.Model) string {
	var color, report string

	color = getStyle(model.IsFullyCovered)

	if model.TotalCount == 1 && model.Bool == nil && model.Enum == nil && model.Properties == nil && model.Variants == nil {
		// leaf property
		report = fmt.Sprintf(`<!-- %[1]v -->
<details>
<summary><span%[2]v>%[3]v</span></summary>

</details>`, model.Identifier, color, name)
	} else {
		childReport := generateReport(model)
		sort.Strings(childReport)
		report = fmt.Sprintf(`<!-- %[1]v -->
<details>
<summary><span%[2]v>%[3]v %[4]v</span></summary>
<blockquote>

%[5]v

</blockquote>
</details>`, model.Identifier, color, name, getCoverageCount(model), strings.Join(childReport, "\n\n"))
	}

	return report
}

func getStyle(isFullyCovered bool) string {
	if isFullyCovered {
		return ""
	}
	return " style=\"color:red\""
}
