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

func TestLoadConfigI18n(t *testing.T) {
	path := writeTempConfig(t, `{
		"i18n": {
			"enabled": true,
			"defaultLocale": "de",
			"locales": ["de", "en"],
			"urlStrategy": "none",
			"detection": ["cookie", "browser", "custom", "default"]
		}
	}`)
	s, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if !s.I18n.Enabled || s.I18n.DefaultLocale != "de" {
		t.Fatalf("unexpected i18n settings: %+v", s.I18n)
	}
	if strings.Join(s.I18n.Locales, ",") != "de,en" {
		t.Errorf("locales = %v", s.I18n.Locales)
	}
}

func TestLoadConfigRejectsInvalidI18n(t *testing.T) {
	cases := []struct {
		name string
		i18n string
		want string
	}{
		{"missing default", `{"enabled":true,"locales":["de"]}`, "defaultLocale is required"},
		{"default unsupported", `{"enabled":true,"defaultLocale":"fr","locales":["de","en"]}`, `defaultLocale "fr" is not listed`},
		{"duplicate locale", `{"enabled":true,"defaultLocale":"de","locales":["de","de"]}`, `duplicate locale "de"`},
		{"canonical duplicate locale", `{"enabled":true,"defaultLocale":"de-DE","locales":["de-DE","de-de"]}`, `duplicate locale "de-DE"`},
		{"invalid locale", `{"enabled":true,"defaultLocale":"de--DE","locales":["de--DE"]}`, `invalid locale "de--DE"`},
		{"invalid strategy", `{"enabled":true,"defaultLocale":"de","locales":["de"],"urlStrategy":"path"}`, `unsupported urlStrategy "path"`},
		{"duplicate detector", `{"enabled":true,"defaultLocale":"de","locales":["de"],"detection":["browser","browser"]}`, `duplicate locale detector "browser"`},
		{"default not last", `{"enabled":true,"defaultLocale":"de","locales":["de"],"detection":["default","browser"]}`, `locale detector "default" must be last`},
		{"unknown detector", `{"enabled":true,"defaultLocale":"de","locales":["de"],"detection":["ip","default"]}`, `unsupported locale detector "ip"`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTempConfig(t, `{"i18n":`+tc.i18n+`}`)
			_, err := LoadConfig(path)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}

func TestLoadSettingsWarnsOnInvalidConfig(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, configFileName), []byte(`{not valid`), 0o600); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	old := slog.Default()
	slog.SetDefault(slog.New(slog.NewTextHandler(&buf, nil)))
	defer slog.SetDefault(old)

	s := loadSettings(root)
	if s != nil {
		t.Fatalf("expected nil Settings for invalid config, got %+v", s)
	}
	if !strings.Contains(buf.String(), configFileName+" is invalid") {
		t.Errorf("expected warning about invalid config, got %q", buf.String())
	}
}
