package main

import (
	"fmt"
	"os"

	"codeberg.org/dreego/dreego/pkg/generate"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dreego <command>")
		fmt.Fprintln(os.Stderr, "  generate    Transpile .dreego files to _dreego.go")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		if err := generate.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
