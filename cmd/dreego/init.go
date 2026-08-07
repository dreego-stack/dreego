package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
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
		content := strings.ReplaceAll(string(data), "§$name$§", moduleName(target))
		dest = strings.TrimSuffix(dest, ".tmpl")
		return os.WriteFile(dest, []byte(content), 0644)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("initialized %s\n", target)
}

// moduleName returns the module path declared in the go.mod of dir (first
// "module " line), or the directory base name as fallback when no go.mod
// exists. Used to qualify the generated package import in blueprints.
func moduleName(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if rest, ok := strings.CutPrefix(line, "module "); ok {
				return strings.TrimSpace(rest)
			}
		}
	}
	return filepath.Base(dir)
}
