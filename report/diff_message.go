package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ms-henglu/armstrong/types"
	"github.com/nsf/jsondiff"
)

func DiffMessageTerraform(diff types.Change) string {
	option := jsondiff.DefaultConsoleOptions()
	_, msg := jsondiff.Compare([]byte(diff.Before), []byte(diff.After), &option)
	return msg
}

func DiffMessageReadable(diff types.Change) string {
	option := jsondiff.Options{
		Added:                 jsondiff.Tag{Begin: "\033[0;32m", End: " is not returned from response\033[0m"},
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
		Added:                 jsondiff.Tag{Begin: "", End: " is not returned from response"},
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

func compare(old interface{}, new interface{}, path string) []string {
	if new == nil {
		return []string{fmt.Sprintf("%s: expect %v, but got null", path, old)}
	}
	switch oldValue := old.(type) {
	case map[string]interface{}:
		if newMap, ok := new.(map[string]interface{}); ok {
			res := make([]string, 0)
			for key, value := range oldValue {
				res = append(res, compare(value, newMap[key], fmt.Sprintf("%s.%s", path, key))...)
			}
			return res
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a map, but got %v", path, old, new)}
		}
	case []interface{}:
		if newArr, ok := new.([]interface{}); ok {
			if len(oldValue) != len(newArr) {
				return []string{fmt.Sprintf("%s: expect %d in length, but got %d", path, len(oldValue), len(newArr))}
			}
			res := make([]string, 0)
			for index := range oldValue {
				res = append(res, compare(oldValue[index], newArr[index], fmt.Sprintf("%s.%d", path, index))...)
			}
			return res
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is an array, but got %v", path, old, new)}
		}
	case bool:
		if newBool, ok := new.(bool); ok {
			if newBool != oldValue {
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, new, old)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a bool, but got %v", path, old, new)}
		}
	case string:
		if newString, ok := new.(string); ok {
			if newString != oldValue {
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, new, old)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a string, but got %v", path, old, new)}
		}
	case float64:
		if newValue, ok := new.(float64); ok {
			if newValue != oldValue {
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, new, old)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a number, but got %v", path, old, new)}
		}
	case int64:
		if newValue, ok := new.(int64); ok {
			if newValue != oldValue {
				return []string{fmt.Sprintf("%s: expect %v, but got %v", path, new, old)}
			}
		} else {
			return []string{fmt.Sprintf("%s: expect %v which is a number, but got %v", path, old, new)}
		}
	}
	return nil
}
