package tf

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

type Terraform struct {
	exec       *tfexec.Terraform
	LogEnabled bool
}

const planfile = "tfplan"

func NewTerraform(logEnabled bool) (*Terraform, error) {
	execPath, err := FindTerraform(context.TODO())
	if err != nil {
		return nil, err
	}
	workingDirectory, _ := os.Getwd()
	tf, err := tfexec.NewTerraform(workingDirectory, execPath)
	if err != nil {
		return nil, err
	}

	t := &Terraform{
		exec:       tf,
		LogEnabled: logEnabled,
	}
	t.SetLogEnabled(true)
	err = t.exec.SetLogPath(path.Join(workingDirectory, "log.txt"))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Terraform) SetLogEnabled(enabled bool) {
	if enabled && t.LogEnabled {
		t.exec.SetStdout(os.Stdout)
		t.exec.SetStderr(os.Stderr)
		t.exec.SetLogger(log.New(os.Stdout, "", 0))
	} else {
		t.exec.SetStdout(ioutil.Discard)
		t.exec.SetStderr(ioutil.Discard)
		t.exec.SetLogger(log.New(ioutil.Discard, "", 0))
	}
}

func (t *Terraform) Init() error {
	if _, err := os.Stat(".terraform"); os.IsNotExist(err) {
		err := t.exec.Init(context.Background(), tfexec.Upgrade(false))
		// ignore the error if can't find azapi
		if err != nil && strings.Contains(err.Error(), "Azure/azapi: provider registry registry.terraform.io does not have") {
			return nil
		}
		return err
	}
	log.Println("[INFO] skip running init command because .terraform folder exist")
	return nil
}

func (t *Terraform) Show() (*tfjson.State, error) {
	return t.exec.Show(context.TODO())
}

func (t *Terraform) Plan() (*tfjson.Plan, error) {
	ok, err := t.exec.Plan(context.TODO(), tfexec.Out(planfile))
	if err != nil {
		return nil, err
	}
	if !ok {
		// no changes
		return nil, nil
	}

	t.SetLogEnabled(false)
	p, err := t.exec.ShowPlanFile(context.TODO(), planfile)
	t.SetLogEnabled(true)
	return p, err
}

func (t *Terraform) Apply() error {
	return t.exec.Apply(context.TODO())
}

func (t *Terraform) Destroy() error {
	return t.exec.Destroy(context.TODO())
}
