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

func TestFindModDirWithoutGoMod(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("GOMODCACHE", cache)
	root := t.TempDir()
	if _, err := findModDir(root, "github.com/dreego-stack/dreego"); err == nil {
		t.Fatal("expected an error when there is no go.mod and no cached module")
	}
}

func TestModuleCacheDirPicksCachedVersion(t *testing.T) {
	cache := t.TempDir()
	t.Setenv("GOMODCACHE", cache)
	versioned := filepath.Join(cache, "github.com", "dreego-stack", "dreego@v0.10.3")
	if err := os.MkdirAll(versioned, 0o755); err != nil {
		t.Fatal(err)
	}
	got, ok := moduleCacheDir("github.com/dreego-stack/dreego", "")
	if !ok || got != versioned {
		t.Fatalf("moduleCacheDir without version = %q, %v; want %q", got, ok, versioned)
	}
	got, ok = moduleCacheDir("github.com/dreego-stack/dreego", "v0.10.3")
	if !ok || got != versioned {
		t.Fatalf("moduleCacheDir with version = %q, %v; want %q", got, ok, versioned)
	}
}

func TestEscapeModulePathUppercase(t *testing.T) {
	if got := escapeModulePath("github.com/Acme/App"); got != "github.com/!acme/!app" {
		t.Fatalf("escapeModulePath = %q", got)
	}
}
