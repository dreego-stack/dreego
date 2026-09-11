package main

import (
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

func TestFindModDirUsesGoReplacement(t *testing.T) {
	root := t.TempDir()
	moduleDir := filepath.Join(root, "core-src")
	if err := os.MkdirAll(moduleDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(moduleDir, "go.mod"), []byte("module github.com/dreego-stack/dreego/core\n\ngo 1.27\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	goMod := "module example.com/app\n\ngo 1.27\n\nrequire github.com/dreego-stack/dreego/core v0.8.0\nreplace github.com/dreego-stack/dreego/core => ./core-src\n"
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte(goMod), 0o644); err != nil {
		t.Fatal(err)
	}
	dir, err := findModDir(root, "github.com/dreego-stack/dreego/core")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dir != moduleDir {
		t.Fatalf("resolved directory = %q, want %q", dir, moduleDir)
	}
}

func TestModuleVersionUsesBuildInformation(t *testing.T) {
	info := &debug.BuildInfo{
		Main: debug.Module{Path: "github.com/dreego-stack/dreego/cmd/dreego", Version: "v0.8.0"},
		Deps: []*debug.Module{{Path: "github.com/dreego-stack/dreego", Version: "v0.8.0"}},
	}
	if got := moduleVersion(info, "github.com/dreego-stack/dreego/cmd/dreego"); got != "v0.8.0" {
		t.Fatalf("main module version = %q", got)
	}
	if got := moduleVersion(info, "github.com/dreego-stack/dreego"); got != "v0.8.0" {
		t.Fatalf("dependency module version = %q", got)
	}
}
