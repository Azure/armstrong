package commands

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-testing-tool/tf"
	"github.com/nsf/jsondiff"
)

type TestCommand struct {
	Ui      cli.Ui
	verbose bool
}

func (command *TestCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("version")
	fs.BoolVar(&command.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { command.Ui.Error(command.Help()) }
	return fs
}

func (command TestCommand) Help() string {
	helpText := `
Usage: azurerm-rest-api-testing-tool test [-v]
` + command.Synopsis() + "\n\n" + helpForFlags(command.flags())

	return strings.TrimSpace(helpText)
}

func (command TestCommand) Synopsis() string {
	return "Update dependencies for tests and run tests"
}

func (command TestCommand) Run(args []string) int {
	f := command.flags()
	if err := f.Parse(args); err != nil {
		command.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	log.Println("[INFO] ----------- run tests ---------")
	terraform, err := tf.NewTerraform(command.verbose)
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
	c, r, u, d := 0, 0, 0, 0
	for _, action := range actions {
		switch action {
		case tf.ActionCreate:
			c++
		case tf.ActionReplace:
			r++
		case tf.ActionUpdate:
			u++
		case tf.ActionDelete:
			d++
		}
	}
	log.Printf("[INFO] found %d changes in total, create: %d, replace: %d, update: %d, delete: %d\n", c+r+u+d, c, r, u, d)

	log.Println("[INFO] running apply command to provision test resource...")
	err = terraform.Apply()
	if err != nil {
		log.Fatalf("[Error] error running terraform apply: %+v\n", err)
	}
	log.Println("[INFO] test resource has been provisioned")

	log.Println("[INFO] running plan command to verify test resource...")
	plan, err = terraform.Plan()
	if err != nil {
		log.Fatalf("[Error] error running terraform plan: %+v\n", err)
	}

	if len(tf.GetChanges(plan)) == 0 {
		log.Println("[INFO] Test passed!")
		return 0
	}

	before, after := tf.GetBodyChange(plan)
	option := jsondiff.DefaultConsoleOptions()
	_, msg := jsondiff.Compare([]byte(before), []byte(after), &option)
	log.Printf("[INFO] found differences between response and configuration:\n%s", msg)
	option = jsondiff.Options{
		Added:                 jsondiff.Tag{Begin: "\033[0;32m", End: " is not returned from response\033[0m"},
		Removed:               jsondiff.Tag{Begin: "\033[0;31m", End: "\033[0m"},
		Changed:               jsondiff.Tag{Begin: "\033[0;33m Got ", End: "\033[0m"},
		Skipped:               jsondiff.Tag{Begin: "\033[0;90m", End: "\033[0m"},
		SkippedArrayElement:   jsondiff.SkippedArrayElement,
		SkippedObjectProperty: jsondiff.SkippedObjectProperty,
		ChangedSeparator:      " in response, expect ",
		Indent:                "    ",
	}
	_, msg = jsondiff.Compare([]byte(before), []byte(after), &option)
	log.Fatalf("[INFO] report:\n%s", msg)
	return 1
}
