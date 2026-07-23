package main

import (
	"fmt"
	"os"

	"codeberg.org/dreego/dreego/pkg/generate"
)

func main() {
	if err := generate.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
