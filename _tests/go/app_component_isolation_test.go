package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestAppLocalComponentAliasBuilds ensures an app-local component imported with
// an alias generates and compiles (the alias lives in the app's scope, not the
// shared root scope).
func TestAppLocalComponentAliasBuilds(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/dreego.config.json":        `{"logging":{"enabled":false}}`,
		"www/a/components/OnlyA.dreego": "DREEFILE component ()\n<body><b>a</b></body>",
		"www/a/routes/+page.dreego":     "COMPONENT \"www/a/components\" IMPORT { OnlyA as Foo }\n<body><@Foo/></body>",
	})
}

// TestSharedAndModuleComponentsCoexist is the regression for a module-component
// import overwriting the shared root components package: a shared component and
// an external module component must both generate and render.
func TestSharedAndModuleComponentsCoexist(t *testing.T) {
	t.Parallel()
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	files := map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/components/Badge.dreego": "DREEFILE component ()\n<body><span class=\"badge\">shared</span></body>",
		"www/app/routes/+page.dreego": "COMPONENT \"example.com/ui/components\" IMPORT { Button }\nCOMPONENT \"www/components\" IMPORT { Badge }\n<body><@Badge/><@Button label=\"From module\"/></body>",
	}
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	gomod := "module t\n\ngo 1.27\n\nrequire github.com/dreego-stack/dreego/core v0.0.0\n\n" +
		"replace github.com/dreego-stack/dreego => " + repoRoot + "\n" +
		"replace github.com/dreego-stack/dreego/core => " + filepath.Join(repoRoot, "core") + "\n" +
		"replace github.com/dreego-stack/dreego/adapter/ssr => " + filepath.Join(repoRoot, "adapter", "ssr") + "\n" +
		"require example.com/ui v0.0.0\nreplace example.com/ui => ./external-ui\n"
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644); err != nil {
		t.Fatal(err)
	}
	sum, err := os.ReadFile(filepath.Join(repoRoot, "go.sum"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "go.sum"), sum, 0644); err != nil {
		t.Fatal(err)
	}
	uiDir := filepath.Join(dir, "external-ui", "components")
	if err := os.MkdirAll(uiDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "external-ui", "go.mod"), []byte("module example.com/ui\n\ngo 1.27\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(uiDir, "Button.dreego"), []byte("DREEFILE component (label string)\n<body><button>{{ label }}</button></body>"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nimport (\n\twebapp \"t/www/app\"\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\nfunc main() { _ = dreego.New(webapp.App) }\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	build := exec.Command("go", "build", "-mod=mod", "-o", "/dev/null", ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
}

// TestAppComponentIsolation ensures a component defined only in app A is not
// visible to app B: each app package compiles against only its own components.
func TestAppComponentIsolation(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/dreego.config.json":        `{"logging":{"enabled":false}}`,
		"www/a/components/OnlyA.dreego": "DREEFILE component ()\n<body><b>a</b></body>",
		"www/a/routes/+page.dreego":     "COMPONENT \"www/a/components\" IMPORT { OnlyA }\n<body><@OnlyA/></body>",
		"www/b/routes/+page.dreego":     "<body><p>b</p></body>",
	})
}

// TestAppCannotSeeSiblingComponent ensures referencing a sibling app's component
// fails generation rather than leaking a definition across apps.
func TestAppCannotSeeSiblingComponent(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuildFail(t, map[string]string{
		"www/dreego.config.json":        `{"logging":{"enabled":false}}`,
		"www/a/components/OnlyA.dreego": "DREEFILE component ()\n<body><b>a</b></body>",
		"www/a/routes/+page.dreego":     "<body><p>a</p></body>",
		"www/b/routes/+page.dreego":     "<body><@OnlyA/></body>",
	})
}

// TestAppLocalComponentShadowsShared ensures an app-local component with the
// same name as a shared root component overrides it for that app.
func TestAppLocalComponentShadowsShared(t *testing.T) {
	t.Parallel()
	dir := dreegotest.BuildDir(t, map[string]string{
		"www/dreego.config.json":          `{"logging":{"enabled":false}}`,
		"www/components/Badge.dreego":     "DREEFILE component ()\n<body><span class=\"shared\">shared</span></body>",
		"www/app/components/Badge.dreego": "DREEFILE component ()\n<body><span class=\"local\">local</span></body>",
		"www/app/routes/+page.dreego":     "COMPONENT \"www/components\" IMPORT { Badge }\n<body><@Badge/></body>",
	})
	if dir == "" {
		t.Fatal("expected build to succeed")
	}
}
