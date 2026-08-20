package transpiler

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

func TestGoLiteralPlain(t *testing.T) {
	out := goLiteral(`<p>hi</p>`)
	if out != "`<p>hi</p>`" {
		t.Errorf("goLiteral plain must wrap in backticks, got: %q", out)
	}
}

func TestGoLiteralContainsBacktick(t *testing.T) {
	out := goLiteral("a ` b")
	if strings.Contains(out, "`") && !strings.Contains(out, "\"") {
		t.Errorf("goLiteral with backtick must use quotes, got: %q", out)
	}
	if !strings.Contains(out, `"a `) {
		t.Errorf("goLiteral with backtick must produce a quoted string, got: %q", out)
	}
}

func TestToPascalCaseVariants(t *testing.T) {
	cases := []struct{ in, want string }{
		{"hello world", "HelloWorld"},
		{"user-id", "UserId"},
		{"my_page", "MyPage"},
		{"123abc", "123abc"},
		{"alreadyPascal", "Alreadypascal"},
	}
	for _, c := range cases {
		if got := toPascalCase(c.in); got != c.want {
			t.Errorf("toPascalCase(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestScopeSelectorCommaList(t *testing.T) {
	out := scopeSelector(".a, .b:hover", "[data-scope=x] ")
	if !strings.Contains(out, "[data-scope=x] .a") {
		t.Errorf("scopeSelector must scope first selector, got: %q", out)
	}
	if !strings.Contains(out, "[data-scope=x] .b:hover") {
		t.Errorf("scopeSelector must scope second selector, got: %q", out)
	}
}

func TestSplitTopLevelCommaNested(t *testing.T) {
	in := ":not(.a), .b, .c:hover, div > span"
	out := splitTopLevelComma(in)
	if len(out) != 4 {
		t.Fatalf("splitTopLevelComma must split 4 top-level selectors, got %d: %#v", len(out), out)
	}
}

func TestMatchBraceUnbalanced(t *testing.T) {
	// "a{b{c}" — opening brace at index 1 never closes before end.
	css := "a{b{c}"
	got := matchBrace(css, 1, len(css))
	if got != len(css) {
		t.Errorf("matchBrace unbalanced must return end index, got %d, want %d", got, len(css))
	}
}

func TestMatchBraceBalanced(t *testing.T) {
	// "a{b{c}d}": brace at index 3 (c) closes at index 5.
	css := "a{b{c}d}"
	if got := matchBrace(css, 3, len(css)); got != 5 {
		t.Errorf("matchBrace nested must return matching close, got %d, want 5", got)
	}
}

// attrVal must resolve a quoted {var} attribute value as a Go expression, not a
// literal string. Feedback (testspace/feedback.md, core v0.0.23): prop="{var}"
// generated []string{"{var}"} for slices and a compile error for bools.
func TestAttrValQuotedExpr(t *testing.T) {
	out := attrVal(`prop="{var}"`)
	if out != "var" {
		t.Errorf("attrVal(prop=\"{var}\") = %q, want expression var", out)
	}
}

// attrVal must resolve an unquoted {var} attribute value as a Go expression.
func TestAttrValUnquotedExpr(t *testing.T) {
	out := attrVal(`prop={var}`)
	if out != "var" {
		t.Errorf("attrVal(prop={var}) = %q, want expression var", out)
	}
}

// attrVal must keep a plain quoted literal as a Go string literal.
func TestAttrValLiteral(t *testing.T) {
	out := attrVal(`prop="literal"`)
	if out != `"literal"` {
		t.Errorf("attrVal(prop=\"literal\") = %q, want string literal \"literal\"", out)
	}
}

// attrVal on a quoted bool-looking string keeps it a string literal. This is
// documented behavior: a bool prop must be passed unquoted (active={true}),
// passing active="false" yields a string and fails to compile at the call site.
func TestAttrValQuotedBoolString(t *testing.T) {
	out := attrVal(`active="false"`)
	if out != `"false"` {
		t.Errorf("attrVal(active=\"false\") = %q, want string literal \"false\"", out)
	}
}

// concatPlaceholders builds a component-call argument for a quoted attribute
// value with mixed literal text and multiple {…} placeholders. Escaping must be
// deferred to the prop-injection point (the component body escapes its own
// placeholders via dreego.SafeText/dreego.SafeAttr), so the call argument must
// NOT itself escape the expression. Otherwise multi-placeholder calls would be
// double-escaped.
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
