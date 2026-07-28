package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

)

func cmdNew(args []string) {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dreego new <name>\n")
		os.Exit(1)
	}
	name := args[0]
	projName := filepath.Base(name)

	target, _ := filepath.Abs(name)
	if _, err := os.Stat(target); err == nil {
		fmt.Fprintf(os.Stderr, "error: %s already exists\n", name)
		os.Exit(1)
	}

	fmt.Printf("Creating %s/\n", name)

	templateRoot := "blueprints/landing"

	err := fs.WalkDir(blueprintsSrc, templateRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(templateRoot, path)
		if rel == "." {
			return nil
		}

		dest := filepath.Join(target, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		data, err := blueprintsSrc.ReadFile(path)
		if err != nil {
			return err
		}

		content := strings.ReplaceAll(string(data), "§$name$§", projName)

		os.MkdirAll(filepath.Dir(dest), 0755)
		return os.WriteFile(dest, []byte(content), 0644)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	dreegoCoreVersion := "v0.0.14"

	c := exec.Command("go", "mod", "init", projName)
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod init failed: %v\n", err)
	}

	c = exec.Command("go", "mod", "edit", "-go=1.22")
	c.Dir = target
	c.Run()

	c = exec.Command("go", "mod", "edit", "-require", "codeberg.org/dreego/dreego@"+dreegoCoreVersion)
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	c.Run()

	fmt.Printf("Done!\n")
	fmt.Printf("  cd %s && dreego generate && go run .\n", name)
}
