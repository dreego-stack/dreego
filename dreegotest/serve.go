package dreegotest

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Serve builds a temp module from files, starts it as a subprocess on a free
// port, and returns a client for HTTP requests. It replaces shell tests that
// start a server and curl it.
func Serve(t *testing.T, files map[string]string) *Client {
	t.Helper()
	return serveSetup(t, files, "")
}

// ServeSetup is like Serve but injects setup code into the generated main
// function before dreego.Listen. It replaces shell tests that customise the
// server main (e.g. dreego.SetCSRF(false)).
func ServeSetup(t *testing.T, files map[string]string, setup string) *Client {
	t.Helper()
	return serveSetup(t, files, setup)
}

func serveSetup(t *testing.T, files map[string]string, setup string) *Client {
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
	mainGo := fmt.Sprintf("package main\nimport (\n\t_ \"t/dreego/gen\"\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\nfunc main() { %sdreego.Listen(\":%d\") }\n", setup, port)
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

	if _, err := RunCLI(t, dir, "generate"); err != nil {
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

	return &Client{base: fmt.Sprintf("http://127.0.0.1:%d", port), jar: newJar()}
}

// Client wraps HTTP requests against a served app. It keeps a cookie jar so
// session/CSRF cookies set by the server are sent back on later requests.
type Client struct {
	base   string
	client *http.Client
	jar    *jar
}

type jar struct{ cookies map[string]*http.Cookie }

func newJar() *jar { return &jar{cookies: map[string]*http.Cookie{}} }

func (j *jar) save(resp *http.Response) {
	for _, c := range resp.Cookies() {
		j.cookies[c.Name] = c
	}
}

func (j *jar) header() []*http.Cookie {
	out := make([]*http.Cookie, 0, len(j.cookies))
	for _, c := range j.cookies {
		out = append(out, c)
	}
	return out
}

// Cookie returns the value of the most recent cookie with the given name set
// by the server, or "" if none.
func (c *Client) Cookie(name string) string {
	if c == nil || c.jar == nil {
		return ""
	}
	if ck, ok := c.jar.cookies[name]; ok {
		return ck.Value
	}
	return ""
}

// Get performs a GET request and returns status code and body.
func (c *Client) Get(t *testing.T, path string) (int, string) {
	t.Helper()
	code, body, _ := c.Request(t, "GET", path, "", nil)
	return code, body
}
// Request performs an HTTP request with the given method, path, body and
// headers. It disables automatic redirect following so tests can assert on
// redirect status codes and Location headers. Returns status, body and headers.
func (c *Client) Request(t *testing.T, method, path, body string, headers map[string]string) (int, string, http.Header) {
	t.Helper()
	req, err := http.NewRequest(method, c.base+path, strings.NewReader(body))
	if err != nil {
		t.Fatalf("Request: new %s %s: %v", method, path, err)
	}
	for _, ck := range c.jar.header() {
		req.AddCookie(ck)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Request: do %s %s: %v", method, path, err)
	}
	c.jar.save(resp)
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(respBody), resp.Header
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
