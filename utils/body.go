package utils

import (
	"fmt"
	"strconv"
)

func UpdatedBody(body interface{}, replacements map[string]string, path string) interface{} {
	if len(replacements) == 0 {
		return body
	}
	switch bodyValue := body.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{})
		for key, value := range bodyValue {
			if temp := UpdatedBody(value, replacements, path+"."+key); temp != nil {
				if replaceKey := replacements[fmt.Sprintf("key:%s.%s", path, key)]; replaceKey != "" {
					key = replaceKey
				}
				res[key] = temp
			}
		}
		return res
	case []interface{}:
		res := make([]interface{}, 0)
		for index, value := range bodyValue {
			if temp := UpdatedBody(value, replacements, path+"."+strconv.Itoa(index)); temp != nil {
				res = append(res, temp)
			}
		}
		return res
	case string:
		for key, replacement := range replacements {
			if key == path {
				return replacement
			}
		}
	default:

	}
	return body
}
