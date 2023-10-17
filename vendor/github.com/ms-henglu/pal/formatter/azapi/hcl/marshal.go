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
			wrapKey := key
			if strings.HasPrefix(key, "${") && strings.HasSuffix(key, "}") {
				wrapKey = fmt.Sprintf("(%s)", key[2:len(key)-1])
			} else if needWrapKey(key) {
				wrapKey = fmt.Sprintf(`"%s"`, key)
			}
			content += fmt.Sprintf("%s%s = %s\n", prefix+indent, wrapKey, MarshalIndent(value, prefix+indent, indent))
		}
		return fmt.Sprintf("{\n%s%s}", content, prefix)
	case []interface{}:
		content := ""
		for _, value := range i {
			content += fmt.Sprintf("%s%s,\n", prefix+indent, MarshalIndent(value, prefix+indent, indent))
		}
		return fmt.Sprintf("[\n%s%s]", content, prefix)
	case string:
		if strings.HasPrefix(i, "${") && strings.HasSuffix(i, "}") {
			return i[2 : len(i)-1]
		}
		i = strings.ReplaceAll(i, "\\", "\\\\") // escape backslashes
		i = strings.ReplaceAll(i, "\n", "\\n")
		i = strings.ReplaceAll(i, "\t", "\\t")
		i = strings.ReplaceAll(i, "\r", "\\r")
		i = strings.ReplaceAll(i, "\"", "\\\"")
		return fmt.Sprintf(`"%s"`, i)
	default:
		return fmt.Sprintf("%v", i)
	}
}

func needWrapKey(input string) bool {
	if len(input) == 0 {
		return false
	}
	if strings.Contains(input, ".") || strings.Contains(input, "/") {
		return true
	}
	if input[0] == '$' {
		return true
	}
	return input[0] >= '0' && input[0] <= '9'
}
