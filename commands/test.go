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

	"github.com/ms-henglu/armstrong/coverage"
	"github.com/ms-henglu/armstrong/report"
	"github.com/ms-henglu/armstrong/tf"
	"github.com/ms-henglu/armstrong/types"
	"github.com/ms-henglu/armstrong/utils"
	"github.com/ms-henglu/pal/formatter"
	"github.com/ms-henglu/pal/trace"
	paltypes "github.com/ms-henglu/pal/types"
	"github.com/sirupsen/logrus"
)

type TestCommand struct {
	verbose          bool
	workingDir       string
	destroyAfterTest bool
	swaggerPath      string
}

func (c *TestCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("test")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.BoolVar(&c.destroyAfterTest, "destroy-after-test", false, "whether to destroy the created resources after each test")
	fs.StringVar(&c.swaggerPath, "swagger", "", "path to the .json swagger which is being test")
	fs.Usage = func() { logrus.Error(c.Help()) }
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
		logrus.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	if c.verbose {
		log.SetOutput(os.Stdout)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Infof("verbose mode enabled")
	}
	return c.Execute()
}

func (c TestCommand) Execute() int {
	const (
		allPassedReportFileName     = "all_passed_report.md"
		partialPassedReportFileName = "partial_passed_report.md"
	)

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
		logrus.Fatalf("error creating terraform executable: %+v\n", err)
	}

	logrus.Infof("prepare working directory\n")
	_ = terraform.Init()

	logrus.Infof("running plan command to check changes...")
	plan, err := terraform.Plan()
	if err != nil {
		logrus.Fatalf("error running terraform plan: %+v\n", err)
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
	logrus.Infof("found %d changes in total, create: %d, replace: %d, update: %d, delete: %d\n", create+replace+update+delete, create, replace, update, delete)
	logrus.Infof("running apply command to provision test resource...")
	applyErr := terraform.Apply()
	if applyErr != nil {
		logrus.Errorf("error running terraform apply: %+v\n", applyErr)
	} else {
		logrus.Infof("test resource has been provisioned")
	}

	logrus.Infof("running plan command to verify test resource...")
	plan, err = terraform.Plan()
	if err != nil {
		logrus.Fatalf("error running terraform plan: %+v\n", err)
	}

	reportDir := fmt.Sprintf("armstrong_reports_%s", time.Now().Format(time.Stamp))
	reportDir = strings.ReplaceAll(reportDir, ":", "")
	reportDir = strings.ReplaceAll(reportDir, " ", "_")
	reportDir = path.Join(wd, reportDir)
	logrus.Infof("creating report directory %s\n", reportDir)
	err = os.Mkdir(reportDir, 0755)
	if err != nil {
		logrus.Fatalf("error creating report dir %s: %+v", reportDir, err)
	}

	logrus.Infof("parsing log.txt...")
	logs, err := trace.RequestTracesFromFile(path.Join(wd, "log.txt"))
	if err != nil {
		logrus.Errorf("parsing log.txt: %+v", err)
	}

	logrus.Infof("generating reports...")
	var passReport types.PassReport
	if applyErr == nil && len(tf.GetChanges(plan)) == 0 {
		if state, err := terraform.Show(); err == nil {
			passReport = tf.NewPassReportFromState(state)
			coverageReport, err := tf.NewCoverageReportFromState(state, c.swaggerPath)
			if err != nil {
				logrus.Errorf("error producing coverage report: %+v", err)
			}
			storePassReport(passReport, coverageReport, reportDir, allPassedReportFileName)
		} else {
			logrus.Fatalf("error showing terraform state: %+v", err)
		}
	} else {
		passReport = tf.NewPassReport(plan)
		coverageReport, err := tf.NewCoverageReport(plan, c.swaggerPath)
		if err != nil {
			logrus.Errorf("error producing coverage report: %+v", err)
		}
		storePassReport(passReport, coverageReport, reportDir, partialPassedReportFileName)
	}

	errorReport := tf.NewErrorReport(applyErr, logs)
	storeErrorReport(errorReport, reportDir)

	diffReport := tf.NewDiffReport(plan, logs)
	storeDiffReport(diffReport, reportDir)

	if applyErr == nil && c.destroyAfterTest {
		logrus.Infof("running destroy command to delete resources...")
		destroyErr := terraform.Destroy()
		if destroyErr != nil {
			logrus.Errorf("error running terraform destroy: %+v\n", destroyErr)
		} else {
			logrus.Infof("test resource has been deleted")
		}
	} else {
		logrus.Warnf("the created resources will not be destroyed because either there is an error or destroy-after-test flag is not set")
	}

	logrus.Infof("generating traces...")
	traceDir := path.Join(wd, "traces")
	if !utils.Exists(traceDir) {
		err = os.Mkdir(traceDir, 0755)
		if err != nil {
			logrus.Errorf("error creating trace dir %s: %+v", traceDir, err)
		}
	}

	storeOavTraffic(logs, traceDir)
	logrus.Infof("copying traces to report directory...")
	if err := utils.Copy(traceDir, path.Join(reportDir, "traces")); err != nil {
		logrus.Errorf("error copying traces: %+v", err)
	}

	if c.swaggerPath != "" {
		logrus.Infof("generating swagger accuracy report...")
		if _, err = report.OavValidateTraffic(traceDir, c.swaggerPath, reportDir); err != nil {
			logrus.Errorf("error storing swagger accuracy report: %+v", err)
		}
	} else {
		logrus.Warnf("no swagger file provided, swagger accuracy report will not be generated")
	}

	logrus.Infof("---------------- Summary ----------------")
	logrus.Infof("%d resources passed the tests.", len(passReport.Resources))
	if len(errorReport.Errors) != 0 {
		logrus.Infof("%d errors when creating the testing resources.", len(errorReport.Errors))
	}
	if len(diffReport.Diffs) != 0 {
		logrus.Infof("%d API issues.", len(diffReport.Diffs))
	}
	logrus.Infof("all reports have been saved in the report directory: %s, please check.", reportDir)
	return 0
}

func storePassReport(passReport types.PassReport, coverageReport coverage.CoverageReport, reportDir string, reportName string) {
	if len(passReport.Resources) != 0 {
		err := os.WriteFile(path.Join(reportDir, reportName), []byte(report.PassedMarkdownReport(passReport, coverageReport)), 0644)
		if err != nil {
			logrus.Warnf("failed to save passed markdown report to %s: %+v", reportName, err)
		} else {
			logrus.Infof("markdown report saved to %s", reportName)
		}
	}
}

func storeErrorReport(errorReport types.ErrorReport, reportDir string) {
	for _, r := range errorReport.Errors {
		logrus.Warnf("found an error when creating %s, address: azapi_resource.%s\n", r.Type, r.Label)
		markdownFilename := fmt.Sprintf("%s_%s.md", strings.ReplaceAll(r.Type, "/", "_"), r.Label)
		err := os.WriteFile(path.Join(reportDir, markdownFilename), []byte(report.ErrorMarkdownReport(r, errorReport.Logs)), 0644)
		if err != nil {
			logrus.Warnf("failed to save markdown report to %s: %+v", markdownFilename, err)
		} else {
			logrus.Infof("markdown report saved to %s", markdownFilename)
		}
	}
}

func storeDiffReport(diffReport types.DiffReport, reportDir string) {
	for _, r := range diffReport.Diffs {
		logrus.Warnf("found differences between response and configuration:\n\naddress: %s\n\n%s\n", r.Address, report.DiffMessageTerraform(r.Change))
		logrus.Infof("report:\n\naddresss: %s\t%s\n", r.Address, report.DiffMessageReadable(r.Change))
		markdownFilename := fmt.Sprintf("%s_%s.md", strings.ReplaceAll(r.Type, "/", "_"), strings.TrimPrefix(r.Address, "azapi_resource."))
		err := os.WriteFile(path.Join(reportDir, markdownFilename), []byte(report.DiffMarkdownReport(r, diffReport.Logs)), 0644)
		if err != nil {
			logrus.Warnf("failed to save markdown report to %s: %+v", markdownFilename, err)
		} else {
			logrus.Infof("markdown report saved to %s", markdownFilename)
		}
	}
}

func storeOavTraffic(traces []paltypes.RequestTrace, output string) {
	format := formatter.OavTrafficFormatter{}
	files, err := os.ReadDir(output)
	if err != nil {
		logrus.Warnf("failed to read trace output directory: %v", err)
	}
	index := len(files)
	for _, t := range traces {
		out := format.Format(t)
		index = index + 1
		outputPath := path.Join(output, fmt.Sprintf("trace-%d.json", index))
		if err := os.WriteFile(outputPath, []byte(out), 0644); err != nil {
			logrus.Warnf("failed to write file: %v", err)
		} else {
			logrus.Debugf("trace saved to %s", outputPath)
		}
	}
}
