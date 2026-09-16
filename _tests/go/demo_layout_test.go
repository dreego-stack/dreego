package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDemoApplicationsUseOneDirectoryEach(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	for _, name := range []string{"demo-ssr", "demo-wailsv3"} {
		for _, file := range []string{"main.go", "go.mod"} {
			path := filepath.Join(repoRoot, "demo", name, file)
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("required demo file %s: %v", path, err)
			}
		}
		for _, path := range []string{"dreego.config.json", "routes", "components"} {
			if _, err := os.Stat(filepath.Join(repoRoot, "demo", name, path)); !os.IsNotExist(err) {
				t.Fatalf("Dreego content must not live beside %s/main.go: %s", name, path)
			}
		}
	}
	for _, path := range []string{
		"demo-ssr/www/dreego.config.json",
		"demo-wailsv3/app/dreego.config.json",
	} {
		if _, err := os.Stat(filepath.Join(repoRoot, "demo", path)); err != nil {
			t.Fatalf("required nested Dreego application %s: %v", path, err)
		}
	}
	for _, path := range []string{"main.go", "go.mod", "wails-v3", "cmd"} {
		if _, err := os.Stat(filepath.Join(repoRoot, "demo", path)); !os.IsNotExist(err) {
			t.Fatalf("legacy demo path %s still exists", path)
		}
	}
}
