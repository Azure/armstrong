package commands

import (
	"flag"
	"io"
	"strings"
)

func defaultFlagSet(cmdName string) *flag.FlagSet {
	f := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	f.SetOutput(io.Discard)

	// Set the default Usage to empty
	f.Usage = func() {}

	return f
}

func helpForFlags(fs *flag.FlagSet) string {
	buf := &strings.Builder{}
	buf.WriteString("Options:\n\n")

	w := fs.Output()
	defer fs.SetOutput(w)
	fs.SetOutput(buf)
	fs.PrintDefaults()

	return buf.String()
}
