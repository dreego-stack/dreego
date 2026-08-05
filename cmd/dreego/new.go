package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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
		dest = strings.TrimSuffix(dest, ".tmpl")
		return os.WriteFile(dest, []byte(content), 0644)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	dreegoCoreVersion := dreegoVersion()

	c := exec.Command("go", "mod", "init", projName)
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod init failed: %v\n", err)
	}

	c = exec.Command("go", "mod", "edit", "-go=1.22")
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod edit -go failed: %v\n", err)
	}

	c = exec.Command("go", "mod", "edit", "-require", "codeberg.org/dreego/dreego/core@"+dreegoCoreVersion)
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod edit -require failed: %v\n", err)
	}

	// When the CLI runs from a repo-local build (or a pre-release version), the
	// required core version is not yet on the remote module proxy. Point the
	// scaffold at the local core module so `go mod tidy` and the build resolve
	// fully offline. For a release-installed binary there is no local core
	// directory, so tidy resolves the published tag instead.
	if coreDir := findLocalCore(); coreDir != "" {
		c = exec.Command("go", "mod", "edit", "-replace=codeberg.org/dreego/dreego/core="+coreDir)
		c.Dir = target
		c.Stdout, c.Stderr = nil, os.Stderr
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: go mod edit -replace failed: %v\n", err)
		}
	}

	c = exec.Command("go", "mod", "tidy")
	c.Env = append(os.Environ(), "GOWORK=off")
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod tidy failed: %v\n", err)
	}

	fmt.Printf("Done!\n")
	fmt.Printf("  cd %s && dreego generate && go run .\n", name)
}

// findLocalCore returns the absolute path to the local core module directory
// (used to replace the remote core dependency in generated scaffolds), or ""
// if it cannot be located. The path is resolved relative to this source file,
// which lives in <repo>/cmd/dreego/ for a repo-local build.
func findLocalCore() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	coreDir := filepath.Join(filepath.Dir(file), "..", "..", "core")
	st, err := os.Stat(filepath.Join(coreDir, "go.mod"))
	if err != nil || st.IsDir() {
		return ""
	}
	return coreDir
}
