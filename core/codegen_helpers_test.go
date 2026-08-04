package core

import (
	"strings"
	"testing"
)

func TestScopeCSSKeepsDeclarationsWithBraces(t *testing.T) {
	css := ".bg-grid { background-image: radial-gradient(circle, #ccfbf1 1px, transparent 1px); }"
	out := scopeCSS(css, "abc")

	if !strings.Contains(out, "background-image") {
		t.Errorf("declaration name lost: got %q", out)
	}
	if !strings.Contains(out, "radial-gradient(circle, #ccfbf1 1px, transparent 1px)") {
		t.Errorf("declaration value with commas/parens lost: got %q", out)
	}
	if !strings.Contains(out, "[data-scope=abc]") {
		t.Errorf("selector not scoped: got %q", out)
	}
}

func TestScopeCSSMediaPreservesDeclarationsAndScopesInnerSelectors(t *testing.T) {
	css := "@media (min-width: 600px) { .box { color: blue; } }"
	out := scopeCSS(css, "abc")

	if !strings.Contains(out, "@media (min-width: 600px)") {
		t.Errorf("@media rule dropped: got %q", out)
	}
	if !strings.Contains(out, "color: blue") {
		t.Errorf("@media declaration dropped: got %q", out)
	}
	if !strings.Contains(out, "[data-scope=abc] .box") {
		t.Errorf("inner selector not scoped: got %q", out)
	}
}

func TestScopeCSSPseudoSelectorKeepsDeclaration(t *testing.T) {
	css := ".box:hover { color: green; }"
	out := scopeCSS(css, "abc")

	if !strings.Contains(out, ":hover") {
		t.Errorf("pseudo selector lost: got %q", out)
	}
	if !strings.Contains(out, "color: green") {
		t.Errorf("declaration lost: got %q", out)
	}
	if !strings.Contains(out, "[data-scope=abc] .box:hover") {
		t.Errorf("pseudo selector not scoped: got %q", out)
	}
}
