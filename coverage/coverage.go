package coverage

import (
	"fmt"
	"log"
)

func MarkCovered(root interface{}, model *Model) {
	if root == nil || model == nil || model.IsReadOnly {
		return
	}

	model.IsAnyCovered = true

	// https://pkg.go.dev/encoding/json#Unmarshal
	switch value := root.(type) {
	case string:
		if model.Enum != nil {
			if _, ok := (*model.Enum)[value]; !ok {
				log.Println(fmt.Errorf("unexpected enum value %s in %s", value, model.Identifier))
			}

			(*model.Enum)[value] = true
		}

	case bool:
		(*model.Bool)[value] = true

	case float64:

	case []interface{}:
		if model.Item == nil {
			log.Println(fmt.Errorf("unexpected array in %s", model.Identifier))
		}

		for _, item := range value {
			MarkCovered(item, model.Item)
		}

	case map[string]interface{}:
		if model.Discriminator != nil {
			for k, v := range value {
				if k == *model.Discriminator {
					if _, ok := (*model.Variants)[v.(string)]; !ok {
						log.Println(fmt.Errorf("unexpected variants %s in %s", v.(string), model.Identifier))
					}
					MarkCovered(value, (*model.Variants)[v.(string)])
					break
				}
			}
		} else {
			for k, v := range value {
				if _, ok := (*model.Properties)[k]; !ok {
					log.Println(fmt.Errorf("unexpected key %s in %s", k, model.Identifier))
				}
				MarkCovered(v, (*model.Properties)[k])
			}
		}

	default:
		panic(fmt.Errorf("unexpect type %T for json unmarshaled value", value))
	}
}

func ComputeCoverage(model *Model) (int, int) {
	if model == nil || model.IsReadOnly {
		return 0, 0
	}

	// first assume is leaf property
	model.TotalCount = 1

	if model.Enum != nil {
		model.TotalCount = len(*model.Enum)
		for _, isCovered := range *model.Enum {
			if isCovered {
				model.CoveredCount++
			}
		}
	}

	if model.Bool != nil {
		model.TotalCount = 2
		for _, isCovered := range *model.Bool {
			if isCovered {
				model.CoveredCount++
			}
		}
	}

	if model.Item != nil {
		covered, total := ComputeCoverage(model.Item)
		model.CoveredCount += covered
		model.TotalCount += total
	}

	if model.Variants != nil {
		for _, v := range *model.Variants {
			covered, total := ComputeCoverage(v)
			model.CoveredCount += covered
			model.TotalCount += total
		}
	}

	if model.Properties != nil {
		for _, v := range *model.Properties {
			covered, total := ComputeCoverage(v)
			model.CoveredCount += covered
			model.TotalCount += total
		}
	}

	if model.TotalCount == 1 && model.IsAnyCovered {
		model.CoveredCount = 1
	}

	model.IsFullyCovered = model.TotalCount > 0 && model.CoveredCount == model.TotalCount

	return model.CoveredCount, model.TotalCount
}

func SplitCovered(model *Model, covered, uncovered *[]string) {
	if model == nil || model.IsReadOnly {
		return
	}

	if model.IsAnyCovered {
		*covered = append(*covered, model.Identifier)
	} else {
		*uncovered = append(*uncovered, model.Identifier)
	}

	if model.Properties != nil {
		for _, v := range *model.Properties {
			SplitCovered(v, covered, uncovered)
		}
	}

	if model.Variants != nil {
		for _, v := range *model.Variants {
			SplitCovered(v, covered, uncovered)
		}
	}

	if model.Item != nil {
		SplitCovered(model.Item, covered, uncovered)
	}

	if model.Enum != nil {
		for k, isCovered := range *model.Enum {
			if isCovered {
				*covered = append(*covered, fmt.Sprintf("%s(%v)", model.Identifier, k))
			} else {
				*uncovered = append(*uncovered, fmt.Sprintf("%s(%v)", model.Identifier, k))
			}
		}
	}

	if model.Bool != nil {
		for k, isCovered := range *model.Bool {
			if isCovered {
				*covered = append(*covered, fmt.Sprintf("%s(%t)", model.Identifier, k))
			} else {
				*uncovered = append(*uncovered, fmt.Sprintf("%s(%t)", model.Identifier, k))
			}
		}
	}

}
