package commands_test

import (
	"os"
	"path"
	"testing"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/commands"
	"github.com/ms-henglu/armstrong/tf"
)

func TestGenerateCommand_multiple(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %+v", err)
	}
	runGenerateCommand([][]string{
		{"-type", "data", "-path", path.Join(wd, "testdata", "case0", "DatabaseSqlVulnerabilityAssessmentGet.json")},
		{"-path", path.Join(wd, "testdata", "case0", "DatabaseSqlVulnerabilityAssessmentBaselineAdd.json")},
	}, t)
}

func TestGenerateCommand_identity(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %+v", err)
	}
	runGenerateCommand([][]string{
		{"-path", path.Join(wd, "testdata", "case1", "ConfigurationStoresCreateWithIdentity.json")},
	}, t)
}

func runGenerateCommand(args [][]string, t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %+v", err)
	}

	wd = path.Join(wd, ".temp", t.Name())
	if err := os.MkdirAll(wd, 0755); err != nil {
		t.Fatalf("failed to create working directory: %+v", err)
	}
	defer os.RemoveAll(wd)

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			Reader:      os.Stdin,
			ErrorWriter: os.Stderr,
		},
	}
	command := commands.GenerateCommand{Ui: ui}

	for _, arg := range args {
		res := command.Run(append([]string{"-working-dir", wd}, arg...))
		if res != 0 {
			t.Fatalf("failed to generate terraform configuration")
		}
	}

	terraform, err := tf.NewTerraform(wd, true)
	if err != nil {
		t.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}

	if err := terraform.Init(); err != nil {
		t.Fatalf("[Error] error initializing terraform configuration: %+v\n", err)
	}

	out, err := terraform.Validate()
	if err != nil {
		t.Fatalf("[Error] error validating terraform configuration: %+v\n", err)
	}
	if out.ErrorCount != 0 {
		t.Fatalf("[Error] terraform configuration is not valid: %+v\n", out)
	}
}
