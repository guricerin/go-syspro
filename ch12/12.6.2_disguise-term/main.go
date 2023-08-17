package main

import (
	"io"
	"os"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

func main() {
	var out io.Writer
	if isatty.IsTerminal(os.Stdout.Fd()) {
		out = colorable.NewColorableStdout()
	} else {
		out = colorable.NewNonColorable(os.Stdout)
	}
}
