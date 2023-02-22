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

type TestCommand struct {
	Ui         cli.Ui
	verbose    bool
	workingDir string
}

func (c *TestCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("test")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c TestCommand) Help() string {
	helpText := `
Usage: armstrong test [-v] [-working-dir <path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c TestCommand) Synopsis() string {
	return "Update dependencies for tests and run tests"
}

func (c TestCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	return c.Execute()
}

func (c TestCommand) Execute() int {
	log.Println("[INFO] ----------- run tests ---------")
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
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}

	log.Printf("[INFO] prepare working directory\n")
	_ = terraform.Init()

	log.Println("[INFO] running plan command to check changes...")
	plan, err := terraform.Plan()
	if err != nil {
		log.Fatalf("[Error] error running terraform plan: %+v\n", err)
	}

	actions := tf.GetChanges(plan)
	create, replace, update, delete := 0, 0, 0, 0
	for _, action := range actions {
		switch action {
		case tf.ActionCreate:
			create++
		case tf.ActionReplace:
			replace++
		case tf.ActionUpdate:
			update++
		case tf.ActionDelete:
			delete++
		}
	}
	log.Printf("[INFO] found %d changes in total, create: %d, replace: %d, update: %d, delete: %d\n", create+replace+update+delete, create, replace, update, delete)
	log.Println("[INFO] running apply command to provision test resource...")
	applyErr := terraform.Apply()
	if applyErr != nil {
		log.Printf("[Error] error running terraform apply: %+v\n", applyErr)
	} else {
		log.Println("[INFO] test resource has been provisioned")
	}

	log.Println("[INFO] running plan command to verify test resource...")
	plan, err = terraform.Plan()
	if err != nil {
		log.Fatalf("[Error] error running terraform plan: %+v\n", err)
	}

	reportDir := fmt.Sprintf("armstrong_reports_%s", time.Now().Format(time.Stamp))
	reportDir = strings.ReplaceAll(reportDir, ":", "")
	reportDir = strings.ReplaceAll(reportDir, " ", "_")
	reportDir = path.Join(wd, reportDir)
	err = os.Mkdir(reportDir, 0755)
	if err != nil {
		log.Fatalf("[Error] error creating report dir %s: %+v", reportDir, err)
	}

	if applyErr == nil && len(tf.GetChanges(plan)) == 0 {
		if state, err := terraform.Show(); err == nil {
			passReport := tf.NewPassReportFromState(state)
			storePassReport(passReport, reportDir, "all_passed_report.md")
			log.Printf("[INFO] %d resources passed the tests.", len(passReport.Resources))
			log.Printf("[INFO] all reports have been saved in the report directory: %s, please check.", reportDir)
		} else {
			log.Fatalf("[Error] error showing terraform state: %+v", err)
		}
		return 0
	}

	logs, err := report.ParseLogs(path.Join(wd, "log.txt"))
	if err != nil {
		log.Printf("[ERROR] parsing log.txt: %+v", err)
	}

	errorReport := types.ErrorReport{}
	if applyErr != nil {
		errorReport = tf.NewErrorReport(applyErr, logs)
		storeErrorReport(errorReport, reportDir)
	}

	diffReport := tf.NewDiffReport(plan, logs)
	storeDiffReport(diffReport, reportDir)

	passReport := tf.NewPassReport(plan)
	storePassReport(passReport, reportDir, "partially_passed_report.md")

	log.Println("[INFO] ---------------- Summary ----------------")
	log.Printf("[INFO] %d resources passed the tests.", len(passReport.Resources))
	if len(errorReport.Errors) != 0 {
		log.Printf("[INFO] %d errors when creating the testing resources.", len(errorReport.Errors))
	}
	if len(diffReport.Diffs) != 0 {
		log.Printf("[INFO] %d API issues.", len(diffReport.Diffs))
	}
	log.Printf("[INFO] all reports have been saved in the report directory: %s, please check.", reportDir)
	return 1
}

func storePassReport(passReport types.PassReport, reportDir string, reportName string) {
	if len(passReport.Resources) != 0 {
		err := os.WriteFile(path.Join(reportDir, reportName), []byte(report.PassedMarkdownReport(passReport)), 0644)
		if err != nil {
			log.Printf("[WARN] failed to save passed markdown report to %s: %+v", reportName, err)
		} else {
			log.Printf("[INFO] markdown report saved to %s", reportName)
		}
	}
}

func storeErrorReport(errorReport types.ErrorReport, reportDir string) {
	for _, r := range errorReport.Errors {
		log.Printf("[WARN] found an error when create %s, address: azapi_resource.%s\n", r.Type, r.Label)
		markdownFilename := fmt.Sprintf("%s_%s.md", strings.ReplaceAll(r.Type, "/", "_"), r.Label)
		err := os.WriteFile(path.Join(reportDir, markdownFilename), []byte(report.ErrorMarkdownReport(r, errorReport.Logs)), 0644)
		if err != nil {
			log.Printf("[WARN] failed to save markdown report to %s: %+v", markdownFilename, err)
		} else {
			log.Printf("[INFO] markdown report saved to %s", markdownFilename)
		}
	}
}

func storeDiffReport(diffReport types.DiffReport, reportDir string) {
	for _, r := range diffReport.Diffs {
		log.Printf("[WARN] found differences between response and configuration:\n\naddress: %s\n\n%s\n", r.Address, report.DiffMessageTerraform(r.Change))
		log.Printf("[INFO] report:\n\naddresss: %s\t%s\n", r.Address, report.DiffMessageReadable(r.Change))
		markdownFilename := fmt.Sprintf("%s_%s.md", strings.ReplaceAll(r.Type, "/", "_"), strings.TrimPrefix(r.Address, "azapi_resource."))
		err := os.WriteFile(path.Join(reportDir, markdownFilename), []byte(report.DiffMarkdownReport(r, diffReport.Logs)), 0644)
		if err != nil {
			log.Printf("[WARN] failed to save markdown report to %s: %+v", markdownFilename, err)
		} else {
			log.Printf("[INFO] markdown report saved to %s", markdownFilename)
		}
	}
}
