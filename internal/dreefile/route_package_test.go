package dreefile

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRoutePackageDirMergesInvalidSegments(t *testing.T) {
	root := filepath.Join("/tmp", "site")
	cases := map[string]string{
		"routes":                 "routes",
		"routes/about":           "routes/about",
		"routes/blog/[id]":       "routes/blog",
		"routes/(auth)":          "routes",
		"routes/users/[id]/edit": "routes/users",
		"routes/admin-panel":     "routes/admin-panel",
		"routes/my-site/v2":      "routes/my-site/v2",
	}
	for in, want := range cases {
		got := routePackageDir(root, filepath.Join(root, filepath.FromSlash(in)))
		wantPath := filepath.Join(root, filepath.FromSlash(want))
		if got != wantPath {
			t.Errorf("routePackageDir(%q) = %q, want %q", in, got, wantPath)
		}
	}
}

func TestRoutePackageNamesComeFromFolder(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/+page.dreego":              "<body>home</body>",
		"routes/about/+page.dreego":        "<body>about</body>",
		"routes/my-site/+page.dreego":      "<body>site</body>",
		"routes/blog/[id]/+page.dreego":    "<body>post</body>",
		"routes/(auth)/login/+page.dreego": "<body>login</body>",
	})
	dirs, _, count, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 5 {
		t.Fatalf("scanRoutes count = %d, want 5", count)
	}
	byRel := map[string]*routePkg{}
	for _, d := range dirs {
		byRel[d.rel] = d
	}
	want := map[string]string{
		"":        "routes",
		"about":   "about",
		"my-site": "my_site",
		"blog":    "blog",
	}
	for rel, pkg := range want {
		d, ok := byRel[rel]
		if !ok {
			t.Fatalf("missing route package for rel %q; got %v", rel, keysOf(byRel))
		}
		if d.pkg != pkg {
			t.Errorf("route package %q has pkg %q, want %q", rel, d.pkg, pkg)
		}
	}
}

func TestRootRoutePackageCollectsSubPackages(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/+page.dreego":           "<body>home</body>",
		"routes/about/+page.dreego":     "<body>about</body>",
		"routes/blog/[id]/+page.dreego": "<body>post</body>",
	})
	gen := NewGenerator()
	gen.Module = "example.com/app"
	gen.RootRel = "www"
	dirs, _, _, err := scanRoutes(gen, root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err != nil {
		t.Fatal(err)
	}
	var rootFile string
	for _, d := range dirs {
		if d.rel == "" {
			rootFile, err = buildRoutePackageFile(gen, d)
			if err != nil {
				t.Fatal(err)
			}
		}
	}
	if rootFile == "" {
		t.Fatal("no root route package generated")
	}
	for _, want := range []string{
		"package routes",
		"about \"example.com/app/www/routes/about\"",
		"blog \"example.com/app/www/routes/blog\"",
		"about.Register(app)",
		"blog.Register(app)",
		"func Register(app *dreego.App) error",
	} {
		if !strings.Contains(rootFile, want) {
			t.Errorf("root route package missing %q:\n%s", want, rootFile)
		}
	}
	if _, err := parser.ParseFile(token.NewFileSet(), "dree.go", rootFile, 0); err != nil {
		t.Fatalf("root route package is not valid Go: %v\n%s", err, rootFile)
	}
}

func TestRoutePackagesAllowDuplicateDeclarationsAcrossFolders(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/registrierung/+page.dreego": "<server>\ntype Form struct{ Name string }\n</server>\n<body>{{ (Form{Name: \"a\"}).Name }}</body>",
		"routes/anmeldung/+page.dreego":     "<server>\ntype Form struct{ Name string }\n</server>\n<body>{{ (Form{Name: \"b\"}).Name }}</body>",
	})
	dirs, _, _, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err != nil {
		t.Fatalf("same-named declarations in different route folders must not conflict: %v", err)
	}
	if len(dirs) < 3 {
		t.Fatalf("expected root plus two sub-packages, got %d", len(dirs))
	}
}

func TestRoutePackagesConflictWithinSameFolder(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/a.dreego": "<server>\ntype Form struct{ Name string }\n</server>\n<body>{{ (Form{Name: \"a\"}).Name }}</body>",
		"routes/b.dreego": "<server>\ntype Form struct{ Name string }\n</server>\n<body>{{ (Form{Name: \"b\"}).Name }}</body>",
	})
	_, _, _, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err == nil || !strings.Contains(err.Error(), "Form") {
		t.Fatalf("same-folder duplicate declaration must conflict with a diagnostic naming it, got %v", err)
	}
}

func TestRunSplitsRouteFoldersIntoPackages(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"www/dreego.config.json":                "{}",
		"www/routes/+page.dreego":               "<body>home</body>",
		"www/routes/registrierung/+page.dreego": "<server>\ntype Form struct{ Email string }\nfunc label() string { return \"register\" }\n</server>\n<body>{{ label() }}{{ (Form{Email: \"a\"}).Email }}</body>",
		"www/routes/anmeldung/+page.dreego":     "<server>\ntype Form struct{ Email string }\nfunc label() string { return \"login\" }\n</server>\n<body>{{ label() }}{{ (Form{Email: \"b\"}).Email }}</body>",
	})
	if err := runInDir(t, dir); err != nil {
		t.Fatalf("generation failed: %v", err)
	}

	read := func(rel string) string {
		t.Helper()
		data, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(rel)))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		return string(data)
	}

	root := read("www/routes/dree.go")
	if !isGeneratedFile(root) {
		t.Fatalf("root routes file must carry a generated marker:\n%s", root)
	}
	if !strings.Contains(root, "\npackage routes\n") {
		t.Fatalf("root routes file must be package routes:\n%s", root)
	}
	for _, want := range []string{"registrierung.Register(app)", "anmeldung.Register(app)"} {
		if !strings.Contains(root, want) {
			t.Errorf("root collector missing %q:\n%s", want, root)
		}
	}
	if strings.Contains(root, "type Form struct") || strings.Contains(root, "func label(") {
		t.Errorf("sub-folder glue must not leak into the root package:\n%s", root)
	}

	for folder, pkg := range map[string]string{"registrierung": "registrierung", "anmeldung": "anmeldung"} {
		src := read("www/routes/" + folder + "/dree.go")
		if !strings.Contains(src, "\npackage "+pkg+"\n") {
			t.Errorf("%s/dree.go must be package %s:\n%s", folder, pkg, src)
		}
		for _, want := range []string{"type Form struct", "func label() string"} {
			if !strings.Contains(src, want) {
				t.Errorf("%s/dree.go missing its own glue %q:\n%s", folder, want, src)
			}
		}
	}
}

func keysOf[T any](m map[string]T) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}
