package dreegotest_test

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestGenerateComponentUsesExplicitNameForDreefile(t *testing.T) {
	t.Parallel()
	src := "DREEFILE component (title string)\n<body><h1>{{ title }}</h1></body>"
	out := dreegotest.GenerateComponent(t, "Card", src)
	if out == "" {
		t.Fatal("GenerateComponent returned empty output for a DREEFILE component")
	}
	if !strings.Contains(out, "func Card(title string) dreego.Component") {
		t.Fatalf("generated output must declare func Card, got:\n%s", out)
	}
}

func TestMustCompileComponentAcceptsDreefileComponent(t *testing.T) {
	t.Parallel()
	dreegotest.MustCompileComponent(t, "Badge", "DREEFILE component (label string)\n<body class=\"badge\">{{ label }}</body>")
}
