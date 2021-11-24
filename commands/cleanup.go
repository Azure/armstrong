package commands

import (
	"log"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

func Cleanup() {
	log.Println("[INFO] ----------- cleanup resources ---------")
	terraform, err := tf.NewTerraform()
	if err != nil {
		log.Fatalf("[Error] error creating terraform executable: %+v\n", err)
	}
	log.Println("[INFO] prepare working directory")
	_ = terraform.Init()
	log.Println("[INFO] running destroy command to cleanup resources...")
	err = terraform.Destroy()
	if err != nil {
		log.Fatalf("[Error] error cleaning up resources: %+v\n", err)
	}
	log.Println("[INFO] all resources have been deleted")
}
