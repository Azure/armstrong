package resource

import (
	"fmt"
	"strings"
)

type Reference struct {
	Label        string
	Type         string
	ResourceType string
	PropertyName string
}

func (r Reference) String() string {
	if r.Type == "data" {
		return fmt.Sprintf("data.%s.%s.%s", r.ResourceType, r.Label, r.PropertyName)
	}
	return fmt.Sprintf("%s.%s.%s", r.ResourceType, r.Label, r.PropertyName)
}

func (r *Reference) IsKnown() bool {
	return r != nil && r.Label != "" && r.PropertyName != "" && r.ResourceType != "" && r.Type != ""
}

func NewReferenceFromAddress(address string) *Reference {
	parts := strings.Split(address, ".")
	switch len(parts) {
	case 3:
		return &Reference{
			Type:         "resource",
			ResourceType: parts[0],
			Label:        parts[1],
			PropertyName: parts[2],
		}
	case 4:
		return &Reference{
			Type:         "data",
			ResourceType: parts[1],
			Label:        parts[2],
			PropertyName: parts[3],
		}
	default:
		return nil
	}
}
