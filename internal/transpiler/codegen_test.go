package transpiler

import (
	"strings"
	"testing"
)

// GenerateErrorHandler for 404 must register a catch-all route and write the
// 404 status.
func TestGenerateErrorHandler404(t *testing.T) {
	file := &File{
		Body: &BodySection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<p>not found</p>"},
			},
		},
	}
	out, _, err := GenerateErrorHandler(NewGenerator(), file, "Site", 404, "/{path...}", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{"renderErrorSite404", "HandleErrorSite404", "w.WriteHeader(404)"} {
		if !strings.Contains(out, want) {
			t.Errorf("404 handler missing %q, got:\n%s", want, out)
		}
	}
}

// GenerateErrorHandler for 500 must register via SetErrorHandler.
func TestGenerateErrorHandler500(t *testing.T) {
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
	for _, want := range []string{"renderErrorSite500", "HandleErrorSite500"} {
		if !strings.Contains(out, want) {
			t.Errorf("500 handler missing %q, got:\n%s", want, out)
		}
	}
	if strings.Contains(out, "w.WriteHeader") {
		t.Errorf("500 handler must not write an explicit header (default 200 path), got:\n%s", out)
	}
}

// GenerateMethodHandler with a non-GET firstMethod must register that method and
// append a method suffix to the handler names.
func TestGenerateMethodHandlerNonGET(t *testing.T) {
	file := &File{
		Server: []ServerSection{
			{Code: "x := 1\n_ = x", Method: "POST"},
		},
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
	for _, want := range []string{"renderHomePOST", "HandleHomePOST"} {
		if !strings.Contains(out, want) {
			t.Errorf("non-GET handler missing %q, got:\n%s", want, out)
		}
	}
}
