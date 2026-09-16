package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dreego-stack/dreego/cmd/dreego/internal/templates"
)

func cmdNew(args []string) {
	flags, err := parseScaffoldFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if flags.list {
		printTemplates()
		return
	}

	if len(flags.args) < 1 {
		fmt.Fprintf(os.Stderr, "usage: dreego new <name> [-t <template>]\n")
		os.Exit(1)
	}
	name := flags.args[0]
	if !validProjectName(name) {
		fmt.Fprintf(os.Stderr, "error: invalid project name %q\n", name)
		fmt.Fprintf(os.Stderr, "  the name must be a Go module path segment: start with a letter, use only letters, digits, '-', '_', '.', and '/'.\n")
		fmt.Fprintf(os.Stderr, "  examples: myapp, github.com/me/myapp\n")
		os.Exit(1)
	}

	meta, err := resolveTemplate(flags.template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !goAvailable() {
		fmt.Fprintf(os.Stderr, "error: 'go' executable not found on PATH.\n")
		fmt.Fprintf(os.Stderr, "  Dreego requires Go 1.27 or newer. Install it from https://go.dev/doc/install and retry.\n")
		os.Exit(1)
	}

	target, _ := filepath.Abs(name)
	if _, err := os.Stat(target); err == nil {
		fmt.Fprintf(os.Stderr, "error: %s already exists\n", name)
		os.Exit(1)
	}

	fmt.Printf("Creating %s/\n", name)

	if err := templates.Install(target, moduleName(target), meta.Name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	dreegoVersion := scaffoldVersion(dreegoVersion())

	c := exec.Command("go", "mod", "init", name)
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod init failed: %v\n", err)
	}

	c = exec.Command("go", "mod", "edit", "-go=1.27")
	c.Dir = target
	c.Stdout, c.Stderr = nil, os.Stderr
	if err := c.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: go mod edit -go failed: %v\n", err)
	}

	for _, module := range requiredModules(meta) {
		c = exec.Command("go", "mod", "edit", "-require", module+"@"+dreegoVersion)
		c.Dir = target
		c.Stdout, c.Stderr = nil, os.Stderr
		if err := c.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "warning: go mod edit -require failed: %v\n", err)
		}
	}

	// When the CLI runs from a repo-local build (or a pre-release version), the
	// required dreego version is not yet on the remote module proxy. Point the
	// scaffold at the local repo root so `go mod tidy` and the build resolve
	// fully offline. For a release-installed binary there is no local repo
	// directory, so tidy resolves the published tag instead.
	if repoDir := findLocalRepo(); repoDir != "" {
		for _, replacement := range localReplacements(repoDir, meta) {
			c = exec.Command("go", "mod", "edit", "-replace="+replacement.module+"="+replacement.dir)
			c.Dir = target
			c.Stdout, c.Stderr = nil, os.Stderr
			if err := c.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "warning: go mod edit -replace failed: %v\n", err)
			}
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

// requiredModules returns the dreego modules a scaffolded project must require:
// core always, the module for the template's adapter when it names one, and any
// extra module@version entries declared by the template.
func requiredModules(meta templates.Meta) []string {
	modules := []string{"github.com/dreego-stack/dreego/core"}
	if module, _ := adapterModule(meta.Adapter); module != "" {
		modules = append(modules, module)
	}
	return append(modules, meta.ExtraRequires...)
}

// adapterModule maps a template adapter name to its dreego module path and its
// path relative to the repo root. An empty adapter means the template needs no
// adapter module.
func adapterModule(adapter string) (module, dir string) {
	if adapter == "" {
		return "", ""
	}
	return "github.com/dreego-stack/dreego/adapter/" + adapter, filepath.Join("adapter", adapter)
}

type moduleReplacement struct {
	module string
	dir    string
}

// localReplacements returns the replace directives that point a scaffolded
// project at the local dreego checkout, derived from the selected template.
func localReplacements(repoDir string, meta templates.Meta) []moduleReplacement {
	replacements := []moduleReplacement{
		{"github.com/dreego-stack/dreego", repoDir},
		{"github.com/dreego-stack/dreego/core", filepath.Join(repoDir, "core")},
	}
	if module, dir := adapterModule(meta.Adapter); module != "" {
		replacements = append(replacements, moduleReplacement{module, filepath.Join(repoDir, dir)})
	}
	return replacements
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
// which lives in <repo>/cmd/dreego/ for a repo-local build.
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
	for line := range strings.SplitSeq(string(data), "\n") {
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
	for seg := range strings.SplitSeq(name, "/") {
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
