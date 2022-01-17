package commands

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-restapi-testing-tool/tf"
)

type SetupCommand struct {
	Ui      cli.Ui
	verbose bool
}

func (c *SetupCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("version")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c SetupCommand) Help() string {
	helpText := `
Usage: azurerm-rest-api-testing-tool setup [-v]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c SetupCommand) Synopsis() string {
	return "Update dependencies for tests"
}

func (c SetupCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	log.Println("[INFO] ----------- update resources ---------")
	terraform, err := tf.NewTerraform(c.verbose)
	if err != nil {
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}
	log.Printf("[INFO] prepare working directory\n")
	_ = terraform.Init()
	log.Println("[INFO] running apply command to update dependency resources...")
	err = terraform.Apply()
	if err != nil {
		log.Fatalf("[Error] error setting up resources: %+v\n", err)
	}
	log.Println("[INFO] all dependencies have been updated")
	return 0
}
