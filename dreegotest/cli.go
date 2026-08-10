package dreegotest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

var (
	cliOnce  sync.Once
	cliBin   string
	cliErr   error
	cliBuild string
)

// ProjectDir creates a temp project directory with a go.mod (module t) that
// replaces the dreego module with the repo root, a main.go importing the
// generated package, and the given dreego/… files. It returns the directory
// path. It replaces shell tests that scaffold a temp module and run the CLI
// inside it.
func ProjectDir(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()

	repoRoot, err := RepoRoot()
	if err != nil {
		t.Fatalf("ProjectDir: %v", err)
	}

	goMod := fmt.Sprintf("module t\ngo 1.22\nrequire github.com/dreego-stack/dreego v0.0.0\nreplace github.com/dreego-stack/dreego => %s\n", repoRoot)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("ProjectDir: write go.mod: %v", err)
	}
	mainGo := "package main\nimport _ \"t/dreego/gen\"\nfunc main() {}\n"
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatalf("ProjectDir: write main.go: %v", err)
	}
	writeFiles(t, dir, files)
	return dir
}

// CLIBin builds the dreego CLI once per test binary and returns its path. It
// replaces shell tests that `go build -o … ./cli/dreego` themselves. The
// version is injected via ldflags (latest git tag) so `dreego new` writes a
// valid require directive.
func CLIBin(t *testing.T) string {
	t.Helper()
	cliOnce.Do(func() {
		repoRoot, err := RepoRoot()
		if err != nil {
			cliErr = err
			return
		}
		dir, err := os.MkdirTemp("", "dreego-cli")
		if err != nil {
			cliErr = err
			return
		}
		cliBuild = dir
		bin := filepath.Join(dir, "dreego")

		tag := latestTag(repoRoot)
		ldflags := "-X main.version=" + tag
		build := exec.Command("go", "build", "-ldflags", ldflags, "-o", bin, "./cli/dreego")
		build.Dir = repoRoot
		if out, err := build.CombinedOutput(); err != nil {
			cliErr = fmt.Errorf("build CLI: %v\n%s", err, out)
			return
		}
		cliBin = bin
	})
	if cliErr != nil {
		t.Fatalf("CLIBin: %v", cliErr)
	}
	return cliBin
}

func latestTag(repoRoot string) string {
	// Prefer DREEGO_VERSION (set by make test / test.sh and the Dockerfile) so
	// tests behave identically inside the container (where git is absent).
	if v := os.Getenv("DREEGO_VERSION"); v != "" {
		return v
	}
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = repoRoot
	if out, err := cmd.Output(); err == nil {
		if tag := strings.TrimSpace(string(out)); tag != "" {
			return tag
		}
	}
	return "dev"
}

// LatestTag returns the latest git tag of the repo (as injected into the CLI
// build), or "dev" if no tag exists. Used to assert the CLI version matches the
// current tag.
func LatestTag(t *testing.T) string {
	t.Helper()
	repoRoot, err := RepoRoot()
	if err != nil {
		t.Fatalf("LatestTag: %v", err)
	}
	return latestTag(repoRoot)
}

// RunCLI runs the dreego CLI with the given arguments in dir and returns the
// combined stdout/stderr output and the exit error. It replaces shell tests
// that call $DREEGO_BIN … and grep the output.
func RunCLI(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	return runIn(t, dir, CLIBin(t), args...)
}

// RunCLIIn runs bin (e.g. dreegotest.CLIBin(t)) with args in dir and returns
// combined output and exit error. Used when the CLI must run in an arbitrary
// directory that is not a generated project dir (e.g. the repo root).
func RunCLIIn(t *testing.T, dir, bin string, args ...string) (string, error) {
	t.Helper()
	return runIn(t, dir, bin, args...)
}

// runIn runs bin with args in dir and returns combined output and exit error.
func runIn(t *testing.T, dir, bin string, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// MustBuildInDir runs `go build` in dir, failing the test on error. It
// replaces shell tests that call `go build -o /dev/null .`.
func MustBuildInDir(t *testing.T, dir string) {
	t.Helper()
	if ok := BuildInDirOK(t, dir); !ok {
		t.Fatal("MustBuildInDir: go build failed")
	}
}

// BuildInDirOK runs `go build` in dir and returns whether it succeeded. It
// replaces shell tests that assert on the go build exit status.
func BuildInDirOK(t *testing.T, dir string) bool {
	t.Helper()
	build := exec.Command("go", "build", "-o", "/dev/null", ".")
	build.Dir = dir
	return build.Run() == nil
}

func writeFiles(t *testing.T, dir string, files map[string]string) {
	t.Helper()
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("writeFiles: mkdir %s: %v", filepath.Dir(full), err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatalf("writeFiles: write %s: %v", path, err)
		}
	}
}
