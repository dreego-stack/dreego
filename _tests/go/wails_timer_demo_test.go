package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestWailsTimerDemoGeneratesAccessibleDocument(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	source, err := os.ReadFile(filepath.Join(repoRoot, "demo", "wails-v3", "routes", "+page.dreego"))
	if err != nil {
		t.Fatalf("read timer demo: %v", err)
	}
	about, err := os.ReadFile(filepath.Join(repoRoot, "demo", "wails-v3", "routes", "about", "+page.dreego"))
	if err != nil {
		t.Fatalf("read timer about page: %v", err)
	}
	files := map[string]string{
		"www/routes/+page.dreego":       string(source),
		"www/routes/about/+page.dreego": string(about),
	}
	generated := dreegotest.Build(t, files)
	routes := generated["www/routes/dree.go"]
	regenerated := dreegotest.Build(t, files)["www/routes/dree.go"]
	if routes != regenerated {
		t.Fatal("restart-based generation produced different route output")
	}
	for _, fragment := range []string{
		`app.RegisterRender("/", PageIndex())`,
		`app.RegisterRender("/about", PageAbout())`,
		`<title>Dreego Timer</title>`,
		`role="timer"`,
		`aria-label="Timer controls"`,
		`role="status"`,
		`button:focus-visible`,
		`bindings/demo/wails-v3/timerservice.js`,
		`toggle.addEventListener`,
		`reset.addEventListener`,
	} {
		dreegotest.MustContain(t, routes, fragment)
	}
	for _, path := range []string{
		"demo/wails-v3/bindings/demo/wails-v3/models.ts",
		"demo/wails-v3/bindings/demo/wails-v3/timerservice.ts",
		"demo/wails-v3/static/bindings/demo/wails-v3/timerservice.js",
	} {
		if _, err := os.Stat(filepath.Join(repoRoot, path)); err != nil {
			t.Fatalf("required generated binding %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "demo", "wails-v3", "package.json")); !os.IsNotExist(err) {
		t.Fatalf("Wails timer must not own package.json: %v", err)
	}
}
