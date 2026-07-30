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
