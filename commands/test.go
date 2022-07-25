package commands

import (
	"flag"
	"fmt"
	"github.com/ms-henglu/armstrong/report"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/tf"
)

type TestCommand struct {
	Ui      cli.Ui
	verbose bool
}

func (command *TestCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("test")
	fs.BoolVar(&command.verbose, "v", false, "whether show terraform logs")
	fs.Usage = func() { command.Ui.Error(command.Help()) }
	return fs
}

func (command TestCommand) Help() string {
	helpText := `
Usage: armstrong test [-v]
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
