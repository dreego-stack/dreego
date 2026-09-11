package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunBuildHooksRunsPreBuildStep(t *testing.T) {
	projDir := t.TempDir()
	cmd := "echo hook-ran > " + filepath.Join(projDir, "hook-ran.txt")
	setupPlugin(t, projDir, "plugin-tailwind", manifestFor(cmd, "pre-build"))
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/dreego":          "v0.0.69",
		"github.com/dreego-stack/plugin-tailwind": "v0.0.1",
	})
	writeApprovals(t, projDir, map[string]bool{
		"github.com/dreego-stack/plugin-tailwind:" + cmd: true,
	})

	if err := runBuildHooks(projDir, false, strings.NewReader("")); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, "hook-ran.txt")); err != nil {
		t.Fatalf("pre-build step did not run: %v", err)
	}
}

func TestRunBuildHooksAutoApproveYes(t *testing.T) {
	projDir := t.TempDir()
	cmd := "echo hook-ran > " + filepath.Join(projDir, "hook-ran.txt")
	setupPlugin(t, projDir, "plugin-tailwind", manifestFor(cmd, "pre-build"))
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-tailwind": "v0.0.1",
	})

	if err := runBuildHooks(projDir, true, strings.NewReader("")); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, "hook-ran.txt")); err != nil {
		t.Fatalf("pre-build step did not run: %v", err)
	}
	af := readApprovalsFile(t, projDir)
	key := "github.com/dreego-stack/plugin-tailwind:" + cmd
	if !af.ApprovedHooks[key] {
		t.Fatalf("approval not saved for key %q, got: %v", key, af.ApprovedHooks)
	}
}

func TestRunBuildHooksPromptsAndApproves(t *testing.T) {
	withTerminal(t, true)
	projDir := t.TempDir()
	cmd := "echo hook-ran > " + filepath.Join(projDir, "hook-ran.txt")
	setupPlugin(t, projDir, "plugin-tailwind", manifestFor(cmd, "pre-build"))
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-tailwind": "v0.0.1",
	})

	if err := runBuildHooks(projDir, false, strings.NewReader("y\n")); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, "hook-ran.txt")); err != nil {
		t.Fatalf("pre-build step did not run: %v", err)
	}
	af := readApprovalsFile(t, projDir)
	key := "github.com/dreego-stack/plugin-tailwind:" + cmd
	if !af.ApprovedHooks[key] {
		t.Fatalf("approval not saved for key %q", key)
	}
}

func TestRunBuildHooksPromptsAndRejects(t *testing.T) {
	withTerminal(t, true)
	projDir := t.TempDir()
	cmd := "echo hook-ran > " + filepath.Join(projDir, "hook-ran.txt")
	setupPlugin(t, projDir, "plugin-tailwind", manifestFor(cmd, "pre-build"))
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-tailwind": "v0.0.1",
	})

	err := runBuildHooks(projDir, false, strings.NewReader("N\n"))
	if err == nil {
		t.Fatal("expected error when user rejects approval")
	}
	if !strings.Contains(err.Error(), "not approved by user") {
		t.Fatalf("expected 'not approved by user' in error, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, "hook-ran.txt")); err == nil {
		t.Fatal("command should not have run after rejection")
	}
}

func TestRunBuildHooksFailsInNonInteractiveWithoutApproval(t *testing.T) {
	projDir := t.TempDir()
	cmd := "echo hook-ran > " + filepath.Join(projDir, "hook-ran.txt")
	setupPlugin(t, projDir, "plugin-tailwind", manifestFor(cmd, "pre-build"))
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-tailwind": "v0.0.1",
	})

	err := runBuildHooks(projDir, false, strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error in non-interactive mode without approval")
	}
	if !strings.Contains(err.Error(), "dreego build --yes") {
		t.Fatalf("expected help message mentioning --yes, got: %v", err)
	}
	if _, err := os.Stat(filepath.Join(projDir, "hook-ran.txt")); err == nil {
		t.Fatal("command should not have run without approval")
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

	if err := runBuildHooks(projDir, false, strings.NewReader("")); err != nil {
		t.Fatalf("runBuildHooks should skip silently: %v", err)
	}
}

func TestRunBuildHooksFailsOnStepError(t *testing.T) {
	projDir := t.TempDir()
	setupPlugin(t, projDir, "plugin-broken", manifestFor("exit 1", "pre-build"))
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-broken": "v0.0.1",
	})
	writeApprovals(t, projDir, map[string]bool{
		"github.com/dreego-stack/plugin-broken:exit 1": true,
	})

	err := runBuildHooks(projDir, false, strings.NewReader(""))
	if err == nil {
		t.Fatal("expected error from failing build step")
	}
	if !strings.Contains(err.Error(), "build step failed") {
		t.Fatalf("expected 'build step failed' in error, got: %v", err)
	}
}

func TestRunBuildHooksSkipsNonPreBuildSteps(t *testing.T) {
	projDir := t.TempDir()
	marker := filepath.Join(projDir, "post-build-ran.txt")
	setupPlugin(t, projDir, "plugin-future", `{"build":{"steps":[{"cmd":"echo ran > `+marker+`","when":"post-build"}]}}`)
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-future": "v0.0.1",
	})

	if err := runBuildHooks(projDir, false, strings.NewReader("")); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("post-build step should not have run")
	}
}

func TestRunBuildHooksIgnoresCoreModule(t *testing.T) {
	projDir := t.TempDir()
	marker := filepath.Join(projDir, "core-ran.txt")
	coreDir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", "dreego")
	if err := os.MkdirAll(coreDir, 0755); err != nil {
		t.Fatal(err)
	}
	manifest := `{"build":{"steps":[{"cmd":"echo ran > ` + marker + `","when":"pre-build"}]}}`
	if err := os.WriteFile(filepath.Join(coreDir, "dreego-plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/dreego": "v0.0.69",
	})

	if err := runBuildHooks(projDir, false, strings.NewReader("")); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	if _, err := os.Stat(marker); err == nil {
		t.Fatal("core module build hook should not run")
	}
}

func TestRunBuildHooksRunsPluginsAlphabetically(t *testing.T) {
	projDir := t.TempDir()
	plugins := []string{"plugin-zeta", "plugin-alpha"}
	log := filepath.Join(projDir, "order.log")
	approvals := map[string]bool{}
	for _, p := range plugins {
		cmd := "echo " + p + " >> " + log
		setupPlugin(t, projDir, p, manifestFor(cmd, "pre-build"))
		approvals["github.com/dreego-stack/"+p+":"+cmd] = true
	}
	writeApprovals(t, projDir, approvals)
	writeGoMod(t, projDir, "myapp", map[string]string{
		"github.com/dreego-stack/plugin-zeta":  "v0.0.1",
		"github.com/dreego-stack/plugin-alpha": "v0.0.1",
	})

	if err := runBuildHooks(projDir, false, strings.NewReader("")); err != nil {
		t.Fatalf("runBuildHooks: %v", err)
	}
	body, err := os.ReadFile(log)
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
