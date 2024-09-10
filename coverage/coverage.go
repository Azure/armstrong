package coverage

import (
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Model struct {
	Bool                    *map[string]bool   `json:"Bool,omitempty"` // key is the Enum value, value is coverage status
	BoolCoveredCount        int                `json:"BoolCoveredCount,omitempty"`
	CoveredCount            int                `json:"CoveredCount,omitempty"`
	Discriminator           *string            `json:"Discriminator,omitempty"`
	Enum                    *map[string]bool   `json:"Enum,omitempty"` // key is the Enum value, value is coverage status
	EnumCoveredCount        int                `json:"EnumCoveredCount,omitempty"`
	EnumTotalCount          int                `json:"EnumTotalCount,omitempty"`
	Format                  *string            `json:"Format,omitempty"`
	HasAdditionalProperties bool               `json:"HasAdditionalProperties,omitempty"`
	Identifier              string             `json:"Identifier,omitempty"` // e.g., #.properties.accessPolicies[].permissions.certificates
	IsAnyCovered            bool               `json:"IsAnyCovered"`
	IsFullyCovered          bool               `json:"IsFullyCovered,omitempty"`
	IsReadOnly              bool               `json:"IsReadOnly,omitempty"`
	IsRequired              bool               `json:"IsRequired,omitempty"`
	IsRoot                  bool               `json:"IsRoot,omitempty"`
	IsSecret                bool               `json:"IsSecret,omitempty"` // related to x-ms-secret
	Item                    *Model             `json:"Item,omitempty"`
	ModelName               string             `json:"ModelName,omitempty"`
	Properties              *map[string]*Model `json:"Properties,omitempty"`
	RootCoveredCount        int                `json:"RootCoveredCount,omitempty"` // only for root model, covered count plus all variant count if any
	RootTotalCount          int                `json:"RootTotalCount,omitempty"`   // only for root model, total count plus all variant count if any
	SourceFile              string             `json:"SourceFile,omitempty"`
	TotalCount              int                `json:"TotalCount,omitempty"`
	Type                    *string            `json:"Type,omitempty"`
	Variants                *map[string]*Model `json:"Variants,omitempty"`    // variant model name is used as key, this may only contains
	VariantType             *string            `json:"VariantType,omitempty"` // the x-ms-discriminator-value of the variant model if exists, otherwise model name
}

// CredScan scans the input payload (root) and extract the secret field and value in the secrets map.
func (m *Model) CredScan(root interface{}, secrets map[string]string) {
	if root == nil || m == nil || m.IsReadOnly {
		return
	}

	// https://pkg.go.dev/encoding/json#Unmarshal
	switch value := root.(type) {
	case string:

	case bool:

	case float64:

	case []interface{}:
		if m.Item == nil {
			logrus.Errorf("unexpected array in %s", m.Identifier)
		}

		for _, item := range value {
			m.Item.CredScan(item, secrets)
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
						variant.CredScan(value, secrets)

						break
					}
					for _, variant := range *m.Variants {
						if variant.VariantType != nil && *variant.VariantType == v.(string) {
							isMatchProperty = false
							variant.CredScan(value, secrets)

							break Loop
						}
					}
					logrus.Errorf("unexpected variant %s in %s", v.(string), m.Identifier)
				}
			}
		}

		if isMatchProperty {
			for k, v := range value {
				if m.Properties == nil {
					if !m.HasAdditionalProperties {
						logrus.Errorf("unexpected key %s in %s", k, m.Identifier)
					}
					continue
				}
				if _, ok := (*m.Properties)[k]; !ok {
					if !m.HasAdditionalProperties {
						logrus.Errorf("unexpected key %s in %s", k, m.Identifier)
						continue
					}
				}
				if (*m.Properties)[k].IsSecret {
					secrets[fmt.Sprintf("%v.%v", m.Identifier, k)] = fmt.Sprintf("%v", v)
				}
				(*m.Properties)[k].CredScan(v, secrets)
			}
		}

	case nil:

	default:
		logrus.Errorf("unexpect type %T for json unmarshaled value", value)
	}
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
				logrus.Errorf("unexpected enum %s in %s", value, m.Identifier)
			}

			(*m.Enum)[strValue] = true
		}

	case bool:
		if m.Bool == nil {
			logrus.Errorf("unexpected bool %v in %v", value, m.Identifier)
		}
		(*m.Bool)[strconv.FormatBool(value)] = true

	case float64:

	case []interface{}:
		if m.Item == nil {
			logrus.Errorf("unexpected array in %s", m.Identifier)
		}

		for _, item := range value {
			m.Item.MarkCovered(item)
		}

	case map[string]interface{}:
		// decide to match property or variant
		isMatchProperty := true
		if m.Discriminator != nil && m.Variants != nil {
		Loop:
			for k, v := range value {
				if k == *m.Discriminator {
					// if model name or variant type is matched, then we match the property
					if m.ModelName == v.(string) {
						break
					}
					if m.VariantType != nil && *m.VariantType == v.(string) {
						break
					}

					// either the discriminator value hit the variant model name or variant type, we match the variant
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
					logrus.Errorf("unexpected variant %s in %s", v.(string), m.Identifier)
				}
			}
		}

		if isMatchProperty {
			for k, v := range value {
				if m.Properties == nil {
					if !m.HasAdditionalProperties {
						logrus.Errorf("unexpected key %s in %s", k, m.Identifier)
					}
					continue
				}
				if _, ok := (*m.Properties)[k]; !ok {
					if !m.HasAdditionalProperties {
						logrus.Errorf("unexpected key %s in %s", k, m.Identifier)
						continue
					}
				}
				(*m.Properties)[k].MarkCovered(v)
			}
		}

	case nil:

	default:
		logrus.Errorf("unexpect type %T for json unmarshaled value", value)
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
			if v.Item != nil && v.Item.IsReadOnly {
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

	if m.IsRoot {
		if m.Variants != nil {
			for _, v := range *m.Variants {
				v.CountCoverage()
			}
		}

		if m.Item != nil && m.Item.Variants != nil {
			for _, v := range *m.Item.Variants {
				v.CountCoverage()
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

	if m.IsRoot {
		m.RootCoveredCount = m.CoveredCount
		m.RootTotalCount = m.TotalCount
		if m.Variants != nil {
			for _, v := range *m.Variants {
				m.RootCoveredCount += v.CoveredCount
				m.RootTotalCount += v.TotalCount
			}
		}
		if m.Item != nil && m.Item.Variants != nil {
			for _, v := range *m.Item.Variants {
				m.RootCoveredCount += v.CoveredCount
				m.RootTotalCount += v.TotalCount
			}
		}
	}

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
