package commands

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

type CleanupCommand struct {
	Ui      cli.Ui
	verbose bool
}

func (c *CleanupCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("version")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c CleanupCommand) Help() string {
	helpText := `
Usage: azurerm-rest-api-testing-tool cleanup [-v]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c CleanupCommand) Synopsis() string {
	return "Clean up dependency"
}

func (c CleanupCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	log.Println("[INFO] ----------- cleanup resources ---------")
	terraform, err := tf.NewTerraform(c.verbose)
	if err != nil {
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}
	log.Println("[INFO] prepare working directory")
	_ = terraform.Init()
	log.Println("[INFO] running destroy command to cleanup resources...")
	err = terraform.Destroy()
	if err != nil {
		log.Fatalf("[Error] error cleaning up resources: %+v\n", err)
	}
	log.Println("[INFO] all resources have been deleted")
	return 0
}
