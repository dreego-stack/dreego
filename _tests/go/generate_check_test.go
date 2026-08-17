package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestGenerateCheckNoGenFileFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hi</p></div>`,
	})
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatalf("expected non-zero exit when no gen files exist, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "missing") {
		t.Fatalf("expected missing-file diagnostic, got: %s", out)
	}
}

func TestGenerateCheckUpToDatePasses(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hello</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err != nil {
		t.Fatalf("check should pass on up-to-date tree: %v\n%s", err, out)
	}
	if !strings.Contains(out, "up-to-date") {
		t.Fatalf("expected up-to-date message, got: %s", out)
	}
}

func TestGenerateCheckDoesNotModifyWorkingTree(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>before</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	routesPath := filepath.Join(dir, "dreego", "gen", "routes.go")
	before, err := os.ReadFile(routesPath)
	if err != nil {
		t.Fatalf("read routes.go: %v", err)
	}
	beforeInfo, err := os.Stat(routesPath)
	if err != nil {
		t.Fatalf("stat routes.go: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err != nil {
		t.Fatalf("check: %v\n%s", err, out)
	}
	after, err := os.ReadFile(routesPath)
	if err != nil {
		t.Fatalf("read routes.go after check: %v", err)
	}
	afterInfo, err := os.Stat(routesPath)
	if err != nil {
		t.Fatalf("stat routes.go after check: %v", err)
	}
	if string(before) != string(after) {
		t.Fatalf("check modified routes.go content")
	}
	if !beforeInfo.ModTime().Equal(afterInfo.ModTime()) {
		t.Fatalf("check modified routes.go mtime: before=%v after=%v", beforeInfo.ModTime(), afterInfo.ModTime())
	}
}

func TestGenerateCheckStaleContentFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>original</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(dir, "dreego", "routes", "get.dreego"), []byte(`<div><p>changed</p></div>`), 0644); err != nil {
		t.Fatalf("edit source: %v", err)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatalf("expected check to fail on stale content, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "stale") && !strings.Contains(strings.ToLower(out), "differ") && !strings.Contains(strings.ToLower(out), "routes.go") {
		t.Fatalf("expected stale/diff diagnostic naming routes.go, got: %s", out)
	}
}

func TestGenerateCheckManipulatedTimestampsNoFalsePass(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>v1</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(dir, "dreego", "routes", "get.dreego"), []byte(`<div><p>v2</p></div>`), 0644); err != nil {
		t.Fatalf("edit source: %v", err)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatalf("check must not false-pass on manipulated timestamps: %s", out)
	}
}

func TestGenerateCheckExtraGenFileFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hi</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(dir, "dreego", "gen", "stale.go"), []byte("package gen\n"), 0644); err != nil {
		t.Fatalf("write stale.go: %v", err)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatalf("expected check to fail on extra gen file, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "stale.go") {
		t.Fatalf("expected diagnostic naming stale.go, got: %s", out)
	}
}

func TestGenerateCheckMissingGenFileFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego":       `<div><p>page</p></div>`,
		"dreego/components/Cmp.dreego":   "Component Card (title string)\n<div><article><h2>{{ title }}</h2></article></div>",
		"dreego/routes/about/get.dreego": `<div><@Card title="A"/></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	compPath := filepath.Join(dir, "dreego", "gen", "components.go")
	if err := os.Remove(compPath); err != nil {
		t.Fatalf("remove components.go: %v", err)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatalf("expected check to fail when components.go missing, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "components.go") {
		t.Fatalf("expected diagnostic naming components.go, got: %s", out)
	}
}

func TestGenerateCheckRemovesStaleGenFiles(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego":      `<div><p>home</p></div>`,
		"dreego/components/Card.dreego": "Component Card ()\n<div><p>card</p></div>",
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	compPath := filepath.Join(dir, "dreego", "gen", "components.go")
	if _, err := os.Stat(compPath); err != nil {
		t.Fatalf("components.go not generated: %v", err)
	}
	if err := os.RemoveAll(filepath.Join(dir, "dreego", "components")); err != nil {
		t.Fatalf("remove components dir: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate after removing last component: %v\n%s", err, out)
	}
	if _, err := os.Stat(compPath); !os.IsNotExist(err) {
		t.Fatalf("stale components.go still present after generate: %v", err)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err != nil {
		t.Fatalf("check should pass after generate cleanup: %v\n%s", err, out)
	}
}
