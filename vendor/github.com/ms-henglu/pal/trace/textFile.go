package trace

import (
	"fmt"
	"log"
	"os"
	"regexp"

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

func requestTracesFromFile(input string) ([]types.RequestTrace, error) {
	data, err := os.ReadFile(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}
	logRegex := regexp.MustCompile(`([\d+.:T\-/ ]{19,28})\s\[([A-Z]+)]`)
	lines := utils.SplitBefore(string(data), logRegex)
	log.Printf("[INFO] total lines: %d", len(lines))

	requestCount, responseCount := 0, 0

	traces := make([]types.RequestTrace, 0)
	for _, line := range lines {
		l, err := rawlog.NewRawLog(line)
		if err != nil {
			log.Printf("[WARN] failed to parse log: %v", err)
		}
		if l == nil {
			continue
		}
		t, err := newRequestTrace(*l)
		if err == nil {
			traces = append(traces, *t)
			if t.Request != nil {
				requestCount++
			}
			if t.Response != nil {
				responseCount++
			}
		}
	}

	log.Printf("[INFO] total traces: %d", len(traces))
	log.Printf("[INFO] request count: %d", requestCount)
	log.Printf("[INFO] response count: %d", responseCount)

	return mergeTraces(traces), nil
}
