package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFetchDocLocalPluginDiscovery verifies that dreego docs resolves a
// plugin doc from the local plugins/<name>/_docs/ directory without any
// HTTP call. Priority: local plugin docs win over embedded/remote.
func TestFetchDocLocalPluginDiscovery(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "sample", "_docs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# Sample Plugin\n\nLocal plugin docs."
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	oldRoot := pluginDocsRoot
	pluginDocsRoot = root
	defer func() { pluginDocsRoot = oldRoot }()

	body, fromLocal, err := fetchDocLocal("/plugins/sample/_docs/index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fromLocal {
		t.Fatal("expected fromLocal=true (local plugin doc should win)")
	}
	if string(body) != content {
		t.Errorf("expected body %q, got %q", content, string(body))
	}
}

// TestFetchDocLocalReadsFilesystem verifies the discovery reads the doc
// from the filesystem (plugins/<name>/_docs/) rather than importing any
// plugin package. A plugin without a matching _docs file is not found.
func TestFetchDocLocalReadsFilesystem(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "auth", "_docs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "# Auth\n\nAuth plugin docs."
	if err := os.WriteFile(filepath.Join(dir, "guide.md"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	oldRoot := pluginDocsRoot
	pluginDocsRoot = root
	defer func() { pluginDocsRoot = oldRoot }()

	body, fromLocal, err := fetchDocLocal("/plugins/auth/_docs/guide.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fromLocal {
		t.Fatal("expected fromLocal=true")
	}
	if string(body) != content {
		t.Errorf("expected %q, got %q", content, string(body))
	}

	// A plugin with no _docs file must not be found (no error).
	_, fromLocal, err = fetchDocLocal("/plugins/missing/_docs/index.md")
	if err != nil {
		t.Fatalf("unexpected error for missing plugin: %v", err)
	}
	if fromLocal {
		t.Fatal("expected fromLocal=false for missing plugin doc")
	}
}

// TestFetchDocLocalFallback verifies that fetchDocLocal does NOT invoke the
// fallback internally. For a non-plugin path (or a plugin with no local
// _docs file) it returns (nil, false, nil); the caller decides whether to
// call the fallback. Priority: local plugin docs → embedded/remote.
func TestFetchDocLocalFallback(t *testing.T) {
	oldRoot := pluginDocsRoot
	pluginDocsRoot = t.TempDir() // empty: no local plugin docs
	defer func() { pluginDocsRoot = oldRoot }()

	fallbackCalls := 0
	oldFallback := fetchDocFallback
	fetchDocFallback = func(path string) ([]byte, error) {
		fallbackCalls++
		return []byte("remote:" + path), nil
	}
	defer func() { fetchDocFallback = oldFallback }()

	body, fromLocal, err := fetchDocLocal("/_docs/index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fromLocal {
		t.Fatal("expected fromLocal=false (no local plugin doc)")
	}
	if body != nil {
		t.Errorf("expected nil body, got %q", string(body))
	}
	if fallbackCalls != 0 {
		t.Errorf("fetchDocLocal must not call the fallback internally, got %d calls", fallbackCalls)
	}
}

// TestCmdDocsFallbackCalledOnce verifies the caller (cmdDocs) invokes the
// fallback exactly once when no local plugin doc exists. This guards against
// the regression where fetchDocLocal called the fallback internally AND the
// caller called it again, producing two HTTP requests.
func TestCmdDocsFallbackCalledOnce(t *testing.T) {
	oldRoot := pluginDocsRoot
	pluginDocsRoot = t.TempDir() // empty: no local plugin docs
	defer func() { pluginDocsRoot = oldRoot }()

	fallbackCalls := 0
	oldFallback := fetchDocFallback
	fetchDocFallback = func(path string) ([]byte, error) {
		fallbackCalls++
		return []byte("remote:" + path), nil
	}
	defer func() { fetchDocFallback = oldFallback }()

	cmdDocs([]string{"/_docs/index.md"})

	if fallbackCalls != 1 {
		t.Errorf("expected fallback called exactly once, got %d calls", fallbackCalls)
	}
}

// TestFetchDocLocalPriorityLocalWins verifies the priority order: when a
// local plugin doc exists, it wins over the embedded/remote fallback.
func TestFetchDocLocalPriorityLocalWins(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "sample", "_docs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "index.md"), []byte("local"), 0644); err != nil {
		t.Fatal(err)
	}

	oldRoot := pluginDocsRoot
	pluginDocsRoot = root
	defer func() { pluginDocsRoot = oldRoot }()

	oldFallback := fetchDocFallback
	fetchDocFallback = func(path string) ([]byte, error) {
		return []byte("remote"), nil
	}
	defer func() { fetchDocFallback = oldFallback }()

	body, fromLocal, err := fetchDocLocal("/plugins/sample/_docs/index.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !fromLocal {
		t.Fatal("expected fromLocal=true (local must win over fallback)")
	}
	if string(body) != "local" {
		t.Errorf("expected local body, got %q", string(body))
	}
}
