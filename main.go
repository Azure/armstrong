package main

import (
	"io"
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/commands"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	log.SetOutput(io.Discard)

	c := &cli.CLI{
		Name:       "armstrong",
		Version:    VersionString(),
		Args:       os.Args[1:],
		HelpWriter: os.Stdout,
	}

	c.Commands = map[string]cli.CommandFactory{
		"generate": func() (cli.Command, error) {
			return &commands.GenerateCommand{}, nil
		},
		"validate": func() (cli.Command, error) {
			return &commands.ValidateCommand{}, nil
		},
		"test": func() (cli.Command, error) {
			return &commands.TestCommand{}, nil
		},
		"cleanup": func() (cli.Command, error) {
			return &commands.CleanupCommand{}, nil
		},
		"report": func() (cli.Command, error) {
			return &commands.ReportCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		logrus.Fatal(err)
	}

	os.Exit(exitStatus)
}
