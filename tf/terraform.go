package tf

import (
	"context"
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/sirupsen/logrus"
)

type Terraform struct {
	exec       *tfexec.Terraform
	LogEnabled bool
}

const (
	planfile            = "tfplan"
	logSensitiveDataEnv = "AZAPI_LOG_SENSITIVE_DATA"
)

func NewTerraform(workingDirectory string, logEnabled bool) (*Terraform, error) {
	execPath, err := FindTerraform(context.TODO())
	if err != nil {
		return nil, err
	}
	tf, err := tfexec.NewTerraform(workingDirectory, execPath)
	if err != nil {
		return nil, err
	}
	// SetEnv replaces the inherited environment, so copy it to preserve variables such as ARM_* and PATH.
	env := make(map[string]string)
	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok && key != "" && !strings.EqualFold(key, logSensitiveDataEnv) {
			env[key] = value
		}
	}
	env[logSensitiveDataEnv] = "true"
	if err := tf.SetEnv(env); err != nil {
		return nil, err
	}

	t := &Terraform{
		exec:       tf,
		LogEnabled: logEnabled,
	}
	t.SetLogEnabled(true)
	logPath := path.Join(workingDirectory, "log.txt")
	_ = os.RemoveAll(logPath)
	err = t.exec.SetLogPath(logPath)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *Terraform) SetLogEnabled(enabled bool) {
	if enabled && t.LogEnabled {
		t.exec.SetStdout(os.Stdout)
		t.exec.SetStderr(os.Stderr)
		t.exec.SetLogger(logrus.StandardLogger())
	} else {
		t.exec.SetStdout(io.Discard)
		t.exec.SetStderr(io.Discard)
		t.exec.SetLogger(log.New(io.Discard, "", 0))
	}
}

func (t *Terraform) Init() error {
	return t.exec.Init(context.Background(), tfexec.Upgrade(false))
}

func (t *Terraform) Show() (*tfjson.State, error) {
	return t.exec.Show(context.TODO())
}

func (t *Terraform) Plan() (*tfjson.Plan, error) {
	ok, err := t.exec.Plan(context.TODO(), tfexec.Out(planfile), tfexec.Refresh(true))
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

func (t *Terraform) Validate() (*tfjson.ValidateOutput, error) {
	return t.exec.Validate(context.TODO())
}
