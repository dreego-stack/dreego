package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

// Testable core of the `dreego dev` file watcher. These functions are pure and
// deterministic so the change-detection + restart logic can be unit-tested
// without a real background daemon.
//
// detectChanges(dir, mtimes) compares the current .dreego file modtimes against
// the previous map and returns the changed files (relative paths) plus the
// updated mtime map. A file is "changed" when it is new, its modtime moved, or
// it disappeared from the previous map.
//
// shouldRestart(changed) reports whether a server restart is required. Any
// .dreego change needs codegen + rebuild, so it is true whenever changed is
// non-empty.

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func mtimeOf(t *testing.T, dir, name string) time.Time {
	t.Helper()
	info, err := os.Stat(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return info.ModTime()
}

func TestDetectChangesNewFile(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "page.dreego", "page / {}\n")

	changed, updated := detectChanges(dir, map[string]time.Time{})

	if !reflect.DeepEqual(changed, []string{"page.dreego"}) {
		t.Errorf("expected new file detected, got %v", changed)
	}
	if _, ok := updated["page.dreego"]; !ok {
		t.Errorf("expected updated map to contain page.dreego, got %v", updated)
	}
}

func TestDetectChangesNoChange(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "page.dreego", "page / {}\n")
	prev := map[string]time.Time{"page.dreego": mtimeOf(t, dir, "page.dreego")}

	changed, updated := detectChanges(dir, prev)

	if len(changed) != 0 {
		t.Errorf("expected no changes, got %v", changed)
	}
	if !reflect.DeepEqual(updated, prev) {
		t.Errorf("expected unchanged map, got %v", updated)
	}
}

func TestDetectChangesModified(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "page.dreego", "page / {}\n")
	prev := map[string]time.Time{"page.dreego": mtimeOf(t, dir, "page.dreego")}

	future := prev["page.dreego"].Add(2 * time.Second)
	if err := os.Chtimes(filepath.Join(dir, "page.dreego"), future, future); err != nil {
		t.Fatal(err)
	}

	changed, updated := detectChanges(dir, prev)

	if !reflect.DeepEqual(changed, []string{"page.dreego"}) {
		t.Errorf("expected modified file detected, got %v", changed)
	}
	if !updated["page.dreego"].Equal(future) {
		t.Errorf("expected updated modtime, got %v", updated["page.dreego"])
	}
}

func TestDetectChangesRemoved(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "page.dreego", "page / {}\n")
	prev := map[string]time.Time{"page.dreego": mtimeOf(t, dir, "page.dreego")}

	if err := os.Remove(filepath.Join(dir, "page.dreego")); err != nil {
		t.Fatal(err)
	}

	changed, updated := detectChanges(dir, prev)

	if !reflect.DeepEqual(changed, []string{"page.dreego"}) {
		t.Errorf("expected removed file detected, got %v", changed)
	}
	if _, ok := updated["page.dreego"]; ok {
		t.Errorf("expected removed file absent from updated map, got %v", updated)
	}
}

func TestDetectChangesIgnoresNonDreego(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "page.dreego", "page / {}\n")
	writeFile(t, dir, "README.md", "# readme\n")

	changed, _ := detectChanges(dir, map[string]time.Time{})

	if !reflect.DeepEqual(changed, []string{"page.dreego"}) {
		t.Errorf("expected only .dreego file, got %v", changed)
	}
}

func TestShouldRestartTrueOnChange(t *testing.T) {
	if !shouldRestart([]string{"page.dreego"}) {
		t.Error("expected restart on .dreego change")
	}
}

func TestShouldRestartFalseOnNoChange(t *testing.T) {
	if shouldRestart(nil) {
		t.Error("expected no restart when nothing changed")
	}
	if shouldRestart([]string{}) {
		t.Error("expected no restart on empty change list")
	}
}

// TestDetectChangesInitialPriming guards Blocker 1: the watcher must prime the
// mtime map before its first tick so a second, unchanged scan reports no
// changes instead of treating every .dreego file as newly added.
func TestDetectChangesInitialPriming(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "page.dreego", "page / {}\n")

	// First scan with an empty map establishes the diff baseline. It must
	// report every file as new once, but the returned map is the baseline.
	changed, baseline := detectChanges(dir, map[string]time.Time{})
	if !reflect.DeepEqual(changed, []string{"page.dreego"}) {
		t.Fatalf("expected first scan to report new file, got %v", changed)
	}

	// Second scan with the baseline, nothing modified: no changes.
	changed2, _ := detectChanges(dir, baseline)
	if len(changed2) != 0 {
		t.Errorf("expected no changes after priming, got %v", changed2)
	}
}

// TestCmdBuildEErrorsInsteadOfExit guards Blocker 2: cmdBuildE must return an
// error rather than calling os.Exit so a build failure in the dev watcher does
// not kill the process. It is invoked in a directory without a valid main.go
// module, where the underlying `go build` fails.
func TestCmdBuildEErrorsInsteadOfExit(t *testing.T) {
	dir := t.TempDir()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(wd)

	// No go.mod / main.go here → the go build step must fail with an error,
	// and cmdBuildE must surface it instead of calling os.Exit(1).
	if err := cmdBuildE(nil); err == nil {
		t.Fatal("expected cmdBuildE to return an error in a dir without a valid module")
	}
}
