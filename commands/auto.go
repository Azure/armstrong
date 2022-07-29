package commands

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mitchellh/cli"
)

type AutoCommand struct {
	Ui                cli.Ui
	path              string
	workingDir        string
	verbose           bool
	useRawJsonPayload bool
	overwrite         bool
}

func (c *AutoCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("auto")
	fs.StringVar(&c.path, "path", "", "filepath of rest api to create arm resource example")
	fs.StringVar(&c.workingDir, "working-dir", "", "output path to Terraform configuration files")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.BoolVar(&c.useRawJsonPayload, "raw", false, "whether use raw json payload in 'body'")
	fs.BoolVar(&c.overwrite, "overwrite", false, "whether overwrite existing terraform configurations")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}
func (c AutoCommand) Help() string {
	helpText := `
Usage: armstrong auto -path <path to a swagger 'Create' example> [-v] [-working-dir <output path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c AutoCommand) Synopsis() string {
	return "Run generate and test, if test passed, run cleanup"
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
	args = make([]string, 0)
	if c.verbose {
		args = append(args, "-v")
	}
	TestCommand{Ui: c.Ui}.Run(args)
	CleanupCommand{Ui: c.Ui}.Run(args)
	log.Println("[INFO] Test passed!")
	return 0
}
