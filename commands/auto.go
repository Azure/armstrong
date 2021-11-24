package commands

import "log"

func Auto(args []string) {
	Generate(args)
	Test([]string{})
	Cleanup()
	log.Println("[INFO] Test passed!")
}
