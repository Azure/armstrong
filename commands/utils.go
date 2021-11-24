package commands

import (
	"os"

	"github.com/ms-henglu/azurerm-rest-api-testing-tool/tf"
)

func GetCommandArgs() (string, []string) {
	args := make([]string, 0)
	for _, arg := range os.Args {
		if arg == "-v" {
			tf.LogEnabled = true
		} else {
			args = append(args, arg)
		}
	}
	if len(args) > 1 {
		return args[1], args[2:]
	}
	return "", []string{}
}
