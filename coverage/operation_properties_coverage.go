package coverage

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/ms-henglu/pal/formatter"
	"github.com/sirupsen/logrus"
)

func NewOperationPropertiesCoverageReport(traceDir string, swaggerPath string) (*CoverageReport, error) {
	files, err := os.ReadDir(traceDir)
	if err != nil {
		return nil, err
	}
	report := &CoverageReport{
		Coverages: make(map[string]*CoverageItem),
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(path.Join(traceDir, file.Name()))
		if err != nil {
			logrus.Warnf("failed to read file %s: %+v", file.Name(), err)
			continue
		}

		var trace formatter.OavTraffic
		if err := json.Unmarshal(data, &trace); err != nil {
			logrus.Warnf("failed to unmarshal file %s: %+v", file.Name(), err)
			continue
		}

		swaggerModel, err := GetModelInfoFromLocalDir(removeQueryParameters(trace.LiveRequest.Url), swaggerPath, trace.LiveRequest.Method)
		if err != nil {
			logrus.Warnf("failed to get model info from local dir: %+v", err)
			continue
		}

		if swaggerModel == nil {
			// the API is not in the swagger file, usually it's an API that out of the testing scope
			continue
		}

		index := fmt.Sprintf("%s-%s", trace.LiveRequest.Method, swaggerModel.ApiPath)
		if swaggerModel.ModelName == "" {
			// the API has no request body, mark it as fully covered
			report.Coverages[index] = &CoverageItem{
				ApiPath:     swaggerModel.ApiPath,
				DisplayName: swaggerModel.OperationID,
				Model: &Model{
					IsFullyCovered: true,
				},
			}
			continue
		}

		if _, ok := report.Coverages[index]; !ok {
			expanded, err := Expand(swaggerModel.ModelName, swaggerModel.SwaggerPath)
			if err != nil {
				logrus.Warnf("failed to expand model %s property: %+v", swaggerModel.ModelName, err)
				continue
			}
			report.Coverages[index] = &CoverageItem{
				ApiPath:     swaggerModel.ApiPath,
				DisplayName: swaggerModel.OperationID,
				Model:       expanded,
			}
		}

		report.Coverages[index].Model.MarkCovered(trace.LiveRequest.Body)
		report.Coverages[index].Model.CountCoverage()
	}

	return report, nil
}

func removeQueryParameters(url string) string {
	return strings.Split(url, "?")[0]
}
