package trace

import (
	"regexp"

	"github.com/ms-henglu/pal/provider"
	"github.com/ms-henglu/pal/rawlog"
)

var providers = []provider.Provider{
	provider.AzureADProvider{},
	provider.AzureRMProvider{},
	provider.AzAPIProvider{},
}

var providerUrlRegex = regexp.MustCompile(`/subscriptions/[a-zA-Z\d\-]+/providers\?`)

func parsePlainTextLine(line string) (*rawlog.RawLog, error) {
	return rawlog.NewRawLog(line)
}
