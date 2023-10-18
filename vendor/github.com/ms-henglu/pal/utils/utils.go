package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

func IsJson(input string) bool {
	var out interface{}
	err := json.Unmarshal([]byte(input), &out)
	return err == nil
}

func JsonPretty(input string) string {
	var out interface{}
	err := json.Unmarshal([]byte(input), &out)
	if err != nil {
		return input
	}
	b, err := json.MarshalIndent(out, "", "    ")
	if err != nil {
		return input
	}
	return string(b)
}

func SplitBefore(s string, re *regexp.Regexp) []string {
	out := make([]string, 0)
	is := re.FindAllStringIndex(s, -1)
	if len(is) == 0 {
		return append(out, s)
	}
	for i := 0; i < len(is)-1; i++ {
		out = append(out, s[is[i][0]:is[i+1][0]])
	}
	return append(out, s[is[len(is)-1][0]:])
}

func ParseHeader(input string) (string, string, error) {
	deliminatorIndex := strings.Index(input, ":")
	if deliminatorIndex == -1 {
		return "", "", fmt.Errorf("failed to parse header, `:` is not found: %s", input)
	}
	key := strings.Trim(input[0:deliminatorIndex], " ")
	value := input[deliminatorIndex+1:]
	if index := strings.LastIndex(value, ": timestamp"); index != -1 {
		value = value[0:index]
	}
	value = strings.Trim(value, " \n\r")
	return key, value, nil
}

func LineAt(input string, index int) string {
	lines := strings.Split(input, "\n")
	if len(lines) > index {
		return lines[index]
	}
	return input
}

func NormalizeUrlPath(input string) string {
	if !strings.HasPrefix(input, "/") {
		input = "/" + input
	}
	input = strings.ReplaceAll(input, "//", "/")
	return input
}
