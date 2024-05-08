package commands_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/azure/armstrong/commands"
)

func TestTestCommand_AllPass(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}
	runTestCommand(t, map[string]string{
		"all_passed_report.md": "Microsoft.Automation/automationAccounts@2023-11-01 (azapi_resource.automationAccount)",
	})
}

func TestTestCommand_Diff(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}
	runTestCommand(t, map[string]string{
		"partial_passed_report.md": "Microsoft.Resources/resourceGroups@2020-06-01 (azapi_resource.resourceGroup)",
		"Microsoft.Automation_automationAccounts@2023-11-01_automationAccount.md": ".properties.sku.name: expect Free, but got Basic",
	})
}

func TestTestCommand_MissingProperties(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}
	runTestCommand(t, map[string]string{
		"partial_passed_report.md": "Microsoft.Resources/resourceGroups@2020-06-01 (azapi_resource.resourceGroup)",
		"Microsoft.Automation_automationAccounts@2023-11-01_automationAccount.md": ".properties.sku.family = bar: not returned from response",
	})
}

func TestTestCommand_BadRequest(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC is not set")
	}
	runTestCommand(t, map[string]string{
		"partial_passed_report.md": "Microsoft.Resources/resourceGroups@2020-06-01 (azapi_resource.resourceGroup)",
		"Microsoft.Automation_automationAccounts@2023-11-01_automationAccount.md": "RESPONSE 400: 400 Bad Request",
	})
}

func runTestCommand(t *testing.T, fileContentMap map[string]string) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %+v", err)
	}

	tfDir := path.Join(wd, ".temp", t.Name())
	if err := os.MkdirAll(tfDir, 0755); err != nil {
		t.Fatalf("failed to create working directory: %+v", err)
	}
	defer os.RemoveAll(tfDir)
	defer commands.CleanupCommand{}.Run([]string{"-working-dir", tfDir})

	if err := copyFile(path.Join(wd, "testdata", t.Name(), "main.tf"), path.Join(tfDir, "main.tf")); err != nil {
		t.Fatalf("failed to copy file: %+v", err)
	}

	command := commands.TestCommand{}
	command.Run([]string{"-working-dir", tfDir})

	reportDir := tfDir
	dirs, err := os.ReadDir(reportDir)
	if err != nil {
		t.Fatalf("failed to read directory %s: %+v", reportDir, err)
	}
	for _, dir := range dirs {
		if dir.IsDir() && strings.HasPrefix(dir.Name(), "armstrong_reports_") {
			reportDir = path.Join(reportDir, dir.Name())
			break
		}
	}

	for file, content := range fileContentMap {
		filePath := path.Join(reportDir, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Fatalf("file %s not found", filePath)
		}
		if content != "" {
			data, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("failed to read file %s: %+v", filePath, err)
			}
			if !strings.Contains(string(data), content) {
				t.Fatalf("file %s does not contain expected content", filePath)
			}
		}
	}
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
