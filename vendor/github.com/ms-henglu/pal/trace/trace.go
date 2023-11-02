package trace

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/ms-henglu/pal/provider"
	"github.com/ms-henglu/pal/rawlog"
	"github.com/ms-henglu/pal/types"
	"github.com/ms-henglu/pal/utils"
)

var providers = []provider.Provider{
	provider.AzureADProvider{},
	provider.AzureRMProvider{},
	provider.AzAPIProvider{},
}

var providerUrlRegex = regexp.MustCompile(`/subscriptions/[a-zA-Z\d\-]+/providers\?`)

func RequestTracesFromFile(input string) ([]types.RequestTrace, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}
	logRegex := regexp.MustCompile(`([\d+.:T\-/ ]{19,28})\s\[([A-Z]+)]`)
	lines := utils.SplitBefore(string(data), logRegex)
	log.Printf("[INFO] total lines: %d", len(lines))

	traces := make([]types.RequestTrace, 0)
	for _, line := range lines {
		l, err := rawlog.NewRawLog(line)
		if err != nil {
			log.Printf("[WARN] failed to parse log: %v", err)
		}
		if l == nil {
			continue
		}
		t, err := NewRequestTrace(*l)
		if err == nil {
			traces = append(traces, *t)
		}
	}
	requestCount, responseCount := 0, 0
	for _, t := range traces {
		if t.Request != nil {
			requestCount++
		}
		if t.Response != nil {
			responseCount++
		}
	}
	log.Printf("[INFO] total traces: %d", len(traces))
	log.Printf("[INFO] request count: %d", requestCount)
	log.Printf("[INFO] response count: %d", responseCount)

	mergedTraces := make([]types.RequestTrace, 0)
	for i := 0; i < len(traces); i++ {
		// skip GET /subscriptions/******/providers
		if traces[i].Method == "GET" && providerUrlRegex.MatchString(traces[i].Url) && traces[i].Provider == "azurerm" {
			continue
		}

		if traces[i].Request != nil && traces[i].Response != nil {
			mergedTraces = append(mergedTraces, traces[i])
			continue
		}

		if traces[i].Request != nil {
			found := false
			for j := i + 1; j < len(traces); j++ {
				if traces[j].Response == nil || traces[i].Url != traces[j].Url || traces[i].Host != traces[j].Host {
					continue
				}
				found = true
				mergedTraces = append(mergedTraces, types.RequestTrace{
					TimeStamp:  traces[i].TimeStamp,
					Url:        traces[i].Url,
					Method:     traces[i].Method,
					Host:       traces[i].Host,
					StatusCode: traces[j].StatusCode,
					Request:    traces[i].Request,
					Response:   traces[j].Response,
				})
				break
			}
			if !found {
				log.Printf("[WARN] failed to find response for request: url %s, method %s", traces[i].Url, traces[i].Method)
				mergedTraces = append(mergedTraces, traces[i])
			}
		}
	}
	log.Printf("[INFO] merged traces: %d", len(mergedTraces))
	return mergedTraces, nil
}

func VerifyRequestTrace(t types.RequestTrace) []string {
	out := make([]string, 0)
	if len(t.Url) == 0 {
		out = append(out, "[ERROR] url is empty")
	}
	if len(t.Host) == 0 {
		out = append(out, "[ERROR] host is empty")
	}
	if len(t.Method) == 0 {
		out = append(out, "[ERROR] method is empty")
	}
	if t.StatusCode == 0 {
		out = append(out, "[ERROR] status code is empty")
	}
	if t.TimeStamp.IsZero() {
		out = append(out, "[ERROR] timestamp is empty")
	}
	switch {
	case t.Request == nil:
		out = append(out, "[ERROR] request is nil")
	case t.Request.Headers == nil:
		out = append(out, "[ERROR] request headers is nil")
	default:
		if contentLength, ok := t.Request.Headers["Content-Length"]; ok {
			length, err := strconv.ParseInt(contentLength, 10, 64)
			if err != nil {
				out = append(out, fmt.Sprintf("[ERROR] failed to parse content length: %v", err))
			}
			if length != 0 && len(t.Request.Body) == 0 {
				out = append(out, "[ERROR] request body is empty")
			}
		}
	}
	switch {
	case t.Response == nil:
		out = append(out, "[ERROR] response is nil")
	case t.Response.Headers == nil:
		out = append(out, "[ERROR] response headers is nil")
	default:
		if contentLength, ok := t.Response.Headers["Content-Length"]; ok {
			length, err := strconv.ParseInt(contentLength, 10, 64)
			if err != nil {
				out = append(out, fmt.Sprintf("[ERROR] failed to parse content length: %v", err))
			}
			if length != 0 && len(t.Response.Body) == 0 {
				out = append(out, "[ERROR] response body is empty")
			}
		}
	}
	return out
}

func NewRequestTrace(l rawlog.RawLog) (*types.RequestTrace, error) {
	for _, p := range providers {
		if p.IsTrafficTrace(l) {
			return p.ParseTraffic(l)
		}
		if p.IsRequestTrace(l) {
			return p.ParseRequest(l)
		}
		if p.IsResponseTrace(l) {
			return p.ParseResponse(l)
		}
	}
	return nil, fmt.Errorf("TODO: implement other providers")
}
