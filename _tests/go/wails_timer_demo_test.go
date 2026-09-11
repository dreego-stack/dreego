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
	source, err := os.ReadFile(filepath.Join(repoRoot, "demo", "demo-wailsv3", "app", "routes", "+page.dreego"))
	if err != nil {
		t.Fatalf("read timer demo: %v", err)
	}
	about, err := os.ReadFile(filepath.Join(repoRoot, "demo", "demo-wailsv3", "app", "routes", "about", "+page.dreego"))
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
		`<title>Focus — Dreego Timer</title>`,
		`role="timer"`,
		`aria-label="Timer controls"`,
		`role="status"`,
		`class="timer-card"`,
		`class="timer-orbit"`,
		`class="shortcut-hint"`,
		`aria-describedby="timer-caption"`,
		`--progress`,
		`prefers-reduced-motion: reduce`,
		`button:focus-visible`,
		`bindings/demo-wailsv3/app/timerservice.js`,
		`toggle.addEventListener`,
		`reset.addEventListener`,
	} {
		dreegotest.MustContain(t, routes, fragment)
	}
	for _, path := range []string{
		"demo/demo-wailsv3/app/bindings/demo-wailsv3/app/models.ts",
		"demo/demo-wailsv3/app/bindings/demo-wailsv3/app/timerservice.ts",
		"demo/demo-wailsv3/app/static/bindings/demo-wailsv3/app/timerservice.js",
	} {
		if _, err := os.Stat(filepath.Join(repoRoot, path)); err != nil {
			t.Fatalf("required generated binding %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "demo", "demo-wailsv3", "package.json")); !os.IsNotExist(err) {
		t.Fatalf("Wails timer must not own package.json: %v", err)
	}
}
