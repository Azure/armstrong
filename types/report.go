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

func (c *CoverageReport) AddCoverageFromState(resourceId, apiVersion string, jsonBody map[string]interface{}) error {
	var apiPath, modelName, modelSwaggerPath, commitId *string
	var err error

	apiPath, modelName, modelSwaggerPath, commitId, err = coverage.GetModelInfoFromIndex(resourceId, apiVersion)
	if err != nil {
		return fmt.Errorf("error find the path for %s from index:%s", resourceId, err)

	}

	c.CommitId = *commitId

	log.Printf("[INFO] matched API path:%s modelSwawggerPath:%s\n", *apiPath, *modelSwaggerPath)

	versionedPath := fmt.Sprintf("%s?api-version=%s", *apiPath, apiVersion)
	if _, ok := c.Coverages[versionedPath]; !ok {
		expanded, err := coverage.Expand(*modelName, *modelSwaggerPath)
		if err != nil {
			return fmt.Errorf("error expand model %s property:%s", *modelName, err)
		}

		c.Coverages[versionedPath] = expanded
	}
	c.Coverages[versionedPath].MarkCovered(jsonBody)
	c.Coverages[versionedPath].CountCoverage()

	return nil
}
