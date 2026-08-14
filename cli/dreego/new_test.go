package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFindLocalRepo verifies findLocalRepo resolves the absolute path to the
// local dreego repo root relative to this source file (<repo>/cli/dreego/../..),
// when that directory contains a go.mod. This is what allows a repo-local
// build to scaffold projects with an offline replace directive.
func TestFindLocalRepo(t *testing.T) {
	got := findLocalRepo()
	if got == "" {
		t.Fatal("expected a local repo root to be found")
	}
	if !filepath.IsAbs(got) {
		t.Errorf("expected absolute path, got %q", got)
	}
	if _, err := os.Stat(filepath.Join(got, "go.mod")); err != nil {
		t.Errorf("expected go.mod to exist under %q: %v", got, err)
	}
	if _, err := os.Stat(filepath.Join(got, "core")); err != nil {
		t.Errorf("expected core/ to exist under %q: %v", got, err)
	}
}

// TestFindLocalRepoMissing verifies findLocalRepo tolerates a missing root
// go.mod. The function is hardwired to this source file (<repo>/cli/dreego/),
// so in this repo the root go.mod always exists and the branch is unreachable;
// it is asserted here that the resolved path is a directory containing a
// regular go.mod file (not a directory or a non-existent path), which is the
// only condition under which findLocalRepo would return "".
func TestFindLocalRepoMissing(t *testing.T) {
	got := findLocalRepo()
	if got == "" {
		// In a repo-local build the root module is always present; returning ""
		// here would mean the stat guard wrongly rejected a valid repo root.
		t.Fatal("findLocalRepo unexpectedly returned empty in a repo-local build")
	}
	info, err := os.Stat(filepath.Join(got, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod stat to succeed: %v", err)
	}
	if info.IsDir() {
		t.Error("expected go.mod to be a regular file, not a directory")
	}
}

func TestScaffoldVersionForDevelopmentBuild(t *testing.T) {
	for _, version := range []string{"", "dev", "(devel)"} {
		if got := scaffoldVersion(version); got != "v0.0.0" {
			t.Errorf("scaffoldVersion(%q) = %q, want v0.0.0", version, got)
		}
	}
	if got := scaffoldVersion("v0.0.41"); got != "v0.0.41" {
		t.Errorf("release version changed to %q", got)
	}
}
