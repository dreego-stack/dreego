package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestModuleComponentImport(t *testing.T) {
	for _, root := range []string{"www", "shop", "app", "gibberish"} {
		t.Run(root, func(t *testing.T) {
			t.Parallel()
			testModuleComponentImport(t, root)
		})
	}
}

func testModuleComponentImport(t *testing.T, root string) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		root + "/routes/+page.dreego": `from "example.com/ui/components" import {
    Button,
}
<body><@Button label="From module"/></body>`,
	})
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nimport (\n\tsite \"t/"+root+"\"\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\nfunc main() { app := dreego.New(); if err := site.Register(app); err != nil { panic(err) } }\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	moduleDir := filepath.Join(dir, "external-ui")
	if err := os.MkdirAll(filepath.Join(moduleDir, "components"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module example.com/ui\n\ngo 1.25\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "components", "Button.dreego"), []byte(`Component Button (label string)
<body><button>{{ label }}</button></body>`), 0o644); err != nil {
		t.Fatal(err)
	}

	goMod, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	goMod = append(goMod, []byte("\nrequire example.com/ui v0.0.0\nreplace example.com/ui => ./external-ui\n")...)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), goMod, 0o644); err != nil {
		t.Fatal(err)
	}
	if output, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate failed: %v\n%s", err, output)
	}
	if _, err := os.Stat(filepath.Join(dir, root, "components", "dree.go")); err != nil {
		t.Fatalf("generated component package for %s: %v", root, err)
	}
	build := exec.Command("go", "test", "./...")
	build.Dir = dir
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("generated project failed to compile: %v\n%s", err, output)
	}
}
