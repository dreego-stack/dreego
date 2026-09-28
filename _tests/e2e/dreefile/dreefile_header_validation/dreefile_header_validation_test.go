package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDreefileGrammarRejectsEmptyGoImportList(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "GOIMPORT {  }\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an empty GOIMPORT list:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "1:", "GOIMPORT"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarRejectsUnquotedGoImportPath(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "GOIMPORT { encoding json }\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an unquoted GOIMPORT path:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "1:", "GOIMPORT", "encoding json"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarRejectsUnquotedLayout(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "LAYOUT www/layouts/x.dreego\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an unquoted LAYOUT path:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "1:", "LAYOUT"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarRejectsEmptyLayout(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "LAYOUT \"\"\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an empty LAYOUT path:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "1:", "LAYOUT"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarRejectsNestedBraceList(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "GOIMPORT { a, { b } }\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted a nested brace list:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "1:", "nested"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarRejectsTrailingContentAfterBraceList(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Card.dreego": "DREEFILE component ()\n<body><p>x</p></body>",
		"www/routes/+page.dreego":    "COMPONENT \"www/components\" IMPORT { Card } extra\n<body><@Card/></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted trailing content after a brace list:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "1:", "after closing"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarRejectsDuplicateDreefile(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "DREEFILE page\nDREEFILE layout\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted duplicate DREEFILE lines:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "2:", "duplicate DREEFILE"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}

func TestDreefileGrammarRejectsDuplicateLayout(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "LAYOUT \"a.dreego\"\nLAYOUT \"b.dreego\"\n<body><p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted duplicate LAYOUT lines:\n%s", out)
	}
	for _, want := range []string{"+page.dreego", "2:", "duplicate LAYOUT"} {
		if !strings.Contains(out, want) {
			t.Fatalf("diagnostic must contain %q, got:\n%s", want, out)
		}
	}
}
