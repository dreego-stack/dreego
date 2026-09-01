package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestBlueprintUsesSSRHost asserts the generated default blueprint main starts
// the app through the explicit SSR host instead of calling app.Listen directly.
func TestBlueprintUsesSSRHost(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, "cli", "dreego", "blueprints", "default", "main.go.tmpl"))
	if err != nil {
		t.Fatalf("read default main.go.tmpl: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "ssr.Listen(app, addr)") {
		t.Fatalf("default blueprint does not use ssr.Listen(app, addr):\n%s", content)
	}
	if strings.Contains(content, "app.Listen(") {
		t.Fatalf("default blueprint still calls app.Listen directly:\n%s", content)
	}
}

// TestSSRTargetServes builds and serves a minimal app through the full pipeline
// to prove the SSR host path still serves requests.
func TestSSRTargetServes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego": `<body><h1>ssr target</h1></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<h1>ssr target</h1>") {
		t.Fatalf("body missing expected content: %s", body)
	}
}
