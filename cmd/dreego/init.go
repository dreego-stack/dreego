package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed all:blueprints
var blueprintsSrc embed.FS

func cmdInit(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dreego init <path>\n")
		os.Exit(1)
	}
	target := args[0]

	if err := os.MkdirAll(target, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err := fs.WalkDir(blueprintsSrc, "blueprints/default", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel("blueprints/default", path)
		if rel == "." {
			return nil
		}
		dest := filepath.Join(target, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		data, err := fs.ReadFile(blueprintsSrc, path)
		if err != nil {
			return err
		}
		return os.WriteFile(dest, data, 0644)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("initialized %s\n", target)
}
