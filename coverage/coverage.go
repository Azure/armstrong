package coverage

import (
	"fmt"
	"log"
	"strconv"
)

type Model struct {
	Bool                    *map[string]bool   `json:"Bool,omitempty"`
	BoolCoveredCount        int                `json:"BoolCoveredCount,omitempty"`
	CoveredCount            int                `json:"CoveredCount,omitempty"`
	Discriminator           *string            `json:"Discriminator,omitempty"`
	Enum                    *map[string]bool   `json:"Enum,omitempty"`
	EnumCoveredCount        int                `json:"EnumCoveredCount,omitempty"`
	EnumTotalCount          int                `json:"EnumTotalCount,omitempty"`
	Format                  *string            `json:"Format,omitempty"`
	HasAdditionalProperties bool               `json:"HasAdditionalProperties,omitempty"`
	Identifier              string             `json:"Identifier,omitempty"`
	IsAnyCovered            bool               `json:"IsAnyCovered"`
	IsFullyCovered          bool               `json:"IsFullyCovered,omitempty"`
	IsReadOnly              bool               `json:"IsReadOnly,omitempty"`
	IsRequired              bool               `json:"IsRequired,omitempty"`
	Item                    *Model             `json:"Item,omitempty"`
	ModelName               string             `json:"ModelName,omitempty"`
	Properties              *map[string]*Model `json:"Properties,omitempty"`
	SourceFile              string             `json:"SourceFile,omitempty"`
	TotalCount              int                `json:"TotalCount,omitempty"`
	Type                    *string            `json:"Type,omitempty"`
	Variants                *map[string]*Model `json:"Variants,omitempty"`
	VariantType             *string            `json:"VariantType,omitempty"`
}

func (m *Model) MarkCovered(root interface{}) {
	if root == nil || m == nil || m.IsReadOnly {
		return
	}

	m.IsAnyCovered = true

	// https://pkg.go.dev/encoding/json#Unmarshal
	switch value := root.(type) {
	case string:
		if m.Enum != nil {
			strValue := fmt.Sprintf("%v", value)
			if _, ok := (*m.Enum)[strValue]; !ok {
				log.Printf("[WARN] unexpected enum %s in %s\n", value, m.Identifier)
			}

			(*m.Enum)[strValue] = true
		}

	case bool:
		if m.Bool == nil {
			log.Printf("[ERROR] unexpected bool %v in %v\n", value, m.Identifier)
		}
		(*m.Bool)[strconv.FormatBool(value)] = true

	case float64:

	case []interface{}:
		if m.Item == nil {
			log.Printf("[ERROR] unexpected array in %s\n", m.Identifier)
		}

		for _, item := range value {
			m.Item.MarkCovered(item)
		}

	case map[string]interface{}:
		isMatchProperty := true
		if m.Discriminator != nil && m.Variants != nil {
		Loop:
			for k, v := range value {
				if k == *m.Discriminator {
					if m.ModelName == v.(string) {
						break
					}
					if m.VariantType != nil && *m.VariantType == v.(string) {
						break
					}
					if variant, ok := (*m.Variants)[v.(string)]; ok {
						isMatchProperty = false
						variant.MarkCovered(value)

						break
					}
					for _, variant := range *m.Variants {
						if variant.VariantType != nil && *variant.VariantType == v.(string) {
							isMatchProperty = false
							variant.MarkCovered(value)

							break Loop
						}
					}
					log.Printf("[ERROR] unexpected variant %s in %s\n", v.(string), m.Identifier)
				}
			}
		}

		if isMatchProperty {
			for k, v := range value {
				if m.Properties == nil {
					if !m.HasAdditionalProperties {
						log.Printf("[WARN] unexpected key %s in %s\n", k, m.Identifier)
					}
					return
				}
				if _, ok := (*m.Properties)[k]; !ok {
					if !m.HasAdditionalProperties {
						log.Printf("[WARN] unexpected key %s in %s\n", k, m.Identifier)
						return
					}
				}
				(*m.Properties)[k].MarkCovered(v)
			}
		}

	case nil:

	default:
		log.Printf("[ERROR] unexpect type %T for json unmarshaled value", value)
	}
}

func (m *Model) CountCoverage() (int, int) {
	if m == nil || m.IsReadOnly {
		return 0, 0
	}

	m.CoveredCount = 0
	m.TotalCount = 0

	if m.Enum != nil {
		m.EnumCoveredCount = 0
		m.EnumTotalCount = len(*m.Enum)
		for _, isCovered := range *m.Enum {
			if isCovered {
				m.EnumCoveredCount++
			}
		}
	}

	if m.Bool != nil {
		m.BoolCoveredCount = 0
		for _, isCovered := range *m.Bool {
			if isCovered {
				m.BoolCoveredCount++
			}
		}
	}

	if m.Item != nil {
		covered, total := m.Item.CountCoverage()
		m.CoveredCount = covered
		m.TotalCount = total
	}

	if m.Properties != nil {
		for _, v := range *m.Properties {
			if v.IsReadOnly {
				continue
			}
			covered, total := v.CountCoverage()
			m.CoveredCount += covered
			m.TotalCount += total
			if v.Variants != nil {
				for _, variant := range *v.Variants {
					covered, total := variant.CountCoverage()
					m.CoveredCount += covered
					m.TotalCount += total
				}
			}
			if v.Item != nil && v.Item.Variants != nil {
				for _, variant := range *v.Item.Variants {
					covered, total := variant.CountCoverage()
					m.CoveredCount += covered
					m.TotalCount += total
				}
			}
		}
	}

	if m.TotalCount == 0 {
		m.TotalCount = 1
	}
	if m.TotalCount == 1 && m.IsAnyCovered {
		m.CoveredCount = 1
	}

	m.IsFullyCovered = m.TotalCount > 0 && m.CoveredCount == m.TotalCount

	return m.CoveredCount, m.TotalCount
}

func (m *Model) SplitCovered(covered, uncovered *[]string) {
	if m == nil || m.IsReadOnly {
		return
	}

	if m.IsAnyCovered {
		*covered = append(*covered, m.Identifier)
	} else {
		*uncovered = append(*uncovered, m.Identifier)
	}

	if m.Properties != nil {
		for _, v := range *m.Properties {
			v.SplitCovered(covered, uncovered)
		}
	}

	if m.Variants != nil {
		for _, v := range *m.Variants {
			v.SplitCovered(covered, uncovered)
		}
	}

	if m.Item != nil {
		m.Item.SplitCovered(covered, uncovered)
	}

	if m.Enum != nil {
		for k, isCovered := range *m.Enum {
			if isCovered {
				*covered = append(*covered, fmt.Sprintf("%s(%v)", m.Identifier, k))
			} else {
				*uncovered = append(*uncovered, fmt.Sprintf("%s(%v)", m.Identifier, k))
			}
		}
	}

	if m.Bool != nil {
		for k, isCovered := range *m.Bool {
			if isCovered {
				*covered = append(*covered, fmt.Sprintf("%s(%v)", m.Identifier, k))
			} else {
				*uncovered = append(*uncovered, fmt.Sprintf("%s(%v)", m.Identifier, k))
			}
		}
	}
}
