package dreegotest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// MustBuild runs the full core pipeline (core.Run) in a temp module and then
// compiles it with `go build`. It replaces shell tests that do
// `dreego generate` + `go build -o /dev/null .`.
//
// files maps relative paths (e.g. "dreego/routes/get.dreego") to their content.
// The temp module gets a go.mod with a replace to the real repo root.
func MustBuild(t *testing.T, files map[string]string) {
	t.Helper()
	_, err := build(t, files, false)
	if err != nil {
		t.Fatalf("MustBuild: %v", err)
	}
}

// MustBuildFail asserts that the full pipeline fails (generate or build error).
func MustBuildFail(t *testing.T, files map[string]string) {
	t.Helper()
	_, err := build(t, files, true)
	if err == nil {
		t.Fatal("MustBuildFail: expected error, got none")
	}
}

// Build runs the full pipeline and returns the generated files (relative path
// → content) plus the temp dir. It replaces shell tests that grep the
// generated output after `dreego generate`.
func Build(t *testing.T, files map[string]string) map[string]string {
	t.Helper()
	dir, err := build(t, files, false)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	generated := map[string]string{}
	filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		generated[rel] = string(data)
		return nil
	})
	return generated
}

func BuildDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir, err := build(t, files, false)
	if err != nil {
		t.Fatalf("BuildDir: %v", err)
	}
	return dir
}

func build(t *testing.T, files map[string]string, expectFail bool) (string, error) {
	t.Helper()
	dir := t.TempDir()

	repoRoot, err := RepoRoot()
	if err != nil {
		return "", err
	}

	goMod := fmt.Sprintf("module t\ngo 1.22\nrequire github.com/dreego-stack/dreego v0.0.0\nreplace github.com/dreego-stack/dreego => %s\n", repoRoot)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		return "", err
	}
	mainGo := "package main\nimport (\n\t\"t/www\"\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\nfunc main() { app := dreego.New(); if err := www.Register(app); err != nil { panic(err) } }\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		return "", err
	}
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			return "", err
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			return "", err
		}
	}
	ensureConfig(t, dir, files)

	// Run codegen in a subprocess (the cached CLI binary) instead of in-process
	// with a global os.Chdir. This makes build() safe for t.Parallel().
	if _, err := RunCLI(t, dir, "generate"); err != nil {
		if expectFail {
			return dir, err
		}
		return "", fmt.Errorf("generate failed: %w", err)
	}

	if expectFail {
		return dir, nil
	}

	cmd := exec.Command("go", "build", "-o", "/dev/null", ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("go build failed: %v\n%s", err, out)
	}
	return dir, nil
}

// ensureConfig writes a dreego.config.json into every website root implied by
// the given files (a directory containing routes/, components/, or layouts/)
// unless the test already provided one. This keeps the new root-marker model
// transparent for tests that only care about routes/components.
func ensureConfig(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	roots := map[string]bool{}
	for path := range files {
		parts := strings.Split(filepath.ToSlash(path), "/")
		for i, p := range parts {
			if p == "routes" || p == "components" || p == "layouts" {
				root := filepath.Join(append([]string{dir}, parts[:i]...)...)
				roots[root] = true
			}
		}
	}
	for root := range roots {
		cfg := filepath.Join(root, "dreego.config.json")
		if _, err := os.Stat(cfg); os.IsNotExist(err) {
			if err := os.WriteFile(cfg, []byte("{}"), 0644); err != nil {
				t.Fatalf("ensureConfig: %v", err)
			}
		}
	}
}

func RepoRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			return "", fmt.Errorf("go.mod not found above %s", wd)
		}
		wd = parent
	}
}
