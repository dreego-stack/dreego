package main

import (
	"strings"
	"testing"
)

func TestDefaultBlueprintGenImport(t *testing.T) {
	data, err := blueprintsSrc.ReadFile("blueprints/default/main.go.tmpl")
	if err != nil {
		t.Fatalf("read default main.go.tmpl: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `"§$name$§/www"`) {
		t.Errorf("default main.go.tmpl must import \"§$name$§/www\" (module-qualified, placeholder), got:\n%s", content)
	}
	if !strings.Contains(content, "www.Register(app)") {
		t.Errorf("default main.go.tmpl must call www.Register(app), got:\n%s", content)
	}
	if !strings.Contains(content, `ssr "github.com/dreego-stack/dreego/target/ssr"`) {
		t.Errorf("default main.go.tmpl must import \"github.com/dreego-stack/dreego/target/ssr\", got:\n%s", content)
	}
	if !strings.Contains(content, "ssr.Listen(app, addr)") {
		t.Errorf("default main.go.tmpl must call ssr.Listen(app, addr), got:\n%s", content)
	}
	if !strings.Contains(content, "ssr.DefaultAddr()") {
		t.Errorf("default main.go.tmpl must use ssr.DefaultAddr() for the default addr, got:\n%s", content)
	}
}

func TestLandingBlueprintGenImport(t *testing.T) {
	data, err := blueprintsSrc.ReadFile("blueprints/landing/main.go.tmpl")
	if err != nil {
		t.Fatalf("read landing main.go.tmpl: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `"§$name$§/www"`) {
		t.Errorf("landing main.go.tmpl must import \"§$name$§/www\", got:\n%s", content)
	}
	if !strings.Contains(content, "www.Register(app)") {
		t.Errorf("landing main.go.tmpl must call www.Register(app), got:\n%s", content)
	}
}
