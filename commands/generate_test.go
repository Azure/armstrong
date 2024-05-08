package commands_test

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/azure/armstrong/commands"
	"github.com/azure/armstrong/resource"
	"github.com/azure/armstrong/tf"
)

func TestGenerateCommand_multiple(t *testing.T) {
	runGenerateCommand(t, [][]string{
		{"-type", "data", "-path", "DatabaseSqlVulnerabilityAssessmentGet.json"},
		{"-path", "DatabaseSqlVulnerabilityAssessmentBaselineAdd.json"},
	})
}

func TestGenerateCommand_identity(t *testing.T) {
	runGenerateCommand(t, [][]string{
		{"-path", "ConfigurationStoresCreateWithIdentity.json"},
	})
}

func TestGenerateCommand_fromSwagger(t *testing.T) {
	runGenerateCommand(t, [][]string{
		{"-swagger", "purview.json"},
	})
}

func runGenerateCommand(t *testing.T, args [][]string) {
	resource.R = rand.New(rand.NewSource(0))

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %+v", err)
	}

	tfDir := path.Join(wd, ".temp", t.Name())
	if err := os.MkdirAll(tfDir, 0755); err != nil {
		t.Fatalf("failed to create working directory: %+v", err)
	}
	defer os.RemoveAll(tfDir)

	command := commands.GenerateCommand{}

	for _, arg := range args {
		for i, a := range arg {
			if a == "-path" || a == "-swagger" {
				arg[i+1] = path.Join(wd, "testdata", t.Name(), arg[i+1])
			}
		}
		res := command.Run(append([]string{"-working-dir", tfDir}, arg...))
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

	expectDir := path.Join(wd, "testdata", t.Name(), "expect")
	if err := requireFoldersEqual(expectDir, tfDir); err != nil {
		t.Fatal(err)
	}
}

func requireFoldersEqual(path1, path2 string) error {
	files, err := os.ReadDir(path1)
	if err != nil {
		return err
	}
	for _, f := range files {
		if f.IsDir() {
			if err := requireFoldersEqual(path.Join(path1, f.Name()), path.Join(path2, f.Name())); err != nil {
				return err
			}
		} else {
			if f.Name() == "provider.tf" || !strings.HasSuffix(f.Name(), ".tf") {
				continue
			}
			actual, err := os.ReadFile(path.Join(path2, f.Name()))
			if err != nil {
				return err
			}
			expect, err := os.ReadFile(path.Join(path1, f.Name()))
			if err != nil {
				return err
			}
			if string(expect) != string(actual) {
				return fmt.Errorf("expected %s to be equal to %s", path.Join(path1, f.Name()), path.Join(path2, f.Name()))
			}
		}
	}
	return nil
}
