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
	result := genTemplateNode(n, 1)

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
	result := genTemplateNode(n, 1)

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
	result := genTemplateNode(n, 1)

	if !strings.Contains(result, "html.EscapeString") {
		t.Errorf("{#each} child expression must use html.EscapeString, got:\n%s", result)
	}
}

func TestGenTemplateNodeTextNoEscape(t *testing.T) {
	n := TemplateNode{
		Type:    NodeText,
		Content: "hello",
	}
	result := genTemplateNode(n, 1)

	if strings.Contains(result, "html.EscapeString") {
		t.Errorf("static text must NOT use html.EscapeString, got:\n%s", result)
	}
}
