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
	generated := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": string(source),
	})
	routes := generated["www/routes/dree.go"]
	for _, fragment := range []string{
		`app.RegisterRender("/", PageIndex())`,
		`<title>Dreego Timer</title>`,
		`role="timer"`,
		`aria-label="Timer controls"`,
		`role="status"`,
		`button:focus-visible`,
		`toggle.addEventListener`,
		`reset.addEventListener`,
	} {
		dreegotest.MustContain(t, routes, fragment)
	}
}
