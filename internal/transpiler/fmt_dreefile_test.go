package transpiler

import (
	"strings"
	"testing"
)

func TestFormatKeepsDreefileDirectivesInHeader(t *testing.T) {
	in := `DREEFILE component (title string)

COMPONENT "www/components" IMPORT {
    Card,
    Button,
}

GOIMPORT { sync, encoding/json }

LAYOUT "www/layouts/admin.dreego"

<body>
  <p>{{  title  }}</p>
</body>
`
	out := Format(in)

	if !strings.HasPrefix(out, "DREEFILE component (title string)") {
		t.Fatalf("DREEFILE line must stay in the header, got:\n%s", out)
	}
	for _, want := range []string{
		`COMPONENT "www/components" IMPORT {`,
		`GOIMPORT { sync, encoding/json }`,
		`LAYOUT "www/layouts/admin.dreego"`,
		"    Card,",
		"    Button,",
		"{{ title }}",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Format missing %q, got:\n%s", want, out)
		}
	}

	bodyIdx := strings.Index(out, "<body>")
	for _, directive := range []string{"DREEFILE", "COMPONENT", "GOIMPORT", "LAYOUT"} {
		if idx := strings.Index(out, directive); idx < 0 || idx > bodyIdx {
			t.Errorf("%s must precede <body> (idx %d, body %d), got:\n%s", directive, idx, bodyIdx, out)
		}
	}
}

func TestFormatDreefileDirectivesIdempotent(t *testing.T) {
	in := `DREEFILE layout

LAYOUT "www/layouts/base.dreego"

GOIMPORT { strings }

COMPONENT "www/components" IMPORT { Nav }
`
	once := Format(in)
	twice := Format(once)
	if once != twice {
		t.Fatalf("Format must be idempotent\nonce:  %q\ntwice: %q", once, twice)
	}
	if !strings.HasPrefix(once, "DREEFILE layout") {
		t.Fatalf("directives must stay in the header, got:\n%s", once)
	}
}

func TestFormatLayoutLineNormalizesSpacing(t *testing.T) {
	out := formatLayoutLine(`LAYOUT   "www/layouts/admin.dreego"`)
	if out != `LAYOUT "www/layouts/admin.dreego"` {
		t.Fatalf("expected normalized LAYOUT line, got %q", out)
	}
}

func TestFormatUnbalancedDirectiveKeepsBody(t *testing.T) {
	in := `DREEFILE component (title string)

COMPONENT "www/components" IMPORT { Card

<body>
  <p>{{ title }}</p>
</body>
`
	out := Format(in)
	if idx := strings.Index(out, "<body>"); idx < 0 {
		t.Fatalf("body must survive an unterminated directive, got:\n%s", out)
	}
	if strings.Count(out, "<body>") != 1 {
		t.Fatalf("body must not be duplicated, got:\n%s", out)
	}
	bodyIdx := strings.Index(out, "<body>")
	if idx := strings.Index(out, "COMPONENT"); idx < 0 || idx > bodyIdx {
		t.Fatalf("COMPONENT directive must stay in the header, got:\n%s", out)
	}
	twice := Format(out)
	if twice != out {
		t.Fatalf("Format must be idempotent for an unbalanced directive\nonce:  %q\ntwice: %q", out, twice)
	}
}

func TestFormatUnbalancedDirectiveDoesNotHoistBodySection(t *testing.T) {
	in := `COMPONENT "www/components" IMPORT { Card
<server>msg := "hi"</server>
<body><p>{{ msg }}</p></body>
`
	out := Format(in)
	bodyIdx := strings.Index(out, "<body>")
	if bodyIdx < 0 {
		t.Fatalf("body section must survive, got:\n%s", out)
	}
	if strings.Count(out, "<server>") != 1 || strings.Count(out, "<body>") != 1 {
		t.Fatalf("sections must not be duplicated, got:\n%s", out)
	}
	if idx := strings.Index(out, "COMPONENT"); idx > bodyIdx {
		t.Fatalf("COMPONENT directive must stay before the sections, got:\n%s", out)
	}
}

func TestFormatBalancedDirectiveStillGroupsMultilineBlock(t *testing.T) {
	in := `COMPONENT "www/components" IMPORT {
    Card,
    Button,
}
<body><@Card/></body>
`
	out := Format(in)
	if !strings.Contains(out, "    Card,") || !strings.Contains(out, "    Button,") {
		t.Fatalf("balanced multiline import block must stay in the header, got:\n%s", out)
	}
	bodyIdx := strings.Index(out, "<body>")
	if idx := strings.Index(out, "Button"); idx > bodyIdx {
		t.Fatalf("import block must precede the body, got:\n%s", out)
	}
}

func TestFormatBareDreefileStaysInHeader(t *testing.T) {
	in := "DREEFILE\n<body><p>x</p></body>\n"
	out := Format(in)
	bodyIdx := strings.Index(out, "<body>")
	if idx := strings.Index(out, "DREEFILE"); idx < 0 || idx > bodyIdx {
		t.Fatalf("bare DREEFILE must stay before the body, got:\n%s", out)
	}
	twice := Format(out)
	if twice != out {
		t.Fatalf("Format must be idempotent for a bare DREEFILE\nonce:  %q\ntwice: %q", out, twice)
	}
}

func TestFormatBareDreefileRoundTripStillRejected(t *testing.T) {
	in := "DREEFILE\n<body><p>x</p></body>\n"
	if _, _, err := ParseFileHeaderStrict(in); err == nil {
		t.Fatal("generate must reject a bare DREEFILE")
	}
	out := Format(in)
	if _, _, err := ParseFileHeaderStrict(out); err == nil {
		t.Fatalf("fmt must not turn a rejected bare DREEFILE into an accepted file, got:\n%s", out)
	}
}

func TestFormatCollapsesDuplicateHeaderBlankLines(t *testing.T) {
	in := "DREEFILE component (title string)\n\n\nLAYOUT \"www/layouts/admin.dreego\"\n\n\n<body><p>x</p></body>\n"
	out := Format(in)
	if strings.Contains(out, "\n\n\n") {
		t.Fatalf("Format must collapse duplicate header blank lines, got:\n%q", out)
	}
	if !strings.HasPrefix(out, "DREEFILE component (title string)\n\nLAYOUT") {
		t.Fatalf("expected a single blank line between header directives, got:\n%q", out)
	}
	twice := Format(out)
	if twice != out {
		t.Fatalf("Format must be idempotent\nonce:  %q\ntwice: %q", out, twice)
	}
}

func TestFormatNormalizesCRLF(t *testing.T) {
	in := "DREEFILE component (title string)\r\n\r\n<body>\r\n  <p>{{ title }}</p>\r\n</body>\r\n"
	out := Format(in)
	if strings.Contains(out, "\r") {
		t.Fatalf("Format must not leave carriage returns, got: %q", out)
	}
	if !strings.Contains(out, "{{ title }}") || !strings.Contains(out, "<body>") {
		t.Fatalf("Format lost content, got: %q", out)
	}
}

func TestParseFileHeaderHandlesCRLF(t *testing.T) {
	header, body, err := ParseFileHeaderStrict("DREEFILE layout\r\n\r\n<body>{#slot}</body>\r\n")
	if err != nil {
		t.Fatalf("CRLF header: %v", err)
	}
	if header.Kind != FileKindLayout {
		t.Fatalf("expected layout kind, got %v", header.Kind)
	}
	if !strings.Contains(body, "<body>") {
		t.Fatalf("expected body after CRLF header, got %q", body)
	}
}

func TestFormatKeepsLegacyComponentAndImport(t *testing.T) {
	in := "Component Card (title string)\n\nimport dreego github.com/dreego-stack/dreego\n\n<body><p>{{ title }}</p></body>\n"
	out := Format(in)
	for _, want := range []string{"Component Card (title string)", "import dreego github.com/dreego-stack/dreego"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format must keep legacy header %q, got:\n%s", want, out)
		}
	}
}
