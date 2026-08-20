package transpiler

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	return path
}

func TestLoadConfig(t *testing.T) {
	path := writeTempConfig(t, `{
		"logging": {"enabled": true},
		"redirects": [{"from": "/old", "to": "/new", "status": 301}],
		"rewrites": [{"from": "/a/*", "to": "/b/*"}]
	}`)

	s, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Settings")
	}
	if !s.Logging.Enabled {
		t.Error("expected Logging.Enabled true")
	}
	if len(s.Redirects) != 1 {
		t.Fatalf("expected 1 redirect, got %d", len(s.Redirects))
	}
	if s.Redirects[0].From != "/old" || s.Redirects[0].To != "/new" || s.Redirects[0].Status != 301 {
		t.Errorf("unexpected redirect: %+v", s.Redirects[0])
	}
	if len(s.Rewrites) != 1 {
		t.Fatalf("expected 1 rewrite, got %d", len(s.Rewrites))
	}
	if s.Rewrites[0].From != "/a/*" || s.Rewrites[0].To != "/b/*" {
		t.Errorf("unexpected rewrite: %+v", s.Rewrites[0])
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")
	s, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if s != nil {
		t.Errorf("expected nil Settings on error, got %+v", s)
	}
}

func TestLoadConfigInvalidJSON(t *testing.T) {
	path := writeTempConfig(t, `{not valid json`)
	s, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if s != nil {
		t.Errorf("expected nil Settings on error, got %+v", s)
	}
}

func TestLoadSettingsWarnsOnInvalidConfig(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "config.json"), []byte(`{not valid`), 0o600); err != nil {
		t.Fatal(err)
	}
	genDir := filepath.Join(root, "gen")
	os.MkdirAll(genDir, 0o755)

	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(old)

	s := loadSettings(genDir)
	if s != nil {
		t.Fatalf("expected nil Settings for invalid config, got %+v", s)
	}
	if !strings.Contains(buf.String(), "config.json is invalid") {
		t.Errorf("expected warning about invalid config, got %q", buf.String())
	}
}
