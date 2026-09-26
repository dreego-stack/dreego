package tests

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// Each route folder is its own Go package. The root routes package collects the
// sub-package Register calls, and the same package-level declaration may be
// reused in two different folders.
func TestPackagePerRouteFolderSeparatesDeclarations(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/+page.dreego": `<body><h1>Home</h1></body>`,
		"www/routes/registrierung/+page.dreego": `<server>
type Form struct{ Email string }
</server>
<body><form g-action="Register" method="post">{{ (Form{Email: "a"}).Email }}</form></body>`,
		"www/routes/anmeldung/+page.dreego": `<server>
type Form struct{ Email string }
</server>
<body><form g-action="Login" method="post">{{ (Form{Email: "b"}).Email }}</form></body>`,
	})

	rootRoutes := gen["www/routes/dree.go"]
	if !strings.HasPrefix(rootRoutes, "package routes\n") {
		t.Fatalf("root routes file must be package routes:\n%s", rootRoutes)
	}
	for _, want := range []string{
		"registrierung \"t/www/routes/registrierung\"",
		"anmeldung \"t/www/routes/anmeldung\"",
		"registrierung.Register(app)",
		"anmeldung.Register(app)",
	} {
		if !strings.Contains(rootRoutes, want) {
			t.Fatalf("root routes collector missing %q:\n%s", want, rootRoutes)
		}
	}
	if strings.Contains(rootRoutes, "type Form struct") {
		t.Fatalf("sub-folder declarations must not leak into the root routes package:\n%s", rootRoutes)
	}

	registrierung := gen["www/routes/registrierung/dree.go"]
	if !strings.HasPrefix(registrierung, "package registrierung\n") {
		t.Fatalf("sub-folder must be its own package named after the folder:\n%s", registrierung)
	}
	if !strings.Contains(registrierung, "type Form struct") {
		t.Fatalf("sub-folder dree.go must contain its own package-level declaration:\n%s", registrierung)
	}
	anmeldung := gen["www/routes/anmeldung/dree.go"]
	if !strings.HasPrefix(anmeldung, "package anmeldung\n") {
		t.Fatalf("second sub-folder must be its own package:\n%s", anmeldung)
	}
	if !strings.Contains(anmeldung, "type Form struct") {
		t.Fatalf("second sub-folder must also keep its declaration:\n%s", anmeldung)
	}
}

// A folder name that is not a valid Go identifier still becomes a valid
// sanitized package name, and dynamic/group folders fold into their nearest
// valid ancestor package.
func TestPackagePerRouteFolderDynamicFoldersFoldIntoAncestor(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/routes/my-site/+page.dreego":      `<body>site</body>`,
		"www/routes/blog/[id]/+page.dreego":    `<body>post</body>`,
		"www/routes/(auth)/login/+page.dreego": `<body>login</body>`,
	})

	if _, ok := gen["www/routes/my-site/dree.go"]; !ok {
		t.Fatalf("my-site folder must get its own package file; got files %v", keys(gen))
	}
	mySite := gen["www/routes/my-site/dree.go"]
	if !strings.HasPrefix(mySite, "package my_site\n") {
		t.Fatalf("my-site must sanitize to package my_site:\n%s", mySite)
	}

	// "[id]" and "(auth)/login" are not importable directories, so their
	// handlers live in the closest valid ancestor (blog and the root routes
	// package respectively).
	blog := gen["www/routes/blog/dree.go"]
	if !strings.Contains(blog, "HandleBlogId") {
		t.Fatalf("dynamic folder must fold into ancestor blog package:\n%s", blog)
	}
	root := gen["www/routes/dree.go"]
	if !strings.Contains(root, "HandleAuthLogin") {
		t.Fatalf("group folder must fold into the root routes package:\n%s", root)
	}
	for _, path := range []string{"www/routes/blog/[id]/dree.go", "www/routes/(auth)/login/dree.go"} {
		if _, ok := gen[path]; ok {
			t.Fatalf("dynamic/group folders must not get their own Go package file: %s", path)
		}
	}
}

// A collision inside one folder is still a dreego diagnostic (not a raw
// compiler error), because one package may not redeclare a name.
func TestPackagePerRouteFolderSameFolderCollisionFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/a.dreego": "<server>\ntype Form struct{ Email string }\n</server>\n<body>{{ (Form{Email: \"a\"}).Email }}</body>",
		"www/routes/b.dreego": "<server>\ntype Form struct{ Email string }\n</server>\n<body>{{ (Form{Email: \"b\"}).Email }}</body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected a duplicate declaration error, got success: %s", out)
	}
	if !strings.Contains(out, "Form") {
		t.Fatalf("diagnostic must name the duplicate declaration, got: %s", out)
	}
}

func TestPackagePerRouteFolderGeneratedTreeBuilds(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"www/routes/+page.dreego":               `<body><h1>Home</h1></body>`,
		"www/routes/registrierung/+page.dreego": "<server>\ntype Form struct{ Email string }\n</server>\n<body><form g-action=\"Register\" method=\"post\">{{ (Form{Email: \"a\"}).Email }}</form></body>",
		"www/routes/[id]/+page.dreego":          `<server>id := c.Param("id")</server><body><p>user:{{ id }}</p></body>`,
	}
	dir := dreegotest.BuildDir(t, files)
	if _, err := os.Stat(filepath.Join(dir, "www", "routes", "registrierung", "dree.go")); err != nil {
		t.Fatalf("expected sub-package file on disk: %v", err)
	}
	if _, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err != nil {
		t.Fatalf("generate --check must be stable after the split: %v", err)
	}
}

func keys[T any](m map[string]T) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
