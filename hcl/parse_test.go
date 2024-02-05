package hcl_test

import (
	"testing"

	"github.com/ms-henglu/armstrong/hcl"
)

func TestParseAzapiResource(t *testing.T) {
	testFileDir := "testdata/"

	tfFiles, err := hcl.FindTfFiles(testFileDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, tfFile := range *tfFiles {
		f, errs := hcl.ParseHclFile(tfFile)
		if errs != nil {
			t.Fatal(errs)
		}

		azapiResources, errs := hcl.ParseAzapiResource(*f)
		if errs != nil {
			t.Fatal(errs)
		}

		for _, ar := range *azapiResources {
			t.Log(ar)
		}
	}
}

func TestParseVariable(t *testing.T) {
	testFileDir := "testdata/"

	tfFiles, err := hcl.FindTfFiles(testFileDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, tfFile := range *tfFiles {
		f, errs := hcl.ParseHclFile(tfFile)
		if errs != nil {
			t.Fatal(errs)
		}

		vars, errs := hcl.ParseVariables(*f)
		if errs != nil {
			t.Fatal(errs)
		}

		for k, v := range *vars {
			t.Log(k, v)
		}
	}
}
