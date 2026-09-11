package transpiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPluginClientBundleSelectsDependenciesDeterministically(t *testing.T) {
	project := t.TempDir()
	plugin := filepath.Join(project, "vendor", "github.com", "dreego-stack", "plugin-auth")
	writePluginClientFixture(t, project, plugin, `{
		"client":{"format":"modules-v1","path":"/_dreego/plugin-auth.js","modules":[
			{"id":"passkeys","path":"client/passkeys.js","dependsOn":["core"]},
			{"id":"core","path":"client/core.js","required":true},
			{"id":"password","path":"client/password.js","dependsOn":["core"]}
		]}
	}`, map[string]string{
		"client/core.js":     "globalThis.core = true;",
		"client/passkeys.js": "globalThis.passkeys = true;",
		"client/password.js": "globalThis.password = true;",
	})

	settings := &Settings{Plugins: map[string]PluginSettings{
		"github.com/dreego-stack/plugin-auth": {Client: []string{"passkeys"}},
	}}
	src, count, err := generatePluginClientAssets(project, settings, map[string]bool{})
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}
	if strings.Contains(src, "password = true") {
		t.Fatal("bundle contains an unselected module")
	}
	if strings.Index(src, "core = true") > strings.Index(src, "passkeys = true") {
		t.Fatal("dependency must precede the selected module")
	}
	if !strings.Contains(src, `app.RegisterStatic("/_dreego/plugin-auth.js"`) {
		t.Fatalf("generated source does not register the manifest path: %s", src)
	}
}

func TestPluginClientBundleRejectsUnsafeAndInvalidManifests(t *testing.T) {
	tests := []struct {
		name     string
		manifest string
		files    map[string]string
		want     string
	}{
		{"unsupported format", `{"client":{"format":"future","path":"/x.js","modules":[]}}`, nil, "unsupported client format"},
		{"unknown module", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"core","path":"client/core.js"}]}}`, map[string]string{"client/core.js": "core();"}, "unknown client module"},
		{"dependency cycle", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"a","path":"a.js","dependsOn":["b"]},{"id":"b","path":"b.js","dependsOn":["a"]}]}}`, map[string]string{"a.js": "a();", "b.js": "b();"}, "dependency cycle"},
		{"traversal", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"bad","path":"../secret.js","required":true}]}}`, nil, "unsafe client module path"},
		{"backslash traversal", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"bad","path":"..\\secret.js","required":true}]}}`, nil, "unsafe client module path"},
		{"redundant separator", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"bad","path":"client//bad.js","required":true}]}}`, nil, "unsafe client module path"},
		{"dot segment", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"bad","path":"client/./bad.js","required":true}]}}`, nil, "unsafe client module path"},
		{"non javascript", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"bad","path":"client/bad.ts","required":true}]}}`, map[string]string{"client/bad.ts": "bad();"}, "must be a .js file"},
		{"unsafe id", `{"client":{"format":"modules-v1","path":"/x.js","modules":[{"id":"bad\\nid","path":"client/bad.js","required":true}]}}`, map[string]string{"client/bad.js": "bad();"}, "invalid client module id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			project := t.TempDir()
			plugin := filepath.Join(project, "vendor", "github.com", "dreego-stack", "plugin-auth")
			writePluginClientFixture(t, project, plugin, tt.manifest, tt.files)
			settings := &Settings{Plugins: map[string]PluginSettings{
				"github.com/dreego-stack/plugin-auth": {Client: []string{"missing"}},
			}}
			if tt.name != "unknown module" {
				settings.Plugins["github.com/dreego-stack/plugin-auth"] = PluginSettings{}
			}
			_, _, err := generatePluginClientAssets(project, settings, map[string]bool{})
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want substring %q", err, tt.want)
			}
		})
	}
}

func TestPluginClientBundleRejectsRouteConflict(t *testing.T) {
	project := t.TempDir()
	plugin := filepath.Join(project, "vendor", "github.com", "dreego-stack", "plugin-auth")
	writePluginClientFixture(t, project, plugin, `{"client":{"format":"modules-v1","path":"/auth.js","modules":[{"id":"core","path":"core.js","required":true}]}}`, map[string]string{"core.js": "core();"})
	settings := &Settings{Plugins: map[string]PluginSettings{"github.com/dreego-stack/plugin-auth": {}}}
	_, _, err := generatePluginClientAssets(project, settings, map[string]bool{"GET /auth.js": true})
	if err == nil || !strings.Contains(err.Error(), "conflicts with existing route") {
		t.Fatalf("error = %v", err)
	}
}

func TestPluginClientBundleRejectsSymlinkEscape(t *testing.T) {
	project := t.TempDir()
	plugin := filepath.Join(project, "vendor", "github.com", "dreego-stack", "plugin-auth")
	writePluginClientFixture(t, project, plugin, `{"client":{"format":"modules-v1","path":"/auth.js","modules":[{"id":"core","path":"client/core.js","required":true}]}}`, nil)
	outside := filepath.Join(project, "outside.js")
	if err := os.WriteFile(outside, []byte("outside();"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(plugin, "client"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(plugin, "client", "core.js")); err != nil {
		t.Fatal(err)
	}
	settings := &Settings{Plugins: map[string]PluginSettings{"github.com/dreego-stack/plugin-auth": {}}}
	if _, _, err := generatePluginClientAssets(project, settings, map[string]bool{}); err == nil {
		t.Fatal("expected symlink escape to fail")
	}
}

func writePluginClientFixture(t *testing.T, project, plugin, manifest string, files map[string]string) {
	t.Helper()
	if err := os.MkdirAll(plugin, 0o755); err != nil {
		t.Fatal(err)
	}
	goMod := "module example.com/app\n\ngo 1.24\n\nrequire github.com/dreego-stack/plugin-auth v0.0.1\n"
	if err := os.WriteFile(filepath.Join(project, "go.mod"), []byte(goMod), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plugin, "dreego-plugin.json"), []byte(manifest), 0o600); err != nil {
		t.Fatal(err)
	}
	for name, content := range files {
		path := filepath.Join(plugin, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}
}
