package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func withTerminal(t *testing.T, isTerm bool) {
	t.Helper()
	old := isTerminalFn
	isTerminalFn = func(io.Reader) bool { return isTerm }
	t.Cleanup(func() { isTerminalFn = old })
}

func writeApprovals(t *testing.T, projDir string, approvals map[string]bool) {
	t.Helper()
	af := approvalsFile{ApprovedHooks: approvals}
	data, err := json.MarshalIndent(af, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(filepath.Join(projDir, "dreego-build.json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func setupPlugin(t *testing.T, projDir, pluginName, manifest string) {
	t.Helper()
	pluginDir := filepath.Join(projDir, "vendor", "github.com", "dreego-stack", pluginName)
	if err := os.MkdirAll(pluginDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(pluginDir, "dreego-plugin.json"), []byte(manifest), 0644); err != nil {
		t.Fatal(err)
	}
}

func manifestFor(cmd, when string) string {
	return `{"build":{"steps":[{"cmd":"` + cmd + `","when":"` + when + `"}]}}`
}

func readApprovalsFile(t *testing.T, projDir string) approvalsFile {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(projDir, "dreego-build.json"))
	if err != nil {
		t.Fatalf("dreego-build.json not written: %v", err)
	}
	var af approvalsFile
	if err := json.Unmarshal(body, &af); err != nil {
		t.Fatal(err)
	}
	return af
}