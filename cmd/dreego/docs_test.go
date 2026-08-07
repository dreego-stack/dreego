package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
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

// TestPrintJSON verifies printJSON emits a JSON document with the extracted
// headings, code blocks and links, plus the source path and text length.
func TestPrintJSON(t *testing.T) {
	body := `# Title

Some text with a [link](/docs) and another [external](https://example.com).

` + "```go\nfmt.Println(\"hi\")\n```" + `

## Section
`
	got := captureStdout(t, func() { printJSON("/_docs/index.md", []byte(body)) })

	var doc map[string]any
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("printJSON output is not valid JSON: %v\n%s", err, got)
	}
	if doc["path"] != "/_docs/index.md" {
		t.Errorf("expected path, got %v", doc["path"])
	}
	if doc["length"] != float64(len(body)) {
		t.Errorf("expected length %d, got %v", len(body), doc["length"])
	}

	headings, ok := doc["headings"].([]any)
	if !ok || len(headings) != 2 || headings[0] != "Title" || headings[1] != "Section" {
		t.Errorf("expected 2 headings [Title Section], got %v", doc["headings"])
	}

	codeBlocks, ok := doc["code_blocks"].([]any)
	if !ok || len(codeBlocks) != 1 {
		t.Errorf("expected 1 code block, got %v", doc["code_blocks"])
	}

	links, ok := doc["links"].([]any)
	if !ok || len(links) != 2 {
		t.Errorf("expected 2 links, got %v", doc["links"])
	}
}

// TestCmdDumpAll verifies cmdDump with "all" iterates the central doc pages and
// prints each preceded by a --- <path> --- separator. Uses the embedded docs.
func TestCmdDumpAll(t *testing.T) {
	oldRoot := pluginDocsRoot
	pluginDocsRoot = t.TempDir() // empty: no local plugin docs override
	defer func() { pluginDocsRoot = oldRoot }()

	got := captureStdout(t, func() { cmdDump("all") })
	if !strings.Contains(got, "--- /_docs/index.md ---") {
		t.Errorf("expected index separator in dump, got:\n%s", got)
	}
	if !strings.Contains(got, "--- /README.md ---") {
		t.Errorf("expected README separator in dump, got:\n%s", got)
	}
	if !strings.Contains(got, "--- /CHANGELOG.md ---") {
		t.Errorf("expected CHANGELOG separator in dump, got:\n%s", got)
	}
}

// TestCmdDocsJSON verifies cmdDocs --json routes through printJSON and emits
// valid JSON with headings for the requested doc path.
func TestCmdDocsJSON(t *testing.T) {
	oldRoot := pluginDocsRoot
	pluginDocsRoot = t.TempDir() // empty: no local plugin docs override
	defer func() { pluginDocsRoot = oldRoot }()

	got := captureStdout(t, func() { cmdDocs([]string{"--json", "/_docs/index.md"}) })
	var doc map[string]any
	if err := json.Unmarshal([]byte(got), &doc); err != nil {
		t.Fatalf("cmdDocs --json output is not valid JSON: %v\n%s", err, got)
	}
	if doc["path"] != "/_docs/index.md" {
		t.Errorf("expected path /_docs/index.md, got %v", doc["path"])
	}
	headings, ok := doc["headings"].([]any)
	if !ok || len(headings) == 0 {
		t.Errorf("expected non-empty headings, got %v", doc["headings"])
	}
}

// captureStdout runs fn and returns everything written to os.Stdout. A reader
// goroutine drains the pipe concurrently so functions that write more than the
// pipe buffer do not deadlock on a full pipe.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	done := make(chan string, 1)
	go func() {
		buf, _ := io.ReadAll(r)
		done <- string(buf)
	}()

	fn()
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() { os.Stdout = old }()
	return <-done
}
