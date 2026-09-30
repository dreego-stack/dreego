package dreefile

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func parseSettings(t *testing.T, raw string) *Settings {
	t.Helper()
	var s Settings
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		t.Fatalf("unmarshal settings: %v", err)
	}
	return &s
}

func TestUnmarshalPresenceFlags(t *testing.T) {
	cases := []struct {
		name        string
		raw         string
		loggingPres bool
		i18nPres    bool
	}{
		{"empty", `{}`, false, false},
		{"logging only", `{"logging":{"enabled":false}}`, true, false},
		{"i18n only", `{"i18n":{"enabled":false}}`, false, true},
		{"both", `{"logging":{"enabled":true},"i18n":{"enabled":false}}`, true, true},
		{"redirects do not set presence", `{"redirects":[]}`, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := parseSettings(t, tc.raw)
			if s.Logging.Present != tc.loggingPres {
				t.Errorf("Logging.Present = %v, want %v", s.Logging.Present, tc.loggingPres)
			}
			if s.I18n.Present != tc.i18nPres {
				t.Errorf("I18n.Present = %v, want %v", s.I18n.Present, tc.i18nPres)
			}
		})
	}
}

func TestMergeSettingsNil(t *testing.T) {
	app := parseSettings(t, `{"logging":{"enabled":false}}`)
	if got := MergeSettings(nil, app); got != app {
		t.Error("MergeSettings(nil, app) must return app")
	}
	root := parseSettings(t, `{"logging":{"enabled":true}}`)
	if got := MergeSettings(root, nil); got != root {
		t.Error("MergeSettings(root, nil) must return root")
	}
	if got := MergeSettings(nil, nil); got != nil {
		t.Error("MergeSettings(nil, nil) must return nil")
	}
}

func TestMergeSettingsFieldOverride(t *testing.T) {
	cases := []struct {
		name       string
		root       string
		app        string
		wantLog    bool
		wantLogSet bool
		wantRedir  int
		wantRewr   int
	}{
		{
			name:      "app absent inherits root",
			root:      `{"logging":{"enabled":false},"redirects":[{"from":"/a","to":"/b","status":301}]}`,
			app:       `{}`,
			wantLog:   false,
			wantRedir: 1,
		},
		{
			name:      "app overrides logging only",
			root:      `{"logging":{"enabled":false},"redirects":[{"from":"/a","to":"/b","status":301}]}`,
			app:       `{"logging":{"enabled":true}}`,
			wantLog:   true,
			wantRedir: 1,
		},
		{
			name:      "app replaces redirects entirely",
			root:      `{"redirects":[{"from":"/a","to":"/b","status":301}]}`,
			app:       `{"redirects":[{"from":"/x","to":"/y","status":302}]}`,
			wantRedir: 1,
		},
		{
			name:      "app empty redirect list clears root",
			root:      `{"redirects":[{"from":"/a","to":"/b","status":301}]}`,
			app:       `{"redirects":[]}`,
			wantRedir: 0,
		},
		{
			name:      "app adds rewrites, root keeps redirects",
			root:      `{"redirects":[{"from":"/a","to":"/b","status":301}]}`,
			app:       `{"rewrites":[{"from":"/x/*","to":"/y/*"}]}`,
			wantRedir: 1,
			wantRewr:  1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			merged := MergeSettings(parseSettings(t, tc.root), parseSettings(t, tc.app))
			if merged.Logging.Enabled != tc.wantLog {
				t.Errorf("Logging.Enabled = %v, want %v", merged.Logging.Enabled, tc.wantLog)
			}
			if len(merged.Redirects) != tc.wantRedir {
				t.Errorf("Redirects = %d, want %d", len(merged.Redirects), tc.wantRedir)
			}
			if len(merged.Rewrites) != tc.wantRewr {
				t.Errorf("Rewrites = %d, want %d", len(merged.Rewrites), tc.wantRewr)
			}
		})
	}
}

func TestMergeSettingsOverridesRedirectContent(t *testing.T) {
	root := parseSettings(t, `{"redirects":[{"from":"/a","to":"/b","status":301}]}`)
	app := parseSettings(t, `{"redirects":[{"from":"/x","to":"/y","status":308}]}`)
	merged := MergeSettings(root, app)
	if merged.Redirects[0].From != "/x" || merged.Redirects[0].Status != 308 {
		t.Fatalf("app redirect should win entirely: %+v", merged.Redirects)
	}
}

func TestMergeSettingsI18nPresence(t *testing.T) {
	root := parseSettings(t, `{"i18n":{"enabled":true,"defaultLocale":"de","locales":["de"]}}`)
	app := parseSettings(t, `{"logging":{"enabled":true}}`)
	merged := MergeSettings(root, app)
	if !merged.I18n.Enabled || merged.I18n.DefaultLocale != "de" {
		t.Fatalf("app without i18n must inherit root i18n: %+v", merged.I18n)
	}
}

func TestMergeSettingsAppI18nWins(t *testing.T) {
	root := parseSettings(t, `{"i18n":{"enabled":true,"defaultLocale":"de","locales":["de"]}}`)
	app := parseSettings(t, `{"i18n":{"enabled":false}}`)
	merged := MergeSettings(root, app)
	if merged.I18n.Enabled {
		t.Fatal("app i18n block must override the root block")
	}
}

func TestMergeSettingsPluginsReplace(t *testing.T) {
	root := parseSettings(t, `{"plugins":{"github.com/dreego-stack/plugin-a":{"client":["x"]}}}`)
	app := parseSettings(t, `{"plugins":{"github.com/dreego-stack/plugin-b":{"client":["y"]}}}`)
	merged := MergeSettings(root, app)
	if _, ok := merged.Plugins["github.com/dreego-stack/plugin-b"]; !ok {
		t.Fatalf("app plugins must replace root plugins: %+v", merged.Plugins)
	}
}

func TestLoadAppSettings(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "www")
	if err := os.MkdirAll(app, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, configFileName), []byte(`{"logging":{"enabled":false}}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(app, configFileName), []byte(`{"logging":{"enabled":true}}`), 0644); err != nil {
		t.Fatal(err)
	}
	s, err := loadAppSettings(root, app)
	if err != nil {
		t.Fatalf("loadAppSettings: %v", err)
	}
	if !s.Logging.Enabled {
		t.Fatal("app config must override the root default")
	}
}

func TestLoadAppSettingsWithoutAppConfig(t *testing.T) {
	root := t.TempDir()
	app := filepath.Join(root, "www")
	if err := os.MkdirAll(app, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, configFileName), []byte(`{"logging":{"enabled":false}}`), 0644); err != nil {
		t.Fatal(err)
	}
	s, err := loadAppSettings(root, app)
	if err != nil {
		t.Fatalf("loadAppSettings: %v", err)
	}
	if s.Logging.Enabled {
		t.Fatal("app without config must inherit the root value")
	}
}

func TestValidateUrlRules(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{"nil settings", ``, false},
		{"empty", `{}`, false},
		{"valid redirect", `{"redirects":[{"from":"/a","to":"/b","status":301}]}`, false},
		{"valid rewrite", `{"rewrites":[{"from":"/a/*","to":"/b/*"}]}`, false},
		{"invalid redirect target", `{"redirects":[{"from":"/a","to":"b","status":301}]}`, true},
		{"invalid redirect status", `{"redirects":[{"from":"/a","to":"/b","status":200}]}`, true},
		{"invalid rewrite", `{"rewrites":[{"from":"/a","to":"b"}]}`, true},
		{"redirect cycle", `{"redirects":[{"from":"/a","to":"/b","status":301},{"from":"/b","to":"/a","status":301}]}`, true},
		{"rewrite cycle", `{"rewrites":[{"from":"/a","to":"/b"},{"from":"/b","to":"/a"}]}`, true},
		{"mixed cycle", `{"redirects":[{"from":"/b","to":"/a","status":301}],"rewrites":[{"from":"/a","to":"/b"}]}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var s *Settings
			if tc.raw != "" {
				s = parseSettings(t, tc.raw)
			}
			err := validateUrlRules(s)
			if tc.wantErr != (err != nil) {
				t.Fatalf("validateUrlRules(%s) = %v, wantErr=%v", tc.raw, err, tc.wantErr)
			}
		})
	}
}
