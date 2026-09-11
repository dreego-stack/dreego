package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestModuleBoundaries(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	modules := map[string]string{
		"go.mod":               "module github.com/dreego-stack/dreego\n",
		"core/go.mod":          "module github.com/dreego-stack/dreego/core\n",
		"adapter/ssr/go.mod":   "module github.com/dreego-stack/dreego/adapter/ssr\n",
		"adapter/wails/go.mod": "module github.com/dreego-stack/dreego/adapter/wails\n",
		"dreegotest/go.mod":    "module github.com/dreego-stack/dreego/dreegotest\n",
		"cmd/dreego/go.mod":    "module github.com/dreego-stack/dreego/cmd/dreego\n",
	}
	for path, declaration := range modules {
		contents, err := os.ReadFile(filepath.Join(repoRoot, path))
		if err != nil {
			t.Errorf("read %s: %v", path, err)
			continue
		}
		if !strings.Contains(string(contents), declaration) {
			t.Errorf("%s does not declare %q", path, strings.TrimSpace(declaration))
		}
		if strings.Contains(string(contents), "replace github.com/dreego-stack/dreego") {
			t.Errorf("published module %s contains a repository-local replacement", path)
		}
	}
	requirements := map[string][]string{
		"core/go.mod":          {"github.com/dreego-stack/dreego v0.8.0"},
		"adapter/ssr/go.mod":   {"github.com/dreego-stack/dreego v0.8.0", "github.com/dreego-stack/dreego/core v0.8.0"},
		"adapter/wails/go.mod": {"github.com/dreego-stack/dreego/core v0.8.0"},
		"dreegotest/go.mod":    {"github.com/dreego-stack/dreego v0.8.0", "github.com/dreego-stack/dreego/core v0.8.0"},
		"cmd/dreego/go.mod":    {"github.com/dreego-stack/dreego v0.8.0"},
	}
	for path, expected := range requirements {
		contents, err := os.ReadFile(filepath.Join(repoRoot, path))
		if err != nil {
			t.Fatal(err)
		}
		for _, requirement := range expected {
			if !strings.Contains(string(contents), requirement) {
				t.Errorf("%s is missing coordinated requirement %q", path, requirement)
			}
		}
	}
	for _, path := range []string{
		"core/_docs",
		"adapter/ssr/_docs",
		"adapter/wails/_docs",
		"dreegotest/_docs",
		"cmd/dreego/_docs",
	} {
		if info, err := os.Stat(filepath.Join(repoRoot, path)); err != nil || !info.IsDir() {
			t.Errorf("required module documentation directory %s is missing", path)
		}
	}
	for _, path := range []string{"core/ssr", "cli/dreego", "target"} {
		if _, err := os.Stat(filepath.Join(repoRoot, path)); !os.IsNotExist(err) {
			t.Errorf("removed v0.7 path %s still exists", path)
		}
	}
	oldImports := []string{
		"github.com/dreego-stack/dreego/core/ssr",
		"github.com/dreego-stack/dreego/cli/dreego",
	}
	if err := filepath.WalkDir(repoRoot, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil || entry.IsDir() || filepath.Ext(path) != ".go" {
			return walkErr
		}
		if filepath.Base(path) == "module_boundaries_test.go" {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		for _, oldImport := range oldImports {
			if strings.Contains(string(contents), oldImport) {
				t.Errorf("removed import %q remains in %s", oldImport, path)
			}
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}

func TestWailsAdapterKeepsApplicationOwnershipExplicit(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	adapterModule, err := os.ReadFile(filepath.Join(repoRoot, "adapter/wails/go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(adapterModule), "github.com/wailsapp/wails") {
		t.Fatal("adapter/wails must not depend on Wails")
	}
	mainSource, err := os.ReadFile(filepath.Join(repoRoot, "demo/demo-wailsv3/main.go"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ownership := range []string{"application.New", "application.NewService", "Window.NewWithOptions", "wailsApp.Run"} {
		if !strings.Contains(string(mainSource), ownership) {
			t.Errorf("demo main.go does not visibly own %s", ownership)
		}
	}
}
