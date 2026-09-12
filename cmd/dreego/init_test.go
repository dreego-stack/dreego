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
	if !strings.Contains(content, `ssr "github.com/dreego-stack/dreego/adapter/ssr"`) {
		t.Errorf("default main.go.tmpl must import \"github.com/dreego-stack/dreego/adapter/ssr\", got:\n%s", content)
	}
	if !strings.Contains(content, "ssr.Listen(app, addr)") {
		t.Errorf("default main.go.tmpl must call ssr.Listen(app, addr), got:\n%s", content)
	}
	if !strings.Contains(content, "ssr.DefaultAddr()") {
		t.Errorf("default main.go.tmpl must use ssr.DefaultAddr() for the default addr, got:\n%s", content)
	}
}

func TestDefaultBlueprintIncludesTaskfile(t *testing.T) {
	data, err := blueprintsSrc.ReadFile("blueprints/default/Taskfile.yml")
	if err != nil {
		t.Fatalf("read default Taskfile.yml: %v", err)
	}
	if !strings.Contains(string(data), "  generate:") {
		t.Fatalf("default Taskfile.yml must define generate task, got:\n%s", data)
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

func TestLandingBlueprintIncludesTaskfile(t *testing.T) {
	data, err := blueprintsSrc.ReadFile("blueprints/landing/Taskfile.yml")
	if err != nil {
		t.Fatalf("read landing Taskfile.yml: %v", err)
	}
	if !strings.Contains(string(data), "  build:") {
		t.Fatalf("landing Taskfile.yml must define build task, got:\n%s", data)
	}
}
