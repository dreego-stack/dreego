package transpiler

import (
	"strings"
	"testing"
)

func TestGenSelfCloseComponentPropagatesError(t *testing.T) {
	n := TemplateNode{
		Type:      NodeComponentCall,
		Tag:       "Card",
		Attrs:     `title="Hi"`,
		SelfClose: true,
	}
	gen := NewGenerator()
	gen.registerDef("Card", &ComponentDef{Name: "Card", Props: []Prop{{Name: "title", Type: "string"}}})
	out, err := genTemplateNode(gen, n, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "if err != nil") {
		t.Errorf("self-closing component call must not discard the Render error, got:\n%s", out)
	}
	if strings.Contains(out, ", _ :=") || strings.Contains(out, "h, _ :=") {
		t.Errorf("self-closing component call must not discard the Render error, got:\n%s", out)
	}
	if !strings.Contains(out, "panic(err)") {
		t.Errorf("self-closing component call must propagate the Render error, got:\n%s", out)
	}
}

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

func TestGenRouteHandlerDoesNotDiscloseInternalError(t *testing.T) {
	file := &File{
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>ok</p>"},
			},
		},
	}
	out, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "main", "home", "/home", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "err.Error()") {
		t.Errorf("generated route handler must not disclose internal errors, got:\n%s", out)
	}
	if !strings.Contains(out, `http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)`) {
		t.Errorf("generated route handler must write a generic 500, got:\n%s", out)
	}
}

func TestGenErrorHandlerDoesNotDiscloseInternalError(t *testing.T) {
	file := &File{
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>oops</p>"},
			},
		},
	}
	out, _, err := GenerateErrorHandler(NewGenerator(), file, "Site", 500, "/", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "err.Error()") {
		t.Errorf("generated error handler must not disclose internal errors, got:\n%s", out)
	}
	if !strings.Contains(out, `http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)`) {
		t.Errorf("generated error handler must write a generic 500, got:\n%s", out)
	}
}

func TestGenRouteHandlerChecksRenderError(t *testing.T) {
	file := &File{
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>ok</p>"},
			},
		},
	}
	out, _, err := GenerateMethodHandler(NewGenerator(), file, nil, "main", "home", "/home", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "html, err := renderHome(c)") {
		t.Errorf("generated route handler must capture the render error, got:\n%s", out)
	}
	if !strings.Contains(out, "if err != nil {") {
		t.Errorf("generated route handler must check the render error, got:\n%s", out)
	}
}
