package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/hashicorp/terraform-exec/tfinstall"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/helper"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/loader"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/resource"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/types"
	"io/ioutil"
	"log"
)

func main() {
	mappingJsonLoader := loader.MappingJsonDependencyLoader{MappingJsonFilepath: "C:\\Users\\henglu\\go\\src\\github.com\\ms-henglu\\azurerm-terraform-mapping-tool\\mappings.json"}
	hardcodeLoader := loader.HardcodeDependencyLoader{}

	deps := make([]types.Dependency, 0)
	if temp, err := mappingJsonLoader.Load(); err == nil {
		deps = append(deps, temp...)
	}
	if temp, err := hardcodeLoader.Load(); err == nil {
		deps = append(deps, temp...)
	}

	exampleFilepath := "C:\\Users\\henglu\\go\\src\\github.com\\Azure\\azure-rest-api-specs\\specification\\synapse\\resource-manager\\Microsoft.Synapse\\stable\\2020-12-01\\examples\\CreateOrUpdateSqlPoolWorkloadGroupMax.json"
	exampleResource, err := resource.NewResourceFromExample(exampleFilepath, "parameters")
	if err != nil {
		panic(err)
	}

	dependencyHcl := exampleResource.GetDependencyHcl(deps)
	testResourceHcl := exampleResource.GetHcl(dependencyHcl)
	finalHcl := helper.GetCombinedHcl(dependencyHcl, testResourceHcl)
	finalHcl = helper.GetCombinedHcl(helper.ProviderHcl, finalHcl)

	fmt.Println(finalHcl)

	err = ioutil.WriteFile("./temp/main.tf", []byte(finalHcl), 0644)
	if err != nil {
		panic(err)
	}

	execPath, err := tfinstall.LookPath().ExecPath(context.TODO())
	if err != nil {
		panic(err)
	}
	tf, err := tfexec.NewTerraform("./temp", execPath)
	if err != nil {
		panic(err)
	}
	err = tf.Init(context.Background(), tfexec.Upgrade(false))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}
	state, err := tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	ok, err := tf.Plan(context.Background(), tfexec.Dir(""))
	if err != nil {
		log.Fatalf("error running Init: %s", err)
	}
	log.Printf("%v\n", ok)

	state, err = tf.Show(context.Background())
	if err != nil {
		log.Fatalf("error running Show: %s", err)
	}

	fmt.Println(state.FormatVersion) // "0.1"
}
