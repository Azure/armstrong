package commands

import (
	"os"
)

func GetCommandArgs() (string, []string) {
	args := os.Args
	if len(args) > 1 {
		return args[1], args[2:]
	}
	return "", []string{}
}
