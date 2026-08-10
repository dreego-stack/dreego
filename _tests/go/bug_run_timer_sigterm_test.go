package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dreego-stack/dreego/dreegotest"
)

// TestBugRunTimerSigterm verifies that `dreego run -t` sends SIGTERM (graceful
// shutdown) instead of SIGKILL (bug B20). The timer logic itself is covered
// deterministically by TestScheduleStopSendsSIGTERM in cli/dreego/main_test.go;
// this test exercises the real CLI subprocess path.
func TestBugRunTimerSigterm(t *testing.T) {
	repoRoot, err := dreegotest.RepoRoot()
	if err != nil {
		t.Fatalf("RepoRoot: %v", err)
	}

	bin := filepath.Join(t.TempDir(), "dreego")
	build := exec.Command("go", "build", "-o", bin, "./cli/dreego")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build CLI: %v\n%s", err, out)
	}

	dir := t.TempDir()
	goMod := fmt.Sprintf("module t\ngo 1.22\nrequire github.com/dreego-stack/dreego v0.0.0\nreplace github.com/dreego-stack/dreego => %s\n", repoRoot)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatal(err)
	}
	mainGo := `package main
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
	dreego.Listen(":0")
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "dreego", "routes"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "dreego", "routes", "get.dreego"), []byte("<div>hello</div>"), 0644); err != nil {
		t.Fatal(err)
	}

	gen := exec.Command(bin, "generate")
	gen.Dir = dir
	if out, err := gen.CombinedOutput(); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}

	outfile := filepath.Join(dir, "run.out")
	f, err := os.Create(outfile)
	if err != nil {
		t.Fatal(err)
	}
	run := exec.Command(bin, "run", "-t", "1")
	run.Dir = dir
	run.Stdout = f
	run.Stderr = f
	if err := run.Start(); err != nil {
		t.Fatalf("run start: %v", err)
	}

	done := make(chan error, 1)
	go func() { done <- run.Wait() }()

	select {
	case <-done:
	case <-time.After(15 * time.Second):
		run.Process.Kill()
		t.Fatal("run -t did not exit in time")
	}
	f.Close()

	out, _ := os.ReadFile(outfile)
	if !strings.Contains(string(out), "SIGTERM received") {
		t.Fatalf("server did not receive SIGTERM (B20)\n%s", out)
	}
}
