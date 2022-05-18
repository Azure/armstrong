package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/helper"
	"github.com/ms-henglu/armstrong/loader"
	"github.com/ms-henglu/armstrong/resource"
	"github.com/ms-henglu/armstrong/types"
)

type GenerateCommand struct {
	Ui   cli.Ui
	path string
}

func (c *GenerateCommand) flags() *flag.FlagSet {
	fs := defaultFlagSet("generate")

	fs.StringVar(&c.path, "path", "", "filepath of rest api to create arm resource example")

	fs.Usage = func() { c.Ui.Error(c.Help()) }

	return fs
}

func (c GenerateCommand) Help() string {
	helpText := `
Usage: armstrong generate -path <filepath to example>
` + c.Synopsis() + "\n\n" + helpForFlags(c.flags())

	return strings.TrimSpace(helpText)
}

func (c GenerateCommand) Synopsis() string {
	return "Generate testing files including terraform configuration for dependencies and testing resource."
}

func (c GenerateCommand) Run(args []string) int {
	f := c.flags()
	if err := f.Parse(args); err != nil {
		c.Ui.Error(fmt.Sprintf("Error parsing command-line flags: %s", err))
		return 1
	}
	if len(c.path) == 0 {
		c.Ui.Error(c.Help())
		return 1
	}

	log.Println("[INFO] ----------- generate dependency and test resource ---------")
	// load dependencies
	log.Println("[INFO] loading dependencies")
	// MappingJsonFilepath: "C:\\Users\\henglu\\go\\src\\github.com\\ms-henglu\\azurerm-terraform-mapping-tool\\mappings.json"
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

	// load example and generate hcl
	//exampleFilepath := "C:\\Users\\henglu\\go\\src\\github.com\\Azure\\azure-rest-api-specs\\specification\\synapse\\resource-manager\\Microsoft.Synapse\\stable\\2020-12-01\\examples\\CreateOrUpdateSqlPoolWorkloadGroupMax.json"
	log.Println("[INFO] generating testing files")
	exampleFilepath := c.path
	exampleResource, err := resource.NewResourceFromExample(exampleFilepath)
	if err != nil {
		log.Fatalf("[Error] error reading example file: %+v\n", err)
	}

	dependencyHcl := exampleResource.GetDependencyHcl(deps)
	finalHcl := helper.GetCombinedHcl(helper.ProviderHcl, dependencyHcl)

	err = ioutil.WriteFile("dependency.tf", []byte(finalHcl), 0644)
	if err != nil {
		log.Fatalf("[Error] error writing dependency.tf: %+v\n", err)
	}
	log.Println("[INFO] dependency.tf generated")

	testResourceHcl := exampleResource.GetHcl(dependencyHcl)
	err = ioutil.WriteFile("testing.tf", []byte(testResourceHcl), 0644)
	if err != nil {
		log.Fatalf("[Error] error writing testing.tf: %+v\n", err)
	}
	log.Println("[INFO] testing.tf generated")
	return 0
}
