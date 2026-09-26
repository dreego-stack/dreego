package lexer

import (
	"errors"
	"strings"
	"testing"
)

func TestParseFileHeaderRejectsUnbalancedBraceList(t *testing.T) {
	for _, src := range []string{
		`GOIMPORT { sync`,
		"COMPONENT \"www/components\" IMPORT { Card\n\n<body>x</body>",
	} {
		_, _, err := ParseFileHeaderStrict(src)
		if err == nil {
			t.Fatalf("expected unbalanced brace list error for %q", src)
		}
		if !strings.Contains(err.Error(), "unbalanced") {
			t.Fatalf("error must mention unbalanced list, got: %v", err)
		}
		assertHeaderPosition(t, err)
	}
}

func TestParseFileHeaderRejectsNestedBraceList(t *testing.T) {
	_, _, err := ParseFileHeaderStrict(`GOIMPORT { a, { b } }
<body>x</body>`)
	if err == nil {
		t.Fatal("expected nested brace list error")
	}
	if !strings.Contains(err.Error(), "nested") {
		t.Fatalf("error must mention nested brace, got: %v", err)
	}
	assertHeaderPosition(t, err)
}

func TestParseFileHeaderRejectsTrailingContentAfterBraceList(t *testing.T) {
	_, _, err := ParseFileHeaderStrict(`COMPONENT "www/components" IMPORT { A, B } extra
<body>x</body>`)
	if err == nil {
		t.Fatal("expected trailing-content error")
	}
	if !strings.Contains(err.Error(), "after closing") {
		t.Fatalf("error must mention trailing content, got: %v", err)
	}
	assertHeaderPosition(t, err)
}

func TestParseFileHeaderGoImportPositionAfterLeadIn(t *testing.T) {
	_, _, err := ParseFileHeaderStrict("DREEFILE page\nGOIMPORT { sync")
	if err == nil {
		t.Fatal("expected unbalanced GOIMPORT error")
	}
	var he *HeaderError
	if !errors.As(err, &he) {
		t.Fatalf("expected *HeaderError, got %T: %v", err, err)
	}
	if he.Line != 2 || he.Col != 10 {
		t.Fatalf("expected 2:10, got %d:%d (%v)", he.Line, he.Col, err)
	}
}

func TestParseFileHeaderRejectsEmptyGoImportList(t *testing.T) {
	_, _, err := ParseFileHeaderStrict(`GOIMPORT {  }
<body>x</body>`)
	if err == nil {
		t.Fatal("expected empty GOIMPORT list error")
	}
	if !strings.Contains(err.Error(), "GOIMPORT") {
		t.Fatalf("error must name GOIMPORT, got: %v", err)
	}
	assertHeaderPosition(t, err)
}

func TestParseFileHeaderRejectsUnquotedGoImportPath(t *testing.T) {
	_, _, err := ParseFileHeaderStrict(`GOIMPORT { encoding json }
<body>x</body>`)
	if err == nil {
		t.Fatal("expected invalid GOIMPORT path error")
	}
	if !strings.Contains(err.Error(), "GOIMPORT") || !strings.Contains(err.Error(), "encoding json") {
		t.Fatalf("error must name GOIMPORT and the offending value, got: %v", err)
	}
	assertHeaderPosition(t, err)
}

func TestParseFileHeaderRejectsInvalidGoImportAlias(t *testing.T) {
	for _, src := range []string{
		`GOIMPORT { 1bad "statuna/auth" }`,
		`GOIMPORT { bad- "statuna/auth" }`,
		`GOIMPORT { alias "statuna/auth`,
		`GOIMPORT { "a b" }`,
	} {
		_, _, err := ParseFileHeaderStrict(src + "\n<body>x</body>")
		if err == nil {
			t.Fatalf("expected invalid GOIMPORT alias/path error for %q", src)
		}
		if !strings.Contains(err.Error(), "GOIMPORT") {
			t.Fatalf("error must name GOIMPORT, got: %v", err)
		}
		assertHeaderPosition(t, err)
	}
}

func TestParseFileHeaderRejectsUnquotedLayout(t *testing.T) {
	_, _, err := ParseFileHeaderStrict("LAYOUT www/layouts/x.dreego\n\n<body>x</body>")
	if err == nil {
		t.Fatal("expected unquoted LAYOUT error")
	}
	if !strings.Contains(err.Error(), "LAYOUT") {
		t.Fatalf("error must name LAYOUT, got: %v", err)
	}
	assertHeaderPosition(t, err)
}

func TestParseFileHeaderRejectsEmptyLayout(t *testing.T) {
	_, _, err := ParseFileHeaderStrict("LAYOUT \"\"\n\n<body>x</body>")
	if err == nil {
		t.Fatal("expected empty LAYOUT error")
	}
	if !strings.Contains(err.Error(), "LAYOUT") {
		t.Fatalf("error must name LAYOUT, got: %v", err)
	}
	assertHeaderPosition(t, err)
}

func TestParseFileHeaderRejectsDuplicateDreefile(t *testing.T) {
	_, _, err := ParseFileHeaderStrict("DREEFILE page\nDREEFILE layout\n\n<body>x</body>")
	if err == nil {
		t.Fatal("expected duplicate DREEFILE error")
	}
	if !strings.Contains(err.Error(), "duplicate DREEFILE") {
		t.Fatalf("error must mention duplicate DREEFILE, got: %v", err)
	}
	var he *HeaderError
	if !errors.As(err, &he) {
		t.Fatalf("expected *HeaderError, got %T: %v", err, err)
	}
	if he.Line != 2 || he.Col != 1 {
		t.Fatalf("expected 2:1, got %d:%d (%v)", he.Line, he.Col, err)
	}
}

func TestParseFileHeaderRejectsConflictingDreefileKind(t *testing.T) {
	_, _, err := ParseFileHeaderStrict("DREEFILE component (title string)\nDREEFILE page\n\n<body>x</body>")
	if err == nil {
		t.Fatal("expected conflicting DREEFILE kind error")
	}
	if !strings.Contains(err.Error(), "duplicate DREEFILE") {
		t.Fatalf("error must mention duplicate DREEFILE, got: %v", err)
	}
	assertHeaderPosition(t, err)
}

func TestParseFileHeaderRejectsDuplicateLayout(t *testing.T) {
	_, _, err := ParseFileHeaderStrict("LAYOUT \"a.dreego\"\nLAYOUT \"b.dreego\"\n\n<body>x</body>")
	if err == nil {
		t.Fatal("expected duplicate LAYOUT error")
	}
	if !strings.Contains(err.Error(), "duplicate LAYOUT") {
		t.Fatalf("error must mention duplicate LAYOUT, got: %v", err)
	}
	var he *HeaderError
	if !errors.As(err, &he) {
		t.Fatalf("expected *HeaderError, got %T: %v", err, err)
	}
	if he.Line != 2 || he.Col != 1 {
		t.Fatalf("expected 2:1, got %d:%d (%v)", he.Line, he.Col, err)
	}
}

func assertHeaderPosition(t *testing.T, err error) {
	t.Helper()
	var he *HeaderError
	if !errors.As(err, &he) {
		t.Fatalf("expected *HeaderError, got %T: %v", err, err)
	}
	if he.Line < 1 || he.Col < 1 {
		t.Fatalf("expected a positive line:col, got %d:%d", he.Line, he.Col)
	}
}
