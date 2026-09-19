package lexer

import "testing"

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

func TestParseHeaderGroupedImports(t *testing.T) {
	_, imports, body := ParseHeader(`COMPONENT "www/components" IMPORT {
    Button,
    Card,
}

COMPONENT "github.com/dreego-stack/dreego-ui/components/dreegoui" IMPORT {
    Navbar,
    PriceCard,
}

<body><@Button/></body>`)

	if len(imports) != 2 {
		t.Fatalf("expected 2 grouped imports, got %d: %+v", len(imports), imports)
	}
	if imports[0].Path != "www/components" || len(imports[0].Names) != 2 {
		t.Fatalf("unexpected local import: %+v", imports[0])
	}
	if imports[1].Path != "github.com/dreego-stack/dreego-ui/components/dreegoui" || len(imports[1].Names) != 2 {
		t.Fatalf("unexpected module import: %+v", imports[1])
	}
	if body == "" {
		t.Fatal("expected body after grouped imports")
	}
}
