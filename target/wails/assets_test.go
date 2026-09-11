package wails

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

func TestAssetHandlerServesLiteralRoutesAndStaticAssets(t *testing.T) {
	app := dreego.New()
	page := dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		return dreego.Result{HTML: []byte(`<!doctype html><title>Timer</title><link rel="stylesheet" href="/assets/app.css"><style>main{display:block}</style><main>Ready <a href="/settings">Settings</a></main><script src="/assets/app.js"></script>`)}, nil
	})
	if err := app.RegisterRender("/timer", page); err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	if err := app.RegisterStatic("/assets/app.js", "application/javascript", []byte("ready()")); err != nil {
		t.Fatalf("RegisterStatic: %v", err)
	}
	if err := app.RegisterStatic("/assets/app.css", "text/css", []byte("main{color:green}")); err != nil {
		t.Fatalf("RegisterStatic CSS: %v", err)
	}
	if err := app.RegisterRender("/settings", dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		return dreego.Result{HTML: []byte("<!doctype html><title>Settings</title><main>Settings</main>")}, nil
	})); err != nil {
		t.Fatalf("RegisterRender settings: %v", err)
	}
	host, err := New(app)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	assertWailsResponse(t, host, "/timer", http.StatusOK, "text/html; charset=utf-8", "<title>Timer</title>")
	assertWailsResponse(t, host, "/settings", http.StatusOK, "text/html; charset=utf-8", "<title>Settings</title>")
	assertWailsResponse(t, host, "/assets/app.js", http.StatusOK, "application/javascript", "ready()")
	assertWailsResponse(t, host, "/assets/app.css", http.StatusOK, "text/css", "main{color:green}")
}

func TestAssetHandlerRejectsUnknownAndTraversalPaths(t *testing.T) {
	host, err := New(dreego.New())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, requestPath := range []string{"/missing", "/../secret", "/assets/%2e%2e/secret", `/assets\secret`} {
		assertWailsResponse(t, host, requestPath, http.StatusNotFound, "text/plain; charset=utf-8", "404 page not found")
	}
}

func TestRunRejectsExternalFrontendDevserver(t *testing.T) {
	t.Setenv("FRONTEND_DEVSERVER_URL", "http://127.0.0.1:5173")
	host, err := New(dreego.New())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := host.Run(Options{Path: "/"}); err == nil || !strings.Contains(err.Error(), "FRONTEND_DEVSERVER_URL") {
		t.Fatalf("Run error = %v", err)
	}
}

func assertWailsResponse(t *testing.T, handler http.Handler, path string, status int, contentType, body string) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != status {
		t.Fatalf("GET %s status = %d, want %d", path, response.Code, status)
	}
	if got := response.Header().Get("Content-Type"); got != contentType {
		t.Fatalf("GET %s Content-Type = %q, want %q", path, got, contentType)
	}
	if got := response.Body.String(); !strings.Contains(got, body) {
		t.Fatalf("GET %s body = %q, want it to contain %q", path, got, body)
	}
}
