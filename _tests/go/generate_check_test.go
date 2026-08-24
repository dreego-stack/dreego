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
		"www/routes/get.dreego": `<body><p>hi</p></body>`,
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
		"www/routes/get.dreego": `<body><p>hello</p></body>`,
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
		"www/routes/get.dreego": `<body><p>before</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	routesPath := filepath.Join(dir, "www", "routes", "dree.go")
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
		"www/routes/get.dreego": `<body><p>original</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(dir, "www", "routes", "get.dreego"), []byte(`<body><p>changed</p></body>`), 0644); err != nil {
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
		"www/routes/get.dreego": `<body><p>v1</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if err := os.WriteFile(filepath.Join(dir, "www", "routes", "get.dreego"), []byte(`<body><p>v2</p></body>`), 0644); err != nil {
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
		"www/routes/get.dreego": `<body><p>hi</p></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	// A dree.go in a directory without .dreego sources is not part of the
	// plan and must be reported as extra.
	compDir := filepath.Join(dir, "www", "components")
	if err := os.MkdirAll(compDir, 0755); err != nil {
		t.Fatalf("mkdir components: %v", err)
	}
	if err := os.WriteFile(filepath.Join(compDir, "dree.go"), []byte("package components\n"), 0644); err != nil {
		t.Fatalf("write extra dree.go: %v", err)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatalf("expected check to fail on extra dree.go file, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "components/dree.go") {
		t.Fatalf("expected diagnostic naming components/dree.go, got: %s", out)
	}
}

func TestGenerateCheckMissingGenFileFails(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego":       `<body><p>page</p></body>`,
		"www/components/Cmp.dreego":   "Component Card (title string)\n<body><article><h2>{{ title }}</h2></article></body>",
		"www/routes/about/get.dreego": `<body><@Card title="A"/></body>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	compPath := filepath.Join(dir, "www", "components", "dree.go")
	if err := os.Remove(compPath); err != nil {
		t.Fatalf("remove components.go: %v", err)
	}
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatalf("expected check to fail when components dree.go missing, got success: %s", out)
	}
	if !strings.Contains(strings.ToLower(out), "components/dree.go") {
		t.Fatalf("expected diagnostic naming components/dree.go, got: %s", out)
	}
}

func TestGenerateCheckRemovesStaleGenFiles(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego":      `<body><p>home</p></body>`,
		"www/components/Card.dreego": "Component Card ()\n<body><p>card</p></body>",
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	compPath := filepath.Join(dir, "www", "components", "dree.go")
	if _, err := os.Stat(compPath); err != nil {
		t.Fatalf("components.go not generated: %v", err)
	}
	if err := os.RemoveAll(filepath.Join(dir, "www", "components")); err != nil {
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
