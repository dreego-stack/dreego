package tests

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestDeploymentCrossCompile(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, map[string]string{
		"www/routes/get.dreego": `<body><h1>hello</h1></body>`,
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
	t.Parallel()
	dir := t.TempDir()
	repoRoot, _ := dreegotest.RepoRoot()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	goMod := "module t\ngo 1.22\nrequire github.com/dreego-stack/dreego v0.0.0\nreplace github.com/dreego-stack/dreego => " + repoRoot + "\n"
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
	os.WriteFile(filepath.Join(dir, "main.go"), []byte(fmt.Sprintf("package main\nimport (\n\t\"t/www\"\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\nfunc main() { app := dreego.New(); if err := app.SetLogging(false); err != nil { panic(err) }; if err := www.Register(app); err != nil { panic(err) }; if err := app.Listen(\":%d\"); err != nil { panic(err) } }\n", port)), 0644)
	os.MkdirAll(filepath.Join(dir, "www", "routes"), 0755)
	os.WriteFile(filepath.Join(dir, "www", "routes", "get.dreego"), []byte("<body><h1>hello</h1></body>"), 0644)
	os.WriteFile(filepath.Join(dir, "www", "dreego.config.json"), []byte("{}"), 0644)

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

	waitReady(t, port)

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

// waitReady polls until a TCP connection to port succeeds, replacing a fixed
// sleep so shutdown tests start immediately once the server is listening.
func waitReady(t *testing.T, port int) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("server on port %d did not start in time", port)
}
