package commands

import (
	"log"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

func Test(args []string) {
	terraform, err := tf.NewTerraform()
	if err != nil {
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}

	log.Printf("[INFO] prepare working directory\n")
	terraform.Init()
	plan, err := terraform.Plan()
	if err != nil {
		log.Fatalf("[Error] error running terraform plan: %+v\n", err)
	}

	log.Printf("[INFO] Running plan completed, found %v changes\n", tf.GetChanges(plan))
	err = terraform.Apply()
	if err != nil {
		log.Fatalf("[Error] error running terraform apply: %+v\n", err)
	}

	if hasDiff := terraform.HasDiff(); hasDiff {
		log.Fatalf("[INFO] Found bugs! TODO: print details\n")
	} else {
		log.Printf("[INFO] Test passed! \n")
	}
}
