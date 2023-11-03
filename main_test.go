package main_test

import (
	"github.com/ms-henglu/armstrong/commands"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"path"
	"testing"
)

func Test_One(t *testing.T) {
	input := "okUP"
	caser := cases.Title(language.Und, cases.NoLower)
	output := caser.String(input)
	log.Printf("output: %s", output)
}

func Test_RandomTest(t *testing.T) {
	t.Skip()
	specPath := "/Users/luheng/go/src/github.com/Azure/azure-rest-api-specs/specification"
	testcases := dfs(specPath)

	log.Printf("[INFO] found %d testcases", len(testcases))
	gen := commands.GenerateCommand{}
	for i, testcase := range testcases {
		log.Printf("[INFO] testcase %d: %s", i, testcase)

		if shouldSkip(testcase) {
			log.Printf("[INFO] skip %s", testcase)
			continue
		}

		output := path.Join(testcase, "terraform")
		err := os.RemoveAll(output)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}
		err = os.MkdirAll(output, 0755)
		if err != nil {
			log.Fatalf("[ERROR] %v", err)
		}

		args := []string{"-swagger", testcase, "-working-dir", output}
		log.Printf("[INFO] args: %v", args)
		gen.Run(args)
	}
}

func shouldSkip(workingDirectory string) bool {
	folder := path.Join(workingDirectory, "terraform")
	cases, _ := os.ReadDir(folder)
	for _, c := range cases {
		if c.IsDir() {
			files, _ := os.ReadDir(path.Join(folder, c.Name()))
			for _, f := range files {
				if f.Name() == "main.tf" {
					return true
				}
			}
		}
	}
	return false
}

func dfs(dir string) []string {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		log.Printf("[ERROR] reading %s", dir)
		return []string{}
	}
	out := make([]string, 0)
	found := false
	for _, d := range dirs {
		if d.IsDir() && d.Name() == "examples" {
			found = true
			break
		}
	}
	if found {
		out = append(out, dir)
	}

	for _, d := range dirs {
		if d.Name() == "data-plane" {
			continue
		}
		if d.IsDir() {
			out = append(out, dfs(path.Join(dir, d.Name()))...)
		}
	}

	return out
}

/*

TODO:

1.



TODO list:
1. Support extension scope resource
Ex:
	1. /{resourceUri}/providers/Microsoft.Advisor/recommendations/{recommendationId}
	2. /{roleId}

2. Key segment is a variable
Ex:
	1. /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Web/sites/{name}/host/default/{keyType}/{keyName}
	{keyType} is an enum

3. Body with all configurable properties (when the example is invalid/empty)

4. Support POST to create a new resource
Ex:
	1. /providers/Microsoft.ADHybridHealthService/addsservices

5. Support the API paths that don't follow the guideline: give warning message or handle them in a special way
	1. "/providers/Microsoft.ADHybridHealthService/reports/DevOps/IsDevOps"
		segment is not a placeholder
	2. /{roleId}
	3. /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupJobs/operationResults/{operationId}"
		swagger issue: resource ID api path should have even number of segments


*/
