package core

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestProject(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for path, content := range files {
		full := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func runInDir(t *testing.T, dir string) error {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	return Run(false)
}

func TestRunAbortsOnSourceReadFailure(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"dreego/routes/get.dreego":    "<div><p>ok</p></div>",
		"dreego/routes/broken.dreego": "",
	})
	target := filepath.Join(dir, "dreego", "routes", "broken.dreego")
	os.Remove(target)
	if err := os.Symlink(filepath.Join(dir, "missing-target"), target); err != nil {
		t.Skipf("cannot create broken symlink: %v", err)
	}

	err := runInDir(t, dir)
	if err == nil {
		t.Fatal("expected generation to abort on unreadable source")
	}
	if !strings.Contains(err.Error(), "broken.dreego") {
		t.Errorf("error must contain the affected path, got %q", err)
	}
}

func TestRunAbortsOnReadDirFailure(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"dreego/routes/get.dreego": "<div><p>ok</p></div>",
	})
	secretDir := filepath.Join(dir, "dreego", "routes", "secret")
	if err := os.MkdirAll(secretDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "missing-target"), filepath.Join(secretDir, "get.dreego")); err != nil {
		t.Fatal(err)
	}

	err := runInDir(t, dir)
	if err == nil {
		t.Fatal("expected generation to abort on unreadable source under secret")
	}
	if !strings.Contains(err.Error(), "secret") {
		t.Errorf("error must contain the affected path, got %q", err)
	}
}

func TestRunAbortsOnLayoutReadFailure(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"dreego/routes/get.dreego": "<div><p>ok</p></div>",
	})
	layoutDir := filepath.Join(dir, "dreego", "layouts")
	if err := os.MkdirAll(layoutDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(dir, "missing-target"), filepath.Join(layoutDir, "default.dreego")); err != nil {
		t.Fatal(err)
	}

	err := runInDir(t, dir)
	if err == nil {
		t.Fatal("expected generation to abort on unreadable layout")
	}
	if !strings.Contains(err.Error(), "default.dreego") {
		t.Errorf("error must contain the affected path, got %q", err)
	}
}

func TestSSRContextSessionValErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.SessionVal("key"); got != "" {
		t.Errorf("expected empty value on store error, got %q", got)
	}
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store read failure")
	}
	if !strings.Contains(c.SessionError().Error(), "store read failure") {
		t.Errorf("SessionError must contain the cause, got %v", c.SessionError())
	}
}

func TestSSRContextSetSessionValErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	c.SetSessionVal("key", "value")
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store write failure")
	}
}

func TestSSRContextDelSessionValErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	c.DelSessionVal("key")
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store delete failure")
	}
}

func TestSSRContextDestroySessionErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	c.DestroySession()
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface store destroy failure")
	}
}

func TestSSRContextCSRFTokenErrorReachable(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r = WithStore(r, failingStore{})
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.CSRFToken(); got != "" {
		t.Errorf("expected empty token on store error, got %q", got)
	}
	if c.SessionError() == nil {
		t.Fatal("expected SessionError to surface CSRF read failure")
	}
}

func TestSSRContextFormErrorReachable(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("a=%zz"))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.FormValue("a"); got != "" {
		t.Errorf("expected empty value on parse error, got %q", got)
	}
	if c.FormError() == nil {
		t.Fatal("expected FormError to surface parse failure")
	}
}

func TestSSRContextFormErrorDistinguishableFromEmpty(t *testing.T) {
	r := httptest.NewRequest("POST", "/", strings.NewReader("a="))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	c := NewSSR(httptest.NewRecorder(), r)

	if got := c.FormValue("a"); got != "" {
		t.Errorf("expected empty value for empty field, got %q", got)
	}
	if c.FormError() != nil {
		t.Errorf("valid empty form must not produce FormError, got %v", c.FormError())
	}
}

func TestRecoveryDoesNotDiscloseCause(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/boom", nil)

	Recovery(nil)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(fmt.Errorf("database: connection to db.example.com:5432 failed"))
	})).ServeHTTP(w, r)

	body := w.Body.String()
	if strings.Contains(body, "database") || strings.Contains(body, "db.example.com") {
		t.Errorf("panic cause disclosed in body: %q", body)
	}
}

func TestCSRFStoreFailureReachesErrorPath(t *testing.T) {
	mw := CSRF(failingStore{})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	hit := false

	mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		hit = true
	})).ServeHTTP(w, r)

	if hit {
		t.Fatal("handler ran after CSRF token persistence failed")
	}
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	if strings.Contains(w.Body.String(), "store read failure") {
		t.Errorf("internal cause disclosed in body: %q", w.Body.String())
	}
}

func TestListenErrorPropagates(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	app := New()
	err = app.Listen(ln.Addr().String())
	if err == nil {
		t.Fatal("expected Listen to return an error when the port is occupied")
	}
}
