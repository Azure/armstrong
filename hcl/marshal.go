package hcl

import (
	"fmt"
	"sort"
	"strings"
)

func MarshalIndent(input interface{}, prefix, indent string) string {
	if input == nil {
		return "null"
	}
	switch i := input.(type) {
	case map[string]interface{}:
		keys := make([]string, 0)
		for key := range i {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		content := ""
		for _, key := range keys {
			value := i[key]
			content += fmt.Sprintf("%s%s = %s\n", prefix+indent, key, MarshalIndent(value, prefix+indent, indent))
		}
		return fmt.Sprintf("{\n%s%s}", content, prefix)
	case []interface{}:
		content := ""
		for _, value := range i {
			content += fmt.Sprintf("%s%s\n", prefix+indent, MarshalIndent(value, prefix+indent, indent))
		}
		return fmt.Sprintf("[\n%s%s]", content, prefix)
	case string:
		if strings.HasPrefix(i, "${") && strings.HasSuffix(i, "}") {
			return i[2 : len(i)-1]
		}
		return fmt.Sprintf(`"%s"`, i)
	default:
		return fmt.Sprintf("%v", i)
	}
}
