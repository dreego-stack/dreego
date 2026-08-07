package core

import (
	"strings"
	"testing"
)

// GenerateErrorHandler must not emit the scope div when the template starts
// with a doctype/comment ("<!"): a scope div before the doctype puts the
// document into quirks mode. Error pages are self-contained documents, so the
// rendered body must start with the doctype.
func TestGenerateErrorHandlerScopeDivNotBeforeDoctype(t *testing.T) {
	file := &File{
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<!doctype html><html><head><title>Not Found</title></head><body><p>Not Found</p></body></html>"},
			},
		},
	}
	out, err := GenerateErrorHandler(file, "Site", 404, "/{p...}", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "data-scope") {
		t.Errorf("scope div must not be emitted when the template starts with a doctype, got:\n%s", out)
	}
}

// GenerateErrorHandler writes the head section before the scope div: the
// rendered body starts with the head content, never with the scope div.
// Control test — pins the ordering for templates without a doctype.
func TestGenerateErrorHandlerHeadBeforeScopeDiv(t *testing.T) {
	file := &File{
		Head: &HeadSection{Content: `<title>Not Found</title>`},
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<div><p>Not Found</p></div>"},
			},
		},
	}
	out, err := GenerateErrorHandler(file, "Site", 404, "/{p...}", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	headIdx := strings.Index(out, "Not Found")
	scopeIdx := strings.Index(out, "data-scope")
	if headIdx < 0 {
		t.Fatalf("head content missing, got:\n%s", out)
	}
	if scopeIdx < 0 {
		t.Fatalf("scope div missing for a template without doctype, got:\n%s", out)
	}
	if headIdx > scopeIdx {
		t.Errorf("head must be emitted before the scope div, got:\n%s", out)
	}
}

// Templates without a doctype keep the scope div: scoping is only dropped
// when the template starts with "<!". Control test — the fix must not remove
// scoping from plain error pages.
func TestGenerateErrorHandlerScopeDivKeptWithoutDoctype(t *testing.T) {
	file := &File{
		Template: &TemplateSection{
			Nodes: []TemplateNode{
				{Type: NodeText, Content: "<div><p>Not Found</p></div>"},
			},
		},
	}
	out, err := GenerateErrorHandler(file, "Site", 404, "/{p...}", "abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "data-scope") {
		t.Errorf("scope div must be kept for templates without doctype, got:\n%s", out)
	}
}
