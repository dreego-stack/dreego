package dreefile

import (
	"strings"
	"testing"
)

func TestFormatExpressionsPipeNormalization(t *testing.T) {
	in := `{{ name | upper }}`
	out := formatExpressions(in)
	if !strings.Contains(out, "{{ name|upper }}") {
		t.Errorf("formatExpressions must normalize pipe spacing, got: %q", out)
	}
}

func TestFormatExpressionsRemovesExtraSpaces(t *testing.T) {
	in := `{{  count   }}`
	out := formatExpressions(in)
	if !strings.Contains(out, "{{ count }}") {
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
	in := "<body>\n\n  <p>hi</p>\n\n</body>"
	out := formatSectionBody("body", in)
	if strings.Contains(out, "\n\n\n") {
		t.Errorf("formatSectionBody must trim surrounding blank lines, got: %q", out)
	}
	if !strings.Contains(out, "<body>\n") || !strings.Contains(out, "</body>") {
		t.Errorf("formatSectionBody must keep tags, got: %q", out)
	}
}

func TestFormatSectionsOrdersKnownSections(t *testing.T) {
	in := "<style>.a{}</style>\n<body><p>x</p></body>"
	out := formatSections(in)
	if !strings.Contains(out, "<style>") || !strings.Contains(out, "<body>") {
		t.Errorf("formatSections must keep both sections, got:\n%s", out)
	}
}

func TestFormatSectionsPreservesExplicitDefaultLanguages(t *testing.T) {
	in := `<client lang="js">console.log("ready")</client>
<body lang="html"><p>{{ value }}</p></body>
<server lang="go">value := "ok"</server>`
	out := formatSections(in)
	server := strings.Index(out, `<server lang="go">`)
	body := strings.Index(out, `<body lang="html">`)
	client := strings.Index(out, `<client lang="js">`)
	if server < 0 || body < server || client < body {
		t.Fatalf("semantic sections not preserved in canonical order:\n%s", out)
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
	in := "Component Button (label string = Hi)\n\nimport dreego github.com/dreego-stack/dreego\n\n<body>\n  <p>{{ label | upper }}</p>\n</body>\n"
	out := Format(in)
	for _, want := range []string{"Component Button (label string = Hi)", "import dreego github.com/dreego-stack/dreego", "{{ label|upper }}"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format missing %q, got:\n%s", want, out)
		}
	}
}
