package gomod

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadParsesModuleAndRequirements(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	source := `module "example.com/app"

go 1.25

require (
	"example.com/direct" v1.2.3
	example.com/indirect v1.0.0 // indirect
)
`
	if err := os.WriteFile(path, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}

	file, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if file.Module != "example.com/app" {
		t.Fatalf("Module = %q", file.Module)
	}
	if file.Requires["example.com/direct"] != "v1.2.3" {
		t.Fatalf("direct requirement = %q", file.Requires["example.com/direct"])
	}
	if file.Requires["example.com/indirect"] != "v1.0.0" {
		t.Fatalf("indirect requirement = %q", file.Requires["example.com/indirect"])
	}
}

func TestReadRejectsInvalidGoMod(t *testing.T) {
	path := filepath.Join(t.TempDir(), "go.mod")
	if err := os.WriteFile(path, []byte("module\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Read(path); err == nil {
		t.Fatal("Read accepted invalid go.mod")
	}
}
