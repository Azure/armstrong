package coverage

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

type ArmResource struct {
	ApiPath string
	Type    string
}

type CoverageReport struct {
	Coverages map[ArmResource]*Model
}

func (c *CoverageReport) AddCoverageFromState(resourceId, resourceType string, jsonBody map[string]interface{}) error {
	var apiPath, modelName, modelSwaggerPath *string
	var err error

	apiVersion := strings.Split(resourceType, "@")[1]
	if !regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}$`).MatchString(apiVersion) {
		return fmt.Errorf("could not parse apiVersion from resourceType: %s", resourceType)
	}

	apiPath, modelName, modelSwaggerPath, err = GetModelInfoFromIndex(resourceId, apiVersion)
	if err != nil {
		return fmt.Errorf("error find the path for %s from index:%s", resourceId, err)

	}

	log.Printf("[INFO] matched API path:%s modelSwawggerPath:%s\n", *apiPath, *modelSwaggerPath)

	resource := ArmResource{
		ApiPath: *apiPath,
		Type:    resourceType,
	}

	if _, ok := c.Coverages[resource]; !ok {
		expanded, err := Expand(*modelName, *modelSwaggerPath)
		if err != nil {
			return fmt.Errorf("error expand model %s property:%s", *modelName, err)
		}

		c.Coverages[resource] = expanded
	}
	c.Coverages[resource].MarkCovered(jsonBody)
	c.Coverages[resource].CountCoverage()

	return nil
}
