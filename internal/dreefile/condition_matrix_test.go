package dreefile

import (
	"strings"
	"testing"
)

func condNode(cond string) TemplateNode {
	return TemplateNode{
		Type:     NodeIf,
		Cond:     cond,
		Children: []TemplateNode{{Type: NodeText, Content: "yes"}},
	}
}

var conditionMatrix = []struct {
	name string
	cond string
	want string
}{
	{"bool variable", "show", "if dreego.Truthy(show) {"},
	{"negation", "!show", "if dreego.Truthy(!show) {"},
	{"comparison", "score >= 80", "if dreego.Truthy(score >= 80) {"},
	{"string method", `c.Errors("email")`, `if dreego.Truthy(c.Errors("email")) {`},
	{"string variable", "title", "if dreego.Truthy(title) {"},
	{"slice length", "len(items) > 0", "if dreego.Truthy(len(items) > 0) {"},
	{"init statement", `v := c.Get("x"); v != ""`, `if v := c.Get("x"); dreego.Truthy(v != "") {`},
	{"semicolon in string", `label == "a;b"`, `if dreego.Truthy(label == "a;b") {`},
}

// Truthiness matrix: every {#if} condition must compile to a Go boolean,
// including strings, numbers, and slices. The regression was
// "{#if c.Errors('x')}" emitting "if c.Errors('x') {" (non-boolean condition).
func TestConditionCodeTruthinessMatrix(t *testing.T) {
	for _, c := range conditionMatrix {
		t.Run(c.name+"/route", func(t *testing.T) {
			out, err := genTemplateNode(NewGenerator(), condNode(c.cond), 1)
			if err != nil {
				t.Fatalf("route codegen: %v", err)
			}
			if !strings.Contains(out, c.want) {
				t.Errorf("route if missing %q, got:\n%s", c.want, out)
			}
			assertGoStmt(t, out)
		})
		t.Run(c.name+"/component", func(t *testing.T) {
			out, err := genTemplateNodeComp(NewGenerator(), condNode(c.cond))
			if err != nil {
				t.Fatalf("component codegen: %v", err)
			}
			if !strings.Contains(out, c.want) {
				t.Errorf("component if missing %q, got:\n%s", c.want, out)
			}
			assertGoStmt(t, out)
		})
	}
}

// A nested {#if} inside {#else} must wrap its own condition too.
func TestConditionCodeElseIfWrapped(t *testing.T) {
	n := TemplateNode{
		Type:     NodeIf,
		Cond:     "a",
		Children: []TemplateNode{{Type: NodeText, Content: "A"}},
		ElseChildren: []TemplateNode{
			{Type: NodeIf, Cond: `c.Errors("x")`, Children: []TemplateNode{{Type: NodeText, Content: "B"}}},
		},
	}
	out, err := genTemplateNode(NewGenerator(), n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, `} else if dreego.Truthy(c.Errors("x")) {`) {
		t.Errorf("else-if string condition must be wrapped, got:\n%s", out)
	}
	assertGoStmt(t, out)
}
