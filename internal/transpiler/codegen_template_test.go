package transpiler

import (
	"strings"
	"testing"
)

func TestGenTemplateNodeExpressionEscapesHTML(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "name",
	}
	result, err := genTemplateNode(NewGenerator(), n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "dreego.SafeText") {
		t.Errorf("expression node must use dreego.SafeText, got:\n%s", result)
	}
}

func TestGenTemplateNodeIfRecursivelyEscapes(t *testing.T) {
	n := TemplateNode{
		Type: NodeIf,
		Cond: "show",
		Children: []TemplateNode{
			{Type: NodeExpression, Content: "name"},
		},
	}
	result, err := genTemplateNode(NewGenerator(), n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "dreego.SafeText") {
		t.Errorf("{#if} child expression must use dreego.SafeText, got:\n%s", result)
	}
}

func TestGenTemplateNodeEachRecursivelyEscapes(t *testing.T) {
	n := TemplateNode{
		Type:  NodeEach,
		Item:  "item",
		Items: "items",
		Children: []TemplateNode{
			{Type: NodeExpression, Content: "item.Name"},
		},
	}
	result, err := genTemplateNode(NewGenerator(), n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "dreego.SafeText") {
		t.Errorf("{#if} child expression must use dreego.SafeText, got:\n%s", result)
	}
}

func TestGenTemplateNodeTextNoEscape(t *testing.T) {
	n := TemplateNode{
		Type:    NodeText,
		Content: "hello",
	}
	result, err := genTemplateNode(NewGenerator(), n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "dreego.SafeText") {
		t.Errorf("static text must NOT use dreego.SafeText, got:\n%s", result)
	}
}

func TestGenTemplateNodeNestedIfInElseNotDropped(t *testing.T) {
	input := `<body>{#if a}A{#else}{#if b}B{#else}C{/if}D{/if}</body>`
	tokens, err := Lex(input)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := genTemplateNode(NewGenerator(), file.Body.Nodes[0], 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Fatal("nested {#if} inside {#else} silently dropped: generated code is empty")
	}
	for _, want := range []string{"if a", "if b", "`A`", "`B`", "`C`", "`D`"} {
		if !strings.Contains(out, want) {
			t.Errorf("generated code missing %q, got:\n%s", want, out)
		}
	}
}

func TestGenTemplateNodeCompNestedIfInElseNotDropped(t *testing.T) {
	input := `<body>{#if a}A{#else}{#if b}B{#else}C{/if}D{/if}</body>`
	tokens, err := Lex(input)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := genTemplateNodeComp(NewGenerator(), file.Body.Nodes[0])
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == "" {
		t.Fatal("nested {#if} inside {#else} in component silently dropped")
	}
	for _, want := range []string{"if a", "if b", "`A`", "`B`", "`C`", "`D`"} {
		if !strings.Contains(out, want) {
			t.Errorf("generated component code missing %q, got:\n%s", want, out)
		}
	}
}

func TestGenTemplateNodeCompReturnsErrorForUnsupportedNode(t *testing.T) {
	_, err := genTemplateNodeComp(NewGenerator(), TemplateNode{Type: TemplateNodeType(999)})
	if err == nil {
		t.Fatal("expected error for unsupported template node type in component codegen")
	}
}

// An attribute expression in a component body must be escaped (XSS-safe).
// `<a href="{{ url }}">` must generate a safe-value call for the url value so a
// quote in the value cannot break out of the attribute.
func TestGenTemplateNodeCompAttrExpressionEscapesValue(t *testing.T) {
	body := `<body><a href="{{ url }}">{{ label }}</a></body>`
	tokens, err := Lex(body)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	var out strings.Builder
	for _, n := range file.Body.Nodes {
		code, err := genTemplateNodeComp(NewGenerator(), n)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		out.WriteString(code)
	}
	got := out.String()
	if strings.Contains(got, "{{ url }}") {
		t.Errorf("attribute expression {{ url }} left literal, must be resolved. got:\n%s", got)
	}
	if !strings.Contains(got, "dreego.SafeURL") {
		t.Errorf("href attribute expression must be scheme-validated (XSS-safe), got:\n%s", got)
	}
}

func TestGenTemplateNodeRouteAttrExpressionEscapesValue(t *testing.T) {
	n := TemplateNode{Type: NodeText, Content: `<a href="{{ url }}">link</a>`}
	got, err := genTemplateNode(NewGenerator(), n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, `fmt.Sprintf("%v", url)`) {
		t.Fatalf("route attribute expression was not generated: %s", got)
	}
	if !strings.Contains(got, "dreego.SafeURL") {
		t.Fatalf("route href attribute expression must be scheme-validated: %s", got)
	}
}
