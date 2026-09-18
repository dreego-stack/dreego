package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestBugLayoutComponentCompiles(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/components/Nav.dreego":  "DREEFILE component (label string)\n<body><nav>{{ label }}</nav></body>",
		"www/layouts/default.dreego": `<body><html><body><@Nav label="Home"/>{#slot}</body></html></body>`,
		"www/routes/+page.dreego":    `<body><p>Page</p></body>`,
	})
}

func TestBugLayoutComponentRenders(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Nav.dreego":  "DREEFILE component (label string)\n<body><nav>{{ label }}</nav></body>",
		"www/layouts/default.dreego": `<body><html><body><@Nav label="Home"/>{#slot}</body></html></body>`,
		"www/routes/+page.dreego":    `<body><p>Page</p></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<nav>Home</nav>") {
		t.Fatalf("layout component output missing in body: %s", body)
	}
	if !strings.Contains(body, "<p>Page</p>") {
		t.Fatalf("route content missing in body: %s", body)
	}
}

func TestBugLayoutComponentImportDirective(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Nav.dreego": "DREEFILE component (label string)\n<body><nav>{{ label }}</nav></body>",
		"www/layouts/default.dreego": `COMPONENT "www/components" IMPORT { Nav }
<body><html><body><@Nav label="Home"/>{#slot}</body></html></body>`,
		"www/routes/+page.dreego": `<body><p>Page</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err != nil {
		t.Fatalf("generate rejected layout header import: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestBugLayoutComponentErrorPointsAtLayout(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/layouts/default.dreego": `<body><html><body><@NoSuchComponent />{#slot}</body></html></body>`,
		"www/routes/+page.dreego":    `<body><p>Page</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatalf("generate accepted an unknown component in a layout:\n%s", out)
	}
	if !strings.Contains(out, "www/layouts/default.dreego") {
		t.Fatalf("diagnostic does not reference the layout file, got:\n%s", out)
	}
	if !strings.Contains(out, "unknown component NoSuchComponent") {
		t.Fatalf("unexpected diagnostic:\n%s", out)
	}
}

func TestBugComponentInRouteStillWorks(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Nav.dreego": "DREEFILE component (label string)\n<body><nav>{{ label }}</nav></body>",
		"www/routes/+page.dreego":   `<body><@Nav label="Home"/></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
	if !strings.Contains(body, "<nav>Home</nav>") {
		t.Fatalf("route component output missing in body: %s", body)
	}
}

func TestBugLayoutComponentGeneratedImport(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/components/Nav.dreego":  "DREEFILE component (label string)\n<body><nav>{{ label }}</nav></body>",
		"www/layouts/default.dreego": `<body><html><body><@Nav label="Home"/>{#slot}</body></html></body>`,
		"www/routes/+page.dreego":    `<body><p>Page</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	src, err := os.ReadFile(filepath.Join(dir, "www", "layouts", "dree.go"))
	if err != nil {
		t.Fatalf("read generated layout: %v", err)
	}
	dreegotest.MustContain(t, string(src), "components.Nav")
	dreegotest.MustContain(t, string(src), "/www/components\"")
}
