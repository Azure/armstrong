package commands

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/cli"
)

type AutoCommand struct {
	Ui      cli.Ui
	path    string
	verbose bool
}

func (c *AutoCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("generate")
	fs.StringVar(&c.path, "path", "", "filepath of rest api to create arm resource example")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}
func (c AutoCommand) Help() string {
	helpText := `
Usage: azurerm-rest-api-testing-tool auto -path <filepath to example> [-v]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c AutoCommand) Synopsis() string {
	return "Run generate and test"
}

func (c AutoCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	if len(c.path) == 0 {
		c.Ui.Error(c.Help())
		return 1
	}
	GenerateCommand{Ui: c.Ui}.Run(args)
	TestCommand{Ui: c.Ui}.Run(args)
	CleanupCommand{Ui: c.Ui}.Run(args)
	log.Println("[INFO] Test passed!")
	return 0
}
