package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestFindMainRoot(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	projDir, pkg, name := findMainIn(dir)
	if projDir != "." || pkg != "." || name != filepath.Base(dir) {
		t.Errorf("expected root main, got %q %q %q", projDir, pkg, name)
	}
}

func TestFindMainCmd(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "cmd"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "cmd", "main.go"), []byte("package main\n"), 0644); err != nil {
		t.Fatal(err)
	}
	projDir, pkg, name := findMainIn(dir)
	if projDir != "cmd" || pkg != "cmd" || name != "cmd" {
		t.Errorf("expected cmd main, got %q %q %q", projDir, pkg, name)
	}
}

func TestFindMainFallback(t *testing.T) {
	dir := t.TempDir()
	projDir, pkg, name := findMainIn(dir)
	if projDir != "." || pkg != "." || name != "server" {
		t.Errorf("expected fallback server, got %q %q %q", projDir, pkg, name)
	}
}

// TestScheduleStopSendsSIGTERM verifies the run-timer behavior (bug B20): after
// the delay, scheduleStop sends SIGTERM (graceful) — the child process must
// exit via its SIGTERM handler, not via SIGKILL. Uses a real child with a
// short delay; the 3s integration wait and the go-run signal chain are gone.
func TestScheduleStopSendsSIGTERM(t *testing.T) {
	cmd := exec.Command("sh", "-c", "trap 'exit 42' TERM; while :; do sleep 0.05; done")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	done := make(chan error, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		done <- cmd.Wait()
	}()

	scheduleStop(cmd.Process, 100*time.Millisecond)

	err := <-done
	wg.Wait()
	if err == nil {
		t.Fatal("expected child to exit after SIGTERM, got nil error")
	}
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("expected ExitError, got %v", err)
	}
	if code := exitErr.ExitCode(); code != 42 {
		t.Errorf("expected exit code 42 (SIGTERM handler ran), got %d", code)
	}
}

// TestScheduleStopFallsBackToKill verifies scheduleStop kills the process when
// signaling fails (e.g. already-exited process), instead of hanging forever.
func TestScheduleStopFallsBackToKill(t *testing.T) {
	cmd := exec.Command("sh", "-c", "exit 0")
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	cmd.Wait()

	scheduleStop(cmd.Process, 10*time.Millisecond)
}
