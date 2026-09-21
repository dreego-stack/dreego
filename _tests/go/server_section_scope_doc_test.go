package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// The <server> scope contract must stay documented: route files share one Go
// package, the leading declaration block is hoisted, and shared declarations
// must be unique across route files.
func TestServerSectionDocScope(t *testing.T) {
	t.Parallel()
	root, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "_docs", "server-section.md"))
	if err != nil {
		t.Fatalf("read server-section.md: %v", err)
	}
	text := strings.ToLower(string(data))
	for _, want := range []string{
		"one go package",
		"leading declaration block",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("server-section.md must document %q", want)
		}
	}
	if !strings.Contains(text, "must be unique") && !strings.Contains(text, "cannot each declare") {
		t.Error("server-section.md must document that shared declarations must be unique across route files")
	}
}
