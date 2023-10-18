package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/ms-henglu/pal/types"
)

var _ Formatter = OavTrafficFormatter{}

type OavTrafficFormatter struct {
}

type OavTraffic struct {
	LiveRequest  LiveRequest  `json:"liveRequest"`
	LiveResponse LiveResponse `json:"liveResponse"`
}

type LiveRequest struct {
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Body    interface{}       `json:"body"`
}

type LiveResponse struct {
	StatusCode string            `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       interface{}       `json:"body"`
}

func (o OavTrafficFormatter) Format(r types.RequestTrace) string {
	var requestBody interface{}
	requestHeaders := make(map[string]string)
	if r.Request != nil {
		err := json.Unmarshal([]byte(r.Request.Body), &requestBody)
		if err != nil {
			requestBody = nil
		}
		requestHeaders = r.Request.Headers
	}

	var responseBody interface{}
	responseHeaders := make(map[string]string)
	if r.Response != nil {
		err := json.Unmarshal([]byte(r.Response.Body), &responseBody)
		if err != nil {
			responseBody = nil
		}
		responseHeaders = r.Response.Headers
	}

	out := OavTraffic{
		LiveRequest: LiveRequest{
			Headers: requestHeaders,
			Method:  r.Method,
			Url:     r.Url,
			Body:    requestBody,
		},
		LiveResponse: LiveResponse{
			StatusCode: fmt.Sprintf("%d", r.StatusCode),
			Headers:    responseHeaders,
			Body:       responseBody,
		},
	}

	content, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return ""
	}
	return string(content)
}
