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
	Coverages map[string]*coverage.Model
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

func (c CoverageReport) AddCoverageFromState(resourceId, swaggerRepoDir, apiVersion string, jsonBody map[string]interface{}, refreshIndex bool) error {
	var apiPath, modelName, modelSwaggerPath *string
	var err error

	apiPath, modelName, modelSwaggerPath, err = coverage.PathPatternFromIdFromIndex(resourceId, apiVersion, swaggerRepoDir, refreshIndex)
	if err != nil {
		return fmt.Errorf("error find the path for %s from index:%s", resourceId, err)

	}

	log.Printf("matched API path:%s modelSwawggerPath:%s\n", *apiPath, *modelSwaggerPath)

	if _, ok := c.Coverages[*apiPath]; !ok {
		expanded, err := coverage.Expand(*modelName, *modelSwaggerPath)
		if err != nil {
			return fmt.Errorf("error expand model %s property:%s", *modelName, err)
		}

		c.Coverages[*apiPath] = expanded
	}
	coverage.MarkCovered(jsonBody, c.Coverages[*apiPath])

	return nil
}
