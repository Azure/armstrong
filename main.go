package main

import (
	"github.com/ms-henglu/azurerm-rest-api-testing-tool/commands"
	"log"
)

func main() {
	command, args := commands.GetCommandArgs()
	log.Printf("[INFO] command: %v, args: %v", command, args)
	switch command {
	case "generate":
		commands.Generate(args)
	case "auto":
		commands.Auto(args)
	case "test":
		commands.Test(args)
	case "setup":
		commands.Setup()
	case "cleanup":
		commands.Cleanup()
	case "help":
	default:
		commands.Help()
	}

}
