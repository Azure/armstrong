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

type CleanupCommand struct {
	Ui         cli.Ui
	verbose    bool
	workingDir string
}

func (c *CleanupCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("cleanup")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c CleanupCommand) Help() string {
	helpText := `
Usage: armstrong cleanup [-v] [-working-dir <path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c CleanupCommand) Synopsis() string {
	return "Clean up dependencies and testing resource"
}

func (c CleanupCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}

	log.Println("[INFO] ----------- cleanup resources ---------")
	wd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}
	if c.workingDir != "" {
		wd = c.workingDir
	}
	terraform, err := tf.NewTerraform(wd, c.verbose)
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
