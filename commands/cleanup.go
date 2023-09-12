package commands

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ms-henglu/armstrong/report"
	"github.com/ms-henglu/armstrong/tf"
	"github.com/ms-henglu/armstrong/types"
	"github.com/sirupsen/logrus"
)

type CleanupCommand struct {
	verbose    bool
	workingDir string
}

func (c *CleanupCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("cleanup")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { logrus.Error(c.Help()) }
	return fs
}

func (c CleanupCommand) Help() string {
	helpText := `
Usage: armstrong cleanup [-v] [-working-dir <path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c CleanupCommand) Synopsis() string {
	return "Clean up dependencies and testing resource"
}

func (c CleanupCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		logrus.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	if c.verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	return c.Execute()
}

func (c CleanupCommand) Execute() int {
	const (
		allPassedReportFileName     = "cleanup_all_passed_report.md"
		partialPassedReportFileName = "cleanup_partial_passed_report.md"
	)

	logrus.Infof("cleaning up resources...")
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
	terraform, err := tf.NewTerraform(wd, c.verbose)
	if err != nil {
		logrus.Fatalf("creating terraform executable: %+v", err)
	}

	state, err := terraform.Show()
	if err != nil {
		logrus.Fatalf("failed to get terraform state: %+v", err)
	}

	passReport := tf.NewPassReportFromState(state)
	idAddressMap := tf.NewIdAdressFromState(state)

	reportDir := fmt.Sprintf("armstrong_cleanup_reports_%s", time.Now().Format(time.Stamp))
	reportDir = strings.ReplaceAll(reportDir, ":", "")
	reportDir = strings.ReplaceAll(reportDir, " ", "_")
	reportDir = path.Join(wd, reportDir)
	err = os.Mkdir(reportDir, 0755)
	if err != nil {
		logrus.Fatalf("failed to create report directory: %+v", err)
	}

	logrus.Infof("running terraform init...")
	_ = terraform.Init()
	logrus.Infof("running terraform destroy...")
	destroyErr := terraform.Destroy()
	if destroyErr != nil {
		logrus.Errorf("failed to destroy resources: %+v", destroyErr)
	} else {
		logrus.Infof("all resources are cleaned up")
		storeCleanupReport(passReport, reportDir, allPassedReportFileName)
	}

	logs, err := report.ParseLogs(path.Join(wd, "log.txt"))
	if err != nil {
		logrus.Errorf("failed to parse log.txt: %+v", err)
	}

	errorReport := types.ErrorReport{}
	if destroyErr != nil {
		errorReport := tf.NewCleanupErrorReport(destroyErr, logs)
		for i := range errorReport.Errors {
			if address, ok := idAddressMap[errorReport.Errors[i].Id]; ok {
				errorReport.Errors[i].Label = address
			}
		}
		storeCleanupErrorReport(errorReport, reportDir)
	}

	resources := make([]types.Resource, 0)
	if state, err := terraform.Show(); err == nil && state != nil && state.Values != nil && state.Values.RootModule != nil && state.Values.RootModule.Resources != nil {
		for _, passRes := range passReport.Resources {
			isDeleted := true
			for _, res := range state.Values.RootModule.Resources {
				if passRes.Address == res.Address {
					isDeleted = false
					break
				}
			}
			if isDeleted {
				resources = append(resources, passRes)
			}
		}
	}

	if len(resources) > 0 {
		passReport.Resources = resources
		storeCleanupReport(passReport, reportDir, partialPassedReportFileName)
	}

	logrus.Infof("---------------- Summary ----------------")
	logrus.Infof("%d resources passed the cleanup tests.", len(passReport.Resources))
	if len(errorReport.Errors) != 0 {
		logrus.Infof("%d errors when cleanup the testing resources.", len(errorReport.Errors))
	}

	return 0
}

func storeCleanupReport(passReport types.PassReport, reportDir string, reportName string) {
	if len(passReport.Resources) != 0 {
		err := os.WriteFile(path.Join(reportDir, reportName), []byte(report.CleanupMarkdownReport(passReport)), 0644)
		if err != nil {
			logrus.Errorf("failed to save passed markdown report to %s: %+v", reportName, err)
		} else {
			logrus.Infof("markdown report saved to %s", reportName)
		}
	}
}

func storeCleanupErrorReport(errorReport types.ErrorReport, reportDir string) {
	for _, r := range errorReport.Errors {
		logrus.Warnf("found an error when deleting %s, address: %s\n", r.Type, r.Label)
		markdownFilename := fmt.Sprintf("%s_%s.md", strings.ReplaceAll(r.Type, "/", "_"), r.Label)
		err := os.WriteFile(path.Join(reportDir, markdownFilename), []byte(report.CleanupErrorMarkdownReport(r, errorReport.Logs)), 0644)
		if err != nil {
			logrus.Errorf("failed to save markdown report to %s: %+v", markdownFilename, err)
		} else {
			logrus.Infof("markdown report saved to %s", markdownFilename)
		}
	}
}
