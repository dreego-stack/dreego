package lexer

import (
	"errors"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func TestParseFileHeaderComponentKindAndProps(t *testing.T) {
	header, body := ParseFileHeader(`DREEFILE component (title string, count int = 5)

<body><h1>{{ title }}</h1></body>`)
	if header.Kind != ir.FileKindComponent {
		t.Fatalf("expected component kind, got %v", header.Kind)
	}
	if len(header.Props) != 2 {
		t.Fatalf("expected 2 props, got %+v", header.Props)
	}
	if header.Props[0].Name != "title" || header.Props[0].Type != "string" {
		t.Errorf("unexpected prop0: %+v", header.Props[0])
	}
	if header.Props[1].Name != "count" || header.Props[1].Type != "int" || header.Props[1].Default != "5" {
		t.Errorf("unexpected prop1: %+v", header.Props[1])
	}
	if body == "" {
		t.Fatal("expected body after DREEFILE line")
	}
}

func TestParseFileHeaderLayoutKind(t *testing.T) {
	header, _ := ParseFileHeader("DREEFILE layout\n\n<body>{#slot}</body>")
	if header.Kind != ir.FileKindLayout {
		t.Fatalf("expected layout kind, got %v", header.Kind)
	}
}

func TestParseFileHeaderExplicitPageKind(t *testing.T) {
	header, body, err := ParseFileHeaderStrict("DREEFILE page\n\n<body>hi</body>")
	if err != nil {
		t.Fatalf("DREEFILE page must be accepted: %v", err)
	}
	if header.Kind != ir.FileKindPage {
		t.Fatalf("expected page kind, got %v", header.Kind)
	}
	if body != "<body>hi</body>" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestParseFileHeaderRejectsUnknownKind(t *testing.T) {
	for _, src := range []string{
		"DREEFILE Component\n\n<body>hi</body>",
		"DREEFILE foo\n\n<body>hi</body>",
		"DREEFILE\n\n<body>hi</body>",
	} {
		_, _, err := ParseFileHeaderStrict(src)
		if err == nil {
			t.Fatalf("expected error for %q", src)
		}
		if !strings.Contains(err.Error(), "page, component, layout") {
			t.Fatalf("error must name accepted values, got: %v", err)
		}
	}
}

func TestParseFileHeaderComponentImportPathContainsImport(t *testing.T) {
	header, body, err := ParseFileHeaderStrict(`COMPONENT "www/IMPORTant/components" IMPORT { Card }

<body><@Card/></body>`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(header.Imports) != 1 {
		t.Fatalf("expected 1 import, got %+v", header.Imports)
	}
	if header.Imports[0].Path != "www/IMPORTant/components" {
		t.Errorf("expected full path, got %q", header.Imports[0].Path)
	}
	if len(header.Imports[0].Names) != 1 || header.Imports[0].Names[0] != "Card" {
		t.Errorf("expected Card name, got %+v", header.Imports[0].Names)
	}
	if body == "" {
		t.Fatal("expected body after COMPONENT import")
	}
}

func TestParseFileHeaderDefaultsToPage(t *testing.T) {
	header, body := ParseFileHeader(`LAYOUT "www/layouts/admin.dreego"

<body>hello</body>`)
	if header.Kind != ir.FileKindPage {
		t.Fatalf("expected page kind without DREEFILE, got %v", header.Kind)
	}
	if header.Layout != "www/layouts/admin.dreego" {
		t.Fatalf("expected layout path, got %q", header.Layout)
	}
	if body != "<body>hello</body>" {
		t.Fatalf("unexpected body: %q", body)
	}
}

func TestParseFileHeaderComponentImportWithAlias(t *testing.T) {
	header, body := ParseFileHeader(`COMPONENT "www/components" IMPORT { Card, Button, Card as ProductCard }

<body><@ProductCard/></body>`)
	if len(header.Imports) != 1 {
		t.Fatalf("expected 1 import, got %+v", header.Imports)
	}
	imp := header.Imports[0]
	if imp.Path != "www/components" {
		t.Errorf("expected path www/components, got %q", imp.Path)
	}
	if len(imp.Names) != 3 {
		t.Fatalf("expected 3 names, got %+v", imp.Names)
	}
	if imp.Aliases["ProductCard"] != "Card" {
		t.Errorf("expected alias ProductCard->Card, got %+v", imp.Aliases)
	}
	if body == "" {
		t.Fatal("expected body after COMPONENT import")
	}
}

func TestParseFileHeaderComponentImportWithoutAlias(t *testing.T) {
	header, _ := ParseFileHeader(`COMPONENT "github.com/dreego-stack/dreego-ui/components/dreegoui" IMPORT { Navbar, PriceCard }`)
	if len(header.Imports) != 1 {
		t.Fatalf("expected 1 import, got %+v", header.Imports)
	}
	if header.Imports[0].Aliases != nil {
		t.Errorf("expected no aliases, got %+v", header.Imports[0].Aliases)
	}
	if len(header.Imports[0].Names) != 2 {
		t.Errorf("expected 2 names, got %+v", header.Imports[0].Names)
	}
}

func TestParseFileHeaderGoImports(t *testing.T) {
	header, _ := ParseFileHeader(`GOIMPORT { sync, encoding/json }

<body>ok</body>`)
	if len(header.GoImports) != 2 {
		t.Fatalf("expected 2 go imports, got %+v", header.GoImports)
	}
	if header.GoImports[0].Path != "sync" || header.GoImports[0].Alias != "" ||
		header.GoImports[1].Path != "encoding/json" || header.GoImports[1].Alias != "" {
		t.Errorf("unexpected go imports: %+v", header.GoImports)
	}
}

func TestParseFileHeaderGoImportAlias(t *testing.T) {
	header, _ := ParseFileHeader(`GOIMPORT { strings, myauth "statuna/auth" }

<body>ok</body>`)
	if len(header.GoImports) != 2 {
		t.Fatalf("expected 2 go imports, got %+v", header.GoImports)
	}
	if header.GoImports[0].Path != "strings" || header.GoImports[0].Alias != "" {
		t.Errorf("unexpected bare import: %+v", header.GoImports[0])
	}
	if header.GoImports[1].Alias != "myauth" || header.GoImports[1].Path != "statuna/auth" {
		t.Errorf("unexpected aliased import: %+v", header.GoImports[1])
	}
}

func TestParseFileHeaderRejectsLegacyForms(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "component line",
			src:  "Component Navbar (title string)\n\n<body>x</body>",
			want: "DREEFILE component",
		},
		{
			name: "bare import",
			src:  "import dreego github.com/dreego-stack/dreego\n\n<body>x</body>",
			want: "GOIMPORT",
		},
		{
			name: "from import",
			src:  "from \"www/components\" import {\n    Button,\n}\n\n<body>x</body>",
			want: "COMPONENT",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ParseFileHeaderStrict(tc.src)
			if err == nil {
				t.Fatalf("expected a hard error for legacy header %q", tc.src)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error must name the replacement %q, got: %v", tc.want, err)
			}
		})
	}
}

func TestParseFileHeaderLegacyErrorPositionAfterLeadIn(t *testing.T) {
	cases := []struct {
		name string
		src  string
		line int
		col  int
	}{
		{
			name: "component after dreefile line",
			src:  "DREEFILE page\nComponent Foo (title string)\n\n<body>x</body>",
			line: 2,
			col:  1,
		},
		{
			name: "from after blank lead-in",
			src:  "\nfrom \"www/components\" import {\n\tFoo,\n}\n\n<body>x</body>",
			line: 2,
			col:  1,
		},
		{
			name: "bare import after blank lead-in",
			src:  "\nimport \"sync\"\n\n<body>x</body>",
			line: 2,
			col:  1,
		},
		{
			name: "indented bare import after dreefile line",
			src:  "DREEFILE component ()\n   import \"sync\"\n\n<body>x</body>",
			line: 2,
			col:  4,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := ParseFileHeaderStrict(tc.src)
			if err == nil {
				t.Fatalf("expected a hard error for legacy header %q", tc.src)
			}
			var he *HeaderError
			if !errors.As(err, &he) {
				t.Fatalf("expected *HeaderError, got %T: %v", err, err)
			}
			if he.Line != tc.line || he.Col != tc.col {
				t.Fatalf("expected position %d:%d, got %d:%d (%v)", tc.line, tc.col, he.Line, he.Col, err)
			}
		})
	}
}
