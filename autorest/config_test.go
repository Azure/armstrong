package autorest

import (
	"log"
	"os"
	"path"
	"testing"
)

func Test_ParseAutoRestConfig(t *testing.T) {
	wd, _ := os.Getwd()
	testcases := []struct {
		Input    string
		Expected []Package
	}{
		{
			Input: path.Join(wd, "testdata", "readme.md"),
			Expected: []Package{
				{
					Tag: "package-2015-10",
				},
				{
					Tag: "package-2017-05-preview",
				},
				{
					Tag: "package-2018-01-preview",
				},
				{
					Tag: "package-2018-06-preview",
				},
				{
					Tag: "package-2019-06",
				},
				{
					Tag: "package-2020-01-13-preview",
				},
				{
					Tag: "package-2021-06-22",
				},
				{
					Tag: "package-2022-01-31",
				},
				{
					Tag: "package-2022-02-22",
				},
				{
					Tag: "package-2022-08-08",
				},
			},
		},
	}

	for _, testcase := range testcases {
		log.Printf("[DEBUG] testcase: %+v", testcase.Input)
		actual := ParseAutoRestConfig(testcase.Input)
		if len(actual) != len(testcase.Expected) {
			t.Errorf("expected %d packages, got %d", len(testcase.Expected), len(actual))
		}
		for i := range actual {
			if actual[i].Tag != testcase.Expected[i].Tag {
				t.Errorf("expected %s, got %s", testcase.Expected[i].Tag, actual[i].Tag)
			}
			if len(actual[i].InputFiles) == 0 {
				t.Errorf("expected non-empty input files")
			}
		}
	}
}

func Test_ParseYamlConfig(t *testing.T) {
	testcases := []struct {
		Input       string
		Expected    *Package
		ExpectError bool
	}{
		{
			Input: `$(tag) == 'package-2015-10'
input-file:
- Microsoft.Automation/stable/2015-10-31/account.json
- Microsoft.Automation/stable/2015-10-31/certificate.json
`,
			Expected: &Package{
				Tag: "package-2015-10",
				InputFiles: []string{
					"Microsoft.Automation/stable/2015-10-31/account.json",
					"Microsoft.Automation/stable/2015-10-31/certificate.json",
				},
			},
			ExpectError: false,
		},
	}

	for _, testcase := range testcases {
		log.Printf("[DEBUG] testcase: %+v", testcase.Input)

		actual, err := ParseYamlConfig(testcase.Input)
		if testcase.ExpectError != (err != nil) {
			t.Errorf("expected error %v, got %v", testcase.ExpectError, err)
			continue
		}
		if actual.Tag != testcase.Expected.Tag {
			t.Errorf("expected %s, got %s", testcase.Expected.Tag, actual.Tag)
		}
		if len(actual.InputFiles) != len(testcase.Expected.InputFiles) {
			t.Errorf("expected %d input files, got %d", len(testcase.Expected.InputFiles), len(actual.InputFiles))
			continue
		}
		for i := range actual.InputFiles {
			if actual.InputFiles[i] != testcase.Expected.InputFiles[i] {
				t.Errorf("expected %s, got %s", testcase.Expected.InputFiles[i], actual.InputFiles[i])
			}
		}
	}

}
