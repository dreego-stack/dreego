package main

import (
	"fmt"

	typescript "github.com/dreego-stack/dreego/internal/transpiler/js/ts"
)

func cmdTools(args []string) error {
	if len(args) != 2 || args[0] != "install" || args[1] != "typescript" {
		return fmt.Errorf("usage: dreego tools install typescript")
	}
	path, err := typescript.Install()
	if err != nil {
		return err
	}
	fmt.Printf("TypeScript %s installed at %s\n", typescript.Version, path)
	return nil
}
