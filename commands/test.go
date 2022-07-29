package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/report"
	"github.com/ms-henglu/armstrong/tf"
)

type TestCommand struct {
	Ui         cli.Ui
	verbose    bool
	workingDir string
}

func (c *TestCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("test")
	fs.BoolVar(&c.verbose, "v", false, "whether show terraform logs")
	fs.StringVar(&c.workingDir, "working-dir", "", "path to Terraform configuration files")
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	return fs
}

func (c TestCommand) Help() string {
	helpText := `
Usage: armstrong test [-v] [-working-dir <path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c TestCommand) Synopsis() string {
	return "Update dependencies for tests and run tests"
}

func (c TestCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	return c.Execute()
}

func (c TestCommand) Execute() int {
	log.Println("[INFO] ----------- run tests ---------")
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

	log.Printf("[INFO] prepare working directory\n")
	_ = terraform.Init()

	log.Println("[INFO] running plan command to check changes...")
	plan, err := terraform.Plan()
	if err != nil {
		log.Fatalf("[Error] error running terraform plan: %+v\n", err)
	}

	actions := tf.GetChanges(plan)
	create, replace, update, delete := 0, 0, 0, 0
	for _, action := range actions {
		switch action {
		case tf.ActionCreate:
			create++
		case tf.ActionReplace:
			replace++
		case tf.ActionUpdate:
			update++
		case tf.ActionDelete:
			delete++
		}
	}
	log.Printf("[INFO] found %d changes in total, create: %d, replace: %d, update: %d, delete: %d\n", create+replace+update+delete, create, replace, update, delete)
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

	reports := tf.NewReports(plan)
	logs, err := report.ParseLogs("./log.txt")
	if err != nil {
		log.Printf("[ERROR] parsing log.txt: %+v", err)
	}
	for _, r := range reports {
		log.Printf("[INFO] found differences between response and configuration:\n\naddress: %s\n\n%s\n",
			r.Address, report.DiffMessageTerraform(r.Change))
		log.Printf("[INFO] report:\n\naddresss: %s\n\n%s\n", r.Address, report.DiffMessageReadable(r.Change))
		markdownFilename := fmt.Sprintf("%s_%s.md", strings.ReplaceAll(r.Type, "/", "_"), time.Now().Format("20060102030405PM"))
		err := ioutil.WriteFile(markdownFilename, []byte(report.MarkdownReport(r, logs)), 0644)
		if err != nil {
			log.Printf("[WARN] failed to save markdown report to %s: %+v", markdownFilename, err)
		} else {
			log.Printf("[INFO] markdown report saved to %s", markdownFilename)
		}
	}

	log.Fatalf("[ERROR] found %v API issues", len(reports))
	return 1
}
