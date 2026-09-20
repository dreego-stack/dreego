package tests

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestCLINewSelectedTemplate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.NewProject(t, "app", "web-minimal")
	for _, f := range []string{"main.go", "www/routes/+page.dreego"} {
		if _, statErr := os.Stat(filepath.Join(dir, f)); statErr != nil {
			t.Fatalf("missing %s after new: %v", f, statErr)
		}
	}
}

func TestCLINewWebAppTemplate(t *testing.T) {
	t.Parallel()
	dir := dreegotest.NewProject(t, "app", "web-app")
	for _, f := range []string{
		"main.go",
		"www/layouts/default.dreego",
		"www/components/Nav.dreego",
		"www/components/Card.dreego",
		"www/routes/+page.dreego",
		"www/routes/dashboard/+page.dreego",
	} {
		if _, statErr := os.Stat(filepath.Join(dir, f)); statErr != nil {
			t.Fatalf("missing %s after new -t web-app: %v", f, statErr)
		}
	}
	if _, statErr := os.Stat(filepath.Join(dir, "template.json")); statErr == nil {
		t.Error("scaffolded web-app tree must not contain template.json")
	}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if strings.Contains(string(data), "§$name$§") {
			t.Errorf("scaffolded web-app tree contains an unresolved §$name$§ placeholder: %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk web-app scaffold: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate in web-app scaffold: %v\n%s", err, out)
	}
	if !dreegotest.BuildInDirOK(t, dir) {
		t.Fatal("web-app scaffold must build")
	}
}

func TestCLINewUnknownTemplateFails(t *testing.T) {
	t.Parallel()
	parent := t.TempDir()
	out, err := dreegotest.RunCLI(t, parent, "new", "app", "-t", "does-not-exist")
	if err == nil {
		t.Fatal("expected non-zero exit for an unknown template")
	}
	if !strings.Contains(out, "does-not-exist") || !strings.Contains(out, "web-minimal") {
		t.Fatalf("unknown-template error must name the value and valid names, got: %s", out)
	}
	if _, statErr := os.Stat(filepath.Join(parent, "app")); statErr == nil {
		t.Error("a failed new must not scaffold files")
	}
}

func TestCLINewScaffoldClean(t *testing.T) {
	t.Parallel()
	dir := dreegotest.NewProject(t, "app", "web-minimal")
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
	for _, name := range []string{"web-app", "web-minimal"} {
		if !strings.Contains(out, name) {
			t.Fatalf("new --list must list %s, got: %s", name, out)
		}
	}
	if _, statErr := os.Stat(filepath.Join(dir, "myapp")); statErr == nil {
		t.Error("new --list must not scaffold a project")
	}
}
