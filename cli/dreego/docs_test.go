package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestTree(t *testing.T, files map[string]string) string {
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
	return root
}

func writeGoMod(t *testing.T, root, module string, reqs map[string]string) {
	t.Helper()
	var b strings.Builder
	b.WriteString("module " + module + "\n\ngo 1.22\n")
	if len(reqs) > 0 {
		b.WriteString("require (\n")
		for m, v := range reqs {
			b.WriteString("\t" + m + " " + v + "\n")
		}
		b.WriteString(")\n")
	}
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(b.String()), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestParseGoMod(t *testing.T) {
	root := t.TempDir()
	writeGoMod(t, root, "example.com/app", map[string]string{
		"github.com/dreego-stack/dreego":     "v0.0.27",
		"github.com/dreego-stack/plugin-sse": "v0.1.0",
	})
	gm, err := parseGoMod(filepath.Join(root, "go.mod"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gm.Module != "example.com/app" {
		t.Errorf("expected module example.com/app, got %q", gm.Module)
	}
	if gm.Requires["github.com/dreego-stack/dreego"] != "v0.0.27" {
		t.Errorf("unexpected core version: %q", gm.Requires["github.com/dreego-stack/dreego"])
	}
	if gm.Requires["github.com/dreego-stack/plugin-sse"] != "v0.1.0" {
		t.Errorf("unexpected plugin version: %q", gm.Requires["github.com/dreego-stack/plugin-sse"])
	}
}

func TestFindModDirSelfRepo(t *testing.T) {
	root := writeTestTree(t, map[string]string{"go.mod": ""})
	writeGoMod(t, root, coreModule, nil)
	dir, err := findModDir(root, coreModule)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != root {
		t.Errorf("expected self repo dir %q, got %q", root, dir)
	}
}

func TestFindModDirVendorWins(t *testing.T) {
	root := writeTestTree(t, map[string]string{
		"vendor/github.com/dreego-stack/plugin-sse/_docs/index.md": "# Plugin\n",
	})
	writeGoMod(t, root, "example.com/app", map[string]string{
		"github.com/dreego-stack/plugin-sse": "v0.1.0",
	})
	dir, err := findModDir(root, "github.com/dreego-stack/plugin-sse")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(dir, "vendor/github.com/dreego-stack/plugin-sse") {
		t.Errorf("expected vendor dir, got %q", dir)
	}
}

func TestFindModDirCacheFallback(t *testing.T) {
	root := writeTestTree(t, map[string]string{})
	writeGoMod(t, root, "example.com/app", map[string]string{
		"github.com/dreego-stack/plugin-sse": "v0.1.0",
	})
	oldCache := modCacheDir
	modCacheDir = filepath.Join(root, "cache")
	defer func() { modCacheDir = oldCache }()
	os.MkdirAll(filepath.Join(root, "cache", "github.com", "dreego-stack", "plugin-sse@v0.1.0"), 0755)

	dir, err := findModDir(root, "github.com/dreego-stack/plugin-sse")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(dir, "cache/github.com/dreego-stack/plugin-sse@v0.1.0") {
		t.Errorf("expected cache dir, got %q", dir)
	}
}

func TestFindModDirNotInGoMod(t *testing.T) {
	root := writeTestTree(t, map[string]string{})
	writeGoMod(t, root, "example.com/app", nil)
	if _, err := findModDir(root, "github.com/dreego-stack/plugin-sse"); err == nil {
		t.Fatal("expected error for module not in go.mod")
	}
}

func TestReadDocFromDir(t *testing.T) {
	root := writeTestTree(t, map[string]string{
		"_docs/cli.md": "# CLI\n",
	})
	body, err := readDocFrom(root, "/_docs/cli.md")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(body) != "# CLI\n" {
		t.Errorf("expected body, got %q", string(body))
	}
}

func TestReadDocFromDirRejectsTraversal(t *testing.T) {
	root := writeTestTree(t, map[string]string{
		"_docs/cli.md": "# CLI\n",
	})
	for _, p := range []string{"../../../etc/passwd", "/../secret", ".."} {
		if _, err := readDocFrom(root, p); err == nil {
			t.Errorf("expected traversal error for path %q", p)
		}
	}
}

func TestReadSitemap(t *testing.T) {
	root := writeTestTree(t, map[string]string{
		"_docs/sitemap.json": `{"module":"github.com/dreego-stack/dreego","pages":[{"path":"/_docs/index.md","title":"Index"}]}`,
	})
	sm, err := readSitemap(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sm.Pages) != 1 || sm.Pages[0].Path != "/_docs/index.md" || sm.Pages[0].Title != "Index" {
		t.Errorf("unexpected sitemap: %+v", sm)
	}
}

func TestSitemapPathsMissingFile(t *testing.T) {
	root := writeTestTree(t, map[string]string{})
	if got := sitemapPaths(root); got != nil {
		t.Errorf("expected nil paths for missing sitemap, got %v", got)
	}
}

func TestPrintJSON(t *testing.T) {
	body := `# Title

Some text with a [link](/docs) and another [external](https://example.com).

` + "```go\nfmt.Println(\"hi\")\n```" + `

## Section
`
	got := captureStdout(t, func() { printJSON("https://github.com/example/repo/blob/main", "/_docs/index.md", []byte(body)) })

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

func TestCmdDumpAll(t *testing.T) {
	root := writeTestTree(t, map[string]string{
		"_docs/sitemap.json": `{"module":"github.com/dreego-stack/dreego","pages":[{"path":"/_docs/index.md"},{"path":"/README.md"},{"path":"/CHANGELOG.md"}]}`,
		"_docs/index.md":     "# Index\n",
		"README.md":          "# Readme\n",
		"CHANGELOG.md":       "# Changelog\n",
	})
	got := captureStdout(t, func() {
		cmdDump(root, "/_docs/index.md", "https://github.com/dreego-stack/dreego/blob/main", "https://raw.githubusercontent.com/dreego-stack/dreego/main")
	})
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

func TestCmdList(t *testing.T) {
	root := writeTestTree(t, map[string]string{
		"go.mod":             "",
		"_docs/sitemap.json": `{"module":"github.com/dreego-stack/dreego","pages":[{"path":"/_docs/index.md","title":"Index"}]}`,
		"_docs/index.md":     "# Index\n",
	})
	writeGoMod(t, root, "github.com/dreego-stack/dreego", map[string]string{
		"github.com/dreego-stack/plugin-sse": "v0.1.0",
	})
	os.MkdirAll(filepath.Join(root, "vendor", "github.com", "dreego-stack", "plugin-sse", "_docs"), 0755)
	os.WriteFile(filepath.Join(root, "vendor", "github.com", "dreego-stack", "plugin-sse", "_docs", "sitemap.json"), []byte(`{"module":"github.com/dreego-stack/plugin-sse","pages":[{"path":"/_docs/index.md","title":"Plugin Index"}]}`), 0644)

	oldWd := wdFunc
	wdFunc = func() string { return root }
	defer func() { wdFunc = oldWd }()

	got := captureStdout(t, cmdList)
	if !strings.Contains(got, "[github.com/dreego-stack/dreego]") {
		t.Errorf("expected core section, got:\n%s", got)
	}
	if !strings.Contains(got, "[github.com/dreego-stack/plugin-sse]") {
		t.Errorf("expected plugin section, got:\n%s", got)
	}
}

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
