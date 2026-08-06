package core

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

// genComponentCall must resolve {expr} attribute values when a component calls
// another component, mirroring how the route path uses extractAttrValues.
// Currently it passes n.Attrs through raw, emitting invalid Go like
// Link(href={url} label="x").Render(ctx).
func TestGenComponentCallResolvesAttrExpression(t *testing.T) {
	n := TemplateNode{
		Type:      NodeComponentCall,
		Tag:       "Nav.Link",
		Attrs:     `href={url} label="x"`,
		SelfClose: true,
	}
	out, err := genComponentCall(n)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "{url}") {
		t.Errorf("genComponentCall left {url} literal, must resolve to expression url. got:\n%s", out)
	}
	if !strings.Contains(out, "Link(url, \"x\")") {
		t.Errorf("genComponentCall must pass resolved args Link(url, \"x\"), got:\n%s", out)
	}
}

// compTextWithAttrs is applied to every NodeText in a component body, including
// text inside <script>/<style> sections. The lexer treats those sections as raw
// text blocks where {…} is NOT a Go expression, so compTextWithAttrs must leave
// them untouched. Currently it replaces {x} inside "<script>const s = "{x}";"
// with html.EscapeString(fmt.Sprintf("%v", x)).
func TestCompTextWithAttrsLeavesScriptStyleBodiesUntouched(t *testing.T) {
	cases := []string{
		`<script>const s = "{x}";</script>`,
		`<style>.a { color: "{x}"; }</style>`,
	}
	for _, in := range cases {
		out := compTextWithAttrs(in)
		if !strings.Contains(out, "{x}") {
			t.Errorf("compTextWithAttrs must leave script/style body {x} literal (not resolve as Go expr). input=%q got=%q", in, out)
		}
	}
}

// GenerateComponent drives the stateful compGen over file.Template.Nodes. Unlike
// the wrapper-only unit tests, this exercises the real production path and must
// (1) keep a <script>/<style> body literal, (2) still resolve a quoted attribute
// placeholder like href="{url}", and (3) produce syntactically valid Go.
func TestGenerateComponentStatefulGenerator(t *testing.T) {
	src := `Component Card (x string, url string)

<div>
    <script>const s = "literal {x}";</script>
    <a href="{url}">go</a>
</div>
`
	_, _, body := ParseHeader(src)
	file := parseFile(t, body)
	file.Component = &ComponentDef{
		Name: "Card",
		Props: []Prop{
			{Name: "x", Type: "string"},
			{Name: "url", Type: "string"},
		},
	}

	out, err := GenerateComponent(file, scopeHashFor(src))
	if err != nil {
		t.Fatalf("GenerateComponent: %v", err)
	}

	if !strings.Contains(out, "literal {x}") {
		t.Errorf("script body {x} must stay literal, got resolved as expression:\n%s", out)
	}
	if strings.Contains(out, "fmt.Sprintf(\"%v\", x)") {
		t.Errorf("script body {x} must NOT become an expression, got:\n%s", out)
	}
	if !strings.Contains(out, `fmt.Sprintf("%v", url)`) {
		t.Errorf("href={url} must resolve to an expression, got literal:\n%s", out)
	}
	if strings.Contains(out, `href="{url}"`) {
		t.Errorf("href must not keep the literal placeholder {url}:\n%s", out)
	}

	if _, err := parser.ParseFile(token.NewFileSet(), "comp.go", "package comp\n"+out, 0); err != nil {
		t.Fatalf("generated component is not valid Go: %v\n--- body ---\n%s", err, out)
	}
}
