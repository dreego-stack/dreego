package codegen

import (
	"strings"
	"testing"

	"codeberg.org/dreego/dreego/pkg/ast"
)

func TestGenTemplateNodeExpressionEscapesHTML(t *testing.T) {
	n := ast.TemplateNode{
		Type:    ast.NodeExpression,
		Content: "name",
	}
	result := genTemplateNode(n, 1)

	if !strings.Contains(result, "html.EscapeString") {
		t.Errorf("expression node must use html.EscapeString, got:\n%s", result)
	}
}

func TestGenTemplateNodeIfRecursivelyEscapes(t *testing.T) {
	n := ast.TemplateNode{
		Type: ast.NodeIf,
		Cond: "show",
		Children: []ast.TemplateNode{
			{Type: ast.NodeExpression, Content: "name"},
		},
	}
	result := genTemplateNode(n, 1)

	if !strings.Contains(result, "html.EscapeString") {
		t.Errorf("{#if} child expression must use html.EscapeString, got:\n%s", result)
	}
}

func TestGenTemplateNodeEachRecursivelyEscapes(t *testing.T) {
	n := ast.TemplateNode{
		Type:  ast.NodeEach,
		Item:  "item",
		Items: "items",
		Children: []ast.TemplateNode{
			{Type: ast.NodeExpression, Content: "item.Name"},
		},
	}
	result := genTemplateNode(n, 1)

	if !strings.Contains(result, "html.EscapeString") {
		t.Errorf("{#each} child expression must use html.EscapeString, got:\n%s", result)
	}
}

func TestGenTemplateNodeTextNoEscape(t *testing.T) {
	n := ast.TemplateNode{
		Type:    ast.NodeText,
		Content: "hello",
	}
	result := genTemplateNode(n, 1)

	if strings.Contains(result, "html.EscapeString") {
		t.Errorf("static text must NOT use html.EscapeString, got:\n%s", result)
	}
}
