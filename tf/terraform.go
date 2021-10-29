package tf

import (
	"context"
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

	tf.SetStdout(os.Stdout)
	tf.SetStderr(os.Stderr)
	tf.SetLogger(log.New(os.Stdout, "", 0))
	return &Terraform{
		exec: *tf,
	}, nil
}

func (t *Terraform) Init() error {
	if _, err := os.Stat(".terraform"); os.IsNotExist(err) {
		err := t.exec.Init(context.Background(), tfexec.Upgrade(false))
		// ignore the error if can't find azurermg
		if err != nil && strings.Contains(err.Error(), "ms-henglu/azurermg: provider registry registry.terraform.io does not have") {
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
	return t.exec.ShowPlanFile(context.TODO(), planfile)
}

func (t *Terraform) Apply() error {
	return t.exec.Apply(context.TODO())
}

func (t *Terraform) Destroy() error {
	return t.exec.Destroy(context.TODO())
}

func (t *Terraform) HasDiff() bool {
	plan, err := t.Plan()
	if err != nil {
		return true
	}
	return plan != nil && GetChanges(plan) > 0
}
