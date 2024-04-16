package trace

import "github.com/ms-henglu/pal/types"

type RequestTraceParser struct {
	format LogFormat
}

type LogFormat string

const (
	JsonParser LogFormat = "json"
	TextParser LogFormat = "text/plain"
)

func NewRequestTraceParser(format LogFormat) *RequestTraceParser {
	return &RequestTraceParser{format: format}
}

func (rtp *RequestTraceParser) ParseFromFile(intput string) ([]types.RequestTrace, error) {
	if rtp.format == JsonParser {
		return requestTracesFromJsonFile(intput)
	}

	return requestTracesFromFile(intput)
}
