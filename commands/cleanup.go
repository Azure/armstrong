package commands

import (
	"log"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

type CleanupCommand struct {
	Ui cli.Ui
}

func (c CleanupCommand) Help() string {
	helpText := `
Usage: azurerm-rest-api-testing-tool cleanup
` + c.Synopsis() + "\n\n"

	return strings.TrimSpace(helpText)
}

func (c CleanupCommand) Synopsis() string {
	return "Clean up dependency"
}

func (c CleanupCommand) Run(args []string) int {
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
	return 0
}
