package main

import (
	"strings"
	"testing"
)

// TestDefaultBlueprintGenImport verifies the default blueprint's main.go.tmpl
// imports the generated package via the §$name$§ placeholder (module-qualified
// path, replaced with the project name at scaffold time), exactly like the
// landing blueprint. A bare `_ "gen"` import would break `dreego build` with
// "package gen is not in std", because generate emits the package into
// dreego/gen/.
func TestDefaultBlueprintGenImport(t *testing.T) {
	data, err := blueprintsSrc.ReadFile("blueprints/default/main.go.tmpl")
	if err != nil {
		t.Fatalf("read default main.go.tmpl: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `_ "§$name$§/dreego/gen"`) {
		t.Errorf("default main.go.tmpl must import _ \"§$name$§/dreego/gen\" (module-qualified, placeholder), got:\n%s", content)
	}
	if strings.Contains(content, `_ "gen"`) {
		t.Errorf("default main.go.tmpl must not contain the bare _ \"gen\" import:\n%s", content)
	}
}

// TestLandingBlueprintGenImport verifies the landing blueprint (the reference
// pattern) uses the same module-qualified placeholder import.
func TestLandingBlueprintGenImport(t *testing.T) {
	data, err := blueprintsSrc.ReadFile("blueprints/landing/main.go.tmpl")
	if err != nil {
		t.Fatalf("read landing main.go.tmpl: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `_ "§$name$§/dreego/gen"`) {
		t.Errorf("landing main.go.tmpl must import _ \"§$name$§/dreego/gen\", got:\n%s", content)
	}
}
