package report

import (
	"encoding/json"
	"fmt"
	"strings"

	paltypes "github.com/ms-henglu/pal/types"
)

func IsUrlMatchWithId(url string, id string) bool {
	return strings.HasPrefix(url, id+"?")
}

func RequestTraceToString(r paltypes.RequestTrace) string {
	return fmt.Sprintf(`%s %s
Status Code: %d
------------ Request ------------
%s
------------ Response ------------
%s

`, r.Method, r.Url, r.StatusCode, HttpRequestToString(r.Request), HttpResponseToString(r.Response))
}

func HttpRequestToString(r *paltypes.HttpRequest) string {
	if r == nil {
		return ""
	}
	headers := ""
	for k, v := range r.Headers {
		headers += fmt.Sprintf("%s: %s\n", k, v)
	}
	bodyContent := r.Body
	var body interface{}
	if err := json.Unmarshal([]byte(bodyContent), &body); err == nil {
		if data, err := json.MarshalIndent(body, "", "  "); err == nil {
			bodyContent = string(data)
		}
	}
	return fmt.Sprintf(`%s
---
%s
`, headers, bodyContent)
}

func HttpResponseToString(r *paltypes.HttpResponse) string {
	if r == nil {
		return ""
	}
	headers := ""
	for k, v := range r.Headers {
		headers += fmt.Sprintf("%s: %s\n", k, v)
	}
	bodyContent := r.Body
	var body interface{}
	if err := json.Unmarshal([]byte(bodyContent), &body); err == nil {
		if data, err := json.MarshalIndent(body, "", "  "); err == nil {
			bodyContent = string(data)
		}
	}
	return fmt.Sprintf(`%s------
%s`, headers, bodyContent)
}
