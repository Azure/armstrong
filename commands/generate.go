package commands

import (
	"io/ioutil"
	"log"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/helper"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/loader"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/resource"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/types"
)

func Generate(args []string) {
	if len(args) == 0 {
		Help()
		return
	}

	// load dependencies
	log.Println("[INFO] loading dependencies")
	// MappingJsonFilepath: "C:\\Users\\henglu\\go\\src\\github.com\\ms-henglu\\azurerm-terraform-mapping-tool\\mappings.json"
	mappingJsonLoader := loader.MappingJsonDependencyLoader{}
	hardcodeLoader := loader.HardcodeDependencyLoader{}

	deps := make([]types.Dependency, 0)
	depsMap := make(map[string]types.Dependency, 0)
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
	exampleFilepath := args[0]
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
}
