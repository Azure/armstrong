package trace

import (
	"encoding/json"
	"fmt"

	"github.com/ms-henglu/pal/rawlog"
)

var jsonLine map[string]interface{}

func parseJsonLine(logLine []byte) (*rawlog.RawLog, error) {
	if err := json.Unmarshal(logLine, &jsonLine); err != nil {
		return nil, fmt.Errorf("could not unmarhal text into json %v - json data %s", err, string(logLine))
	}

	return rawlog.NewRawLogJson(jsonLine)
}
