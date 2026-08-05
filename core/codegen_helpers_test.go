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

// extractAttrValues must resolve {expr} inside a quoted attribute value as a
// Go expression (unquoted), not a literal string. Currently attrVal trims the
// quotes and returns "{url}" as a literal.
func TestExtractAttrValuesResolvesExprInQuotedValue(t *testing.T) {
	out := extractAttrValues(`href="{url}" label="x"`)
	if strings.Contains(out, "{url}") {
		t.Errorf("extractAttrValues left {url} literal, must resolve to url. got: %s", out)
	}
	if !strings.Contains(out, "url") {
		t.Errorf("extractAttrValues must contain url expression, got: %s", out)
	}
}

// attrVal must resolve multiple placeholders inside one quoted attribute value
// (`href="{a}-{b}"`) into a concatenated Go expression, not treat the whole
// inner "a}-{b" as a single expression. Currently attrVal returns "a}-{b"
// verbatim because it only handles the exact `{...}` (single placeholder) shape.
func TestAttrValResolvesMultiplePlaceholdersToConcatenation(t *testing.T) {
	out := attrVal(`href="{a}-{b}"`)
	if strings.Contains(out, "a}-{b") {
		t.Errorf("attrVal must split {a}-{b} into separate expressions, got: %s", out)
	}
	if !strings.Contains(out, "a") || !strings.Contains(out, "b") {
		t.Errorf("attrVal must keep both placeholder expressions a and b, got: %s", out)
	}
	if !strings.Contains(out, `"-"`) {
		t.Errorf("attrVal must join placeholders with a literal separator \"-\", got: %s", out)
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

// concatPlaceholders builds a component-call argument for a quoted attribute
// value with mixed literal text and multiple {…} placeholders. Escaping must be
// deferred to the prop-injection point (the component body escapes its own
// placeholders), so the call argument must NOT itself wrap the expression in
// html.EscapeString. Otherwise multi-placeholder calls would be double-escaped.
func TestConcatPlaceholdersDoesNotEscape(t *testing.T) {
	out := concatPlaceholders(`href="{a}-{b}"`)
	if strings.Contains(out, "html.EscapeString") {
		t.Errorf("concatPlaceholders must not escape, escaping is deferred to prop injection. got: %s", out)
	}
	if !strings.Contains(out, "fmt.Sprintf(\"%v\", a)") {
		t.Errorf("concatPlaceholders must emit fmt.Sprintf for placeholder a, got: %s", out)
	}
	if !strings.Contains(out, "fmt.Sprintf(\"%v\", b)") {
		t.Errorf("concatPlaceholders must emit fmt.Sprintf for placeholder b, got: %s", out)
	}
	if !strings.Contains(out, `"-"`) {
		t.Errorf("concatPlaceholders must keep the literal separator \"-\", got: %s", out)
	}
}
