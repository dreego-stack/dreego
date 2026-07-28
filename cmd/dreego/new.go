package main

import (
	"fmt"
	"io/fs"
	"os"
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

	fmt.Printf("Done! Next steps:\n")
	fmt.Printf("  cd %s\n", name)
	fmt.Printf("  go mod init %s\n", projName)
	fmt.Printf("  go mod edit -require codeberg.org/dreego/dreego@v0.0.13\n")
	fmt.Printf("  go mod tidy\n")
	fmt.Printf("  dreego generate\n")
	fmt.Printf("  go run .\n")
	fmt.Printf("  docker build -t %s .  # production build\n", projName)
}
