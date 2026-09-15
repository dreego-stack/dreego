package tests

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestCLIInitListTemplates(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "init", ".", "-l")
	if err != nil {
		t.Fatalf("init -l: %v\n%s", err, out)
	}
	if !strings.Contains(out, "web-minimal") {
		t.Fatalf("init -l must list web-minimal, got: %s", out)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "Taskfile.yml")); statErr == nil {
		t.Error("init -l must not scaffold files")
	}
}

func TestCLIInitSelectedTemplate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "init", ".", "-t", "web-minimal")
	if err != nil {
		t.Fatalf("init -t web-minimal: %v\n%s", err, out)
	}
	for _, f := range []string{"main.go", "www/routes/+page.dreego"} {
		if _, statErr := os.Stat(filepath.Join(dir, f)); statErr != nil {
			t.Fatalf("missing %s after init: %v", f, statErr)
		}
	}
}

func TestCLIInitUnknownTemplateFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "init", ".", "-t", "does-not-exist")
	if err == nil {
		t.Fatal("expected non-zero exit for an unknown template")
	}
	if !strings.Contains(out, "does-not-exist") || !strings.Contains(out, "web-minimal") {
		t.Fatalf("unknown-template error must name the value and valid names, got: %s", out)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "Taskfile.yml")); statErr == nil {
		t.Error("a failed init must not scaffold files")
	}
}

func TestCLIInitScaffoldClean(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", ".", "-t", "web-minimal"); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}
		if filepath.Base(path) == "template.json" {
			t.Errorf("scaffolded tree must not contain template.json: %s", path)
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), "§$name$§") {
			t.Errorf("scaffolded tree contains an unresolved §$name$§ placeholder: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk scaffold: %v", err)
	}
}

func TestCLINewListTemplates(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "new", "myapp", "--list")
	if err != nil {
		t.Fatalf("new --list: %v\n%s", err, out)
	}
	if !strings.Contains(out, "web-minimal") {
		t.Fatalf("new --list must list web-minimal, got: %s", out)
	}
	if _, statErr := os.Stat(filepath.Join(dir, "myapp")); statErr == nil {
		t.Error("new --list must not scaffold a project")
	}
}
