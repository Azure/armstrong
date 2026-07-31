package trace

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/ms-henglu/pal/rawlog"
	"github.com/ms-henglu/pal/types"
)

type RequestTraceParser struct{}

const (
	plainTextRegex string = `([\d+.:TZ\-/ ]{19,28})\s\[([A-Z]+)]`
	jsonTextRegex  string = `({"@{1}).*}{1}`
)

func NewRequestTraceParser() *RequestTraceParser {
	return &RequestTraceParser{}
}

func (rtp *RequestTraceParser) ParseFromFile(input string) ([]types.RequestTrace, error) {
	fileData, err := os.Open(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file %s: %v", input, err)
	}

	defer func() {
		_ = fileData.Close()
	}()

	reader := bufio.NewReader(fileData)

	traces := make([]types.RequestTrace, 0)

	traceLines, requestCount, responseCount, plines, jlines := 0, 0, 0, 0, 0

	plainRegex := regexp.MustCompile(plainTextRegex)
	jsonRegex := regexp.MustCompile(jsonTextRegex)

	for {
		lineData, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("could not read line %v", err)
		}

		text := string(lineData)

		fPR := plainRegex.FindString(text)
		jPR := jsonRegex.FindString(text)

		var line *rawlog.RawLog
		var pErr error

		if len(fPR) > 0 {
			if line, pErr = parsePlainTextLine(text); pErr != nil {
				log.Printf("[WARN] failed to parse plain log line: %v", pErr)
			}
			plines++
		}

		if len(jPR) > 0 {
			if line, pErr = parseJsonLine(lineData); pErr != nil {
				log.Printf("[WARN] failed to parse json log line: %v", pErr)
			}
			jlines++
		}

		if line == nil {
			continue
		}

		t, err := newRequestTrace(*line)
		if err == nil {
			traces = append(traces, *t)

			if t.Request != nil {
				requestCount++
			}
			if t.Response != nil {
				responseCount++
			}
		}

		traceLines++
	}

	log.Printf("[INFO] total traces: %d", traceLines)
	log.Printf("[INFO] total plain log lines: %d", plines)
	log.Printf("[INFO] total json log lines: %d", jlines)
	log.Printf("[INFO] request count: %d", requestCount)
	log.Printf("[INFO] response count: %d", responseCount)

	return mergeTraces(traces), nil
}

func read(reader *bufio.Reader) ([]byte, error) {
	lineData := make([]byte, 0)

	for {
		line, prefix, err := reader.ReadLine()
		if err != nil {
			return line, err
		}

		lineData = append(lineData, line...)

		if !prefix {
			break
		}
	}

	return lineData, nil
}
