package transpiler

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

func TestScanFormActionsNested(t *testing.T) {
	nodes := []TemplateNode{
		{
			Type: NodeIf,
			Cond: "show",
			Children: []TemplateNode{
				{
					Type:    NodeText,
					Content: `<form g-action="save">`,
				},
			},
		},
	}
	actions := scanFormActions(nodes)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action in nested {#if}, got %d", len(actions))
	}
	if actions[0] != "save" {
		t.Errorf("expected 'save', got '%s'", actions[0])
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

func TestScanFormActionsNoAction(t *testing.T) {
	nodes := []TemplateNode{
		{
			Type:    NodeText,
			Content: `<form method="post">`,
		},
	}
	actions := scanFormActions(nodes)
	if len(actions) != 0 {
		t.Fatalf("expected 0 actions, got %d", len(actions))
	}
}

func TestFindFormStruct(t *testing.T) {
	serverSections := []ServerSection{
		{Code: "func save(w http.ResponseWriter, f MyForm) {}"},
	}
	result := findFormStruct(serverSections, "save")
	if result != "MyForm" {
		t.Errorf("expected 'MyForm', got '%s'", result)
	}
}

func TestFindFormStructMissing(t *testing.T) {
	serverSections := []ServerSection{
		{Code: "func save(w http.ResponseWriter) {}"},
	}
	result := findFormStruct(serverSections, "save")
	if result != "" {
		t.Errorf("expected '', got '%s'", result)
	}
}

func TestFindFormHandler(t *testing.T) {
	serverSections := []ServerSection{
		{Code: "func save(w http.ResponseWriter, f MyForm) {}"},
	}
	if !findFormHandler(serverSections, "save") {
		t.Error("expected handler to be found")
	}
}

func TestFindFormHandlerMissing(t *testing.T) {
	serverSections := []ServerSection{
		{Code: "func other(w http.ResponseWriter) {}"},
	}
	if findFormHandler(serverSections, "save") {
		t.Error("expected handler to not be found")
	}
}

func TestHasValidateTagScoped(t *testing.T) {
	serverSections := []ServerSection{
		{Code: "type MyForm struct {\n\tName string `validate:\"required\"`\n}"},
	}
	if !hasValidateTag(serverSections, "MyForm") {
		t.Error("expected validate tag to be found")
	}
}

func TestHasFormTagScoped(t *testing.T) {
	serverSections := []ServerSection{
		{Code: "type MyForm struct {\n\tName string `form:\"name\"`\n}"},
	}
	if !hasFormTag(serverSections, "MyForm") {
		t.Error("expected form tag to be found")
	}
}
