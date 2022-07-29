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

func (command *ValidateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("validate")
	fs.StringVar(&command.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { command.Ui.Error(command.Help()) }
	return fs
}

func (command ValidateCommand) Help() string {
	helpText := `
Usage: armstrong validate [-v] [-working-dir <path to Terraform configuration files>]
` + command.Synopsis() + "\n\n" + helpForFlags(command.flags())

	return strings.TrimSpace(helpText)
}

func (command ValidateCommand) Synopsis() string {
	return "Generates a speculative execution plan, showing what actions Terraform would take to apply the current configuration."
}

func (command ValidateCommand) Run(args []string) int {
	f := command.flags()
	if err := f.Parse(args); err != nil {
		command.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	wd, err := os.Getwd()
	if err != nil {
		command.Ui.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}
	if command.workingDir != "" {
		wd = command.workingDir
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
