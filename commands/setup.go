package commands

import (
	"log"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

type SetupCommand struct {
	Ui cli.Ui
}

func (c SetupCommand) Help() string {
	helpText := `
Usage: azurerm-rest-api-testing-tool setup
` + c.Synopsis() + "\n\n"

	return strings.TrimSpace(helpText)
}

func (c SetupCommand) Synopsis() string {
	return "Update dependency for tests"
}

func (c SetupCommand) Run(args []string) int {
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
	return 0
}
