package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindMainRoot(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	projDir, pkg, name := findMainIn(dir)
	if projDir != "." || pkg != "." || name != filepath.Base(dir) {
		t.Errorf("expected root main, got %q %q %q", projDir, pkg, name)
	}
}

func TestFindMainCmd(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "cmd"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cmd", "main.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	projDir, pkg, name := findMainIn(dir)
	if projDir != "cmd" || pkg != "cmd" || name != "cmd" {
		t.Errorf("expected cmd main, got %q %q %q", projDir, pkg, name)
	}
}

func TestFindMainFallback(t *testing.T) {
	dir := t.TempDir()
	projDir, pkg, name := findMainIn(dir)
	if projDir != "." || pkg != "." || name != "server" {
		t.Errorf("expected fallback server, got %q %q %q", projDir, pkg, name)
	}
}
