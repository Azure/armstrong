package provider

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ms-henglu/pal/rawlog"
	"github.com/ms-henglu/pal/types"
)

var _ Provider = AzAPIProvider{}

var r1 = regexp.MustCompile(`Live traffic: (.+): timestamp`)
var r2 = regexp.MustCompile(`Live traffic: (.+)`)

type AzAPIProvider struct {
}

func (a AzAPIProvider) IsTrafficTrace(l rawlog.RawLog) bool {
	return l.Level == "DEBUG" && strings.Contains(l.Message, "Live traffic:")
}

func (a AzAPIProvider) ParseTraffic(l rawlog.RawLog) (*types.RequestTrace, error) {
	matches := r1.FindAllStringSubmatch(l.Message, -1)
	if len(matches) == 0 || len(matches[0]) != 2 {
		matches = r2.FindAllStringSubmatch(l.Message, -1)
		if len(matches) == 0 || len(matches[0]) != 2 {
			return nil, fmt.Errorf("failed to parse request trace, no matches found")
		}
	}
	trafficJson := matches[0][1]
	var liveTraffic traffic
	err := json.Unmarshal([]byte(trafficJson), &liveTraffic)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request trace, %v", err)
	}
	parsedUrl, err := url.Parse(liveTraffic.LiveRequest.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse request trace, %v", err)
	}

	if liveTraffic.LiveRequest.Headers == nil {
		liveTraffic.LiveRequest.Headers = map[string]string{}
	}
	if liveTraffic.LiveResponse.Headers == nil {
		liveTraffic.LiveResponse.Headers = map[string]string{}
	}

	return &types.RequestTrace{
		TimeStamp:  l.TimeStamp,
		Method:     liveTraffic.LiveRequest.Method,
		Host:       parsedUrl.Host,
		Url:        parsedUrl.Path + "?" + parsedUrl.RawQuery,
		StatusCode: liveTraffic.LiveResponse.StatusCode,
		Provider:   "azapi",
		Request: &types.HttpRequest{
			Headers: liveTraffic.LiveRequest.Headers,
			Body:    liveTraffic.LiveRequest.Body,
		},
		Response: &types.HttpResponse{
			Headers: liveTraffic.LiveResponse.Headers,
			Body:    liveTraffic.LiveResponse.Body,
		},
	}, nil
}

func (a AzAPIProvider) IsRequestTrace(l rawlog.RawLog) bool {
	return false
}

func (a AzAPIProvider) IsResponseTrace(l rawlog.RawLog) bool {
	return false
}

func (a AzAPIProvider) ParseRequest(l rawlog.RawLog) (*types.RequestTrace, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a AzAPIProvider) ParseResponse(l rawlog.RawLog) (*types.RequestTrace, error) {
	return nil, fmt.Errorf("not implemented")
}

type traffic struct {
	LiveRequest  liveRequest  `json:"request"`
	LiveResponse liveResponse `json:"response"`
}

type liveRequest struct {
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Body    string            `json:"body"`
}

type liveResponse struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}
