package dreegotest

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"testing"
)

// Fixture copies the reference app _tests/fixtures/<name> into a temp dir and
// rewrites its go.mod replace directive to point at the repo root, so the copy
// builds offline. It returns the temp dir.
func Fixture(t *testing.T, name string) string {
	t.Helper()
	repoRoot, err := RepoRoot()
	if err != nil {
		t.Fatalf("Fixture: %v", err)
	}
	src := filepath.Join(repoRoot, "_tests", "fixtures", name)
	if _, err := os.Stat(filepath.Join(src, "go.mod")); err != nil {
		t.Fatalf("Fixture: reference app %q not found at %s", name, src)
	}
	dir := t.TempDir()
	err = filepath.WalkDir(src, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dest := filepath.Join(dir, rel)
		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(dest, data, 0644)
	})
	if err != nil {
		t.Fatalf("Fixture: copy %s: %v", name, err)
	}
	gomod := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(gomod)
	if err != nil {
		t.Fatalf("Fixture: read go.mod: %v", err)
	}
	re := regexp.MustCompile(`(?m)^replace github\.com/dreego-stack/dreego => .*$`)
	rewritten := re.ReplaceAllString(string(data), "replace github.com/dreego-stack/dreego => "+repoRoot)
	if err := os.WriteFile(gomod, []byte(rewritten), 0644); err != nil {
		t.Fatalf("Fixture: write go.mod: %v", err)
	}
	return dir
}

// ServeFixture runs the full CLI-to-HTTP workflow for a reference app: dreego
// generate, go build, start the server on a free port, and return a Client.
func ServeFixture(t *testing.T, name string) *Client {
	t.Helper()
	dir := Fixture(t, name)
	if _, err := RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("ServeFixture: generate %s: %v", name, err)
	}
	port := FreePort(t)
	bin := filepath.Join(dir, "server")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("ServeFixture: go build %s: %v\n%s", name, err, out)
	}
	proc := exec.Command(bin)
	proc.Dir = dir
	proc.Env = append(os.Environ(), "PORT="+fmt.Sprintf("%d", port))
	if err := proc.Start(); err != nil {
		t.Fatalf("ServeFixture: start %s: %v", name, err)
	}
	t.Cleanup(func() {
		proc.Process.Kill()
		proc.Wait()
	})
	WaitForPort(t, port)
	return &Client{base: fmt.Sprintf("http://127.0.0.1:%d", port), jar: newJar()}
}
