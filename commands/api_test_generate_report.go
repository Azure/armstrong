package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ms-henglu/armstrong/report"
	"github.com/sirupsen/logrus"
)

type ApiTestGenerateReportCommand struct {
	workingDir  string
	swaggerPath string
}

func (c *ApiTestGenerateReportCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("api-test-generate-report")
	fs.StringVar(&c.workingDir, "working-dir", "", "path that contains all the test cases")
	fs.StringVar(&c.swaggerPath, "swagger", "", "path to the .json swagger which is being test")
	fs.Usage = func() { logrus.Error(c.Help()) }
	return fs
}

func (c ApiTestGenerateReportCommand) Help() string {
	helpText := `
Usage: armstrong api-test-generate-report
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c ApiTestGenerateReportCommand) Synopsis() string {
	return "Generate test report for a set of test results"
}

func (c ApiTestGenerateReportCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		logrus.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	return c.Execute()
}

func (c ApiTestGenerateReportCommand) Execute() int {
	wd, err := os.Getwd()
	if err != nil {
		logrus.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}

	if c.workingDir != "" {
		wd, err = filepath.Abs(c.workingDir)
		if err != nil {
			logrus.Error(fmt.Sprintf("working directory is invalid: %+v", err))
			return 1
		}
	}

	if err := report.GenerateApiTestReports(wd, c.swaggerPath); err != nil {
		logrus.Errorf("failed to generate API Test Report: %+v", err)
		return 1
	}

	return 0
}
