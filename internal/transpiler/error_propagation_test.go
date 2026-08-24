package transpiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestProject(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func runInDir(t *testing.T, dir string) error {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	return Run(false)
}

func TestRunAbortsOnSourceReadFailure(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"www/dreego.config.json":   "{}",
		"www/routes/get.dreego":    "<body><p>ok</p></body>",
		"www/routes/broken.dreego": "",
	})
	target := filepath.Join(dir, "www", "routes", "broken.dreego")
	os.Remove(target)
	if err := os.Symlink(filepath.Join(dir, "missing-target"), target); err != nil {
		t.Skipf("cannot create broken symlink: %v", err)
	}

	err := runInDir(t, dir)
	if err == nil {
		t.Fatal("expected generation to abort on unreadable source")
	}
	if !strings.Contains(err.Error(), "broken.dreego") {
		t.Errorf("error must contain the affected path, got %q", err)
	}
}

func TestRunAbortsOnReadDirFailure(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"www/dreego.config.json": "{}",
		"www/routes/get.dreego":  "<body><p>ok</p></body>",
	})
	secretDir := filepath.Join(dir, "www", "routes", "secret")
	if err := os.MkdirAll(secretDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "missing-target"), filepath.Join(secretDir, "get.dreego")); err != nil {
		t.Fatal(err)
	}

	err := runInDir(t, dir)
	if err == nil {
		t.Fatal("expected generation to abort on unreadable source under secret")
	}
	if !strings.Contains(err.Error(), "secret") {
		t.Errorf("error must contain the affected path, got %q", err)
	}
}

func TestRunAbortsOnLayoutReadFailure(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"www/dreego.config.json": "{}",
		"www/routes/get.dreego":  "<body><p>ok</p></body>",
	})
	layoutDir := filepath.Join(dir, "www", "layouts")
	if err := os.MkdirAll(layoutDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "missing-target"), filepath.Join(layoutDir, "default.dreego")); err != nil {
		t.Fatal(err)
	}

	err := runInDir(t, dir)
	if err == nil {
		t.Fatal("expected generation to abort on unreadable layout")
	}
	if !strings.Contains(err.Error(), "default.dreego") {
		t.Errorf("error must contain the affected path, got %q", err)
	}
}
