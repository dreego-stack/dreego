package transpiler

import (
	"os"
	"path/filepath"
	"strings"
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

func TestRouteFileRelSupportsFlatAndLegacyRoutes(t *testing.T) {
	root := filepath.Join("/tmp", "site")
	cases := map[string]string{
		"routes/about.dreego":            "about",
		"routes/page.dreego":             "",
		"routes/users/[id]/page.dreego":   "users/[id]",
		"routes/+page.dreego":            "",
		"routes/users/[id]/+page.dreego":  "users/[id]",
		"routes/(auth)/login.dreego":     "(auth)/login",
		"routes/post.dreego":             "",
		"routes/admin/post-users.dreego": "admin",
	}
	for path, want := range cases {
		dir := filepath.Dir(filepath.Join(root, path))
		name := filepath.Base(path)
		if got := routeFileRel(root, dir, name); got != want {
			t.Errorf("routeFileRel(%q) = %q, want %q", path, got, want)
		}
	}
}

func TestScanRoutesGeneratesFlatPatternsAndRejectsDuplicates(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/about.dreego":            "<body>about</body>",
		"routes/page.dreego":             "<body>home</body>",
		"routes/users/[id]/page.dreego":  "<body>user</body>",
		"routes/(auth)/login.dreego":     "<body>login</body>",
	})
	dirs, _, count, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Fatalf("scanRoutes count = %d, want 4", count)
	}
	src := strings.Join(dirs[0].regs, "")
	for _, want := range []string{"app.Register(\"GET\", \"/about\"", "app.Register(\"GET\", \"/{$}\"", "app.Register(\"GET\", \"/users/{id}\"", "app.Register(\"GET\", \"/login\""} {
		if !strings.Contains(src, want) {
			t.Errorf("generated source missing %q", want)
		}
	}

	duplicateRoot := writeTestProject(t, map[string]string{
		"routes/about.dreego":        "<body>one</body>",
		"routes/(auth)/about.dreego": "<body>two</body>",
	})
	_, _, _, err = scanRoutes(NewGenerator(), duplicateRoot, map[string]*layoutEntry{})
	if err == nil || !strings.Contains(err.Error(), "about.dreego") {
		t.Fatalf("expected duplicate source paths in error, got %v", err)
	}

	pageLegacyConflict := writeTestProject(t, map[string]string{
		"routes/page.dreego": "<body>page</body>",
		"routes/get.dreego":  "<body>legacy</body>",
	})
	_, _, _, err = scanRoutes(NewGenerator(), pageLegacyConflict, map[string]*layoutEntry{})
	if err == nil || !strings.Contains(err.Error(), "duplicate route") {
		t.Fatalf("expected duplicate route error for page.dreego + get.dreego, got %v", err)
	}
}
