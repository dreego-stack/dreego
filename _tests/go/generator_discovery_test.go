package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDiscoveryIgnoresRoutesOutsideProjectRoot(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":                    `{}`,
		"www/routes/get.dreego":                     `<div><p>real</p></div>`,
		"vendor/somepkg/dreego/routes/get.dreego":   `<div><p>vendor</p></div>`,
		"node_modules/foo/dreego/routes/get.dreego": `<div><p>nm</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	gen, err := os.ReadFile(filepath.Join(dir, "www", "routes", "dree.go"))
	if err != nil {
		t.Fatalf("read routes dree.go: %v", err)
	}
	if strings.Contains(string(gen), "vendor") || strings.Contains(string(gen), "node_modules") {
		t.Fatalf("discovery leaked out-of-root routes into gen: %s", gen)
	}
	if !strings.Contains(string(gen), "real") {
		t.Fatalf("expected in-root route content, got: %s", gen)
	}
}

func TestDiscoveryIgnoresNestedDreegoProjectRoots(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":          `{}`,
		"www/routes/get.dreego":           `<div><p>outer</p></div>`,
		"subapp/dreego/routes/get.dreego": `<div><p>inner</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	gen, err := os.ReadFile(filepath.Join(dir, "www", "routes", "dree.go"))
	if err != nil {
		t.Fatalf("read routes dree.go: %v", err)
	}
	if strings.Contains(string(gen), "inner") {
		t.Fatalf("nested subapp routes leaked into outer project: %s", gen)
	}
	if !strings.Contains(string(gen), "outer") {
		t.Fatalf("expected outer route content, got: %s", gen)
	}
}

func TestDiscoveryIgnoresComponentsOutsideProjectRoot(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":                    `{}`,
		"www/components/Inner.dreego":               "Component Inner ()\n<div><p>inner</p></div>",
		"www/routes/get.dreego":                     `<div><@Inner/></div>`,
		"vendor/lib/dreego/components/Outer.dreego": "Component Outer ()\n<div><p>outer</p></div>",
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	gen, err := os.ReadFile(filepath.Join(dir, "www", "components", "dree.go"))
	if err != nil {
		t.Fatalf("read components dree.go: %v", err)
	}
	if strings.Contains(string(gen), "outer") || strings.Contains(string(gen), "Outer") {
		t.Fatalf("out-of-root component leaked: %s", gen)
	}
	if !strings.Contains(string(gen), "Inner") {
		t.Fatalf("expected in-root component, got: %s", gen)
	}
}

func TestDiscoveryLayoutLocalCascades(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":                 `{}`,
		"www/layouts/default.dreego":             `<div><html><body><nav>Root</nav>{#slot}</body></html></div>`,
		"www/routes/get.dreego":                  `<div><p>home</p></div>`,
		"www/routes/blog/get.dreego":             `<div><p>blog</p></div>`,
		"www/routes/blog/layouts/default.dreego": `<div><html><body><nav>Blog</nav>{#slot}</body></html></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	gen, err := os.ReadFile(filepath.Join(dir, "www", "layouts", "dree.go"))
	if err != nil {
		t.Fatalf("read layouts dree.go: %v", err)
	}
	if !strings.Contains(string(gen), "Root") {
		t.Fatalf("root layout not applied to root route: %s", gen)
	}
	if !strings.Contains(string(gen), "Blog") {
		t.Fatalf("route-local layout not applied to nested route: %s", gen)
	}
}

func TestDiscoveryAmbiguousLayoutsFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":     `{}`,
		"www/layouts/default.dreego": `<div><html><body>{#slot}</body></html></div>`,
		"www/layouts/layout.dreego":  `<div><html><body><nav>X</nav>{#slot}</body></html></div>`,
		"www/routes/get.dreego":      `<div><p>home</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected failure on ambiguous layout files, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "ambiguous") && !strings.Contains(strings.ToLower(out), "layout") {
		t.Fatalf("expected ambiguous layout diagnostic, got: %s", out)
	}
}

func TestDiscoveryAmbiguousErrorPagesFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":        `{}`,
		"www/routes/404.dreego":         `<div><p>root 404</p></div>`,
		"www/routes/(group)/404.dreego": `<div><p>group 404</p></div>`,
		"www/routes/get.dreego":         `<div><p>home</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("expected error for duplicate catch-all 404 pattern via group, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "duplicate") {
		t.Fatalf("expected duplicate diagnostic, got: %s", out)
	}
}

func TestDiscovery404And500SameDirGenerateAndCheck(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json": `{}`,
		"www/routes/404.dreego":  `<div><p>not found</p></div>`,
		"www/routes/500.dreego":  `<div><p>server error</p></div>`,
		"www/routes/get.dreego":  `<div><p>home</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate with 404+500 in same dir: %v\n%s", err, out)
	}
	gen, err := os.ReadFile(filepath.Join(dir, "www", "routes", "dree.go"))
	if err != nil {
		t.Fatalf("read routes dree.go: %v", err)
	}
	if !strings.Contains(string(gen), "SetErrorHandler(500") {
		t.Fatalf("expected SetErrorHandler(500) registration, got: %s", gen)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err != nil {
		t.Fatalf("check should pass with 404+500 in same dir: %v\n%s", err, out)
	}
}

func TestDiscoveryDeterministicAddMoveDelete(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json": `{}`,
		"www/routes/get.dreego":  `<div><p>home</p></div>`,
	})
	gen1 := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{}`,
		"www/routes/get.dreego":       `<div><p>home</p></div>`,
		"www/routes/about/get.dreego": `<div><p>about</p></div>`,
	})
	if err := os.MkdirAll(filepath.Join(dir, "www", "routes", "about"), 0755); err != nil {
		t.Fatalf("mkdir about: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "www", "routes", "about", "get.dreego"), []byte(`<div><p>about</p></div>`), 0644); err != nil {
		t.Fatalf("write about/get.dreego: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate after add: %v\n%s", err, out)
	}
	gen2, err := os.ReadFile(filepath.Join(dir, "www", "routes", "dree.go"))
	if err != nil {
		t.Fatalf("read routes dree.go: %v", err)
	}
	if !strings.Contains(string(gen2), "about") {
		t.Fatalf("added route not discovered: %s", gen2)
	}
	if !strings.Contains(string(gen1["www/routes/dree.go"]), "home") {
		t.Fatalf("baseline missing home: %s", gen1["www/routes/dree.go"])
	}
	if err := os.Rename(filepath.Join(dir, "www", "routes", "about"), filepath.Join(dir, "www", "routes", "info")); err != nil {
		t.Fatalf("move about -> info: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate after move: %v\n%s", err, out)
	}
	gen3, err := os.ReadFile(filepath.Join(dir, "www", "routes", "dree.go"))
	if err != nil {
		t.Fatalf("read routes dree.go after move: %v", err)
	}
	if strings.Contains(string(gen3), "HandleAbout") {
		t.Fatalf("stale About handler remained after move: %s", gen3)
	}
	if !strings.Contains(string(gen3), "HandleInfo") {
		t.Fatalf("moved route not discovered as Info: %s", gen3)
	}
	if err := os.RemoveAll(filepath.Join(dir, "www", "routes", "info")); err != nil {
		t.Fatalf("remove info: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate after delete: %v\n%s", err, out)
	}
	gen4, err := os.ReadFile(filepath.Join(dir, "www", "routes", "dree.go"))
	if err != nil {
		t.Fatalf("read routes dree.go after delete: %v", err)
	}
	if strings.Contains(string(gen4), "Info") || strings.Contains(string(gen4), "About") {
		t.Fatalf("deleted route handlers remained: %s", gen4)
	}
}

func TestDiscoveryGeneratesRoutesNamedStaticAndGen(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":            `{}`,
		"www/routes/get.dreego":             `<div><p>home</p></div>`,
		"www/routes/static/get.dreego":      `<div><p>static route</p></div>`,
		"www/routes/gen/get.dreego":         `<div><p>gen route</p></div>`,
		"www/routes/blog/static/get.dreego": `<div><p>blog static</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	gen, err := os.ReadFile(filepath.Join(dir, "www", "routes", "dree.go"))
	if err != nil {
		t.Fatalf("read routes dree.go: %v", err)
	}
	for _, want := range []string{"HandleStatic", "HandleGen", "HandleBlogStatic"} {
		if !strings.Contains(string(gen), want) {
			t.Fatalf("expected %s handler in generated routes dree.go: %s", want, gen)
		}
	}
}
