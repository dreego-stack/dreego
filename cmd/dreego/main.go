package main

import (
	"fmt"
	"os"

	"codeberg.org/dreego/dreego/pkg/generate"
)

func main() {
	dir := "routes"
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if err := generate.Run(dir); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
