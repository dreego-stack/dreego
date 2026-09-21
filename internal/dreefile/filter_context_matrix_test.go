package dreefile

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

type filterContextCase struct {
	name    string
	text    bool
	pattern string
	safe    string
}

var filterContexts = []filterContextCase{
	{"text", true, `{{ v%s }}`, "dreego.SafeText"},
	{"attr", false, `<a title="{{ v%s }}">`, "dreego.SafeAttr"},
	{"url", false, `<a href="{{ v%s }}">`, "dreego.SafeURL"},
	{"script", false, `<button onclick="{{ v%s }}">`, "dreego.SafeScript"},
	{"style", false, `<div style="{{ v%s }}">`, "dreego.SafeStyle"},
}

var filterMatrix = []string{"", "|raw", "|upper", "|raw|upper"}

// Filter x context matrix: every filter must behave identically in text and
// attribute contexts. The regression was |raw in an attribute emitting
// "undefined: raw" because attribute expressions bypassed the filter pipeline.
func TestFilterContextMatrix(t *testing.T) {
	for _, ctx := range filterContexts {
		for _, filter := range filterMatrix {
			t.Run(ctx.name+"/"+strings.Trim(filter, "|"), func(t *testing.T) {
				code := generateFilterContext(t, ctx, filter)
				if strings.Contains(code, "|raw") || strings.Contains(code, "|upper") {
					t.Fatalf("filter leaked into the Go expression:\n%s", code)
				}
				isRaw := strings.Contains(filter, "raw")
				isUpper := strings.Contains(filter, "upper")
				if isUpper != strings.Contains(code, "strings.ToUpper") {
					t.Errorf("upper mismatch (want=%v):\n%s", isUpper, code)
				}
				if isRaw == strings.Contains(code, ctx.safe) {
					t.Errorf("escaping mismatch (raw=%v, want wrapper %s):\n%s", isRaw, ctx.safe, code)
				}
			})
		}
	}
}

func generateFilterContext(t *testing.T, ctx filterContextCase, filter string) string {
	t.Helper()
	expr := strings.Replace(ctx.pattern, "%s", filter, 1)
	if ctx.text {
		tokens, err := Lex("<body>" + expr + "</body>")
		if err != nil {
			t.Fatalf("lex: %v", err)
		}
		file, err := NewParser(tokens).Parse()
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		out, err := genTemplateNode(NewGenerator(), file.Body.Nodes[0], 1)
		if err != nil {
			t.Fatalf("codegen: %v", err)
		}
		assertGoStmt(t, out)
		return out
	}
	out := compTextWithAttrs(expr)
	assertGoExpr(t, out)
	return out
}

// The exact feedback case: <a href="{{ webcal|raw }}"> must compile to raw
// output, not "dreego.SafeURL(fmt.Sprintf("%v", webcal|raw))".
func TestAttributeRawMatchesTextRaw(t *testing.T) {
	attr := compTextWithAttrs(`<a href="{{ webcal|raw }}">go</a>`)
	text, err := genTemplateNode(NewGenerator(), TemplateNode{
		Type:    NodeExpression,
		Content: "webcal",
		Filters: []string{"raw"},
	}, 1)
	if err != nil {
		t.Fatalf("text codegen: %v", err)
	}
	if strings.Contains(attr, "|raw") {
		t.Fatalf("attribute |raw was not applied:\n%s", attr)
	}
	if strings.Contains(attr, "SafeURL") || strings.Contains(attr, "dreego.") {
		t.Fatalf("attribute |raw must bypass SafeURL, got:\n%s", attr)
	}
	if !strings.Contains(text, `WriteString(fmt.Sprintf("%v", webcal))`) {
		t.Fatalf("text |raw baseline changed, got:\n%s", text)
	}
	assertGoExpr(t, attr)
}

func assertGoExpr(t *testing.T, src string) {
	t.Helper()
	if _, err := parser.ParseFile(token.NewFileSet(), "x.go", "package p\n\nvar _ = "+src+"\n", 0); err != nil {
		t.Fatalf("generated expression is not valid Go: %v\n--- code ---\n%s", err, src)
	}
}

func assertGoStmt(t *testing.T, src string) {
	t.Helper()
	if _, err := parser.ParseFile(token.NewFileSet(), "x.go", "package p\n\nfunc f() {\n"+src+"\n}\n", 0); err != nil {
		t.Fatalf("generated code is not valid Go: %v\n--- code ---\n%s", err, src)
	}
}
