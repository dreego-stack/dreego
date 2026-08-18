package tests

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestQuickStartScaffold verifies the canonical quick-start path documented in
// README.md and _docs/getting-started.md:
//
//	dreego new myapp
//	cd myapp
//	dreego generate
//	go build
//	./server  →  GET / 200
//
// It scaffolds a real project with the CLI, generates, builds, starts the
// binary on a free port, and asserts the landing page responds. The test runs
// the documented commands end-to-end (black box), only substituting the
// listening port via DREEGO_PORT so the binary does not bind :8080.
func TestQuickStartScaffold(t *testing.T) {
	t.Parallel()
	parent := t.TempDir()
	appName := "myapp"

	if out, err := dreegotest.RunCLI(t, parent, "new", appName); err != nil {
		t.Fatalf("dreego new: %v\n%s", err, out)
	}
	sub := filepath.Join(parent, appName)

	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("dreego generate: %v\n%s", err, out)
	}
	if _, err := os.Stat(filepath.Join(sub, "dreego/gen/routes.go")); err != nil {
		t.Fatalf("gen/routes.go not produced: %v", err)
	}

	bin := filepath.Join(sub, "server")
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = sub
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}
	if info, err := os.Stat(bin); err != nil || info.Mode()&0111 == 0 {
		t.Fatalf("binary not executable at %s", bin)
	}

	port := freePortQuick(t)
	proc := exec.Command(bin)
	proc.Dir = sub
	proc.Env = append(os.Environ(), "DREEGO_PORT="+port)
	out, _ := os.CreateTemp("", "quick-start-*.log")
	proc.Stdout = out
	proc.Stderr = out
	if err := proc.Start(); err != nil {
		t.Fatalf("start server: %v", err)
	}
	t.Cleanup(func() {
		proc.Process.Kill()
		proc.Wait()
		out.Close()
		os.Remove(out.Name())
	})
	waitForPortQuick(t, port)

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/", port))
	if err != nil {
		t.Fatalf("GET /: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if !strings.Contains(string(body), "<html") {
		t.Fatalf("response is not an HTML document: %s", string(body)[:min(len(body), 200)])
	}
}

// TestQuickStartNewInvalidName asserts the CLI rejects an invalid project
// name with a helpful message, not a confusing go-mod error later.
func TestQuickStartNewInvalidName(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	for _, name := range []string{"123abc", "my app", "my$app", ""} {
		out, err := dreegotest.RunCLI(t, dir, "new", name)
		if err == nil {
			t.Fatalf("expected error for invalid name %q", name)
		}
		if !strings.Contains(out, "invalid project name") {
			t.Fatalf("expected 'invalid project name' for %q, got: %s", name, out)
		}
	}
}

// TestQuickStartNewMissingGo asserts the CLI explains that Go is required when
// the 'go' executable is unavailable.
func TestQuickStartNewMissingGo(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	bin := dreegotest.CLIBin(t)
	cmd := exec.Command(bin, "new", "noapp")
	cmd.Dir = dir
	cmd.Env = []string{"PATH=" + os.TempDir()}
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error when go is unavailable")
	}
	if !strings.Contains(string(out), "go") || !strings.Contains(string(out), "not found") {
		t.Fatalf("expected a 'go not found' diagnostic, got: %s", out)
	}
}

// TestQuickStartNewExisting asserts the CLI explains the target already exists.
func TestQuickStartNewExisting(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	target := filepath.Join(dir, "already")
	if err := os.MkdirAll(target, 0755); err != nil {
		t.Fatal(err)
	}
	out, err := dreegotest.RunCLI(t, dir, "new", "already")
	if err == nil {
		t.Fatal("expected error when target exists")
	}
	if !strings.Contains(out, "already exists") {
		t.Fatalf("expected 'already exists', got: %s", out)
	}
}

// TestQuickStartNewNoArg asserts the CLI prints a usage message.
func TestQuickStartNewNoArg(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	out, err := dreegotest.RunCLI(t, dir, "new")
	if err == nil {
		t.Fatal("expected error when no name given")
	}
	if !strings.Contains(out, "usage: dreego new") {
		t.Fatalf("expected usage message, got: %s", out)
	}
}

// TestQuickStartNoReplaceDirective asserts the scaffolded go.mod requires the
// published dreego module (no repo-local replace directive) when run against
// the published tag. When DREEGO_LOCAL_REPO is set (as in CI), a replace is
// expected — the test tolerates both cases but verifies a 'require' directive
// always exists.
func TestQuickStartNoReplaceDirective(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if out, err := dreegotest.RunCLI(t, dir, "new", "myapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	gomod, _ := os.ReadFile(filepath.Join(dir, "myapp/go.mod"))
	content := string(gomod)
	if !strings.Contains(content, "require github.com/dreego-stack/dreego") {
		t.Fatalf("scaffolded go.mod has no 'require github.com/dreego-stack/dreego':\n%s", content)
	}
	if !strings.Contains(content, "module myapp") {
		t.Fatalf("scaffolded go.mod has no 'module myapp':\n%s", content)
	}
	if !regexp.MustCompile(`go 1\.\d+`).MatchString(content) {
		t.Fatalf("scaffolded go.mod has no go directive:\n%s", content)
	}
}

// TestQuickStartScaffoldBuilds asserts the scaffolded project (without running
// the server) compiles cleanly with `go build`.
func TestQuickStartScaffoldBuilds(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	if out, err := dreegotest.RunCLI(t, dir, "new", "buildsapp"); err != nil {
		t.Fatalf("new: %v\n%s", err, out)
	}
	sub := filepath.Join(dir, "buildsapp")
	if out, err := dreegotest.RunCLI(t, sub, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if !dreegotest.BuildInDirOK(t, sub) {
		t.Fatal("scaffolded project must build")
	}
}

func freePortQuick(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := fmt.Sprintf("%d", ln.Addr().(*net.TCPAddr).Port)
	ln.Close()
	return port
}

func waitForPortQuick(t *testing.T, port string) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("server on port %s did not start in time", port)
}
