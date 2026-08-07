package core

import "testing"

func TestParseComponentHeaderNoParens(t *testing.T) {
	comp := parseComponentHeader("Component Navbar")
	if comp == nil {
		t.Fatal("expected non-nil ComponentDef")
	}
	if comp.Name != "Navbar" {
		t.Errorf("expected Name Navbar, got %q", comp.Name)
	}
	if len(comp.Props) != 0 {
		t.Errorf("expected no props, got %d", len(comp.Props))
	}
	if len(comp.Slots) != 0 {
		t.Errorf("expected no slots, got %d", len(comp.Slots))
	}
}

func TestParseComponentHeaderWithSlots(t *testing.T) {
	// Slots are declared as a second parenthesized group after the props.
	comp := parseComponentHeader("Component Card (title string) (header, footer)")
	if comp == nil {
		t.Fatal("expected non-nil ComponentDef")
	}
	if comp.Name != "Card" {
		t.Errorf("expected Name Card, got %q", comp.Name)
	}
	if len(comp.Props) != 1 {
		t.Fatalf("expected 1 prop, got %d", len(comp.Props))
	}
	if comp.Props[0].Name != "title" || comp.Props[0].Type != "string" {
		t.Errorf("unexpected prop: %+v", comp.Props[0])
	}
	if len(comp.Slots) != 2 {
		t.Fatalf("expected 2 slots, got %d", len(comp.Slots))
	}
	if comp.Slots[0] != "header" || comp.Slots[1] != "footer" {
		t.Errorf("unexpected slots: %v", comp.Slots)
	}
}

func TestParsePropsWithDefault(t *testing.T) {
	props := parseProps("title string, count int = 5, flag bool")
	if len(props) != 3 {
		t.Fatalf("expected 3 props, got %d", len(props))
	}
	if props[0].Name != "title" || props[0].Type != "string" || props[0].Default != "" {
		t.Errorf("unexpected prop0: %+v", props[0])
	}
	if props[1].Name != "count" || props[1].Type != "int" || props[1].Default != "5" {
		t.Errorf("unexpected prop1: %+v", props[1])
	}
	if props[2].Name != "flag" || props[2].Type != "bool" || props[2].Default != "" {
		t.Errorf("unexpected prop2: %+v", props[2])
	}
}

func TestParsePropsDefaultType(t *testing.T) {
	props := parseProps("name")
	if len(props) != 1 {
		t.Fatalf("expected 1 prop, got %d", len(props))
	}
	if props[0].Name != "name" {
		t.Errorf("expected Name name, got %q", props[0].Name)
	}
	if props[0].Type != "string" {
		t.Errorf("expected default Type string, got %q", props[0].Type)
	}
}

func TestParseImportLineWithAlias(t *testing.T) {
	imp := parseImportLine(`import Navbar "components/navbar.dreego"`)
	if imp == nil {
		t.Fatal("expected non-nil Import")
	}
	if imp.Alias != "Navbar" {
		t.Errorf("expected Alias Navbar, got %q", imp.Alias)
	}
	if imp.Path != "components/navbar.dreego" {
		t.Errorf("expected Path components/navbar.dreego, got %q", imp.Path)
	}
}

func TestParseImportLineWithoutAlias(t *testing.T) {
	imp := parseImportLine(`import "components/navbar.dreego"`)
	if imp == nil {
		t.Fatal("expected non-nil Import")
	}
	if imp.Alias != "" {
		t.Errorf("expected empty Alias, got %q", imp.Alias)
	}
	if imp.Path != "components/navbar.dreego" {
		t.Errorf("expected Path components/navbar.dreego, got %q", imp.Path)
	}
}

func TestParseImportLineTooShort(t *testing.T) {
	if imp := parseImportLine("import"); imp != nil {
		t.Errorf("expected nil Import for too-short line, got %+v", imp)
	}
}
