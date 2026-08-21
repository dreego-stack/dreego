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
		"www/routes/get.dreego": "<div>{#if true}<p>x</p></div>",
	})
	out, err := dreegotest.RunCLI(t, dir, "generate")
	if err == nil {
		t.Fatal("expected generate failure for unclosed {#if}")
	}
	if strings.Contains(out, "\x1b[") {
		t.Fatalf("error output contains ANSI color codes: %q", out)
	}
	if !strings.Contains(out, "www/routes/get.dreego") {
		t.Fatalf("error must lead with the source file, got: %q", out)
	}
	if !regexp.MustCompile(`www/routes/get\.dreego:\d+:\d+`).MatchString(out) {
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
		"www/routes/get.dreego": `<head><title>T</title></head>
<div><p>check me</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	src := filepath.Join(dir, "www/routes/get.dreego")
	if err := os.WriteFile(src, []byte(`<head><title>T</title></head>
<div><p>changed content</p></div>`), 0644); err != nil {
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
		"www/dreego.config.json": `{}`,
		"www/routes/get.dreego":  `<div><p>hi</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatal("expected failure without generated files")
	}
	if !strings.Contains(out, "Fix:") {
		t.Fatalf("no-gen diagnostic must name the next action, got: %q", out)
	}
}

func TestCLIBlueprintSemanticHTML(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	sub := filepath.Join(dir, "testapp")

	layout, err := os.ReadFile(filepath.Join(sub, "www/layouts/default.dreego"))
	if err != nil {
		t.Fatalf("read layout: %v", err)
	}
	lay := string(layout)
	for _, want := range []string{"<main", "{#slot}", "<nav", "skip"} {
		if !strings.Contains(lay, want) {
			t.Errorf("landing layout missing %q (semantic landmarks + skip link required)", want)
		}
	}

	route, err := os.ReadFile(filepath.Join(sub, "www/routes/get.dreego"))
	if err != nil {
		t.Fatalf("read route: %v", err)
	}
	if strings.Contains(string(route), "<img") && !strings.Contains(string(route), "alt=") {
		t.Error("landing route must give every <img> an alt attribute")
	}
	mainGo, err := os.ReadFile(filepath.Join(sub, "main.go"))
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	if !strings.Contains(string(mainGo), "script-src 'self' 'unsafe-inline' https://cdn.tailwindcss.com") {
		t.Error("landing CSP must allow its Tailwind development script")
	}

	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("generate in scaffold: %v\n%s", err, out)
	}
	if !dreegotest.BuildInDirOK(t, sub) {
		t.Fatal("accessible scaffold must still build")
	}
}

func TestCLIBlueprintDefaultRouteAccessible(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	route, err := os.ReadFile(filepath.Join(dir, "www/routes/get.dreego"))
	if err != nil {
		t.Fatalf("read route: %v", err)
	}
	if strings.Contains(string(route), "<img") && !strings.Contains(string(route), "alt=") {
		t.Error("default blueprint must give every <img> an alt attribute")
	}
}
