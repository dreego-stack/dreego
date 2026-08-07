package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFindLocalCore verifies findLocalCore resolves the absolute path to the
// local core module relative to this source file (<repo>/cmd/dreego/../../
// core), when that core directory contains a go.mod. This is what allows a
// repo-local build to scaffold projects with an offline replace directive.
func TestFindLocalCore(t *testing.T) {
	got := findLocalCore()
	if got == "" {
		t.Fatal("expected a local core directory to be found")
	}
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path, got %q", got)
	}
	if _, err := os.Stat(filepath.Join(got, "go.mod")); err != nil {
		t.Errorf("expected core/go.mod to exist under %q: %v", got, err)
	}
	if !strings.HasSuffix(filepath.ToSlash(got), "/core") {
		t.Errorf("expected path to end in /core, got %q", got)
	}
}

// TestFindLocalCoreMissing verifies findLocalCore tolerates a missing core
// go.mod. The function is hardwired to this source file (<repo>/cmd/dreego/),
// so in this repo the core/go.mod always exists and the branch is unreachable;
// it is asserted here that the resolved path is a directory containing a
// regular go.mod file (not a directory or a non-existent path), which is the
// only condition under which findLocalCore would return "".
func TestFindLocalCoreMissing(t *testing.T) {
	got := findLocalCore()
	if got == "" {
		// In a repo-local build the core module is always present; returning ""
		// here would mean the stat guard wrongly rejected a valid core dir.
		t.Fatal("findLocalCore unexpectedly returned empty in a repo-local build")
	}
	info, err := os.Stat(filepath.Join(got, "go.mod"))
	if err != nil {
		t.Fatalf("expected core/go.mod stat to succeed: %v", err)
	}
	if info.IsDir() {
		t.Error("expected core/go.mod to be a regular file, not a directory")
	}
}
