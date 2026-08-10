package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDeploymentCrossCompile(t *testing.T) {
	dir := dreegotest.ProjectDir(t, map[string]string{
		"dreego/routes/get.dreego": `<div><h1>hello</h1></div>`,
	})
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "build", "--target", "linux/amd64"); err != nil {
		t.Fatalf("build --target: %v\n%s", err, out)
	}
	matches, _ := filepath.Glob(filepath.Join(dir, "build/bin/*-linux-amd64"))
	if len(matches) == 0 {
		t.Fatal("binary not found")
	}
	info, err := os.Stat(matches[0])
	if err != nil {
		t.Fatalf("stat binary: %v", err)
	}
	if info.Mode()&0111 == 0 {
		t.Fatalf("binary not executable")
	}
}

func TestDeploymentGracefulShutdown(t *testing.T) {
	dir := t.TempDir()
	repoRoot, _ := dreegotest.RepoRoot()
	goMod := "module t\ngo 1.22\nrequire github.com/dreego-stack/dreego v0.0.0\nreplace github.com/dreego-stack/dreego => " + repoRoot + "\n"
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\nimport (\n\t_ \"t/dreego/gen\"\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\nfunc main() { dreego.SetLogging(false); dreego.Listen(\":0\") }\n"), 0644)
	os.MkdirAll(filepath.Join(dir, "dreego", "routes"), 0755)
	os.WriteFile(filepath.Join(dir, "dreego", "routes", "get.dreego"), []byte("<div><h1>hello</h1></div>"), 0644)

	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}

	srv := filepath.Join(dir, "srv")
	build := exec.Command("go", "build", "-o", srv, ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	proc := exec.Command(srv)
	proc.Dir = dir
	if err := proc.Start(); err != nil {
		t.Fatalf("start: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- proc.Wait() }()

	time.Sleep(2 * time.Second)

	if err := proc.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("signal: %v", err)
	}

	select {
	case <-done:
		// http.Server.Shutdown returns nil on graceful close; Wait() succeeds.
	case <-time.After(10 * time.Second):
		proc.Process.Kill()
		t.Fatal("server did not exit after signal")
	}
}
