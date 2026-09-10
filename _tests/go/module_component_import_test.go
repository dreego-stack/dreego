package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestModuleComponentImport(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `from "example.com/ui/components" import {
    Button,
}
<body><@Button label="From module"/></body>`,
	})

	moduleDir := filepath.Join(dir, "external-ui")
	if err := os.MkdirAll(filepath.Join(moduleDir, "components"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module example.com/ui\n\ngo 1.22\n"), 0o644); err != nil {
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
	build := exec.Command("go", "test", "./...")
	build.Dir = dir
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("generated project failed to compile: %v\n%s", err, output)
	}
}
