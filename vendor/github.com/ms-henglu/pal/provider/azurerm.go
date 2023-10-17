package provider

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ms-henglu/pal/rawlog"
	"github.com/ms-henglu/pal/types"
	"github.com/ms-henglu/pal/utils"
)

var _ Provider = AzureRMProvider{}

type AzureRMProvider struct {
}

func (a AzureRMProvider) IsTrafficTrace(l rawlog.RawLog) bool {
	return false
}

func (a AzureRMProvider) ParseTraffic(l rawlog.RawLog) (*types.RequestTrace, error) {
	return nil, fmt.Errorf("not implemented")
}

func (a AzureRMProvider) IsRequestTrace(l rawlog.RawLog) bool {
	return l.Level == "DEBUG" && strings.Contains(l.Message, "AzureRM Request:")
}

func (a AzureRMProvider) IsResponseTrace(l rawlog.RawLog) bool {
	return l.Level == "DEBUG" && strings.Contains(l.Message, "AzureRM Response for")
}

func (a AzureRMProvider) ParseRequest(l rawlog.RawLog) (*types.RequestTrace, error) {
	urlPath := ""
	method := ""
	headers := make(map[string]string)
	body := ""

	lines := strings.Split(l.Message, "\n")
	i := 0
	foundBodySegment := false
	for ; i < len(lines); i++ {
		line := lines[i]
		switch {
		case strings.TrimSpace(line) == "":
			foundBodySegment = true
		case strings.Contains(line, ": "):
			key, value, err := utils.ParseHeader(line)
			if strings.HasPrefix(key, "provider.terraform-provider-azurerm") {
				continue
			}
			if key == "AzureRM Request" {
				continue
			}
			if err != nil {
				return nil, err
			}
			headers[key] = value
		default:
			if parts := strings.Split(line, " "); len(parts) == 3 {
				method = parts[0]
				urlPath = parts[1]
			}
		}
		if foundBodySegment {
			break
		}
	}

	if i+1 < len(lines) {
		line := strings.Join(lines[i+1:], "\n")
		if strings.Contains(line, ": timestamp") {
			index := strings.LastIndex(line, ": timestamp")
			if utils.IsJson(line[0:index]) {
				body = line[0:index]
			} else {
				lineTrimTimestamp := line[0:index]
				key, value, err := utils.ParseHeader(lineTrimTimestamp)
				if err == nil {
					headers[key] = value
				}
			}
		} else {
			body = line
		}
	}
	return &types.RequestTrace{
		TimeStamp: l.TimeStamp,
		Url:       utils.NormalizeUrlPath(urlPath),
		Method:    method,
		Host:      headers["Host"],
		Provider:  "azurerm",
		Request: &types.HttpRequest{
			Headers: headers,
			Body:    body,
		},
	}, nil
}

func (a AzureRMProvider) ParseResponse(l rawlog.RawLog) (*types.RequestTrace, error) {
	urlPath := ""
	host := ""
	method := "" // TODO: this is not available in the response
	body := ""
	headers := make(map[string]string)
	statusCode := 0

	lines := strings.Split(l.Message, "\n")
	i := 0
	foundBodySegment := false
	for ; i < len(lines); i++ {
		line := lines[i]
		switch {
		case strings.TrimSpace(line) == "":
			foundBodySegment = true
		case strings.Contains(line, "AzureRM Response for "):
			urlLine := line[strings.Index(line, "AzureRM Response for ")+len("AzureRM Response for "):]
			urlLine = strings.Trim(urlLine, " \n\r:")
			parsedUrl, err := url.Parse(urlLine)
			if err != nil {
				return nil, err
			}
			host = parsedUrl.Host
			urlPath = fmt.Sprintf("%s?%s", parsedUrl.Path, parsedUrl.RawQuery)
		case strings.Contains(line, ": "):
			key, value, err := utils.ParseHeader(line)
			if err != nil {
				return nil, err
			}
			headers[key] = value
		default:
			if matches := statusCodeRegex.FindAllStringSubmatch(line, -1); len(matches) > 0 && len(matches[0]) == 2 {
				fmt.Sscanf(matches[0][1], "%d", &statusCode)
			}
		}
		if foundBodySegment {
			break
		}
	}

	if i+1 < len(lines) {
		line := strings.Join(lines[i+1:], "\n")
		if strings.Contains(line, ": timestamp") {
			index := strings.LastIndex(line, ": timestamp")
			if utils.IsJson(line[0:index]) {
				body = line[0:index]
			}
		} else {
			body = line
		}
	}

	return &types.RequestTrace{
		TimeStamp:  l.TimeStamp,
		Url:        utils.NormalizeUrlPath(urlPath),
		Host:       host,
		Method:     method,
		StatusCode: statusCode,
		Provider:   "azurerm",
		Response: &types.HttpResponse{
			Headers: headers,
			Body:    body,
		},
	}, nil
}
