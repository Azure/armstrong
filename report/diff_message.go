package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/azure/armstrong/types"
	"github.com/nsf/jsondiff"
)

func DiffMessageTerraform(diff types.Change) string {
	option := jsondiff.DefaultConsoleOptions()
	_, msg := jsondiff.Compare([]byte(diff.Before), []byte(diff.After), &option)
	return msg
}

func DiffMessageReadable(diff types.Change) string {
	option := jsondiff.Options{
		Added:                 jsondiff.Tag{Begin: "\033[0;32m", End: " is not returned from response, if it's on purpose, please follow https://github.com/Azure/armstrong?tab=readme-ov-file#troubleshooting\033[0m"},
		Removed:               jsondiff.Tag{Begin: "\033[0;31m", End: "\033[0m"},
		Changed:               jsondiff.Tag{Begin: "\033[0;33m Got ", End: "\033[0m"},
		Skipped:               jsondiff.Tag{Begin: "\033[0;90m", End: "\033[0m"},
		SkippedArrayElement:   jsondiff.SkippedArrayElement,
		SkippedObjectProperty: jsondiff.SkippedObjectProperty,
		ChangedSeparator:      " in response, expect ",
		Indent:                "    ",
	}
	_, msg := jsondiff.Compare([]byte(diff.Before), []byte(diff.After), &option)
	return msg
}

func DiffMessageMarkdown(diff types.Change) string {
	option := jsondiff.Options{
		Added:                 jsondiff.Tag{Begin: "", End: " is not returned from response, if it's on purpose, please follow https://github.com/Azure/armstrong?tab=readme-ov-file#troubleshooting"},
		Removed:               jsondiff.Tag{Begin: "", End: ""},
		Changed:               jsondiff.Tag{Begin: "Got ", End: ""},
		Skipped:               jsondiff.Tag{Begin: "", End: ""},
		SkippedArrayElement:   jsondiff.SkippedArrayElement,
		SkippedObjectProperty: jsondiff.SkippedObjectProperty,
		ChangedSeparator:      " in response, expect ",
		Indent:                "    ",
	}
	_, msg := jsondiff.Compare([]byte(diff.Before), []byte(diff.After), &option)
	return msg
}

func DiffMessageDescription(diff types.Change) string {
	var before, after interface{}
	_ = json.Unmarshal([]byte(diff.Before), &before)
	_ = json.Unmarshal([]byte(diff.After), &after)
	diffs := compare(before, after, "- ")
	return strings.Join(diffs, "\n")
}

// compare two json objects, return the difference in string array
// path is the path of the json object
// got is the value returned from the api
// expect is the expected value which is defined in the config file
func compare(got interface{}, expect interface{}, path string) []string {
	if expect == nil && got == nil {
		return []string{}
	}
	if expect == nil {
		return []string{fmt.Sprintf("%s: expect null, but got %v", path, got)}
	}
	if got == nil {
		return []string{fmt.Sprintf("%s = %v: not returned from response", path, expect)}
	}
	switch expectValue := expect.(type) {
	case map[string]interface{}:
		if gotMap, ok := got.(map[string]interface{}); ok {
			res := make([]string, 0)
			for key, value := range expectValue {
				res = append(res, compare(gotMap[key], value, fmt.Sprintf("%s.%s", path, key))...)
			}
			return res
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a map, but got %v", path, expect, got)}
		}
	case []interface{}:
		if gotArr, ok := got.([]interface{}); ok {
			if len(gotArr) != len(expectValue) {
				return []string{fmt.Sprintf("%s: expect %d in length, but got %d", path, len(expectValue), len(gotArr))}
			}
			res := make([]string, 0)
			for index := range expectValue {
				res = append(res, compare(gotArr[index], expectValue[index], fmt.Sprintf("%s.%d", path, index))...)
			}
			return res
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is an array, but got %v", path, expect, got)}
		}
	case bool:
		if gotBool, ok := got.(bool); ok {
			if gotBool != expectValue {
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, expect, got)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a bool, but got %v", path, expect, got)}
		}
	case string:
		if gotString, ok := got.(string); ok {
			if gotString != expectValue {
				if strings.EqualFold(gotString, expectValue) {
					return []string{fmt.Sprintf("%s: the values are not equal case-sensitively, expect %v, but got %v", path, expect, got)}
				}
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, expect, got)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a string, but got %v", path, expect, got)}
		}
	case float64:
		if gotFloat, ok := got.(float64); ok {
			if gotFloat != expectValue {
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, expect, got)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a number, but got %v", path, expect, got)}
		}
	case int64:
		if gotInt, ok := got.(int64); ok {
			if gotInt != expectValue {
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, expect, got)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a number, but got %v", path, expect, got)}
		}
	}
	return nil
}
