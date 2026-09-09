package main

import (
	"fmt"
	"io"

	"github.com/dreego-stack/dreego/internal/transpiler"
)

func cmdI18n(args []string, output io.Writer) error {
	if len(args) != 1 || args[0] != "extract" {
		return fmt.Errorf("usage: dreego i18n extract")
	}
	return transpiler.ExtractI18n(output)
}
