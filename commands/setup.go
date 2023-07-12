package commands

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/tf"
)

type SetupCommand struct {
	Ui         cli.Ui
	verbose    bool
	workingDir string
}

func (c *SetupCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("setup")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c SetupCommand) Help() string {
	helpText := `
Usage: armstrong setup [-v] [-working-dir <path to Terraform configuration files>]
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
	return c.Execute()
}

func (c SetupCommand) Execute() int {
	log.Println("[INFO] ----------- update resources ---------")
	wd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}
	if c.workingDir != "" {
		wd, err = filepath.Abs(c.workingDir)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("working directory is invalid: %+v", err))
			return 1
		}
	}
	terraform, err := tf.NewTerraform(wd, c.verbose)
	if err != nil {
		log.Fatalf("[ERROR] error creating terraform executable: %+v\n", err)
	}
	log.Printf("[INFO] prepare working directory\n")
	_ = terraform.Init()
	log.Println("[INFO] running apply command to update dependency resources...")
	err = terraform.Apply()
	if err != nil {
		log.Fatalf("[ERROR] error setting up resources: %+v\n", err)
	}
	log.Println("[INFO] all dependencies have been updated")
	return 0
}
