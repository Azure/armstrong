package commands

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/hcl"
	"github.com/ms-henglu/armstrong/loader"
	"github.com/ms-henglu/armstrong/resource"
	"github.com/ms-henglu/armstrong/types"
)

type GenerateCommand struct {
	Ui                cli.Ui
	path              string
	workingDir        string
	useRawJsonPayload bool
	overwrite         bool
	resourceType      string
}

func (c *GenerateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("generate")

	fs.StringVar(&c.path, "path", "", "path to a swagger 'Create' example")
	fs.StringVar(&c.workingDir, "working-dir", "", "output path to Terraform configuration files")
	fs.BoolVar(&c.useRawJsonPayload, "raw", false, "whether use raw json payload in 'body'")
	fs.BoolVar(&c.overwrite, "overwrite", false, "whether overwrite existing terraform configurations")
	fs.StringVar(&c.resourceType, "type", "resource", "the type of the resource to be generated, allowed values: 'resource'(supports CRUD) and 'data'(read-only). Defaults to 'resource'")
	fs.Usage = func() { c.Ui.Error(c.Help()) }

	return fs
}

func (c GenerateCommand) Help() string {
	helpText := `
Usage: armstrong generate -path <path to a swagger 'Create' example> [-working-dir <output path to Terraform configuration files>]
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c GenerateCommand) Synopsis() string {
	return "Generate testing files including terraform configuration for dependencies and testing resource."
}

func (c GenerateCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %+v", err))
		return 1
	}
	if len(c.path) == 0 {
		c.Ui.Error(c.Help())
		return 1
	}
	return c.Execute()
}

func (c GenerateCommand) Execute() int {
	wd, err := os.Getwd()
	if err != nil {
		c.Ui.Error(fmt.Sprintf("failed to get working directory: %+v", err))
		return 1
	}
	if c.workingDir != "" {
		wd, err = filepath.Abs(c.workingDir)
		if err != nil {
			c.Ui.Error(fmt.Sprintf("working directory is invalid: %+v", err))
			return 1
		}
	}
	if c.overwrite {
		_ = os.RemoveAll(path.Join(wd, "testing.tf"))
		_ = os.RemoveAll(path.Join(wd, "dependency.tf"))
	}
	err = os.WriteFile(path.Join(wd, "provider.tf"), hclwrite.Format([]byte(hcl.ProviderHcl)), 0644)
	if err != nil {
		log.Fatalf("[ERROR] error writing provider.tf: %+v\n", err)
	}

	log.Println("[INFO] ----------- generate dependencies and testing resource ---------")
	// load dependencies
	log.Println("[INFO] loading dependencies")
	existDeps, deps := loadDependencies(wd)

	// load example
	log.Println("[INFO] generating testing files")
	exampleFilepath := c.path
	var base resource.Base
	if c.resourceType == "data" {
		base, err = resource.NewDataSourceFromExample(exampleFilepath)
		if err != nil {
			log.Fatalf("[ERROR] error loading data source: %+v\n", err)
		}
	} else {
		base, err = resource.NewResourceFromExample(exampleFilepath)
		if err != nil {
			log.Fatalf("[ERROR] error loading resource: %+v\n", err)
		}
	}

	// generate dependency.tf
	requiredDependencies := base.RequiredDependencies(existDeps, deps)
	dependencyHcl := ""
	for _, dep := range requiredDependencies {
		dependencyHcl = hcl.Combine(dependencyHcl, hcl.RenameLabel(dep.ExampleConfiguration))
	}
	err = appendFile(path.Join(wd, "dependency.tf"), dependencyHcl)
	if err != nil {
		log.Fatalf("[ERROR] error writing dependency.tf: %+v\n", err)
	}
	log.Println("[INFO] dependency.tf generated")

	// generate testing.tf
	refs := make([]resource.Reference, 0)
	for _, dep := range hcl.LoadExistingDependencies(wd) {
		refs = append(refs, *resource.NewReferenceFromAddress(fmt.Sprintf("%s.%s", dep.Address, dep.ReferredProperty)))
	}
	base.UpdatePropertyDependencyMappingsReference(append(deps, existDeps...), refs)
	base.GenerateLabel(refs)
	testResourceHcl := base.Hcl(c.useRawJsonPayload)
	err = appendFile(path.Join(wd, "testing.tf"), testResourceHcl)
	if err != nil {
		log.Fatalf("[ERROR] error writing testing.tf: %+v\n", err)
	}
	log.Println("[INFO] testing.tf generated")
	return 0
}

func appendFile(filename string, hclContent string) error {
	content := hclContent
	if _, err := os.Stat(filename); err == nil {
		existingHcl, err := os.ReadFile(filename)
		if err != nil {
			log.Printf("[WARN] reading %s: %+v", filename, err)
		} else {
			content = hcl.Combine(string(existingHcl), content)
		}
	}
	return os.WriteFile(filename, hclwrite.Format([]byte(content)), 0644)
}

func loadDependencies(workingDir string) ([]types.Dependency, []types.Dependency) {
	mappingJsonLoader := loader.MappingJsonDependencyLoader{}
	hardcodeLoader := loader.HardcodeDependencyLoader{}

	deps := make([]types.Dependency, 0)
	depsMap := make(map[string]types.Dependency)
	if temp, err := mappingJsonLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType+"."+dep.ReferredProperty] = dep
		}
	}
	if temp, err := hardcodeLoader.Load(); err == nil {
		for _, dep := range temp {
			depsMap[dep.ResourceType+"."+dep.ReferredProperty] = dep
		}
	}
	for _, dep := range depsMap {
		deps = append(deps, dep)
	}
	existDeps := hcl.LoadExistingDependencies(workingDir)
	for i := range existDeps {
		ref := existDeps[i].ResourceType + "." + existDeps[i].ReferredProperty
		if dep, ok := depsMap[ref]; ok {
			existDeps[i].Pattern = dep.Pattern
		}
	}
	return existDeps, deps
}
