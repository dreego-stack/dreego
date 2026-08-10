package dreegotest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

var buildMu sync.Mutex

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
	mainGo := "package main\nimport _ \"t/dreego/gen\"\nfunc main() {}\n"
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

	buildMu.Lock()
	defer buildMu.Unlock()
	old, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if err := os.Chdir(dir); err != nil {
		return "", err
	}
	defer os.Chdir(old)

	if err := dreego.Run(false); err != nil {
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
