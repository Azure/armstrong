package rawlog

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type RawLog struct {
	TimeStamp time.Time
	Level     string
	Message   string
}

var regLayoutMap = map[*regexp.Regexp]string{
	regexp.MustCompile(`([\d+.:T\-]{28})\s\[([A-Z]+)]`):  "2006-01-02T15:04:05.999-0700",
	regexp.MustCompile(`([\d+.:T\- ]{19})\s\[([A-Z]+)]`): "2006-01-02 15:04:05",
	regexp.MustCompile(`([\d+.:T/ ]{19})\s\[([A-Z]+)]`):  "2006/01/02 15:04:05",
}

func NewRawLog(message string) (*RawLog, error) {
	for reg, layout := range regLayoutMap {
		matches := reg.FindAllStringSubmatch(message, -1)
		if len(matches) == 0 || len(matches[0]) != 3 {
			continue
		}
		t, err := time.Parse(layout, matches[0][1])
		if err != nil {
			continue
		}
		m := message[len(matches[0][0]):]
		m = strings.Trim(m, " \n")
		return &RawLog{
			TimeStamp: t,
			Level:     matches[0][2],
			Message:   m,
		}, nil
	}
	return nil, fmt.Errorf("failed to parse log message: %s", message)
}
