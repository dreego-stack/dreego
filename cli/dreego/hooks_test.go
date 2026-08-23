package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBuildHooksRunsPreBuildStep(t *testing.T) {
	projDir := t.TempDir()
	pluginDir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", "plugin-tailwind")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}

	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/dreego":          "v0.0.69",
		"github.com/dreego-stack/plugin-tailwind": "v0.0.1",
	})

	marker := filepath.Join(projDir, "hook-ran.txt")
	manifest := `{"build":{"steps":[{"cmd":"echo hook-ran > ` + marker + `","when":"pre-build"}]}}`
	if err := os.WriteFile(filepath.Join(pluginDir, "dreego-plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	if err := runBuildHooks(projDir); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("pre-build step did not run: %v", err)
	}
}

func TestRunBuildHooksSkipsPluginWithoutManifest(t *testing.T) {
	projDir := t.TempDir()
	pluginDir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", "plugin-plain")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}

	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-plain": "v0.0.1",
	})

	if err := runBuildHooks(projDir); err != nil {
		t.Fatalf("runBuildHooks should skip silently: %v", err)
	}
}

func TestRunBuildHooksFailsOnStepError(t *testing.T) {
	projDir := t.TempDir()
	pluginDir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", "plugin-broken")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}

	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-broken": "v0.0.1",
	})

	manifest := `{"build":{"steps":[{"cmd":"exit 1","when":"pre-build"}]}}`
	if err := os.WriteFile(filepath.Join(pluginDir, "dreego-plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	err := runBuildHooks(projDir)
	if err == nil {
		t.Fatal("expected error from failing build step")
	}
	if !strings.Contains(err.Error(), "build step failed") {
		t.Fatalf("expected 'build step failed' in error, got: %v", err)
	}
}

func TestRunBuildHooksSkipsNonPreBuildSteps(t *testing.T) {
	projDir := t.TempDir()
	pluginDir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", "plugin-future")
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}

	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-future": "v0.0.1",
	})

	marker := filepath.Join(projDir, "post-build-ran.txt")
	manifest := `{"build":{"steps":[{"cmd":"echo ran > ` + marker + `","when":"post-build"}]}}`
	if err := os.WriteFile(filepath.Join(pluginDir, "dreego-plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	if err := runBuildHooks(projDir); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("post-build step should not have run")
	}
}

func TestRunBuildHooksIgnoresCoreModule(t *testing.T) {
	projDir := t.TempDir()
	coreDir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", "dreego")
	if err := os.MkdirAll(coreDir, 0755); err != nil {
		t.Fatal(err)
	}

	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/dreego": "v0.0.69",
	})

	marker := filepath.Join(projDir, "core-ran.txt")
	manifest := `{"build":{"steps":[{"cmd":"echo ran > ` + marker + `","when":"pre-build"}]}}`
	if err := os.WriteFile(filepath.Join(coreDir, "dreego-plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}

	if err := runBuildHooks(projDir); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("core module build hook should not run")
	}
}

func TestRunBuildHooksRunsPluginsAlphabetically(t *testing.T) {
	projDir := t.TempDir()

	var plugins = []string{"plugin-zeta", "plugin-alpha"}
	for _, p := range plugins {
		dir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", p)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		log := filepath.Join(projDir, "order.log")
		manifest := `{"build":{"steps":[{"cmd":"echo ` + p + ` >> ` + log + `","when":"pre-build"}]}}`
		if err := os.WriteFile(filepath.Join(dir, "dreego-plugin.json"), []byte(manifest), 0644); err != nil {
			t.Fatal(err)
		}
	}
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-zeta":  "v0.0.1",
		"github.com/dreego-stack/plugin-alpha": "v0.0.1",
	})

	if err := runBuildHooks(projDir); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	body, err := os.ReadFile(filepath.Join(projDir, "order.log"))
	if err != nil {
		t.Fatalf("order.log not written: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(string(body)), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d: %v", len(lines), lines)
	}
	if lines[0] != "plugin-alpha" || lines[1] != "plugin-zeta" {
		t.Fatalf("expected alpha then zeta, got: %v", lines)
	}
}