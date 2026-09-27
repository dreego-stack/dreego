package tests

import (
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestMarkdownFenceLanguageCannotInjectAttribute(t *testing.T) {
	out, err := dreego.MarkdownToHTML("```0 onerror=\n```")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(strings.ToLower(out), "onerror=") {
		t.Fatalf("generated code fence contains an event attribute: %s", out)
	}
}
