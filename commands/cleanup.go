package commands

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/report"
	"github.com/ms-henglu/armstrong/tf"
	"github.com/ms-henglu/armstrong/types"
)

type CleanupCommand struct {
	Ui         cli.Ui
	verbose    bool
	workingDir string
}

func (c *CleanupCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("cleanup")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
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
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	return c.Execute()
}

func (c CleanupCommand) Execute() int {
	const (
		allPassedReportFileName     = "cleanup_all_passed_report.md"
		partialPassedReportFileName = "cleanup_partial_passed_report.md"
	)

	log.Println("[INFO] ----------- cleanup resources ---------")
	wd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}
	if c.workingDir != "" {
		wd, err = filepath.Abs(c.workingDir)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("working directory is invalid: %+v", err))
			return 1
		}
	}
	terraform, err := tf.NewTerraform(wd, c.verbose)
	if err != nil {
		log.Fatalf("[ERROR] error creating terraform executable: %+v\n", err)
	}
	state, err := terraform.Show()
	if err != nil {
		log.Fatalf("[ERROR] error getting state: %+v\n", err)
	}

	passReport := tf.NewPassReportFromState(state)

	reportDir := fmt.Sprintf("armstrong_cleanup_reports_%s", time.Now().Format(time.Stamp))
	reportDir = strings.ReplaceAll(reportDir, ":", "")
	reportDir = strings.ReplaceAll(reportDir, " ", "_")
	reportDir = path.Join(wd, reportDir)
	err = os.Mkdir(reportDir, 0755)
	if err != nil {
		log.Fatalf("[ERROR] error creating report dir %s: %+v", reportDir, err)
	}

	log.Println("[INFO] prepare working directory")
	_ = terraform.Init()
	log.Println("[INFO] running destroy command to cleanup resources...")
	destroyErr := terraform.Destroy()
	if destroyErr != nil {
		log.Printf("[ERROR] error cleaning up resources: %+v\n", destroyErr)
	} else {
		log.Println("[INFO] all resources are cleaned up")
		storeCleanupReport(passReport, reportDir, allPassedReportFileName)
	}

	logs, err := report.ParseLogs(path.Join(wd, "log.txt"))
	if err != nil {
		log.Printf("[ERROR] parsing log.txt: %+v", err)
	}

	errorReport := types.ErrorReport{}
	if destroyErr != nil {
		errorReport := tf.NewCleanupErrorReport(destroyErr, logs)
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

	passReport.Resources = resources
	storeCleanupReport(passReport, reportDir, partialPassedReportFileName)

	log.Println("[INFO] ---------------- Summary ----------------")
	log.Printf("[INFO] %d resources passed the cleanup tests.", len(passReport.Resources))
	if len(errorReport.Errors) != 0 {
		log.Printf("[INFO] %d errors when cleanup the testing resources.", len(errorReport.Errors))
	}

	return 0
}

func storeCleanupReport(passReport types.PassReport, reportDir string, reportName string) {
	if len(passReport.Resources) != 0 {
		err := os.WriteFile(path.Join(reportDir, reportName), []byte(report.CleanupMarkdownReport(passReport)), 0644)
		if err != nil {
			log.Printf("[WARN] failed to save passed markdown report to %s: %+v", reportName, err)
		} else {
			log.Printf("[INFO] markdown report saved to %s", reportName)
		}
	}
}

func storeCleanupErrorReport(errorReport types.ErrorReport, reportDir string) {
	for _, r := range errorReport.Errors {
		log.Printf("[WARN] found an error when deleting %s, address: azapi_resource.%s\n", r.Type, r.Label)
		markdownFilename := fmt.Sprintf("%s_%s.md", strings.ReplaceAll(r.Type, "/", "_"), r.Label)
		err := os.WriteFile(path.Join(reportDir, markdownFilename), []byte(report.CleanupErrorMarkdownReport(r, errorReport.Logs)), 0644)
		if err != nil {
			log.Printf("[WARN] failed to save markdown report to %s: %+v", markdownFilename, err)
		} else {
			log.Printf("[INFO] markdown report saved to %s", markdownFilename)
		}
	}
}
