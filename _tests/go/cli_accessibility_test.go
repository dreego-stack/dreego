package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestCLIOutputNoColor(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	for _, args := range [][]string{{"--help"}, {"version"}, {"generate", "--check"}} {
		out, _ := dreegotest.RunCLI(t, dir, args...)
		if strings.Contains(out, "\x1b[") {
			t.Fatalf("CLI output for %v contains ANSI color codes; meaning must not rely on color: %q", args, out)
		}
	}
}

func TestCLIHelpLinearScreenReader(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "--help")
	if err != nil {
		t.Fatalf("--help: %v", err)
	}
	if !strings.HasPrefix(out, "dreego — Go Web Framework CLI") {
		t.Fatalf("help must start with the program name for screen readers, got: %q", out)
	}
	if !strings.Contains(out, "usage: dreego <command> [flags]") {
		t.Fatalf("help must contain the usage line, got: %q", out)
	}
}

func TestCLIErrorLeadsWithFilePositionCauseAction(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": "<body>{#if true}<p>x</p></body>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatal("expected generate failure for unclosed {#if}")
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("error output contains ANSI color codes: %q", out)
	}
	if !strings.Contains(out, "www/routes/+page.dreego") {
		t.Fatalf("error must lead with the source file, got: %q", out)
	}
	if !regexp.MustCompile(`www/routes/\+page\.dreego:\d+:\d+`).MatchString(out) {
		t.Fatalf("error must contain file:line:col, got: %q", out)
	}
	if !strings.Contains(out, "unclosed {#if") {
		t.Fatalf("error must state the cause, got: %q", out)
	}
	if !strings.Contains(out, "Fix:") {
		t.Fatalf("error must end with a practical next action (Fix:), got: %q", out)
	}
}

func TestCLICheckStaleActionable(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/+page.dreego": `<head><title>T</title></head>
<body><p>check me</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	src := filepath.Join(dir, "www/routes/+page.dreego")
	if err := os.WriteFile(src, []byte(`<head><title>T</title></head>
<body><p>changed content</p></body>`), 0644); err != nil {
		t.Fatalf("write source: %v", err)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatal("expected stale to fail --check")
	}
	if !strings.Contains(out, "stale:") && !strings.Contains(out, "missing:") && !strings.Contains(out, "out of date") {
		t.Fatalf("expected a stale diagnostic, got: %q", out)
	}
	if !strings.Contains(out, "Fix:") {
		t.Fatalf("stale diagnostic must name the next action, got: %q", out)
	}
}

func TestCLICheckNoGenActionable(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":  `{}`,
		"www/routes/+page.dreego": `<body><p>hi</p></body>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatal("expected failure without generated files")
	}
	if !strings.Contains(out, "Fix:") {
		t.Fatalf("no-gen diagnostic must name the next action, got: %q", out)
	}
}

func TestCLITemplateSemanticHTML(t *testing.T) {
	t.Parallel()
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}

	layout, err := os.ReadFile(filepath.Join(repoRoot, "cmd", "dreego", "internal", "templates", "web-minimal", "www", "layouts", "default.dreego"))
	if err != nil {
		t.Fatalf("read template layout: %v", err)
	}
	lay := string(layout)
	for _, want := range []string{"<main", "{#slot}", "skip to content", `href="#main"`, `lang="en"`} {
		if !strings.Contains(lay, want) {
			t.Errorf("template layout missing %q (semantic landmarks + skip link required)", want)
		}
	}

	route, err := os.ReadFile(filepath.Join(repoRoot, "cmd", "dreego", "internal", "templates", "web-minimal", "www", "routes", "+page.dreego"))
	if err != nil {
		t.Fatalf("read template route: %v", err)
	}
	if strings.Contains(string(route), "<img") && !strings.Contains(string(route), "alt=") {
		t.Error("template route must give every <img> an alt attribute")
	}

	mainTmpl, err := os.ReadFile(filepath.Join(repoRoot, "cmd", "dreego", "internal", "templates", "_common", "main.go.tmpl"))
	if err != nil {
		t.Fatalf("read _common/main.go.tmpl: %v", err)
	}
	if strings.Contains(string(mainTmpl), "cdn.tailwindcss.com") {
		t.Error("web-minimal must not depend on a Tailwind CDN script")
	}

	dir := dreegotest.NewProject(t, "app", "")
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate in scaffold: %v\n%s", err, out)
	}
	if !dreegotest.BuildInDirOK(t, dir) {
		t.Fatal("accessible scaffold must still build")
	}
}

func TestCLITemplateRouteAccessible(t *testing.T) {
	t.Parallel()
	dir := dreegotest.NewProject(t, "app", "")
	route, err := os.ReadFile(filepath.Join(dir, "www/routes/+page.dreego"))
	if err != nil {
		t.Fatalf("read route: %v", err)
	}
	if strings.Contains(string(route), "<img") && !strings.Contains(string(route), "alt=") {
		t.Error("scaffolded route must give every <img> an alt attribute")
	}
	layout, err := os.ReadFile(filepath.Join(dir, "www/layouts/default.dreego"))
	if err != nil {
		t.Fatalf("read scaffolded layout: %v", err)
	}
	lay := string(layout)
	for _, want := range []string{`<main id="main">`, "skip to content", `lang="en"`} {
		if !strings.Contains(lay, want) {
			t.Errorf("scaffolded layout missing %q (skip link + main landmark + lang required)", want)
		}
	}
}
