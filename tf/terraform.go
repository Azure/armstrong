package tf

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	tfjson "github.com/hashicorp/terraform-json"
)

type Terraform struct {
	exec tfexec.Terraform
}

const planfile = "tfplan"

var LogEnabled = false

func NewTerraform() (*Terraform, error) {
	execPath, err := tfinstall.LookPath().ExecPath(context.TODO())
	if err != nil {
		return nil, err
	}
	workingDirectory, err := os.Getwd()
	tf, err := tfexec.NewTerraform(workingDirectory, execPath)
	if err != nil {
		return nil, err
	}

	t := &Terraform{
		exec: *tf,
	}
	t.SetLogEnabled(true)
	return t, nil
}

func (t *Terraform) SetLogEnabled(enabled bool) {
	if enabled && LogEnabled {
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
		// ignore the error if can't find azurerm-restapi
		if err != nil && strings.Contains(err.Error(), "Azure/azurerm-restapi: provider registry registry.terraform.io does not have") {
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
