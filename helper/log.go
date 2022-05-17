package helper

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	"github.com/ms-henglu/azurerm-restapi-testing-tool/types"
)

func ParseLogs(filepath string) ([]types.RequestTrace, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	logs := make([]types.RequestTrace, 0)
	logPrefixReg, _ := regexp.Compile(`^\d{4}-\d{2}-\d{2}`)
	temp := ""
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if logPrefixReg.MatchString(line) {
			if (strings.Contains(temp, "OUTGOING REQUEST") || strings.Contains(temp, "REQUEST/RESPONSE")) && strings.Contains(temp, "management.azure.com") {
				logs = append(logs, NewRequestTrace(temp))
			}
			temp = ""
		}
		temp += line + "\n"
	}

	return logs, nil
}

func NewRequestTrace(raw string) types.RequestTrace {
	trace := types.RequestTrace{}

	methodReg, _ := regexp.Compile(`([A-Z]+)\shttps`)
	if matches := methodReg.FindAllStringSubmatch(raw, -1); len(matches) > 0 && len(matches[0]) == 2 {
		trace.HttpMethod = matches[0][1]
	}

	statusCodeReg, _ := regexp.Compile(`RESPONSE\sStatus:\s(\d+)`)
	if matches := statusCodeReg.FindAllStringSubmatch(raw, -1); len(matches) > 0 && len(matches[0]) == 2 {
		trace.StatusCode, _ = strconv.ParseInt(matches[0][1], 10, 32)
	}

	idReg, _ := regexp.Compile(`management\.azure\.com(.+)\?api-version`)
	if matches := idReg.FindAllStringSubmatch(raw, -1); len(matches) > 0 && len(matches[0]) == 2 {
		trace.ID = matches[0][1]
	}

	trace.Content = raw
	return trace
}
