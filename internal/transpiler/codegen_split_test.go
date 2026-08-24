package transpiler

import (
	"strings"
	"testing"
)

// splitServerSections with hasFormActions=true must move type/func declarations into
// pkgCode and inline the rest, skipping typed/custom sections.
func TestSplitServerSectionsDeclarationWithFormActions(t *testing.T) {
	sections := []ServerSection{
		{Code: "type Item struct {\n\tName string\n}", ContentType: ""},
		{Code: "x := 1\n_ = x", ContentType: ""},
		{Code: "c.W.Write([]byte(\"{}\"))", ContentType: "json"},
		{Code: "// custom", ContentType: "custom"},
	}
	pkg, inline := splitServerSections(sections, true)

	if !strings.Contains(pkg, "type Item struct") {
		t.Errorf("declaration must go to pkgCode, got pkg:\n%s", pkg)
	}
	if strings.Contains(pkg, "x := 1") {
		t.Errorf("non-declaration must not go to pkgCode, got pkg:\n%s", pkg)
	}
	if !strings.Contains(inline, "x := 1") {
		t.Errorf("non-declaration must go to inlineCode, got inline:\n%s", inline)
	}
	// json sections are typed and skipped; custom sections pass through as inline.
	if strings.Contains(inline, "json") {
		t.Errorf("typed sections must be skipped, got inline:\n%s", inline)
	}
	if !strings.Contains(inline, "// custom") {
		t.Errorf("custom sections must be kept in inline, got inline:\n%s", inline)
	}
}

// splitServerSections without form actions must inline declarations too.
func TestSplitServerSectionsNoFormActions(t *testing.T) {
	sections := []ServerSection{
		{Code: "func helper() {}", ContentType: ""},
	}
	pkg, inline := splitServerSections(sections, false)
	if pkg != "" {
		t.Errorf("without form actions no pkgCode expected, got: %q", pkg)
	}
	if !strings.Contains(inline, "func helper()") {
		t.Errorf("declaration must be inlined without form actions, got: %q", inline)
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
