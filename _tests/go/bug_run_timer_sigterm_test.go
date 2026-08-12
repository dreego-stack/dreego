package tests

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestBugRunTimerSigterm verifies that the server started via `dreego run`
// shuts down gracefully on SIGTERM instead of dying without cleanup (bug B20).
// The 1s auto-stop timer itself is covered deterministically by
// TestScheduleStopSendsSIGTERM in cli/dreego/main_test.go; a fixed -t 1 second
// would race server startup under parallel load, so this test waits for the
// server to be ready and then signals it directly.
func TestBugRunTimerSigterm(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	goMod := fmt.Sprintf("module t\ngo 1.22\nrequire github.com/dreego-stack/dreego v0.0.0\nreplace github.com/dreego-stack/dreego => %s\n", repoRoot)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}
	mainGo := fmt.Sprintf(`package main
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	_ "t/dreego/gen"
	dreego "github.com/dreego-stack/dreego/core"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-done
		fmt.Println("SIGTERM received")
		return
	}()
	dreego.SetLogging(false)
	dreego.Listen(":%d")
}
`, port)
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "dreego", "routes"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dreego", "routes", "get.dreego"), []byte("<div>hello</div>"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v", err)
	}

	srv := filepath.Join(dir, "srv")
	build := exec.Command("go", "build", "-o", srv, ".")
	build.Dir = dir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	outfile := filepath.Join(dir, "run.out")
	f, err := os.Create(outfile)
	if err != nil {
		t.Fatal(err)
	}
	run := exec.Command(srv)
	run.Dir = dir
	run.Stdout = f
	run.Stderr = f
	if err := run.Start(); err != nil {
		t.Fatalf("run start: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- run.Wait() }()

	waitReady(t, port)

	if err := run.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("signal: %v", err)
	}

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		run.Process.Kill()
		t.Fatal("server did not exit after signal")
	}
	f.Close()

	out, _ := os.ReadFile(outfile)
	if !strings.Contains(string(out), "SIGTERM received") {
		t.Fatalf("server did not receive SIGTERM (B20)\n%s", out)
	}
}
