package coverage

import (
	"fmt"
	"log"
	"strconv"
)

type Model struct {
	Bool                    *map[string]bool   `json:"Bool,omitempty"`
	Discriminator           *string            `json:"Discriminator,omitempty"`
	Enum                    *map[string]bool   `json:"Enum,omitempty"`
	Format                  *string            `json:"Format,omitempty"`
	Identifier              string             `json:"Identifier,omitempty"`
	IsAnyCovered            bool               `json:"IsAnyCovered"`
	IsFullyCovered          bool               `json:"IsFullyCovered,omitempty"`
	HasAdditionalProperties bool               `json:"HasAdditionalProperties,omitempty"`
	CoveredCount            int                `json:"CoveredCount,omitempty"`
	TotalCount              int                `json:"TotalCount,omitempty"`
	IsReadOnly              bool               `json:"IsReadOnly,omitempty"`
	IsRequired              bool               `json:"IsRequired,omitempty"`
	Item                    *Model             `json:"Item,omitempty"`
	Properties              *map[string]*Model `json:"Properties,omitempty"`
	Type                    *string            `json:"Type,omitempty"`
	Variants                *map[string]*Model `json:"Variants,omitempty"`
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
			if m.Enum == nil {
				log.Printf("[Error] unexpected enum %s in %s\n", value, m.Identifier)
			}

			strValue := fmt.Sprintf("%v", value)
			if _, ok := (*m.Enum)[strValue]; !ok {
				log.Printf("[WARN] unexpected enum %s in %s\n", value, m.Identifier)
			}

			(*m.Enum)[strValue] = true
		}

	case bool:
		if m.Bool == nil {
			log.Printf("[Error] unexpected bool %v in %v\n", value, m.Identifier)
		}
		(*m.Bool)[strconv.FormatBool(value)] = true

	case float64:

	case []interface{}:
		if m.Item == nil {
			log.Printf("[Error] unexpected array in %s\n", m.Identifier)
		}

		for _, item := range value {
			m.Item.MarkCovered(item)
		}

	case map[string]interface{}:
		if m.Discriminator != nil {
			for k, v := range value {
				if k == *m.Discriminator {
					if m.Variants == nil {
						log.Printf("[Error] unexpected discriminator %s in %s\n", k, m.Identifier)
					}
					if _, ok := (*m.Variants)[v.(string)]; !ok {
						log.Printf("[Error] unexpected variant %s in %s\n", v.(string), m.Identifier)
					}
					(*m.Variants)[v.(string)].MarkCovered(value)
					break
				}
			}
		}
		for k, v := range value {
			if m.Properties == nil {
				if !m.HasAdditionalProperties && m.Discriminator == nil {
					log.Printf("[Error] unexpected key %s in %s\n", k, m.Identifier)
				}
				return
			}
			if _, ok := (*m.Properties)[k]; !ok {
				if !m.HasAdditionalProperties && m.Discriminator == nil {
					log.Printf("[Error] unexpected key %s in %s\n", k, m.Identifier)
				}
			}
			(*m.Properties)[k].MarkCovered(v)
		}

	default:
		panic(fmt.Errorf("unexpect type %T for json unmarshaled value", value))
	}
}

func (m *Model) CountCoverage() (int, int) {
	if m == nil || m.IsReadOnly {
		return 0, 0
	}

	// first assume is leaf property
	m.TotalCount = 1

	if m.Enum != nil {
		m.TotalCount = len(*m.Enum)
		for _, isCovered := range *m.Enum {
			if isCovered {
				m.CoveredCount++
			}
		}
	}

	if m.Bool != nil {
		m.TotalCount = 2
		for _, isCovered := range *m.Bool {
			if isCovered {
				m.CoveredCount++
			}
		}
	}

	if m.Item != nil {
		covered, total := m.Item.CountCoverage()
		m.CoveredCount += covered
		m.TotalCount += total
	}

	if m.Variants != nil {
		for _, v := range *m.Variants {
			covered, total := v.CountCoverage()
			m.CoveredCount += covered
			m.TotalCount += total
		}
	}

	if m.Properties != nil {
		for _, v := range *m.Properties {
			covered, total := v.CountCoverage()
			m.CoveredCount += covered
			m.TotalCount += total
		}
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
