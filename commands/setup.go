package commands

import (
	"log"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

func Setup() {
	terraform, err := tf.NewTerraform()
	if err != nil {
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}
	log.Printf("[INFO] prepare working directory\n")
	terraform.Init()
	err = terraform.Apply()
	if err != nil {
		log.Fatalf("[Error] error setting up resources: %+v\n", err)
	}
}
