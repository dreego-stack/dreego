package dreefile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateStaticUsesStandardContentType(t *testing.T) {
	root := t.TempDir()
	staticDir := filepath.Join(root, "static")
	if err := os.MkdirAll(staticDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staticDir, "robots.txt"), []byte("User-agent: *\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(staticDir, "sitemap.xml"), []byte("<urlset/>"), 0o644); err != nil {
		t.Fatal(err)
	}

	src, _, err := generateStaticAssets(root, "", map[string]bool{})
	if err != nil {
		t.Fatalf("generateStaticAssets: %v", err)
	}

	if !strings.Contains(src, `app.RegisterStatic("/robots.txt", "text/plain; charset=utf-8"`) {
		t.Fatalf("robots.txt not registered as text/plain:\n%s", src)
	}
	if !strings.Contains(src, `app.RegisterStatic("/sitemap.xml", "text/xml; charset=utf-8"`) {
		t.Fatalf("sitemap.xml not registered as text/xml:\n%s", src)
	}
}
