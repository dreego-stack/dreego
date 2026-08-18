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
	if !validProjectName(name) {
		fmt.Fprintf(os.Stderr, "error: invalid project name %q\n", name)
		fmt.Fprintf(os.Stderr, "  the name must be a Go module path segment: start with a letter, use only letters, digits, '-', '_', '.', and '/'.\n")
		fmt.Fprintf(os.Stderr, "  examples: myapp, github.com/me/myapp\n")
		os.Exit(1)
	}
	projName := filepath.Base(name)

	if !goAvailable() {
		fmt.Fprintf(os.Stderr, "error: 'go' executable not found on PATH.\n")
		fmt.Fprintf(os.Stderr, "  Dreego requires Go 1.22 or newer. Install it from https://go.dev/doc/install and retry.\n")
		os.Exit(1)
	}

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

	dreegoCoreVersion := scaffoldVersion(dreegoVersion())

	c := exec.Command("go", "mod", "init", name)
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

	c = exec.Command("go", "mod", "edit", "-require", "github.com/dreego-stack/dreego@"+dreegoCoreVersion)
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod edit -require failed: %v\n", err)
	}

	// When the CLI runs from a repo-local build (or a pre-release version), the
	// required dreego version is not yet on the remote module proxy. Point the
	// scaffold at the local repo root so `go mod tidy` and the build resolve
	// fully offline. For a release-installed binary there is no local repo
	// directory, so tidy resolves the published tag instead.
	if repoDir := findLocalRepo(); repoDir != "" {
		c = exec.Command("go", "mod", "edit", "-replace=github.com/dreego-stack/dreego="+repoDir)
		c.Dir = target
		c.Stdout, c.Stderr = nil, os.Stderr
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: go mod edit -replace failed: %v\n", err)
		}
	}

	c = exec.Command("go", "mod", "tidy")
	c.Env = append(os.Environ(), "GOWORK=off")
	c.Dir = target
	out, err := c.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod tidy failed: %v\n", err)
		if s := strings.TrimSpace(string(out)); s != "" {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
		fmt.Fprintf(os.Stderr, "  if the dreego module cannot be resolved (no network, or the CLI was built from an untagged checkout), set DREEGO_LOCAL_REPO=/path/to/dreego to use a local checkout.\n")
	}

	fmt.Printf("Done!\n")
	fmt.Printf("  cd %s && dreego generate && go run .\n", name)
}

func scaffoldVersion(version string) string {
	if version == "" || version == "dev" || version == "(devel)" {
		return "v0.0.0"
	}
	return version
}

// findLocalRepo returns the absolute path to the local dreego repo root
// (used to replace the remote dreego dependency in generated scaffolds), or ""
// if it cannot be located. The path is resolved relative to this source file,
// which lives in <repo>/cli/dreego/ for a repo-local build.
func findLocalRepo() string {
	if repoDir := os.Getenv("DREEGO_LOCAL_REPO"); repoDir != "" {
		if info, err := os.Stat(filepath.Join(repoDir, "go.mod")); err == nil && !info.IsDir() {
			return repoDir
		}
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return findRepoFromWorkingModule()
	}
	repoDir := filepath.Join(filepath.Dir(file), "..", "..")
	st, err := os.Stat(filepath.Join(repoDir, "go.mod"))
	if err != nil || st.IsDir() {
		return findRepoFromWorkingModule()
	}
	return repoDir
}

func findRepoFromWorkingModule() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(filepath.Join(wd, "go.mod"))
	if err != nil {
		return ""
	}
	prefix := "replace github.com/dreego-stack/dreego => "
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, prefix) {
			continue
		}
		path := strings.TrimSpace(strings.TrimPrefix(line, prefix))
		if !filepath.IsAbs(path) {
			path = filepath.Join(wd, path)
		}
		if info, statErr := os.Stat(filepath.Join(path, "go.mod")); statErr == nil && !info.IsDir() {
			return filepath.Clean(path)
		}
	}
	return ""
}

// validProjectName reports whether name is a usable Go module path segment
// for a scaffolded project. It allows multi-segment paths like
// github.com/me/myapp but rejects names that would produce an invalid module
// statement (digits-first, spaces, shell metacharacters, empty).
func validProjectName(name string) bool {
	if name == "" || name == "." || name == ".." {
		return false
	}
	if strings.ContainsAny(name, " \t\"'\\`$;|&<>(){}[]!*?") {
		return false
	}
	for _, seg := range strings.Split(name, "/") {
		if seg == "" || seg == "." || seg == ".." {
			return false
		}
		if !isValidFirstChar(seg[0]) {
			return false
		}
		for i := 0; i < len(seg); i++ {
			if !isValidNameChar(seg[i]) {
				return false
			}
		}
	}
	return true
}

func isValidFirstChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func isValidNameChar(b byte) bool {
	if isValidFirstChar(b) {
		return true
	}
	if b >= '0' && b <= '9' {
		return true
	}
	return b == '-' || b == '_' || b == '.'
}

// goAvailable reports whether the 'go' executable can be found on PATH.
func goAvailable() bool {
	_, err := exec.LookPath("go")
	return err == nil
}
