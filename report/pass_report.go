package report

import (
	_ "embed"
	"fmt"
	"sort"
	"strings"

	"github.com/ms-henglu/armstrong/coverage"
	"github.com/ms-henglu/armstrong/types"
)

//go:embed pass_report.md
var passedReportTemplate string

func PassedMarkdownReport(passReport types.PassReport, coverageReport types.CoverageReport) string {
	resourceTypes := make([]string, 0)
	for _, resource := range passReport.Resources {
		resourceTypes = append(resourceTypes, fmt.Sprintf("%s (%s)", resource.Type, resource.Address))
	}

	content := passedReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", strings.Join(resourceTypes, "\n"))

	var coverages []string
	count := 0
	for k, v := range coverageReport.Coverages {
		count++

		reportDetail := generateReport(v)
		sort.Strings(reportDetail)

		coverages = append(coverages, fmt.Sprintf(`### <!-- %[1]v -->
<details open>
<summary>%[1]v (%[2]v) [%[3]v]</summary>

[swagger](%[4]v)
<blockquote>
<details open>
<summary><span%[5]v>body(%[6]v/%[7]v)</span></summary>
<blockquote>

%[8]v

</blockquote>
</details>

</blockquote>
</details>

---
`, k.Type, k.Address, k.ApiPath, v.SourceFile, getStyle(v.IsFullyCovered), v.CoveredCount, v.TotalCount, strings.Join(reportDetail, "\n\n")))
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

	if model.Properties != nil {
		for k, v := range *model.Properties {
			if v.IsReadOnly {
				continue
			}

			if v.Variants != nil {
				for variantType, variant := range *v.Variants {
					out = append(out, getChildReport(fmt.Sprintf("%s{%s}", k, variantType), variant))
				}
			}

			if v.Item != nil && v.Item.Variants != nil {
				for variantType, variant := range *v.Item.Variants {
					out = append(out, getChildReport(fmt.Sprintf("%s{%s}", k, variantType), variant))
				}
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
	var style, report string

	style = getStyle(model.IsFullyCovered)

	if model.TotalCount == 1 && model.Bool == nil && model.Enum == nil && model.Properties == nil && model.Variants == nil && model.Item == nil {
		// leaf property
		report = fmt.Sprintf(`<!-- %[1]v -->
<details>
<summary><span%[2]v>%[3]v</span></summary>

</details>`, model.Identifier, style, name)
	} else {
		childReport := generateReport(model)
		sort.Strings(childReport)
		report = fmt.Sprintf(`<!-- %[1]v -->
<details>
<summary><span%[2]v>%[3]v %[4]v</span></summary>
<blockquote>

%[5]v

</blockquote>
</details>`, model.Identifier, style, name, getCoverageCount(model), strings.Join(childReport, "\n\n"))
	}

	return report
}

func getStyle(isFullyCovered bool) string {
	if isFullyCovered {
		return ""
	}
	return " style=\"color:red\""
}
