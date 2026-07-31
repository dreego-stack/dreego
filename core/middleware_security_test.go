package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecurityHeadersIncludesCSP(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	csp := w.Result().Header.Get("Content-Security-Policy")
	if csp == "" {
		t.Fatal("expected Content-Security-Policy header to be set")
	}
}

func TestSecurityHeadersCSPDefaultHasSelf(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	csp := w.Result().Header.Get("Content-Security-Policy")
	if !containsStr(csp, "self") {
		t.Error("expected default CSP to allow 'self'")
	}
}

func TestSecurityHeadersCSPOverride(t *testing.T) {
	custom := "default-src 'none'"
	SetCSP(custom)
	defer SetCSP("")

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	csp := w.Result().Header.Get("Content-Security-Policy")
	if csp != custom {
		t.Errorf("expected custom CSP %q, got %q", custom, csp)
	}
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
