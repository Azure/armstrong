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

// for now, we don't display bool and enum detail in coverage detail
const isBoolEnumDisplayed = false

func PassedMarkdownReport(passReport types.PassReport, coverageReport coverage.CoverageReport) string {
	resourceTypes := make([]string, 0)
	for _, resource := range passReport.Resources {
		resourceTypes = append(resourceTypes, fmt.Sprintf("%s (%s)", resource.Type, resource.Address))
	}

	content := passedReportTemplate
	content = strings.ReplaceAll(content, "${resource_type}", strings.Join(resourceTypes, "\n"))

	content = addCoverageReport(content, coverageReport)

	return content
}

func addCoverageReport(content string, coverageReport coverage.CoverageReport) string {
	fullyCoveredPath := make([]string, 0)
	partiallyCoveredPath := make([]string, 0)
	for k, v := range coverageReport.Coverages {
		if v.IsFullyCovered {
			fullyCoveredPath = append(fullyCoveredPath, k.Type)
		} else {
			partiallyCoveredPath = append(partiallyCoveredPath, fmt.Sprintf("%v (%v/%v)", k.Type, v.CoveredCount, v.TotalCount))
		}
	}

	summary := ""
	if len(fullyCoveredPath) > 0 {
		summary += fmt.Sprintf("Congratulations! The following resource types are 100%% covered:\n\n- %s\n\n", strings.Join(fullyCoveredPath, "\n- "))
	}
	if len(partiallyCoveredPath) > 0 {
		summary += fmt.Sprintf("The following resource types are partially covered, please help add more test cases:\n\n- %s\n\n", strings.Join(partiallyCoveredPath, "\n- "))
	}

	content = strings.ReplaceAll(content, "${coverage_summary}", summary)

	var coverages []string
	count := 0
	for k, v := range coverageReport.Coverages {
		count++

		reportDetail := getReport(v)
		sort.Strings(reportDetail)

		coverages = append(coverages, fmt.Sprintf(`##### <!-- %[1]v -->
<details open>
<summary>%[1]v  %[2]v</summary>

[swagger](%[3]v)
<blockquote>
<details open>
<summary><span%[4]v>%[5]v(%[6]v/%[7]v)</span></summary>
<blockquote>

%[8]v

</blockquote>
</details>

</blockquote>
</details>

---
`, k.Type, k.ApiPath, v.SourceFile, getStyle(v.IsFullyCovered), v.ModelName, v.CoveredCount, v.TotalCount, strings.Join(reportDetail, "\n\n")))
	}

	sort.Strings(coverages)
	content = strings.ReplaceAll(content, "${coverage_details}", strings.Join(coverages, "\n"))

	return content
}

func getReport(model *coverage.Model) []string {
	out := make([]string, 0)

	if isBoolEnumDisplayed {
		if model.Enum != nil {
			for k, isCovered := range *model.Enum {
				out = append(out, getEnumBoolReport(k, isCovered))
			}
			return out
		}

		if model.Bool != nil {
			for k, isCovered := range *model.Bool {
				out = append(out, getEnumBoolReport(fmt.Sprintf("%v", k), isCovered))
			}
			return out
		}
	}

	if model.Item != nil {
		return getReport(model.Item)
	}

	if model.Properties != nil {
		for k, v := range *model.Properties {
			if v.IsReadOnly {
				continue
			}

			if v.Variants != nil {
				variantType := v.ModelName
				if v.VariantType != nil {
					variantType = *v.VariantType
				}
				variantKey := fmt.Sprintf("%s{%s}", k, variantType)

				out = append(out, getChildReport(variantKey, v))

				for variantType, variant := range *v.Variants {
					variantType := variantType
					if variant.VariantType != nil {
						variantType = *variant.VariantType
					}
					variantKey := fmt.Sprintf("%s{%s}", k, variantType)
					out = append(out, getChildReport(variantKey, variant))
				}
				continue
			}

			if v.Item != nil && v.Item.Variants != nil {
				variantType := v.Item.ModelName
				if v.Item.VariantType != nil {
					variantType = *v.Item.VariantType
				}
				variantKey := fmt.Sprintf("%s{%s}", k, variantType)
				out = append(out, getChildReport(variantKey, v))

				for variantType, variant := range *v.Item.Variants {
					variantType := variantType
					if variant.VariantType != nil {
						variantType = *variant.VariantType
					}
					variantKey := fmt.Sprintf("%s{%s}", k, variantType)
					out = append(out, getChildReport(variantKey, variant))
				}
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
	var style, report string

	style = getStyle(model.IsFullyCovered)

	if hasNoDetail(model) {
		// leaf property
		report = fmt.Sprintf(`<!-- %[1]v -->
<details>
<summary><span%[2]v>%[3]v</span></summary>

</details>`, model.Identifier, style, name)
	} else {
		childReport := getReport(model)
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

func hasNoDetail(model *coverage.Model) bool {
	if model.Properties == nil && model.Variants == nil && model.Item == nil && (!isBoolEnumDisplayed || (model.Bool == nil && model.Enum == nil)) {
		return true
	}

	// array inside array is regarded as no detail
	if model.Item != nil {
		return hasNoDetail(model.Item)
	}

	return false
}

func getStyle(isFullyCovered bool) string {
	if isFullyCovered {
		return ""
	}
	return " style=\"color:red\""
}
