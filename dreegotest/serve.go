package dreegotest

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"
	"time"

	dreego "github.com/dreego-stack/dreego/core"
)

var serveMu sync.Mutex

// Serve builds a temp module from files, starts it as a subprocess on a free
// port, and returns a client for HTTP requests. It replaces shell tests that
// start a server and curl it.
func Serve(t *testing.T, files map[string]string) *Client {
	t.Helper()
	dir := t.TempDir()

	repoRoot, err := RepoRoot()
	if err != nil {
		t.Fatalf("Serve: %v", err)
	}

	port := freePort(t)

	goMod := fmt.Sprintf("module t\ngo 1.22\nrequire github.com/dreego-stack/dreego v0.0.0\nreplace github.com/dreego-stack/dreego => %s\n", repoRoot)
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("Serve: write go.mod: %v", err)
	}
	mainGo := fmt.Sprintf("package main\nimport (\n\t_ \"t/dreego/gen\"\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\nfunc main() { dreego.Listen(\":%d\") }\n", port)
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainGo), 0644); err != nil {
		t.Fatalf("Serve: write main.go: %v", err)
	}
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatalf("Serve: mkdir %s: %v", filepath.Dir(full), err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatalf("Serve: write %s: %v", path, err)
		}
	}

	serveMu.Lock()
	defer serveMu.Unlock()
	old, err := os.Getwd()
	if err != nil {
		t.Fatalf("Serve: getwd: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("Serve: chdir: %v", err)
	}
	defer os.Chdir(old)

	if err := dreego.Run(false); err != nil {
		t.Fatalf("Serve: generate failed: %v", err)
	}

	bin := filepath.Join(dir, "server")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Serve: go build failed: %v\n%s", err, out)
	}

	proc := exec.Command(bin)
	proc.Dir = dir
	if err := proc.Start(); err != nil {
		t.Fatalf("Serve: start server: %v", err)
	}
	t.Cleanup(func() {
		proc.Process.Kill()
		proc.Wait()
	})

	waitForPort(t, port)

	return &Client{base: fmt.Sprintf("http://127.0.0.1:%d", port)}
}

// Client wraps HTTP requests against a served app.
type Client struct {
	base string
}

// Get performs a GET request and returns status code and body.
func (c *Client) Get(t *testing.T, path string) (int, string) {
	t.Helper()
	resp, err := http.Get(c.base + path)
	if err != nil {
		t.Fatalf("Get %s: %v", path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(body)
}

func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()
	return port
}

func waitForPort(t *testing.T, port int) {
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
	t.Fatalf("waitForPort: server on port %d did not start in time", port)
}
