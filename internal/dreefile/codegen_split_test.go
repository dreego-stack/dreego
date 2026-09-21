package dreefile

import (
	"strings"
	"testing"
)

// splitServerSections must move the leading declaration block to package level
// and keep statements inline, skipping typed sections and passing custom
// sections through as inline. hasFormActions must not change that decision.
func TestSplitServerSectionsDeclarationWithStatements(t *testing.T) {
	sections := []ServerSection{
		{Code: "type Item struct {\n\tName string\n}", ContentType: ""},
		{Code: "x := 1\n_ = x", ContentType: ""},
		{Code: "c.W.Write([]byte(\"{}\"))", ContentType: "json"},
		{Code: "// custom", ContentType: "custom"},
	}
	pkg, inline := splitServerSections(sections, map[string]bool{})

	if !strings.Contains(pkg, "type Item struct") {
		t.Errorf("declaration must go to pkgCode, got pkg:\n%s", pkg)
	}
	if strings.Contains(pkg, "x := 1") {
		t.Errorf("non-declaration must not go to pkgCode, got pkg:\n%s", pkg)
	}
	if !strings.Contains(inline, "x := 1") {
		t.Errorf("non-declaration must go to inlineCode, got inline:\n%s", inline)
	}
	if strings.Contains(inline, "json") {
		t.Errorf("typed sections must be skipped, got inline:\n%s", inline)
	}
	if !strings.Contains(inline, "// custom") {
		t.Errorf("custom sections must be kept in inline, got inline:\n%s", inline)
	}
}

func TestUnindentMixed(t *testing.T) {
	in := "    func a() {\n        x := 1\n    }"
	out := unindent(in)
	want := "func a() {\n    x := 1\n}"
	if out != want {
		t.Errorf("unindent mixed = %q, want %q", out, want)
	}
}

func TestUnindentNoIndent(t *testing.T) {
	in := "func a() {\n\tx := 1\n}"
	out := unindent(in)
	if out != in {
		t.Errorf("unindent no-indent must return unchanged, got: %q", out)
	}
}
