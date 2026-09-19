package main

import (
	"fmt"
	"io"

	dreefile "github.com/dreego-stack/dreego/internal/dreefile"
)

func cmdI18n(args []string, output io.Writer) error {
	if len(args) != 1 || args[0] != "extract" {
		return fmt.Errorf("usage: dreego i18n extract")
	}
	return dreefile.ExtractI18n(output)
}
