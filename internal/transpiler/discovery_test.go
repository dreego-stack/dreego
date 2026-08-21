package transpiler

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsSkippedDir(t *testing.T) {
	cases := map[string]bool{
		"vendor":        true,
		"node_modules":  true,
		".git":          true,
		".worktrees":    true,
		"dreego":        false,
		"dreego/routes": false,
		"build":         false,
		"tmp":           false,
		".tmp":          true,
	}
	for in, want := range cases {
		if got := isSkippedDir(in); got != want {
			t.Errorf("isSkippedDir(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestIsWebsiteRoot(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "www")
	if err := os.MkdirAll(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if isWebsiteRoot(sub) {
		t.Error("dir without config must not be a website root")
	}
	if err := os.WriteFile(filepath.Join(sub, configFileName), []byte("{}"), 0644); err != nil {
		t.Fatal(err)
	}
	if !isWebsiteRoot(sub) {
		t.Error("dir with config must be a website root")
	}
}

func TestSanitizePkgName(t *testing.T) {
	cases := map[string]string{
		"www":     "www",
		"my-site": "my_site",
		"my.site": "my_site",
		"123app":  "pkg123app",
		"":        "app",
		"über":    "",
	}
	for in, want := range cases {
		got := sanitizePkgName(in)
		if want == "" {
			if got == "" {
				t.Errorf("sanitizePkgName(%q) must not be empty", in)
			}
			continue
		}
		if got != want {
			t.Errorf("sanitizePkgName(%q) = %q, want %q", in, got, want)
		}
	}
}
