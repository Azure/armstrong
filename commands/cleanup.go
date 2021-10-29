package commands

import (
	"log"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

func Cleanup() {
	terraform, err := tf.NewTerraform()
	if err != nil {
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}
	log.Printf("[INFO] prepare working directory\n")
	terraform.Init()
	err = terraform.Destroy()
	if err != nil {
		log.Fatalf("[Error] error cleaning up resources: %+v\n", err)
	}
}
