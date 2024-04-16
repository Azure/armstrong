package trace

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ms-henglu/pal/rawlog"
	"github.com/ms-henglu/pal/types"
)

func requestTracesFromJsonFile(input string) ([]types.RequestTrace, error) {
	fileData, err := os.Open(input)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %v", err)
	}

	defer fileData.Close()

	reader := bufio.NewReader(fileData)

	var jsonLine map[string]interface{}

	traces := make([]types.RequestTrace, 0)

	traceLines, requestCount, responseCount := 0, 0, 0

	for {
		lineData, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("could not read line %v", err)
		}

		if err := json.Unmarshal(lineData, &jsonLine); err != nil {
			return nil, fmt.Errorf("could not unmarhal text into json %v - json data %s", err, string(lineData))
		}

		l, err := rawlog.NewRawLogJson(jsonLine)
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

		traceLines++
	}

	log.Printf("[INFO] total traces: %d", traceLines)
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