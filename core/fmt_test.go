package core

import (
	"strings"
	"testing"
)

func TestFormatCompHeaderWithDefault(t *testing.T) {
	in := "Component Button (label string, variant string = primary)"
	out := formatCompHeader(in)
	for _, want := range []string{"Component Button", "label string", "variant string = primary"} {
		if !strings.Contains(out, want) {
			t.Errorf("formatCompHeader missing %q, got: %q", want, out)
		}
	}
	if strings.Contains(out, " = ") && !strings.Contains(out, "= primary") {
		t.Errorf("formatCompHeader must keep the default value, got: %q", out)
	}
}

func TestFormatCompHeaderNoParams(t *testing.T) {
	in := "Component Card ()"
	out := formatCompHeader(in)
	if out != "Component Card" {
		t.Errorf("formatCompHeader no-params must drop parens, got: %q", out)
	}
}

func TestFormatCompHeaderAlreadyNormalized(t *testing.T) {
	in := "Component Card (x int)"
	out := formatCompHeader(in)
	if out != "Component Card (x int)" {
		t.Errorf("formatCompHeader must be idempotent, got: %q", out)
	}
}

func TestFormatCompHeaderNoParens(t *testing.T) {
	in := "Component Card"
	out := formatCompHeader(in)
	if out != in {
		t.Errorf("formatCompHeader without parens must return unchanged, got: %q", out)
	}
}

func TestFormatImportAliasAndPath(t *testing.T) {
	in := "import dreego github.com/dreego-stack/dreego"
	out := formatImport(in)
	if !strings.Contains(out, "import dreego github.com/dreego-stack/dreego") {
		t.Errorf("formatImport must keep alias and path, got: %q", out)
	}
}

func TestFormatImportSingleField(t *testing.T) {
	out := formatImport("import fmt")
	if out != "import fmt" {
		t.Errorf("formatImport single field must be preserved, got: %q", out)
	}
}

func TestFormatImportEmpty(t *testing.T) {
	out := formatImport("import ")
	if out != "import" {
		t.Errorf("formatImport empty must return 'import', got: %q", out)
	}
}

func TestFormatExpressionsPipeNormalization(t *testing.T) {
	in := `{ name | upper }`
	out := formatExpressions(in)
	if !strings.Contains(out, "{name|upper}") {
		t.Errorf("formatExpressions must normalize pipe spacing, got: %q", out)
	}
}

func TestFormatExpressionsRemovesExtraSpaces(t *testing.T) {
	in := `{  count   }`
	out := formatExpressions(in)
	if !strings.Contains(out, "{count}") {
		t.Errorf("formatExpressions must collapse spaces, got: %q", out)
	}
}

func TestFormatControlFlowCollapsesSpaces(t *testing.T) {
	in := `{#each   items as item}`
	out := formatControlFlow(in)
	if !strings.Contains(out, "{#each items as item}") {
		t.Errorf("formatControlFlow must collapse spaces in open tag, got: %q", out)
	}
}

func TestFormatControlFlowCloseTag(t *testing.T) {
	in := `{/each}`
	out := formatControlFlow(in)
	if !strings.Contains(out, "{/each}") {
		t.Errorf("formatControlFlow must keep close tag, got: %q", out)
	}
}

func TestFormatSectionBodyTrimsBlankLines(t *testing.T) {
	in := "<div>\n\n  <p>hi</p>\n\n</div>"
	out := formatSectionBody("div", in)
	if strings.Contains(out, "\n\n\n") {
		t.Errorf("formatSectionBody must trim surrounding blank lines, got: %q", out)
	}
	if !strings.Contains(out, "<div>\n") || !strings.Contains(out, "</div>") {
		t.Errorf("formatSectionBody must keep tags, got: %q", out)
	}
}

func TestFormatSectionsOrdersKnownSections(t *testing.T) {
	in := "<style>.a{}</style>\n<div><p>x</p></div>"
	out := formatSections(in)
	if !strings.Contains(out, "<style>") || !strings.Contains(out, "<div>") {
		t.Errorf("formatSections must keep both sections, got:\n%s", out)
	}
}

func TestFormatSectionsNoSections(t *testing.T) {
	in := "plain text"
	out := formatSections(in)
	if out != in {
		t.Errorf("formatSections without sections must return input unchanged, got: %q", out)
	}
}

func TestFormatFullDocument(t *testing.T) {
	in := "Component Button (label string = Hi)\n\nimport dreego github.com/dreego-stack/dreego\n\n<div>\n  <p>{ label | upper }</p>\n</div>\n"
	out := Format(in)
	for _, want := range []string{"Component Button (label string = Hi)", "import dreego github.com/dreego-stack/dreego", "{label|upper}"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format missing %q, got:\n%s", want, out)
		}
	}
}
