package tests

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

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
	if _, err := os.Stat(filepath.Join(dir, "dreego/routes/get.dreego")); err != nil {
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
	if !strings.Contains(string(mainGo), `"t/dreego/gen"`) {
		t.Fatalf("main.go does not import \"t/dreego/gen\": %s", mainGo)
	}
	if !strings.Contains(string(mainGo), "gen.Register(app)") {
		t.Fatalf("main.go does not call gen.Register(app): %s", mainGo)
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
		"dreego/routes/get.dreego": `<head><title>T</title></head>
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
	src := filepath.Join(dir, "dreego/routes/get.dreego")
	future := time.Now().Add(2 * time.Second)
	os.Chtimes(src, future, future)
	if _, err := dreegotest.RunCLI(t, dir, "generate", "--check"); err == nil {
		t.Fatal("expected stale but got up-to-date")
	}
}

func TestCLICheckNoGen(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	out, err := dreegotest.RunCLI(t, dir, "generate", "--check")
	if err == nil {
		t.Fatal("expected non-zero exit when no generated files exist")
	}
	if !strings.Contains(out, "no generated files found") {
		t.Fatalf("expected 'no generated files found' message, got: %s", out)
	}
}

func TestCLINew(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "new", "testapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	for _, f := range []string{
		"testapp/main.go",
		"testapp/go.mod",
		"testapp/dreego/routes/get.dreego",
		"testapp/dreego/layouts/default.dreego",
		"testapp/dreego/components/Hero.dreego",
		"testapp/dreego/components/FeatureCard.dreego",
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
	for _, f := range []string{"main.go", "dreego/layouts/default.dreego", "dreego/config.json"} {
		data, err := os.ReadFile(filepath.Join(sub, f))
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		if strings.Contains(string(data), "§$name$§") {
			t.Fatalf("unreplaced placeholder in %s", f)
		}
	}
	if !strings.Contains(string(readMust(t, filepath.Join(sub, "dreego/config.json"))), `"logging"`) {
		t.Fatal("config.json invalid/missing logging")
	}
	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("generate failed in scaffold: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(sub, "dreego/gen/routes.go")); err != nil {
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
	layouts, _ := filepath.Glob(filepath.Join(sub, "dreego/layouts/*.dreego"))
	if len(layouts) == 0 {
		t.Fatal("layouts/ directory exists but contains no .dreego layout file")
	}
	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	routes, _ := os.ReadFile(filepath.Join(sub, "dreego/gen/routes.go"))
	if !strings.Contains(string(routes), "<html>") {
		t.Fatal("layout exists but generated route does not produce a complete HTML document (no <html> found)")
	}
}

func TestCLIBuildTarget(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hello</p></div>`,
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
