package types

import (
	"fmt"
	"log"

	"github.com/ms-henglu/armstrong/coverage"
)

type PassReport struct {
	Resources []Resource
}

type Resource struct {
	Type    string
	Address string
}

type CoverageReport struct {
	CommitId  string
	Coverages map[string]map[string]bool
}

type DiffReport struct {
	Diffs []Diff
	Logs  []RequestTrace
}

type Diff struct {
	Id      string
	Type    string
	Address string
	Change  Change
}

type Change struct {
	Before string
	After  string
}

type ErrorReport struct {
	Errors []Error
	Logs   []RequestTrace
}

type Error struct {
	Id      string
	Type    string
	Label   string
	Message string
}

func (c CoverageReport) AddCoverageFromState(resourceId, swaggerPath, apiVersion string, jsonBody map[string]interface{}) error {
	var apiPath, modelName, modelSwaggerPath *string
	var err error
	if swaggerPath != "" {
		apiPath, modelName, modelSwaggerPath, err = coverage.PathPatternFromId(resourceId, swaggerPath)
		if err != nil {
			return fmt.Errorf("error find the path for %s in swagger file %s:%s", resourceId, swaggerPath, err)
		}
	} else {
		apiPath, modelName, modelSwaggerPath, err = coverage.PathPatternFromIdFromIndex(resourceId, apiVersion)
		if err != nil {
			return fmt.Errorf("error find the path for %s from index:%s", resourceId, err)
		}
	}

	log.Printf("matched API path %s\n", *apiPath)

	expanded, err := coverage.Expand(*modelName, *modelSwaggerPath)
	if err != nil {
		return fmt.Errorf("error expand model %s property:%s", *modelName, err)
	}

	lookupTable := map[string]bool{}
	if coverageTable, ok := c.Coverages[*apiPath]; ok {
		lookupTable = coverageTable
	}
	discriminatorTable := map[string]string{}
	coverage.Flatten(*expanded, "", lookupTable, discriminatorTable)
	coverage.MarkCovered(jsonBody, "", lookupTable, discriminatorTable)
	c.Coverages[*apiPath] = lookupTable

	return nil
}
