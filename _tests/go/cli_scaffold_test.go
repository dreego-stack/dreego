package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestCLINew(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	for _, f := range []string{
		"testapp/main.go",
		"testapp/go.mod",
		"testapp/www/routes/get.dreego",
		"testapp/www/layouts/default.dreego",
		"testapp/www/components/Hero.dreego",
		"testapp/www/components/FeatureCard.dreego",
	} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Fatalf("missing %s: %v", f, err)
		}
	}
}

func TestCLINewNoArg(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "new")
	if err == nil {
		t.Fatal("expected non-zero exit when no name given")
	}
	if !strings.Contains(out, "usage: dreego new") {
		t.Fatalf("expected usage message, got: %s", out)
	}
}

func TestCLINewExists(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "myapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	out, err := dreegotest.RunCLI(t, dir, "new", "myapp")
	if err == nil {
		t.Fatal("expected non-zero exit when target exists")
	}
	if !strings.Contains(out, "already exists") {
		t.Fatalf("expected 'already exists' message, got: %s", out)
	}
}

func TestCLINewGoSum(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	gomod, _ := os.ReadFile(filepath.Join(dir, "testapp/go.mod"))
	if !strings.Contains(string(gomod), "replace github.com/dreego-stack/dreego => ") {
		t.Fatalf("go.mod has no replace directive for the local dreego: %s", gomod)
	}
	// generate + build inside the scaffold
	sub := filepath.Join(dir, "testapp")
	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if !dreegotest.BuildInDirOK(t, sub) {
		t.Fatal("scaffold must build")
	}
}

func TestCLINewGitignore(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	gi, err := os.ReadFile(filepath.Join(dir, "testapp/.gitignore"))
	if err != nil {
		t.Fatalf("missing .gitignore: %v", err)
	}
	lines := []string{}
	for _, line := range strings.Split(string(gi), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		lines = append(lines, line)
	}
	bad := regexp.MustCompile(`^(/dreego|dreego|dreego/)$|^/dreego\b|\bdreego\s*$`)
	for _, l := range lines {
		if bad.MatchString(l) {
			t.Fatalf(".gitignore contains top-level 'dreego' ignore pattern: %q", l)
		}
	}
}

func TestCLINewBlueprintValid(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	sub := filepath.Join(dir, "testapp")
	// no unreplaced placeholder
	for _, f := range []string{"main.go", "www/layouts/default.dreego", "www/dreego.config.json"} {
		data, err := os.ReadFile(filepath.Join(sub, f))
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		if strings.Contains(string(data), "§$name$§") {
			t.Fatalf("unreplaced placeholder in %s", f)
		}
	}
	if !strings.Contains(string(readMust(t, filepath.Join(sub, "www/dreego.config.json"))), `"logging"`) {
		t.Fatal("config.json invalid/missing logging")
	}
	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("generate failed in scaffold: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(sub, "www/dree.go")); err != nil {
		t.Fatalf("gen/routes.go not produced: %v", err)
	}
}

func TestCLINewLayoutExists(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	sub := filepath.Join(dir, "testapp")
	layouts, _ := filepath.Glob(filepath.Join(sub, "www/layouts/*.dreego"))
	if len(layouts) == 0 {
		t.Fatal("layouts/ directory exists but contains no .dreego layout file")
	}
	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	routes, _ := os.ReadFile(filepath.Join(sub, "www/layouts/dree.go"))
	if !strings.Contains(string(routes), "<html>") {
		t.Fatal("layout exists but generated layout does not produce a complete HTML document (no <html> found)")
	}
}

func TestCLIBuildTarget(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<div><p>hello</p></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "build", "--target", "linux/amd64"); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}
	name := filepath.Base(dir)
	bin := filepath.Join(dir, "build/bin", name+"-linux-amd64")
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("binary not found at %s: %v", bin, err)
	}
	info, _ := os.Stat(bin)
	if info.Mode()&0111 == 0 {
		t.Fatalf("binary not executable at %s", bin)
	}
}

func TestCLIDocs(t *testing.T) {
	t.Parallel()
	bin := dreegotest.CLIBin(t)
	repoRoot, _ := dreegotest.RepoRoot()

	out, _ := dreegotest.RunCLIIn(t, repoRoot, bin, "docs")
	if !strings.Contains(out, "Dreego Documentation") {
		t.Fatalf("docs did not show core index, got: %s", out)
	}
	out, _ = dreegotest.RunCLIIn(t, repoRoot, bin, "docs", "--json")
	if !strings.Contains(out, `"headings"`) {
		t.Fatalf("docs --json missing headings, got: %s", out)
	}
	out, _ = dreegotest.RunCLIIn(t, repoRoot, bin, "docs", "--list")
	if !strings.Contains(out, "github.com/dreego-stack/dreego") {
		t.Fatalf("--list missing core, got: %s", out)
	}

	// plugin docs via vendor/
	dir := dreegotest.ProjectDir(t, nil)
	vendorDocs := filepath.Join(dir, "vendor/github.com/dreego-stack/plugin-sse/_docs")
	os.MkdirAll(vendorDocs, 0755)
	gomod := "module example.com/myapp\n\ngo 1.22\n\nrequire (\n    github.com/dreego-stack/dreego v0.0.27\n    github.com/dreego-stack/plugin-sse v0.1.0\n)\n"
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(gomod), 0644)
	os.WriteFile(filepath.Join(vendorDocs, "index.md"), []byte("# Plugin SSE\n\nVendor-local plugin docs.\n"), 0644)

	out, _ = dreegotest.RunCLI(t, dir, "docs", "-p", "plugin-sse")
	if !strings.Contains(out, "Vendor-local plugin docs") {
		t.Fatalf("docs -p did not show vendor plugin docs, got: %s", out)
	}
	out, _ = dreegotest.RunCLI(t, dir, "docs", "--list")
	if !strings.Contains(out, "github.com/dreego-stack/plugin-sse") {
		t.Fatalf("--list missing plugin, got: %s", out)
	}
}

func TestCLIFmt(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	messy := `<head>
    <title>test</title>
</head>


<go>
    msg   :=    "hello"
</go>

<div>
  <p>{{  msg  }}</p>
    {#if  show}
        <span>visible</span>
    {/if}
</div>
`
	expected := `<head>
    <title>test</title>
</head>

<go>
    msg   :=    "hello"
</go>

<div>
  <p>{{ msg }}</p>
    {#if show}
        <span>visible</span>
    {/if}
</div>
`
	os.WriteFile(filepath.Join(dir, "messy.dreego"), []byte(messy), 0644)
	out, err := dreegotest.RunCLI(t, dir, "fmt", "--stdout", "messy.dreego")
	if err != nil {
		t.Fatalf("fmt: %v\n%s", err, out)
	}
	if out != expected {
		t.Fatalf("basic formatting mismatch\nwant: %q\ngot:  %q", expected, out)
	}
	os.WriteFile(filepath.Join(dir, "actual.dreego"), []byte(out), 0644)
	twice, err := dreegotest.RunCLI(t, dir, "fmt", "--stdout", "actual.dreego")
	if err != nil {
		t.Fatalf("fmt (idempotency): %v\n%s", err, twice)
	}
	if twice != out {
		t.Fatal("not idempotent")
	}
}

func readMust(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return data
}
