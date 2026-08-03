package core

import (
	"strings"
	"testing"
)

func TestGenTemplateNodeExpressionEscapesHTML(t *testing.T) {
	n := TemplateNode{
		Type:    NodeExpression,
		Content: "name",
	}
	result, err := genTemplateNode(n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "html.EscapeString") {
		t.Errorf("expression node must use html.EscapeString, got:\n%s", result)
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
	result, err := genTemplateNode(n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "html.EscapeString") {
		t.Errorf("{#if} child expression must use html.EscapeString, got:\n%s", result)
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
	result, err := genTemplateNode(n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(result, "html.EscapeString") {
		t.Errorf("{#each} child expression must use html.EscapeString, got:\n%s", result)
	}
}

func TestGenTemplateNodeTextNoEscape(t *testing.T) {
	n := TemplateNode{
		Type:    NodeText,
		Content: "hello",
	}
	result, err := genTemplateNode(n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if strings.Contains(result, "html.EscapeString") {
		t.Errorf("static text must NOT use html.EscapeString, got:\n%s", result)
	}
}

func TestGenTemplateNodeNestedIfInElseNotDropped(t *testing.T) {
	input := `{#if a}A{#else}{#if b}B{#else}C{/if}D{/if}`
	tokens, err := Lex(input)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := genTemplateNode(file.Template.Nodes[0], 0)
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
	input := `{#if a}A{#else}{#if b}B{#else}C{/if}D{/if}`
	tokens, err := Lex(input)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	out, err := genTemplateNodeComp(file.Template.Nodes[0])
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
	_, err := genTemplateNodeComp(TemplateNode{Type: TemplateNodeType(999)})
	if err == nil {
		t.Fatal("expected error for unsupported template node type in component codegen")
	}
}
