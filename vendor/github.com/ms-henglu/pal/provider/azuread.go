package provider

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/ms-henglu/pal/rawlog"
	"github.com/ms-henglu/pal/types"
	"github.com/ms-henglu/pal/utils"
)

var _ Provider = AzureADProvider{}
var statusCodeRegex = regexp.MustCompile(`HTTP/\d.\d\s(\d{3})\s.+`)

type AzureADProvider struct {
}

func (a AzureADProvider) IsTrafficTrace(l rawlog.RawLog) bool {
	return false
}

func (a AzureADProvider) ParseTraffic(l rawlog.RawLog) (*types.RequestTrace, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a AzureADProvider) IsRequestTrace(l rawlog.RawLog) bool {
	return l.Level == "INFO" && strings.Contains(l.Message, "============================ Begin AzureAD Request")
}

func (a AzureADProvider) IsResponseTrace(l rawlog.RawLog) bool {
	return l.Level == "INFO" && strings.Contains(l.Message, "============================ Begin AzureAD Response")
}

func (a AzureADProvider) ParseRequest(l rawlog.RawLog) (*types.RequestTrace, error) {
	urlLine := ""
	headers := make(map[string]string)
	method := ""
	url := ""
	for _, line := range strings.Split(l.Message, "\n") {
		switch {
		case line == "" || strings.Contains(line, "======"):
			continue
		case strings.Contains(line, ": "):
			key, value, err := utils.ParseHeader(line)
			if err != nil {
				return nil, err
			}
			headers[key] = value
		default:
			urlLine = line
			if parts := strings.Split(urlLine, " "); len(parts) == 3 {
				method = parts[0]
				url = parts[1]
			}
		}
	}

	return &types.RequestTrace{
		TimeStamp: l.TimeStamp,
		Url:       utils.NormalizeUrlPath(url),
		Method:    method,
		Host:      headers["Host"],
		Provider:  "azuread",
		Request: &types.HttpRequest{
			Headers: headers,
		},
	}, nil
}

func (a AzureADProvider) ParseResponse(l rawlog.RawLog) (*types.RequestTrace, error) {
	headers := make(map[string]string)
	statusCode := 0
	method := ""
	rawUrl := ""
	host := ""
	body := ""
	for _, line := range strings.Split(l.Message, "\n") {
		switch {
		case line == "" || strings.Contains(line, "======"):
			continue
		case statusCodeRegex.FindAllStringSubmatch(line, -1) != nil:
			matches := statusCodeRegex.FindAllStringSubmatch(line, -1)
			if len(matches) > 0 && len(matches[0]) == 2 {
				fmt.Sscanf(matches[0][1], "%d", &statusCode)
			}
		case utils.IsJson(line):
			body = line
		case strings.Contains(line, ": "):
			key, value, err := utils.ParseHeader(line)
			if err != nil {
				return nil, err
			}
			headers[key] = value
		default:
			parts := strings.Split(line, " ")
			if len(parts) == 2 {
				method = parts[0]
				parsedUrl, err := url.Parse(parts[1])
				if err == nil {
					host = parsedUrl.Host
					rawUrl = fmt.Sprintf("%s?%s", parsedUrl.Path, parsedUrl.RawQuery)
				}
			}
		}
	}

	return &types.RequestTrace{
		TimeStamp:  l.TimeStamp,
		Url:        utils.NormalizeUrlPath(rawUrl),
		Method:     method,
		Host:       host,
		StatusCode: statusCode,
		Provider:   "azuread",
		Response: &types.HttpResponse{
			Headers: headers,
			Body:    body,
		},
	}, nil
}
