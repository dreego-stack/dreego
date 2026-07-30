package core

import (
	"testing"
)

func TestScanFormActionsSingle(t *testing.T) {
	nodes := []TemplateNode{
		{
			Type:    NodeText,
			Content: `<form g-action="save">`,
		},
	}
	actions := scanFormActions(nodes)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
	if actions[0] != "save" {
		t.Errorf("expected 'save', got '%s'", actions[0])
	}
}

func TestScanFormActionsMultiple(t *testing.T) {
	nodes := []TemplateNode{
		{
			Type:    NodeText,
			Content: `<form g-action="save"><input g-action="validate">`,
		},
	}
	actions := scanFormActions(nodes)
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions, got %d", len(actions))
	}
	if actions[0] != "save" {
		t.Errorf("expected 'save', got '%s'", actions[0])
	}
	if actions[1] != "validate" {
		t.Errorf("expected 'validate', got '%s'", actions[1])
	}
}

func TestScanFormActionsDeduplicate(t *testing.T) {
	nodes := []TemplateNode{
		{
			Type:    NodeText,
			Content: `<form g-action="save"><input g-action="save">`,
		},
	}
	actions := scanFormActions(nodes)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action after dedup, got %d", len(actions))
	}
	if actions[0] != "save" {
		t.Errorf("expected 'save', got '%s'", actions[0])
	}
}
