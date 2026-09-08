package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDocumentationSitemapReferencesExistingUniqueFiles(t *testing.T) {
	t.Parallel()
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(repoRoot, "_docs", "sitemap.json"))
	if err != nil {
		t.Fatal(err)
	}
	var sitemap struct {
		Pages []struct {
			Path  string `json:"path"`
			Title string `json:"title"`
		} `json:"pages"`
	}
	if err := json.Unmarshal(data, &sitemap); err != nil {
		t.Fatalf("parse sitemap: %v", err)
	}
	seen := make(map[string]bool, len(sitemap.Pages))
	for _, page := range sitemap.Pages {
		if page.Path == "" || page.Title == "" {
			t.Fatalf("sitemap page needs path and title: %+v", page)
		}
		if seen[page.Path] {
			t.Fatalf("duplicate sitemap path %q", page.Path)
		}
		seen[page.Path] = true
		path := filepath.Join(repoRoot, filepath.FromSlash(strings.TrimPrefix(page.Path, "/")))
		if info, err := os.Stat(path); err != nil || info.IsDir() {
			t.Errorf("sitemap path %q does not reference a file", page.Path)
		}
	}
}
