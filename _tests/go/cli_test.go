package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestCLIHelp(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "--help")
	if err != nil {
		t.Fatalf("--help: %v", err)
	}
	if !strings.Contains(out, "usage:") {
		t.Fatalf("--help did not show usage, got: %s", out)
	}
	out, err = dreegotest.RunCLI(t, dir, "-h")
	if err != nil {
		t.Fatalf("-h: %v", err)
	}
	if !strings.Contains(out, "usage:") {
		t.Fatalf("-h did not show usage, got: %s", out)
	}
}

func TestCLINoArgs(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, _ := dreegotest.RunCLI(t, dir)
	if !strings.Contains(out, "usage:") {
		t.Fatalf("no help shown, got: %s", out)
	}
}

func TestCLIUnknownCmd(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	_, err := dreegotest.RunCLI(t, dir, "bogus")
	if err == nil {
		t.Fatal("expected non-zero exit for unknown command")
	}
}

func TestCLIVersion(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	out = strings.TrimSpace(out)
	if out == "" {
		t.Fatal("version output is empty")
	}
	if out == "(devel)" {
		t.Fatal("version is the build-info placeholder (devel)")
	}
}

func TestCLIVersionDrift(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	out = strings.TrimSpace(out)
	tag := dreegotest.LatestTag(t)
	if tag == "dev" {
		t.Skip("no git tag present; version-drift check requires a tag")
	}
	if !strings.Contains(out, tag) {
		t.Fatalf("dreego version is %q, want it to contain %q (latest git tag)", out, tag)
	}
}

func TestCLIVersionFlag(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	flagOut, err := dreegotest.RunCLI(t, dir, "--version")
	if err != nil {
		t.Fatalf("--version: %v", err)
	}
	shortOut, err := dreegotest.RunCLI(t, dir, "-v")
	if err != nil {
		t.Fatalf("-v: %v", err)
	}
	subOut, err := dreegotest.RunCLI(t, dir, "version")
	if err != nil {
		t.Fatalf("version: %v", err)
	}
	flagOut, shortOut, subOut = strings.TrimSpace(flagOut), strings.TrimSpace(shortOut), strings.TrimSpace(subOut)
	if flagOut == "" || shortOut == "" || subOut == "" {
		t.Fatal("version output empty")
	}
	if flagOut == "(devel)" || shortOut == "(devel)" || subOut == "(devel)" {
		t.Fatal("version is build-info placeholder (devel)")
	}
	if flagOut != subOut {
		t.Fatalf("--version and version outputs differ: %q vs %q", flagOut, subOut)
	}
	if shortOut != subOut {
		t.Fatalf("-v and version outputs differ: %q vs %q", shortOut, subOut)
	}
}

func TestCLIInit(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(dir, "main.go")); err != nil {
		t.Fatalf("missing main.go: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "www/routes/get.dreego")); err != nil {
		t.Fatalf("missing get.dreego: %v", err)
	}
}

func TestCLIInitNoArg(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "init")
	if err == nil {
		t.Fatal("expected non-zero exit when no path given")
	}
	if !strings.Contains(out, "usage: dreego init") {
		t.Fatalf("expected usage message, got: %s", out)
	}
}

func TestCLIInitImport(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	mainGo, _ := os.ReadFile(filepath.Join(dir, "main.go"))
	if !strings.Contains(string(mainGo), `"t/www"`) {
		t.Fatalf("main.go does not import \"t/www\": %s", mainGo)
	}
	if !strings.Contains(string(mainGo), "www.Register(app)") {
		t.Fatalf("main.go does not call www.Register(app): %s", mainGo)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if !dreegotest.BuildInDirOK(t, dir) {
		t.Fatal("init project must build")
	}
}

func TestCLICheck(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err != nil {
		t.Fatalf("generate --check: %v", err)
	}
	if !strings.Contains(out, "up-to-date") {
		t.Fatalf("check failed, got: %s", out)
	}
}

func TestCLICheckStale(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<head><title>T</title></head>
<div><p>check me</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err != nil {
		t.Fatalf("initial check: %v\n%s", err, out)
	}
	if !strings.Contains(out, "up-to-date") {
		t.Fatalf("initial check failed, got: %s", out)
	}
	src := filepath.Join(dir, "www/routes/get.dreego")
	if err := os.WriteFile(src, []byte(`<head><title>T2</title></head>
<div><p>changed</p></div>`), 0644); err != nil {
		t.Fatalf("edit source: %v", err)
	}
	if _, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err == nil {
		t.Fatal("expected stale content to fail check, got up-to-date")
	}
}

func TestCLICheckNoGen(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json": `{}`,
		"www/routes/get.dreego":  `<div><p>hi</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatal("expected non-zero exit when no generated files exist")
	}
	if !strings.Contains(strings.ToLower(out), "missing") {
		t.Fatalf("expected 'missing' diagnostic, got: %s", out)
	}
}
