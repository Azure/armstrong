package coverage

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type ArmResource struct {
	ApiPath string
	Type    string
}

type CoverageReport struct {
	Coverages map[ArmResource]*Model
}

func (c *CoverageReport) AddCoverageFromState(resourceId, resourceType string, jsonBody map[string]interface{}, swaggerPath string) error {
	apiVersion := strings.Split(resourceType, "@")[1]

	var swaggerModel *SwaggerModel
	if swaggerPath != "" {
		swaggerModelFromLocal, err := GetModelInfoFromLocalDir(resourceId, apiVersion, swaggerPath)
		if err != nil {
			logrus.Warnf("error find the path for %s from local dir: %+v", resourceId, err)
		}
		swaggerModel = swaggerModelFromLocal
	}
	if swaggerModel == nil {
		swaggerModelFromIndex, err := GetModelInfoFromIndex(resourceId, apiVersion)
		if err != nil {
			return fmt.Errorf("error find the path for %s from index: %+v", resourceId, err)
		}
		swaggerModel = swaggerModelFromIndex
	}

	logrus.Infof("matched API path: %s; modelSwawggerPath: %s\n", swaggerModel.ApiPath, swaggerModel.SwaggerPath)

	resource := ArmResource{
		ApiPath: swaggerModel.ApiPath,
		Type:    resourceType,
	}

	if _, ok := c.Coverages[resource]; !ok {
		expanded, err := Expand(swaggerModel.ModelName, swaggerModel.SwaggerPath)
		if err != nil {
			return fmt.Errorf("error expand model %s property: %+v", swaggerModel.ModelName, err)
		}

		c.Coverages[resource] = expanded
	}
	c.Coverages[resource].MarkCovered(jsonBody)
	c.Coverages[resource].CountCoverage()

	return nil
}
