package commands

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/tf"
)

type ValidateCommand struct {
	Ui         cli.Ui
	workingDir string
}

func (c *ValidateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("validate")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
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
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	return c.Execute()
}

func (c ValidateCommand) Execute() int {
	wd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}
	if c.workingDir != "" {
		wd = c.workingDir
	}
	terraform, err := tf.NewTerraform(wd, true)
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

	_ = tf.GetChanges(plan)
	return 0
}
