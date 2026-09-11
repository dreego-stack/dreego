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

func TestHostDiagnosesDynamicRoute(t *testing.T) {
	app := dreego.New()
	if err := app.Register(http.MethodGet, "/users/{id}", func(http.ResponseWriter, *http.Request) {}); err != nil {
		t.Fatalf("Register: %v", err)
	}
	host, err := New(app)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	request := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	response := httptest.NewRecorder()
	host.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d for %v", response.Code, http.StatusInternalServerError, dreego.ErrDynamicRenderRoute)
	}
}

func TestAssetHandlerEnforcesReadOnlyLiteralRequests(t *testing.T) {
	app := dreego.New()
	if err := app.RegisterRender("/", dreego.ComponentFunc(func(dreego.RenderContext) (dreego.Result, error) {
		return dreego.Result{HTML: []byte("<main>Ready</main>")}, nil
	})); err != nil {
		t.Fatalf("RegisterRender: %v", err)
	}
	host, err := New(app)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	head := httptest.NewRecorder()
	host.ServeHTTP(head, httptest.NewRequest(http.MethodHead, "/", nil))
	if head.Code != http.StatusOK || head.Body.Len() != 0 {
		t.Fatalf("HEAD response = status %d body %q", head.Code, head.Body.String())
	}
	query := httptest.NewRecorder()
	host.ServeHTTP(query, httptest.NewRequest(http.MethodGet, "/?private=value", nil))
	if query.Code != http.StatusNotFound {
		t.Fatalf("query status = %d, want %d", query.Code, http.StatusNotFound)
	}
	post := httptest.NewRecorder()
	host.ServeHTTP(post, httptest.NewRequest(http.MethodPost, "/", nil))
	if post.Code != http.StatusMethodNotAllowed || post.Header().Get("Allow") != "GET, HEAD" {
		t.Fatalf("POST response = status %d Allow %q", post.Code, post.Header().Get("Allow"))
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
