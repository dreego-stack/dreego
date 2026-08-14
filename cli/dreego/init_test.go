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

	if !strings.Contains(content, `"§$name$§/dreego/gen"`) {
		t.Errorf("default main.go.tmpl must import \"§$name$§/dreego/gen\" (module-qualified, placeholder), got:\n%s", content)
	}
	if !strings.Contains(content, "gen.Register(app)") {
		t.Errorf("default main.go.tmpl must call gen.Register(app), got:\n%s", content)
	}
	if !strings.Contains(content, "app.Listen") {
		t.Errorf("default main.go.tmpl must call app.Listen, got:\n%s", content)
	}
}

func TestLandingBlueprintGenImport(t *testing.T) {
	data, err := blueprintsSrc.ReadFile("blueprints/landing/main.go.tmpl")
	if err != nil {
		t.Fatalf("read landing main.go.tmpl: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `"§$name$§/dreego/gen"`) {
		t.Errorf("landing main.go.tmpl must import \"§$name$§/dreego/gen\", got:\n%s", content)
	}
	if !strings.Contains(content, "gen.Register(app)") {
		t.Errorf("landing main.go.tmpl must call gen.Register(app), got:\n%s", content)
	}
}
