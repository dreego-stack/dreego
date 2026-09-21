package dreefile

import (
	"strings"
	"testing"
)

// The formatter normalizes expression and control-flow spacing, but must never
// rewrite string-literal contents. Feedback 3.1: fmt silently changed
// `{{ "a  b" }}` and `{{ "a | b" }}`, altering program behavior.

func TestFormatKeepsStringLiteralSpaces(t *testing.T) {
	in := "<server>\nx := 1\n</server>\n\n<body><p>{{ \"a  b\" }}</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, `"a  b"`) {
		t.Errorf("formatter changed spaces inside a string literal: got %q", out)
	}
}

func TestFormatKeepsStringLiteralPipe(t *testing.T) {
	in := "<server>\nx := 1\n</server>\n\n<body><p>{{ \"a | b\" }}</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, `"a | b"`) {
		t.Errorf("formatter changed a pipe inside a string literal: got %q", out)
	}
}

func TestFormatKeepsConditionStringLiteral(t *testing.T) {
	in := "<server>\ncond := true\n</server>\n\n<body>{#if cond == \"a  b\"}<p>x</p>{/if}</body>\n"
	out := Format(in)
	if !strings.Contains(out, `"a  b"`) {
		t.Errorf("formatter changed a string literal in a condition: got %q", out)
	}
}

// A server section is Go source and must survive formatting byte for byte,
// including raw-string blank lines and alignment spacing.
func TestFormatKeepsServerRawStringBlankLines(t *testing.T) {
	in := "<server>\nraw := `line1\n\n\nline3`\n</server>\n\n<body><p>ok</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, "line1\n\n\nline3") {
		t.Errorf("formatter collapsed blank lines inside a server raw string: got %q", out)
	}
}

func TestFormatKeepsServerSpacing(t *testing.T) {
	in := "<server>\nmsg   :=   \"a  b\"\n</server>\n\n<body><p>ok</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, "msg   :=   \"a  b\"") {
		t.Errorf("formatter rewrote server spacing: got %q", out)
	}
}

func TestFormatStillNormalizesExpressionSpacing(t *testing.T) {
	in := "<server>\nmsg := \"hi\"\n</server>\n\n<body><p>{{  msg  }}</p></body>\n"
	out := Format(in)
	if !strings.Contains(out, "{{ msg }}") {
		t.Errorf("formatter must still normalize expression spacing, got %q", out)
	}
}
