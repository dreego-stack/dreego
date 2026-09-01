package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestGeneratedComponentRendersWithoutHTTP(t *testing.T) {
	t.Parallel()
	files := map[string]string{
		"www/components/Badge.dreego": `Component Badge (label string)
<body><span class="badge">{{ label }}</span></body>
<style>.badge { font-weight: bold; }</style>`,
		"www/routes/get.dreego": `<head><title>Shop</title></head>

<body><@Badge label={"<b>hot</b>"}/></body>`,
	}

	served := dreegotest.Serve(t, files)
	code, body := served.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	wantFragment := `<span class="badge">&lt;b&gt;hot&lt;/b&gt;</span>`
	if !strings.Contains(body, wantFragment) {
		t.Fatalf("served body missing escaped badge fragment %q, got: %s", wantFragment, body)
	}

	dir := dreegotest.BuildDir(t, files)
	testSource := `package components

import (
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestGeneratedBadgeRendersWithoutHTTP(t *testing.T) {
	result, err := dreego.Render(Badge("<b>hot</b>"))
	if err != nil {
		t.Fatal(err)
	}
	html := string(result.HTML)
	if !strings.Contains(html, "<span class=\"badge\">&lt;b&gt;hot&lt;/b&gt;</span>") {
		t.Fatalf("unexpected HTML: %s", html)
	}
	if !strings.Contains(html, "<style>") {
		t.Fatalf("component style missing from HTML: %s", html)
	}
}
`
	testFile := filepath.Join(dir, "www", "components", "render_non_http_test.go")
	if err := os.WriteFile(testFile, []byte(testSource), 0644); err != nil {
		t.Fatal(err)
	}
	command := exec.Command("go", "test", "./www/components", "-run", "TestGeneratedBadgeRendersWithoutHTTP")
	command.Dir = dir
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated non-HTTP render test failed: %v\n%s", err, output)
	}

	pageTest := `package routes_test

import (
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
	"t/www/routes"
)

func TestGeneratedPageRendersWithoutHTTP(t *testing.T) {
	result, err := dreego.Render(routes.PageIndex())
	if err != nil {
		t.Fatal(err)
	}
	html := string(result.HTML)
	if !strings.Contains(html, "<title>Shop</title>") {
		t.Fatalf("page head missing: %s", html)
	}
	if !strings.Contains(html, "&lt;b&gt;hot&lt;/b&gt;") {
		t.Fatalf("component output missing: %s", html)
	}
}
`
	pageTestFile := filepath.Join(dir, "www", "routes", "render_non_http_test.go")
	if err := os.WriteFile(pageTestFile, []byte(pageTest), 0644); err != nil {
		t.Fatal(err)
	}
	command = exec.Command("go", "test", "./www/routes", "-run", "TestGeneratedPageRendersWithoutHTTP")
	command.Dir = dir
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("generated page non-HTTP render test failed: %v\n%s", err, output)
	}
}
