package transpiler

import (
	"strings"
	"testing"
)

func TestGenFormBindErrorDoesNotDiscardRenderError(t *testing.T) {
	file := &File{
		FormActions: []string{"Save"},
		Server: []ServerSection{
			{Code: "type SaveForm struct {\n\tName string\n}\nfunc Save(c dreego.Context, form SaveForm) error {\n\treturn nil\n}"},
		},
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<form g-action=\"Save\"><input name=\"name\"></form>"},
			},
		},
	}
	out, err := generateFormPostHandler(file, "renderIndex", "HandleIndexPost", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "html, _ :=") {
		t.Errorf("form bind error path must not discard the render error, got:\n%s", out)
	}
	if !strings.Contains(out, "html, err :=") {
		t.Errorf("form bind error path must capture the render error, got:\n%s", out)
	}
	if !strings.Contains(out, "if err != nil") {
		t.Errorf("form bind error path must check the render error, got:\n%s", out)
	}
}

func TestGenFormValidationErrorDoesNotDiscardRenderError(t *testing.T) {
	file := &File{
		FormActions: []string{"Save"},
		Server: []ServerSection{
			{Code: "type SaveForm struct {\n\tName string `validate:\"required\"`\n}\nfunc Save(c dreego.Context, form SaveForm) error {\n\treturn nil\n}"},
		},
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<form g-action=\"Save\"><input name=\"name\"></form>"},
			},
		},
	}
	out, err := generateFormPostHandler(file, "renderIndex", "HandleIndexPost", "/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "html, _ :=") {
		t.Errorf("form validation error path must not discard the render error, got:\n%s", out)
	}
	if !strings.Contains(out, "html, err :=") {
		t.Errorf("form validation error path must capture the render error, got:\n%s", out)
	}
}
