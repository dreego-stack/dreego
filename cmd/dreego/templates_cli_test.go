package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/cmd/dreego/internal/templates"
)

func TestTemplatesList(t *testing.T) {
	metas := templates.List()
	if len(metas) == 0 {
		t.Fatal("List() returned no templates")
	}
	found := false
	for _, meta := range metas {
		if meta.Name == templates.DefaultName {
			found = true
		}
		if meta.Title == "" || meta.Description == "" {
			t.Errorf("template %q has empty title or description: %+v", meta.Name, meta)
		}
		if !templates.Exists(meta.Name) {
			t.Errorf("Exists(%q) = false for a listed template", meta.Name)
		}
	}
	if !found {
		t.Errorf("List() does not contain the default template %q", templates.DefaultName)
	}
}

func TestCommonMainEntrypointSSRHost(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("internal", "templates", "_common", "main.go.tmpl"))
	if err != nil {
		t.Fatalf("read _common/main.go.tmpl: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `"§$name$§/www"`) {
		t.Errorf("_common/main.go.tmpl must import \"§$name$§/www\" (module-qualified, placeholder), got:\n%s", content)
	}
	if !strings.Contains(content, "www.Register(app)") {
		t.Errorf("_common/main.go.tmpl must call www.Register(app), got:\n%s", content)
	}
	if !strings.Contains(content, `ssr "github.com/dreego-stack/dreego/adapter/ssr"`) {
		t.Errorf("_common/main.go.tmpl must import \"github.com/dreego-stack/dreego/adapter/ssr\", got:\n%s", content)
	}
	if !strings.Contains(content, "ssr.Listen(app, addr)") {
		t.Errorf("_common/main.go.tmpl must call ssr.Listen(app, addr), got:\n%s", content)
	}
	if !strings.Contains(content, `const port = "8080"`) {
		t.Errorf("_common/main.go.tmpl must declare the listening port as a constant, got:\n%s", content)
	}
}

func TestCommonTaskfileDefinesGenerateAndBuild(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("internal", "templates", "_common", "Taskfile.yml"))
	if err != nil {
		t.Fatalf("read _common/Taskfile.yml: %v", err)
	}
	content := string(data)
	for _, task := range []string{"  generate:", "  build:"} {
		if !strings.Contains(content, task) {
			t.Errorf("_common/Taskfile.yml must define %q, got:\n%s", strings.TrimSpace(task), content)
		}
	}
}
