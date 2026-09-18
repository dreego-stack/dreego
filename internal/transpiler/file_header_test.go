package transpiler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseRouteFileExposesDreefileHeader(t *testing.T) {
	src := `DREEFILE layout

LAYOUT "www/layouts/admin.dreego"

GOIMPORT { sync, encoding/json }

COMPONENT "www/components" IMPORT { Card, Button, Card as ProductCard }

<body><@ProductCard/></body>`

	file, _, err := parseRouteFile(NewGenerator(), "www/routes/index.dreego", []byte(src))
	if err != nil {
		t.Fatalf("parseRouteFile: %v", err)
	}
	if file.Kind != FileKindLayout {
		t.Errorf("expected layout kind, got %v", file.Kind)
	}
	if file.Layout != "www/layouts/admin.dreego" {
		t.Errorf("expected layout path, got %q", file.Layout)
	}
	if len(file.GoImports) != 2 || file.GoImports[0] != "sync" {
		t.Errorf("expected go imports, got %+v", file.GoImports)
	}
	if len(file.Imports) != 1 || file.Imports[0].Aliases["ProductCard"] != "Card" {
		t.Errorf("expected component import with alias, got %+v", file.Imports)
	}
}

func TestParseRouteFileDefaultsToPageKind(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/index.dreego", []byte("<body>hi</body>"))
	if err != nil {
		t.Fatalf("parseRouteFile: %v", err)
	}
	if file.Kind != FileKindPage {
		t.Errorf("expected page kind, got %v", file.Kind)
	}
}

func TestLoadComponentDerivesNameFromFilename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Card.dreego")
	if err := os.WriteFile(path, []byte("DREEFILE component (title string)\n\n<body><h1>{{ title }}</h1></body>\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	component, err := loadComponent(path)
	if err != nil {
		t.Fatalf("loadComponent: %v", err)
	}
	if component.def == nil {
		t.Fatal("expected a component def for DREEFILE component")
	}
	if component.def.Name != "Card" {
		t.Errorf("expected name Card from filename, got %q", component.def.Name)
	}
	if len(component.def.Props) != 1 || component.def.Props[0].Name != "title" {
		t.Errorf("expected title prop, got %+v", component.def.Props)
	}
}

func TestLoadComponentRejectsInvalidFilename(t *testing.T) {
	for _, name := range []string{"product-card.dreego", "2Card.dreego", "Product.Card.dreego", "card.dreego"} {
		dir := t.TempDir()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("DREEFILE component ()\n\n<body><p>x</p></body>\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		_, err := loadComponent(path)
		if err == nil {
			t.Fatalf("expected error for filename %q", name)
		}
		if !strings.Contains(err.Error(), name) {
			t.Fatalf("error must name the file %q, got: %v", name, err)
		}
		if !strings.Contains(err.Error(), "Card.dreego") {
			t.Fatalf("error must explain a valid component filename, got: %v", err)
		}
	}
}

func TestLoadComponentKeepsValidPascalCaseFilename(t *testing.T) {
	for _, name := range []string{"Card.dreego", "ProductCard.dreego", "Hero2.dreego"} {
		dir := t.TempDir()
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte("DREEFILE component ()\n\n<body><p>x</p></body>\n"), 0o644); err != nil {
			t.Fatal(err)
		}
		component, err := loadComponent(path)
		if err != nil {
			t.Fatalf("loadComponent(%q): %v", name, err)
		}
		want := strings.TrimSuffix(name, ".dreego")
		if component.def == nil || component.def.Name != want {
			t.Fatalf("expected name %q, got %+v", want, component.def)
		}
	}
}

func TestLoadComponentIgnoresNonComponentKind(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Card.dreego")
	if err := os.WriteFile(path, []byte("DREEFILE layout\n\n<body>{#slot}</body>\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	component, err := loadComponent(path)
	if err != nil {
		t.Fatalf("loadComponent: %v", err)
	}
	if component.def != nil {
		t.Errorf("a layout file must not be loaded as a component, got %+v", component.def)
	}
}

func TestLoadComponentKeepsLegacyComponentLine(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "Legacy.dreego")
	if err := os.WriteFile(path, []byte("Component Legacy (title string)\n\n<body>{{ title }}</body>\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	component, err := loadComponent(path)
	if err != nil {
		t.Fatalf("loadComponent: %v", err)
	}
	if component.def == nil || component.def.Name != "Legacy" {
		t.Fatalf("legacy Component line must still work, got %+v", component.def)
	}
}
