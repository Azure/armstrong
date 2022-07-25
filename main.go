package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/ms-henglu/armstrong/commands"
)

func main() {
	c := &cli.CLI{
		Name:       "armstrong",
		Version:    VersionString(),
		Args:       os.Args[1:],
		HelpWriter: os.Stdout,
	}

	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			Reader:      os.Stdin,
			ErrorWriter: os.Stderr,
		},
	}

	c.Commands = map[string]cli.CommandFactory{
		"auto": func() (cli.Command, error) {
			return &commands.AutoCommand{
				Ui: ui,
			}, nil
		},
		"generate": func() (cli.Command, error) {
			return &commands.GenerateCommand{
				Ui: ui,
			}, nil
		},
		"cleanup": func() (cli.Command, error) {
			return &commands.CleanupCommand{
				Ui: ui,
			}, nil
		},
		"setup": func() (cli.Command, error) {
			return &commands.SetupCommand{
				Ui: ui,
			}, nil
		},
		"test": func() (cli.Command, error) {
			return &commands.TestCommand{
				Ui: ui,
			}, nil
		},
		"validate": func() (cli.Command, error) {
			return &commands.ValidateCommand{
				Ui: ui,
			}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		ui.Error("Error: " + err.Error())
	}

	os.Exit(exitStatus)
}
