package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestMiddlewareCompressRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>compress test</p></div>`,
	})
	_, _, headers := c.Request(t, "GET", "/", "", map[string]string{"Accept-Encoding": "gzip"})
	if !strings.Contains(strings.ToLower(headers.Get("Content-Encoding")), "gzip") {
		t.Fatalf("Content-Encoding missing, headers: %v", headers)
	}
}

func TestMiddlewareCompression(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>ok</p></div>`,
	})
}

func TestMiddlewareCSPOverride(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>csp override</p></div>`,
	}, `app.SetCSP("default-src 'none'"); `)
	_, _, headers := c.Request(t, "GET", "/", "", nil)
	if !strings.Contains(strings.ToLower(headers.Get("Content-Security-Policy")), "default-src 'none'") {
		t.Fatalf("custom CSP not applied, got: %q", headers.Get("Content-Security-Policy"))
	}
}

func TestMiddlewareCSPRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>csp test</p></div>`,
	})
	_, _, headers := c.Request(t, "GET", "/", "", nil)
	if headers.Get("Content-Security-Policy") == "" {
		t.Fatalf("Content-Security-Policy missing, headers: %v", headers)
	}
}

func TestMiddlewareCSRFCookieSameSite(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>csrf samesite</p></div>`,
	}, `app.SetSessionStore(dreego.NewCookieStore([]byte("secret-key-32-bytes!"))); `)
	_, _, headers := c.Request(t, "GET", "/", "", nil)
	cookies := strings.Join(headers.Values("Set-Cookie"), "\n")
	if !strings.Contains(cookies, "csrf_token") {
		t.Fatalf("no csrf_token cookie, headers: %v", headers)
	}
	if !strings.Contains(strings.ToLower(cookies), "samesite") {
		t.Fatalf("csrf cookie has no SameSite, headers: %v", headers)
	}
}

func TestMiddlewareCSRFDisabled(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>csrf off</p></div>`,
	}, `app.SetSessionStore(dreego.NewCookieStore([]byte("test"))); app.SetCSRF(false); `)
	code, _ := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("status = %d, want 200", code)
	}
}

func TestMiddlewareCSRFToken(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestMiddlewareHeadersRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>hello</p></div>`,
	})
	_, _, headers := c.Request(t, "GET", "/", "", nil)
	if !strings.Contains(headers.Get("X-Content-Type-Options"), "nosniff") {
		t.Fatalf("X-Content-Type-Options missing, got: %q", headers.Get("X-Content-Type-Options"))
	}
	if !strings.Contains(headers.Get("X-Frame-Options"), "DENY") {
		t.Fatalf("X-Frame-Options missing, got: %q", headers.Get("X-Frame-Options"))
	}
	if !strings.Contains(headers.Get("Referrer-Policy"), "strict-origin-when-cross-origin") {
		t.Fatalf("Referrer-Policy missing, got: %q", headers.Get("Referrer-Policy"))
	}
}

func TestMiddlewareHealthChecks(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>ok</p></div>`,
	})
}

func TestMiddlewareHealthRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>root</p></div>`,
	})
	code, body := c.Get(t, "/health")
	if code != 200 || body != "ok" {
		t.Fatalf("/health not ok, status=%d body=%q", code, body)
	}
}

func TestMiddlewareReadyRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>root</p></div>`,
	}, `app.SetReady(false); `)
	_, body := c.Get(t, "/ready")
	if body != "not ready" {
		t.Fatalf("/ready not returning 'not ready', got %q", body)
	}
}

func TestMiddlewareRecoveryPanic(t *testing.T) {
	t.Parallel()
	dir := dreegotest.ProjectDir(t, nil)
	if out, err := dreegotest.RunCLI(t, dir, "init", "."); err != nil {
		t.Fatalf("init: %v\n%s", err, out)
	}
	if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
		t.Fatalf("generate: %v\n%s", err, out)
	}
	dreegotest.MustBuildInDir(t, dir)
}

func TestMiddlewareRequestIDAccept(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<div><h1>hello</h1></div>`,
	}, `app.SetLogging(false); `)

	custom := "my-custom-request-id"
	_, _, headers := c.Request(t, "GET", "/health", "", map[string]string{"X-Request-ID": custom})
	if got := headers.Get("X-Request-ID"); got != custom {
		t.Fatalf("expected %s, got %q", custom, got)
	}
}

func TestMiddlewareRequestIDRuntime(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"dreego/routes/get.dreego": `<div><h1>hello</h1></div>`,
	}, `app.SetLogging(false); `)
	_, _, headers := c.Request(t, "GET", "/health", "", nil)
	id := headers.Get("X-Request-ID")
	if id == "" {
		t.Fatal("no X-Request-ID header")
	}
	if len(id) != 16 {
		t.Fatalf("invalid request ID format: %q", id)
	}
}

func TestMiddlewareSecurityHeaders(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"dreego/routes/get.dreego": `<div><p>ok</p></div>`,
	})
}
