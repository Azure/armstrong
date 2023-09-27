package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ms-henglu/armstrong/tf"
	"github.com/sirupsen/logrus"
)

type ValidateCommand struct {
	verbose    bool
	workingDir string
}

func (c *ValidateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("validate")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { logrus.Error(c.Help()) }
	return fs
}

func (c ValidateCommand) Help() string {
	helpText := `
Usage: armstrong validate [-v] [-working-dir <path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c ValidateCommand) Synopsis() string {
	return "Generates a speculative execution plan, showing what actions Terraform would take to apply the current configuration."
}

func (c ValidateCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		logrus.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	if c.verbose {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Infof("verbose mode enabled")
	}
	return c.Execute()
}

func (c ValidateCommand) Execute() int {
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
	terraform, err := tf.NewTerraform(wd, true)
	if err != nil {
		logrus.Fatalf("creating terraform executable: %+v\n", err)
	}

	logrus.Infof("running terraform init...")
	_ = terraform.Init()

	logrus.Infof("running terraform plan to check the changes...")
	plan, err := terraform.Plan()
	if err != nil {
		logrus.Fatalf("running terraform plan: %+v", err)
	}

	_ = tf.GetChanges(plan)
	return 0
}
