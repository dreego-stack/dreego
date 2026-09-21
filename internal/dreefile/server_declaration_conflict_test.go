package dreefile

import (
	"strings"
	"testing"
)

// One route file with two method sections repeats the same leading declaration.
// It must be emitted once, not once per method.
func TestServerDeclarationEmittedOncePerFile(t *testing.T) {
	src := `<server>
type T struct{ V int }
</server>
<body><p>{{ (T{V: 1}).V }}</p></body>
<server method="post">
type T struct{ V int }
</server>
<body method="post"><p>x</p></body>`
	file := parseFile(t, src)
	out, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "routes", "index", "/", "abc")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if got := strings.Count(out, "type T struct"); got != 1 {
		t.Fatalf("declaration emitted %d times, want 1:\n%s", got, out)
	}
}

// Two route files that hoist the same package-level name must fail generation
// with a dreego diagnostic instead of a raw compiler redeclaration error.
func TestServerDeclarationConflictIsDiagnostic(t *testing.T) {
	err := serverDeclarationConflict("Product", "www/routes/a.dreego", "www/routes/b.dreego")
	if err == nil {
		t.Fatal("expected a conflict error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Product") || !strings.Contains(msg, "a.dreego") || !strings.Contains(msg, "Fix:") {
		t.Fatalf("diagnostic must name the symbol, the first file and a fix: %s", msg)
	}
}

func TestDeclarationName(t *testing.T) {
	cases := map[string]string{
		"type Product struct{ Name string }": "Product",
		"func Helper() int { return 1 }":     "Helper",
		"func (t T) Get() int { return 1 }":  "",
		"var store = 1":                      "store",
		"const answer = 42":                  "answer",
		"const ( A = 1 )":                    "",
		"x := 1":                             "",
	}
	for src, want := range cases {
		if got := declarationName(src); got != want {
			t.Errorf("declarationName(%q) = %q, want %q", src, got, want)
		}
	}
}
