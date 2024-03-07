package hcl

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var ResourceBlockSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "resource",
			LabelNames: []string{"type", "name"},
		},
	},
}

var VarBlockSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "variable",
			LabelNames: []string{"name"},
		},
	},
}

var ProviderBlockSchema = hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "provider",
			LabelNames: []string{"type"},
		},
	},
}

var evalContext = &hcl.EvalContext{
	Functions: map[string]function.Function{
		"abs":        stdlib.AbsoluteFunc,
		"coalesce":   stdlib.CoalesceFunc,
		"concat":     stdlib.ConcatFunc,
		"hasindex":   stdlib.HasIndexFunc,
		"int":        stdlib.IntFunc,
		"jsondecode": stdlib.JSONDecodeFunc,
		"jsonencode": stdlib.JSONEncodeFunc,
		"length":     stdlib.LengthFunc,
		"lower":      stdlib.LowerFunc,
		"max":        stdlib.MaxFunc,
		"min":        stdlib.MinFunc,
		"reverse":    stdlib.ReverseFunc,
		"strlen":     stdlib.StrlenFunc,
		"substr":     stdlib.SubstrFunc,
		"upper":      stdlib.UpperFunc,
	},
}

// could be azurerm or azapi
type AzureProvider struct {
	Type                      string
	Alias                     string
	SubscriptionId            string
	TenantId                  string
	AuxiliaryTenantIds        []string
	AuxiliaryTenantIdsString  string
	ClientId                  string
	ClientCertificate         string
	ClientCertificatePassword string
	ClientSecret              string
	OidcRequestToken          string
	OidcToken                 string
	FileName                  string
	LineNumber                int
}

func (p AzureProvider) Name() string {
	if p.Alias != "" {
		return fmt.Sprintf("%q.%q", p.Type, p.Alias)
	}
	return p.Type
}

type AzapiResource struct {
	Name       string
	Type       string
	Body       string
	FileName   string
	LineNumber int
}

type Variable struct {
	Name        string
	HasDefault  bool
	FileName    string
	LineNumber  int
	IsSensitive bool
}

// mockVariables returns a map of variables that are used in the given traversals.
// The mocked variables are prefixed with the given prefix and of  type string.
func mockVariables(traversals []hcl.Traversal) map[string]cty.Value {
	const variablePrefix = "$"

	ret := make(map[string]cty.Value)
	for _, traversal := range traversals {
		for k, v := range mockVariable(traversal, 0, variablePrefix) {
			ret[k] = v
		}
	}

	return ret
}

// one hcl.Traversal corresponds to one reference
// e.g.,[{{} azapi_resource testdata/test.tf:121,18-32} {{} networkInterface testdata/test.tf:121,32-49} {{} id testdata/test.tf:121,49-52}]
func mockVariable(steps hcl.Traversal, index int, placeholder string) map[string]cty.Value {
	if index >= len(steps) {
		return map[string]cty.Value{}
	}

	step := steps[index]
	result := map[string]cty.Value{}

	switch stepValue := step.(type) {
	case hcl.TraverseRoot:
		placeholder += stepValue.Name
		result[stepValue.Name] = cty.ObjectVal(mockVariable(steps, index+1, placeholder))
	case hcl.TraverseAttr:
		placeholder += "." + stepValue.Name
		if index < len(steps)-1 {
			result[stepValue.Name] = cty.ObjectVal(mockVariable(steps, index+1, placeholder))
		} else {
			result[stepValue.Name] = cty.StringVal(placeholder)
		}
	}
	return result
}

// mockExpression evaluates the given expression with variables replaced by their mock values.
func mockExpression(expr hcl.Expression) (*cty.Value, []error) {
	evalContext.Variables = mockVariables(expr.Variables())

	v, diags := expr.Value(evalContext)
	if diags.HasErrors() {
		return nil, diags.Errs()
	}

	return &v, nil
}

func FindTfFiles(path string) (*[]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	tfFiles := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			// We only care about terraform configuration files.
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".tf") {
			tfFiles = append(tfFiles, path+"/"+name)
		}
	}

	return &tfFiles, nil
}

func ParseHclFile(path string) (*hcl.File, []error) {
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		return nil, diags.Errs()
	}

	return f, nil
}

func ParseAzapiResource(f hcl.File) (*[]AzapiResource, []error) {
	content, _, diags := f.Body.PartialContent(&ResourceBlockSchema)
	if diags.HasErrors() {
		logrus.Error(diags)
	}

	results := make([]AzapiResource, 0)
	for _, block := range content.Blocks {
		if block.Type == "resource" && len(block.Labels) > 1 && block.Labels[0] == "azapi_resource" {

			attrs, diags := block.Body.JustAttributes()
			if diags.HasErrors() {
				diags := skipJustAttributesDiags(diags)
				if diags.HasErrors() {
					return nil, diags.Errs()
				}
			}

			resourceTypeRaw, ok := attrs["type"]
			if !ok {
				return nil, []error{fmt.Errorf("resource type is not specified for azapi_resource.%s", block.Labels[1])}
			}

			resourceType, errs := mockExpression(resourceTypeRaw.Expr)
			if errs != nil {
				return nil, errs
			}

			r := AzapiResource{
				Name:       block.Labels[1],
				Type:       resourceType.AsString(),
				FileName:   block.DefRange.Filename,
				LineNumber: block.DefRange.Start.Line,
			}

			if p := attrs["body"]; p != nil {
				body, errs := mockExpression(p.Expr)
				if errs != nil {
					return nil, errs
				}

				r.Body = body.AsString()
			}

			results = append(results, r)
		}
	}

	return &results, nil
}

func ParseVariables(f hcl.File) (*map[string]Variable, []error) {
	content, _, diags := f.Body.PartialContent(&VarBlockSchema)
	if diags.HasErrors() {
		logrus.Error(diags)
	}

	results := make(map[string]Variable, 0)
	for _, block := range content.Blocks {
		attrs, diags := block.Body.JustAttributes()
		if diags.HasErrors() {
			diags := skipJustAttributesDiags(diags)
			if diags.HasErrors() {
				return nil, diags.Errs()
			}
		}

		var hasDefault bool
		if p := attrs["default"]; p != nil {
			hasDefault = true
		}

		var isSensitive bool
		if p := attrs["sensitive"]; p != nil {
			value, diags := p.Expr.Value(nil)
			if diags.HasErrors() {
				return nil, diags.Errs()
			}
			isSensitive = value.True()
		}

		results[block.Labels[0]] = Variable{
			Name:        block.Labels[0],
			FileName:    block.DefRange.Filename,
			LineNumber:  block.DefRange.Start.Line,
			IsSensitive: isSensitive,
			HasDefault:  hasDefault,
		}
	}

	return &results, nil
}

func ParseAzureProvider(f hcl.File) (*[]AzureProvider, []error) {
	content, _, diags := f.Body.PartialContent(&ProviderBlockSchema)
	if diags.HasErrors() {
		logrus.Error(diags)
	}

	results := make([]AzureProvider, 0)
	for _, block := range content.Blocks {
		if block.Type == "provider" && len(block.Labels) > 0 && (block.Labels[0] == "azapi" || block.Labels[0] == "azurerm") {
			attrs, diags := block.Body.JustAttributes()
			if diags.HasErrors() {
				diags := skipJustAttributesDiags(diags)
				if diags.HasErrors() {
					return nil, diags.Errs()
				}
			}

			p := AzureProvider{
				Type:       block.Labels[0],
				FileName:   block.DefRange.Filename,
				LineNumber: block.DefRange.Start.Line,
			}

			if raw, ok := attrs["alias"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.Alias = v.AsString()
			}

			if raw, ok := attrs["subscription_id"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.SubscriptionId = v.AsString()
			}

			if raw, ok := attrs["tenant_id"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.TenantId = v.AsString()
			}

			if raw, ok := attrs["auxiliary_tenant_ids"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				// mocked variable is string type, but this field could be list or string
				if v.Type().IsPrimitiveType() {
					p.AuxiliaryTenantIdsString = v.AsString()
				}

				if v.CanIterateElements() {
					slice := v.AsValueSlice()
					for _, v := range slice {
						p.AuxiliaryTenantIds = append(p.AuxiliaryTenantIds, v.AsString())
					}
				}
			}

			if raw, ok := attrs["client_id"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.ClientId = v.AsString()
			}

			if raw, ok := attrs["client_certificate"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.ClientCertificate = v.AsString()
			}

			if raw, ok := attrs["client_certificate_password"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.ClientCertificatePassword = v.AsString()
			}

			if raw, ok := attrs["client_secret"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.ClientSecret = v.AsString()
			}

			if raw, ok := attrs["oidc_request_token"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.OidcRequestToken = v.AsString()
			}

			if raw, ok := attrs["oidc_token"]; ok {
				v, errs := mockExpression(raw.Expr)
				if errs != nil {
					return nil, errs
				}

				p.OidcToken = v.AsString()
			}

			results = append(results, p)
		}
	}

	return &results, nil
}

// skip the diags when blocks are found
func skipJustAttributesDiags(diags hcl.Diagnostics) hcl.Diagnostics {
	if diags == nil {
		return nil
	}

	const BlockNotAllowed = "Blocks are not allowed here"

	result := hcl.Diagnostics{}
	for _, diag := range diags {
		if !regexp.MustCompile(BlockNotAllowed).MatchString(diag.Error()) {
			result.Append(diag)
		}
	}

	return result
}
