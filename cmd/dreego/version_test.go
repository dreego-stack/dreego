package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestReadVersionFile verifies the VERSION-file parsing helper: empty, "(devel)"
// and "dev" contents are treated as absent (ok=true, value=""), while real
// versions are returned trimmed. A missing file reports ok=false.
func TestReadVersionFile(t *testing.T) {
	dir := t.TempDir()

	cases := []struct {
		name    string
		content string
		wantVal string
		wantOK  bool
	}{
		{"empty", "", "", true},
		{"whitespace", "  \n", "", true},
		{"devel", "(devel)", "", true},
		{"dev", "dev", "", true},
		{"version", "v0.0.25", "v0.0.25", true},
		{"version-whitespace", "  v1.2.3\n", "v1.2.3", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := filepath.Join(dir, tc.name+".txt")
			if err := os.WriteFile(p, []byte(tc.content), 0644); err != nil {
				t.Fatal(err)
			}
			got, ok := readVersionFile(p)
			if ok != tc.wantOK || got != tc.wantVal {
				t.Errorf("readVersionFile(%q) = (%q,%v), want (%q,%v)", tc.name, got, ok, tc.wantVal, tc.wantOK)
			}
		})
	}

	// Missing file: ok=false.
	if _, ok := readVersionFile(filepath.Join(dir, "nope.txt")); ok {
		t.Error("expected ok=false for missing file")
	}
}

// TestVersionFromWalkUp verifies the VERSION lookup that climbs from the current
// working directory to the filesystem root. A VERSION file in the cwd is found;
// an empty directory yields "".
func TestVersionFromWalkUp(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("found-in-cwd", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.WriteFile(filepath.Join(dir, "VERSION"), []byte("v9.9.9"), 0644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(wd)

		if got := versionFromWalkUp("VERSION"); got != "v9.9.9" {
			t.Errorf("expected v9.9.9, got %q", got)
		}
	})

	t.Run("not-found", func(t *testing.T) {
		dir := t.TempDir()
		if err := os.Chdir(dir); err != nil {
			t.Fatal(err)
		}
		defer os.Chdir(wd)

		if got := versionFromWalkUp("VERSION"); got != "" {
			t.Errorf("expected empty, got %q", got)
		}
	})
}

// TestVersionFromSourceRoot verifies the VERSION lookup that climbs from this
// source file's directory up toward the repo root. In a repo-local build the
// repo-root VERSION file must be resolved.
func TestVersionFromSourceRoot(t *testing.T) {
	got := versionFromSourceRoot("VERSION")
	if got == "" {
		t.Fatal("expected a non-empty version resolved from the repo source root")
	}
	if got == "(devel)" || got == "dev" {
		t.Errorf("expected a real version, got %q", got)
	}
}

// TestDreegoVersionFallbackDev verifies the top-level version resolution. The
// dev fallback ("dev") is reached only when no injected version, no build-info
// version and no VERSION file is found; in a repo-local test the VERSION file is
// resolved, so the result must be a real non-devel version rather than "(devel)".
func TestDreegoVersionFallbackDev(t *testing.T) {
	version = ""
	got := dreegoVersion()
	if got == "" {
		t.Fatal("expected a non-empty version")
	}
	if got == "(devel)" {
		t.Errorf("expected resolved version, got build-info placeholder %q", got)
	}
}
