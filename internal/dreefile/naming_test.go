package dreefile

import (
	"path/filepath"
	"testing"
)

func TestLayoutScopeKey(t *testing.T) {
	root := "/site"
	cases := []struct {
		name string
		dir  string
		want string
		ok   bool
	}{
		{"root layouts", filepath.Join(root, "layouts"), "root", true},
		{"app layouts", filepath.Join(root, "blog", "layouts"), "blog", true},
		{"app routes layouts", filepath.Join(root, "www", "routes", "admin", "layouts"), "www/admin", true},
		{"nested route layouts", filepath.Join(root, "www", "routes", "a", "b", "layouts"), "www/a/b", true},
		{"unrelated dir", filepath.Join(root, "www", "components"), "", false},
		{"layouts deeper", filepath.Join(root, "www", "x", "layouts", "y"), "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := layoutScopeKey(root, tc.dir)
			if ok != tc.ok || got != tc.want {
				t.Fatalf("layoutScopeKey(%q) = (%q,%v), want (%q,%v)", tc.dir, got, ok, tc.want, tc.ok)
			}
		})
	}
}

func TestLayoutFuncName(t *testing.T) {
	cases := []struct {
		scopeKey string
		file     string
		want     string
	}{
		{"root", "default.dreego", "Default"},
		{"", "default.dreego", "Default"},
		{"root", "layout.dreego", "Layout"},
		{"blog", "default.dreego", "DefaultBlog"},
		{"www/admin", "default.dreego", "DefaultWwwAdmin"},
		{"www/admin", "layout.dreego", "LayoutWwwAdmin"},
	}
	for _, tc := range cases {
		if got := layoutFuncName(tc.scopeKey, tc.file); got != tc.want {
			t.Errorf("layoutFuncName(%q,%q) = %q, want %q", tc.scopeKey, tc.file, got, tc.want)
		}
	}
}

func TestComponentPkgName(t *testing.T) {
	base := "/site"
	cases := []struct {
		name   string
		pkgDir string
		want   string
	}{
		{"root components", filepath.Join(base, "components"), "components"},
		{"nested root components", filepath.Join(base, "components", "ui"), "components_ui"},
		{"app components", filepath.Join(base, "blog", "components"), "blog_components"},
		{"module components", filepath.Join(base, "components", "ui"), "components_ui"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := componentPkgName(base, tc.pkgDir); got != tc.want {
				t.Fatalf("componentPkgName(%q) = %q, want %q", tc.pkgDir, got, tc.want)
			}
		})
	}
}

func TestComponentPkgNameDistinctForSharedAndAppLocal(t *testing.T) {
	base := "/site"
	shared := componentPkgName(base, filepath.Join(base, "components"))
	appLocal := componentPkgName(base, filepath.Join(base, "blog", "components"))
	if shared == appLocal {
		t.Fatalf("shared and app-local component packages collide: %q", shared)
	}
}

func TestSanitizePkgNameEdgeCases(t *testing.T) {
	cases := map[string]string{
		"www":       "www",
		"a-b":       "a_b",
		"a.b":       "a_b",
		"a b":       "a_b",
		"2go":       "pkg2go",
		"":          "app",
		"___":       "___",
		"a_b":       "a_b",
		"CamelCase": "CamelCase",
	}
	for in, want := range cases {
		if got := sanitizePkgName(in); got != want {
			t.Errorf("sanitizePkgName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsSkippedDirExtra(t *testing.T) {
	skipped := []string{"vendor", "node_modules", ".git", ".worktrees", ".tmp", ".hidden"}
	for _, name := range skipped {
		if !isSkippedDir(name) {
			t.Errorf("isSkippedDir(%q) = false, want true", name)
		}
	}
	kept := []string{"www", "blog", "dreego", "app", "src", "cmd"}
	for _, name := range kept {
		if isSkippedDir(name) {
			t.Errorf("isSkippedDir(%q) = true, want false", name)
		}
	}
}

func TestHasDirAndIsWebsiteRoot(t *testing.T) {
	dir := t.TempDir()
	mustMkdir(t, filepath.Join(dir, "routes"))
	if !hasDir(filepath.Join(dir, "routes")) {
		t.Fatal("hasDir must be true for an existing directory")
	}
	if hasDir(filepath.Join(dir, "missing")) {
		t.Fatal("hasDir must be false for a missing directory")
	}
	mustWrite(t, filepath.Join(dir, configFileName), "{}")
	if !isWebsiteRoot(dir) {
		t.Fatal("isWebsiteRoot must be true with a config file")
	}
	if isWebsiteRoot(filepath.Join(dir, "routes")) {
		t.Fatal("isWebsiteRoot must be false without a config file")
	}
}
