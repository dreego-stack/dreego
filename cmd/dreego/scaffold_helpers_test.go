package main

import (
	"testing"

	"github.com/dreego-stack/dreego/cmd/dreego/internal/templates"
)

func TestValidProjectName(t *testing.T) {
	valid := []string{"myapp", "my-app", "my_app", "github.com/me/myapp", "a.b.c", "app1", "A1"}
	for _, name := range valid {
		if !validProjectName(name) {
			t.Errorf("validProjectName(%q) = false, want true", name)
		}
	}
	invalid := []string{
		"", ".", "..", "1app", "my app", "my$app", "my;app",
		"a//b", "/a", "a/", "a/../b", "a\\b", "a'b", "a\"b",
		"a*b", "a?b", "a|b", "a(b)", "a[b]", "a{b}", "a!b",
	}
	for _, name := range invalid {
		if validProjectName(name) {
			t.Errorf("validProjectName(%q) = true, want false", name)
		}
	}
}

func TestRequiredModules(t *testing.T) {
	web := templates.Meta{Adapter: "ssr"}
	mods := requiredModules(web)
	if len(mods) != 2 || mods[0] != "github.com/dreego-stack/dreego/core" {
		t.Fatalf("requiredModules(ssr) = %v", mods)
	}
	if mods[1] != "github.com/dreego-stack/dreego/adapter/ssr" {
		t.Fatalf("ssr adapter module missing: %v", mods)
	}

	none := templates.Meta{}
	if mods := requiredModules(none); len(mods) != 1 {
		t.Fatalf("requiredModules(no adapter) = %v, want core only", mods)
	}

	extra := templates.Meta{Adapter: "ssr", ExtraRequires: []string{"github.com/dreego-stack/plugin-auth"}}
	if mods := requiredModules(extra); len(mods) != 3 {
		t.Fatalf("requiredModules(extra) = %v, want core+adapter+extra", mods)
	}
}

func TestAdapterModule(t *testing.T) {
	module, dir := adapterModule("ssr")
	if module != "github.com/dreego-stack/dreego/adapter/ssr" || dir != "adapter/ssr" {
		t.Fatalf("adapterModule(ssr) = (%q, %q)", module, dir)
	}
	module, dir = adapterModule("wails")
	if module != "github.com/dreego-stack/dreego/adapter/wails" || dir != "adapter/wails" {
		t.Fatalf("adapterModule(wails) = (%q, %q)", module, dir)
	}
	module, dir = adapterModule("")
	if module != "" || dir != "" {
		t.Fatalf("adapterModule(\"\") = (%q, %q), want empty", module, dir)
	}
}

func TestLocalReplacements(t *testing.T) {
	meta := templates.Meta{Adapter: "ssr"}
	reps := localReplacements("/repo", meta)
	if len(reps) != 3 {
		t.Fatalf("localReplacements = %d entries, want 3", len(reps))
	}
	byModule := map[string]string{}
	for _, r := range reps {
		byModule[r.module] = r.dir
	}
	if byModule["github.com/dreego-stack/dreego"] != "/repo" {
		t.Errorf("root module replace = %q", byModule["github.com/dreego-stack/dreego"])
	}
	if byModule["github.com/dreego-stack/dreego/core"] != "/repo/core" {
		t.Errorf("core replace = %q", byModule["github.com/dreego-stack/dreego/core"])
	}
	if byModule["github.com/dreego-stack/dreego/adapter/ssr"] != "/repo/adapter/ssr" {
		t.Errorf("ssr replace = %q", byModule["github.com/dreego-stack/dreego/adapter/ssr"])
	}
}

func TestLocalReplacementsWithoutAdapter(t *testing.T) {
	reps := localReplacements("/repo", templates.Meta{})
	if len(reps) != 2 {
		t.Fatalf("localReplacements(no adapter) = %d, want 2", len(reps))
	}
}
