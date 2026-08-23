package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func TestStopServerKillsProcessAfterTimeout(t *testing.T) {
	ready := t.TempDir() + "/ready"
	cmd := exec.Command(os.Args[0], "-test.run=TestDevServerHelperProcess")
	cmd.Env = append(os.Environ(), "DREEGO_DEV_HELPER=1", "DREEGO_DEV_HELPER_READY="+ready)
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for {
		if _, err := os.Stat(ready); err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("helper process did not become ready")
		}
		time.Sleep(time.Millisecond)
	}

	started := time.Now()
	if err := stopServer(cmd, 50*time.Millisecond); err != nil {
		t.Fatalf("stopServer: %v", err)
	}
	if elapsed := time.Since(started); elapsed > time.Second {
		t.Fatalf("stopServer took %s, want at most 1s", elapsed)
	}
}

func TestDevServerHelperProcess(t *testing.T) {
	if os.Getenv("DREEGO_DEV_HELPER") != "1" {
		return
	}
	signal.Ignore(syscall.SIGTERM)
	if err := os.WriteFile(os.Getenv("DREEGO_DEV_HELPER_READY"), nil, 0600); err != nil {
		os.Exit(2)
	}
	for {
		time.Sleep(time.Second)
	}
}
