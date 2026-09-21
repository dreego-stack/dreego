package dreefile

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

type serverCase struct {
	name         string
	declKind     string
	shape        string
	formAction   bool
	decl         string
	statement    string
	formSupport  string
	wantPackage  []string
	wantInRender []string
}

func matrixCases() []serverCase {
	decls := map[string]struct {
		decl      string
		statement string
		wantDecl  string
	}{
		"type": {
			decl:      "type Item struct {\n\tV int\n}",
			statement: "item := Item{V: 1}\n_ = item",
			wantDecl:  "type Item struct",
		},
		"func": {
			decl:      "func helper() int {\n\treturn 42\n}",
			statement: "n := helper()\n_ = n",
			wantDecl:  "func helper() int",
		},
		"const": {
			decl:      "const answer = 42",
			statement: "v := answer\n_ = v",
			wantDecl:  "const answer = 42",
		},
		"var": {
			decl:      "var counter int",
			statement: "counter = 1\n_ = counter",
			wantDecl:  "var counter int",
		},
	}
	shapes := []string{"decl-only", "statements-only", "mixed"}
	kinds := []string{"type", "func", "const", "var"}
	var cases []serverCase
	for _, shape := range shapes {
		for _, kind := range kinds {
			for _, action := range []bool{false, true} {
				d := decls[kind]
				cases = append(cases, serverCase{
					name:        shape + "/" + kind + "/action=" + boolName(action),
					declKind:    kind,
					shape:       shape,
					formAction:  action,
					decl:        d.decl,
					statement:   d.statement,
					wantPackage: []string{d.wantDecl},
				})
			}
		}
	}
	return cases
}

func boolName(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func (c serverCase) source() string {
	var parts []string
	includeDecl := c.shape != "statements-only"
	includeStatement := c.shape != "decl-only"
	if includeDecl {
		parts = append(parts, c.decl)
	}
	if c.formAction {
		parts = append(parts, "type Form struct {\n\tA string `form:\"a\"`\n}", "func Act(c dreego.Context, form Form) error {\n\treturn nil\n}")
	}
	if includeStatement {
		parts = append(parts, c.statement)
	}
	code := strings.Join(parts, "\n\n")
	body := "<body><p>ok</p></body>"
	if c.formAction {
		body = `<body><form g-action="Act" method="post"><input name="a"><button>Go</button></form></body>`
	}
	return "<server>\n" + code + "\n</server>\n" + body
}

func TestServerSectionMatrixGeneratesValidGo(t *testing.T) {
	for _, c := range matrixCases() {
		t.Run(c.name, func(t *testing.T) {
			file := parseFile(t, c.source())
			if c.formAction {
				file.FormActions = []string{"Act"}
			}
			out, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "routes", "index", "/", "abc")
			if err != nil {
				t.Fatalf("GenerateMethodHandler: %v", err)
			}
			if _, err := parser.ParseFile(token.NewFileSet(), "dree.go", "package routes\n"+out, 0); err != nil {
				t.Fatalf("generated Go does not parse: %v\n---\n%s", err, out)
			}
			if c.shape != "statements-only" {
				renderAt := strings.Index(out, "func renderIndex")
				pkgAt := strings.Index(out, c.wantPackage[0])
				if pkgAt < 0 || renderAt < 0 || pkgAt > renderAt {
					t.Fatalf("declaration must be emitted at package level before render, got:\n%s", out)
				}
			}
			if c.shape != "decl-only" {
				renderAt := strings.Index(out, "func renderIndex")
				firstStatement := strings.Index(out[renderAt:], strings.SplitN(c.statement, "\n", 2)[0])
				if firstStatement < 0 {
					t.Fatalf("statement must be emitted inside render, got:\n%s", out)
				}
			}
		})
	}
}

func TestServerSectionMixedDeclarationAndStatementSplit(t *testing.T) {
	sections := []ServerSection{{
		Code: "type T struct{ V int }\nfunc (t T) Get() int { return t.V }\nx := T{V: 5}\n_ = x",
	}}
	pkg, inline := splitServerSections(sections)
	for _, want := range []string{"type T struct", "func (t T) Get() int"} {
		if !strings.Contains(pkg, want) {
			t.Errorf("mixed section declaration %q must go to package level, got:\n%s", want, pkg)
		}
	}
	if strings.Contains(pkg, "x := T{V: 5}") {
		t.Errorf("mixed section statement must not be package level, got:\n%s", pkg)
	}
	if !strings.Contains(inline, "x := T{V: 5}") {
		t.Errorf("mixed section statement must be inline, got:\n%s", inline)
	}
}

func TestServerSectionStatementsOnlyStayInline(t *testing.T) {
	pkg, inline := splitServerSections([]ServerSection{{Code: "a := 1\nb := 2\n_ = a + b"}})
	if pkg != "" {
		t.Fatalf("statements-only must not emit package code, got:\n%s", pkg)
	}
	if !strings.Contains(inline, "a := 1") || !strings.Contains(inline, "b := 2") {
		t.Fatalf("statements must stay inline, got:\n%s", inline)
	}
}

func TestServerSectionConstBlockIsPackageLevel(t *testing.T) {
	sections := []ServerSection{{Code: "const (\n\tA = 1\n\tB = 2\n)\n_ = A"}}
	pkg, inline := splitServerSections(sections)
	if !strings.Contains(pkg, "const (") || !strings.Contains(pkg, "B = 2") {
		t.Fatalf("const block must be package level, got:\n%s", pkg)
	}
	if strings.Contains(inline, "B = 2") {
		t.Fatalf("const block must not be inline, got:\n%s", inline)
	}
}
