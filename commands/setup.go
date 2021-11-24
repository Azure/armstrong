package commands

import (
	"log"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

func Setup() {
	log.Println("[INFO] ----------- update resources ---------")
	terraform, err := tf.NewTerraform()
	if err != nil {
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}
	log.Printf("[INFO] prepare working directory\n")
	_ = terraform.Init()
	log.Println("[INFO] running apply command to update dependency resources...")
	err = terraform.Apply()
	if err != nil {
		log.Fatalf("[Error] error setting up resources: %+v\n", err)
	}
	log.Println("[INFO] all dependencies have been updated")
}
