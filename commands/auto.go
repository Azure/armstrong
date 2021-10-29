package commands

func Auto(args []string) {
	Generate(args)
	Test([]string{})
	Cleanup()
}
