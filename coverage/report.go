package coverage

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

// for now, we don't display bool and enum detail in coverage detail
const isBoolEnumDisplayed = false

type CoverageReport struct {
	Coverages map[string]*CoverageItem
}

type CoverageItem struct {
	ApiPath     string
	DisplayName string
	Model       *Model
}

func (c *CoverageReport) AddCoverageFromState(resourceId, resourceType string, jsonBody map[string]interface{}, swaggerPath string) error {
	apiVersion := strings.Split(resourceType, "@")[1]

	var swaggerModel *SwaggerModel
	if swaggerPath != "" {
		swaggerModelFromLocal, err := GetModelInfoFromLocalDir(resourceId, swaggerPath, "PUT")
		if err != nil {
			logrus.Warnf("error find the path for %s from local dir: %+v", resourceId, err)
		}
		swaggerModel = swaggerModelFromLocal
	}
	if swaggerModel == nil {
		swaggerModelFromIndex, err := GetModelInfoFromIndex(resourceId, apiVersion, "PUT", "")
		if err != nil {
			return fmt.Errorf("error find the path for %s from index: %+v", resourceId, err)
		}
		swaggerModel = swaggerModelFromIndex
	}

	logrus.Infof("matched API path: %s; modelSwawggerPath: %s\n", swaggerModel.ApiPath, swaggerModel.SwaggerPath)

	if _, ok := c.Coverages[swaggerModel.ApiPath]; !ok {
		expanded, err := Expand(swaggerModel.ModelName, swaggerModel.SwaggerPath)
		if err != nil {
			return fmt.Errorf("error expand model %s property: %+v", swaggerModel.ModelName, err)
		}

		c.Coverages[swaggerModel.ApiPath] = &CoverageItem{
			ApiPath:     swaggerModel.ApiPath,
			DisplayName: resourceType,
			Model:       expanded,
		}
	}
	c.Coverages[swaggerModel.ApiPath].Model.MarkCovered(jsonBody)
	c.Coverages[swaggerModel.ApiPath].Model.CountCoverage()

	return nil
}

func (c *CoverageReport) MarkdownContent() string {
	template := `
### Coverage Status:

#### Summary

${coverage_summary}

#### Details

${coverage_details}
`

	fullyCoveredPath := make([]string, 0)
	partiallyCoveredPath := make([]string, 0)
	for _, v := range c.Coverages {
		if v.Model.IsFullyCovered {
			fullyCoveredPath = append(fullyCoveredPath, v.DisplayName)
		} else {
			partiallyCoveredPath = append(partiallyCoveredPath, fmt.Sprintf("%v (%v/%v)", v.DisplayName, v.Model.RootCoveredCount, v.Model.RootTotalCount))
		}
	}

	summary := ""
	if len(fullyCoveredPath) > 0 {
		summary += fmt.Sprintf("Congratulations! The following resource types are 100%% covered:\n\n- %s\n\n", strings.Join(fullyCoveredPath, "\n- "))
	}
	if len(partiallyCoveredPath) > 0 {
		summary += fmt.Sprintf("The following resource types are partially covered, please help add more test cases:\n\n- %s\n\n", strings.Join(partiallyCoveredPath, "\n- "))
	}

	content := strings.ReplaceAll(template, "${coverage_summary}", summary)

	var coverages []string
	count := 0
	for _, v := range c.Coverages {
		count++

		reportDetail := getReport(v.Model.ModelName, v.Model)
		sort.Strings(reportDetail)

		coverages = append(coverages, fmt.Sprintf(`##### <!-- %[1]v -->
<details open>
<summary>%[1]v  %[2]v</summary>

[swagger](%[3]v)
<blockquote>

%[4]v

</blockquote>
</details>

---
`, v.DisplayName, v.ApiPath, v.Model.SourceFile, strings.Join(reportDetail, "\n\n")))
	}

	sort.Strings(coverages)
	content = strings.ReplaceAll(content, "${coverage_details}", strings.Join(coverages, "\n"))

	return content
}

func (c *CoverageReport) MarkdownContentCompact() string {
	template := `
## Coverage Report
|Operation|Tested properties|Total properties|Coverage|
|---|---|---|---|
`
	content := ""
	total := 0
	covered := 0

	for _, v := range c.Coverages {
		coverage := 100.0
		if v.Model.TotalCount > 0 {
			coverage = float64(v.Model.RootCoveredCount * 100 / v.Model.RootTotalCount)
		}
		content += fmt.Sprintf("|%s|%d|%d|%.1f%%|\n", v.DisplayName, v.Model.RootCoveredCount, v.Model.RootTotalCount, coverage)
		total += v.Model.RootTotalCount
		covered += v.Model.RootCoveredCount
	}

	coverage := 100.0
	if total > 0 {
		coverage = float64(covered * 100 / total)
	}

	content = fmt.Sprintf("%s|%d|%d|%.1f%%|\n", "All", covered, total, coverage) + content
	return template + content
}

func getReport(displayName string, model *Model) []string {
	out := make([]string, 0)
	style := getStyle(model.IsFullyCovered)

	if hasNoDetail(model) {
		// leaf property
		out = append(out,
			fmt.Sprintf(`<!-- %[1]v -->
<details>
<summary><span%[2]v>%[3]v</span></summary>

</details>`, model.Identifier, style, displayName),
		)
		return out
	}

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
		return getReport(displayName, model.Item)
	}

	if model.Properties != nil {
		for k, v := range *model.Properties {
			if v.IsReadOnly {
				continue
			}

			if v.Item != nil && v.Item.IsReadOnly {
				continue
			}

			if v.Variants != nil {
				variantType := v.ModelName
				if v.VariantType != nil {
					variantType = *v.VariantType
				}
				variantKey := getDiscriminatorKey(k, variantType)
				out = append(out, getReport(variantKey, v)...)

				for variantType, variant := range *v.Variants {
					variantType := variantType
					if variant.VariantType != nil {
						variantType = *variant.VariantType
					}
					variantKey := getDiscriminatorKey(k, variantType)
					out = append(out, getReport(variantKey, variant)...)
				}
				continue
			}

			if v.Item != nil && v.Item.Variants != nil {
				variantType := v.Item.ModelName
				if v.Item.VariantType != nil {
					variantType = *v.Item.VariantType
				}
				variantKey := getDiscriminatorKey(k, variantType)
				out = append(out, getReport(variantKey, v)...)

				for variantType, variant := range *v.Item.Variants {
					variantType := variantType
					if variant.VariantType != nil {
						variantType = *variant.VariantType
					}
					variantKey := getDiscriminatorKey(k, variantType)
					out = append(out, getReport(variantKey, variant)...)
				}
				continue
			}

			out = append(out, getReport(k, v)...)
		}
	}

	sort.Strings(out)

	outWithoutVariant := []string{
		fmt.Sprintf(`<!-- %[1]v -->
<details>
<summary><span%[2]v>%[3]v %[4]v</span></summary>
<blockquote>

%[5]v

</blockquote>
</details>`, model.Identifier, style, displayName, getCoverageCount(model), strings.Join(out, "\n\n")),
	}

	if model.IsRoot {
		var variants *map[string]*Model
		if model.Variants != nil {
			variants = model.Variants
		}
		if model.Item != nil && model.Item.Variants != nil {
			variants = model.Item.Variants
		}

		if variants != nil {
			outWithVariant := make([]string, 0)
			outWithVariant = append(outWithVariant, outWithoutVariant...)

			for variantType, variant := range *variants {
				variantType := variantType
				if variant.VariantType != nil {
					variantType = *variant.VariantType
				}
				variantKey := getDiscriminatorKey(displayName, variantType)
				outWithVariant = append(outWithVariant, getReport(variantKey, variant)...)
			}

			sort.Strings(outWithVariant)

			return outWithVariant
		}
	}

	return outWithoutVariant

}

func getEnumBoolReport(name string, isCovered bool) string {
	return fmt.Sprintf("- <span %v>value=%v</span>", getStyle(isCovered), name)
}

func getCoverageCount(model *Model) string {
	if model.Bool != nil {
		return fmt.Sprintf("(bool=%v/%v)", model.BoolCoveredCount, 2)
	}
	if model.Enum != nil {
		return fmt.Sprintf("(enum=%v/%v)", model.EnumCoveredCount, model.EnumTotalCount)
	}
	return fmt.Sprintf("(%v/%v)", model.CoveredCount, model.TotalCount)
}

func hasNoDetail(model *Model) bool {
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

func getDiscriminatorKey(modelName, variantType string) string {
	return fmt.Sprintf("%s{%s}", modelName, variantType)
}
