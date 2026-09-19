package lua

import (
	"strings"
	"testing"
)

func TestCompileLuaStringEscapes(t *testing.T) {
	artifact, err := Compile(`local value = "line\nhex:\x41 decimal:\065 unicode:\u{1F680}\z
    trimmed"`)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`line\nhex:A decimal:A unicode:🚀trimmed`} {
		if !strings.Contains(artifact.Code, want) {
			t.Fatalf("JavaScript missing %q:\n%s", want, artifact.Code)
		}
	}
}

func TestCompileLuaEscapedNewline(t *testing.T) {
	artifact, err := Compile("local value = \"first\\\nsecond\"")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Code, `"first\nsecond"`) {
		t.Fatalf("JavaScript = %s", artifact.Code)
	}
}

func TestCompileLuaEscapedCRLF(t *testing.T) {
	artifact, err := Compile("local value = \"first\\\r\nsecond\"")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(artifact.Code, `"first\nsecond"`) {
		t.Fatalf("JavaScript = %s", artifact.Code)
	}
}

func TestCompileRejectsInvalidLuaStringEscapes(t *testing.T) {
	for _, source := range []string{
		`local value = "\q"`,
		`local value = "\x4"`,
		`local value = "\300"`,
		`local value = "\u{110000}"`,
		`local value = "\u{1234567}"`,
	} {
		if _, err := Compile(source); err == nil {
			t.Fatalf("Compile(%q) succeeded", source)
		}
	}
}
