package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// newTestEmbeddedFS replaces the package-level embeddedDocs with a
// t.TempDir()-backed fs.FS built from the given path->content map. This keeps
// the read-logic tests deterministic and independent of the production
// //go:embed copy step (which lives in cmd/dreego/embedded/).
func newTestEmbeddedFS(t *testing.T, files map[string]string) {
	t.Helper()
	root := t.TempDir()
	for p, content := range files {
		full := filepath.Join(root, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	old := embeddedDocs
	embeddedDocs = os.DirFS(root)
	t.Cleanup(func() { embeddedDocs = old })
}

// TestFetchDocEmbeddedReadsEmbeddedFS verifies fetchDocEmbedded reads a doc
// from the embedded filesystem (no network). The leading-slash URL path is
// mapped onto the embedded FS root.
func TestFetchDocEmbeddedReadsEmbeddedFS(t *testing.T) {
	newTestEmbeddedFS(t, map[string]string{
		"_docs/index.md": "# Dreego Documentation\n",
	})
	body, err := fetchDocEmbedded("/_docs/index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "# Dreego Documentation\n" {
		t.Errorf("expected embedded body, got %q", string(body))
	}
}

// TestFetchDocEmbeddedMissingFile verifies a missing embedded file returns an
// error (so the caller can surface it) rather than silently succeeding.
func TestFetchDocEmbeddedMissingFile(t *testing.T) {
	newTestEmbeddedFS(t, map[string]string{})
	if _, err := fetchDocEmbedded("/_docs/nope.md"); err == nil {
		t.Fatal("expected error for missing embedded file")
	}
}

// TestFetchDocFallbackPointsToEmbedded verifies the fallback loader is wired
// to the embedded loader (fetchDocEmbedded), NOT to the HTTP fetchDoc. This is
// the core guarantee that `dreego docs` works offline.
func TestFetchDocFallbackPointsToEmbedded(t *testing.T) {
	fallbackPtr := reflect.ValueOf(fetchDocFallback).Pointer()
	embeddedPtr := reflect.ValueOf(fetchDocEmbedded).Pointer()
	if fallbackPtr != embeddedPtr {
		t.Fatal("fetchDocFallback must point to fetchDocEmbedded, not fetchDoc (HTTP)")
	}
}

// TestFetchDocFallbackUsesEmbedded verifies the fallback returns embedded
// content. Because the content only exists in the local test FS, a non-zero
// body proves the fallback reads embedded docs instead of doing http.Get.
func TestFetchDocFallbackUsesEmbedded(t *testing.T) {
	newTestEmbeddedFS(t, map[string]string{
		"_docs/index.md": "# Offline Docs\n",
	})
	body, err := fetchDocFallback("/_docs/index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "# Offline Docs\n" {
		t.Errorf("expected embedded body from fallback, got %q", string(body))
	}
}

// TestFetchDocEmbeddedCentralFiles verifies the central documentation files
// (_docs/index.md, README.md, CHANGELOG.md) are reachable offline through the
// embedded loader. The production embedded FS mirrors this layout.
func TestFetchDocEmbeddedCentralFiles(t *testing.T) {
	newTestEmbeddedFS(t, map[string]string{
		"_docs/index.md": "# Dreego Documentation\n",
		"README.md":      "# dreego — Go Web Framework\n",
		"CHANGELOG.md":   "# Changelog\n",
	})
	for _, p := range []string{"/_docs/index.md", "/README.md", "/CHANGELOG.md"} {
		body, err := fetchDocEmbedded(p)
		if err != nil {
			t.Errorf("fetchDocEmbedded(%q) error: %v", p, err)
			continue
		}
		if len(body) == 0 {
			t.Errorf("fetchDocEmbedded(%q) returned empty body", p)
		}
	}
}
