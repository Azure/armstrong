package resource

import (
	"net/url"
	"strconv"
	"strings"
)

func GetKeyValueMappings(parameters interface{}, path string) []PropertyDependencyMapping {
	if parameters == nil {
		return []PropertyDependencyMapping{}
	}
	results := make([]PropertyDependencyMapping, 0)
	switch parameters.(type) {
	case map[string]interface{}:
		for key, value := range parameters.(map[string]interface{}) {
			results = append(results, GetKeyValueMappings(value, path+"."+key)...)
		}
	case []interface{}:
		for index, value := range parameters.([]interface{}) {
			results = append(results, GetKeyValueMappings(value, path+"."+strconv.Itoa(index))...)
		}
	case string:
		results = append(results, PropertyDependencyMapping{
			ValuePath: path,
			Value:     parameters.(string),
		})
	default:

	}
	return results
}

func GetUpdatedBody(body interface{}, replacements map[string]string, path string) interface{} {
	if len(replacements) == 0 {
		return body
	}
	switch body.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{}, 0)
		for key, value := range body.(map[string]interface{}) {
			if temp := GetUpdatedBody(value, replacements, path+"."+key); temp != nil {
				res[key] = temp
			}
		}
		return res
	case []interface{}:
		res := make([]interface{}, 0)
		for index, value := range body.([]interface{}) {
			if temp := GetUpdatedBody(value, replacements, path+"."+strconv.Itoa(index)); temp != nil {
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

func GetIdFromResponseExample(response interface{}) string {
	if response != nil {
		if responseMap, ok := response.(map[string]interface{}); ok && responseMap["body"] != nil {
			if bodyMap, ok := responseMap["body"].(map[string]interface{}); ok && bodyMap["id"] != nil {
				if id, ok := bodyMap["id"].(string); ok {
					return id
				}
			}
		}
	}
	return ""
}

func GetParentIdFromId(id string) string {
	idURL, err := url.ParseRequestURI(id)
	if err != nil {
		return ""
	}
	path := idURL.Path

	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")

	components := strings.Split(path, "/")
	parentId := ""
	length := len(components) - 2
	if length-2 >= 0 && components[length-2] == "providers" {
		length -= 2
	}
	for current := 0; current <= length-2; current += 2 {
		key := components[current]
		value := components[current+1]
		parentId += "/" + key + "/" + value
	}

	return parentId
}
