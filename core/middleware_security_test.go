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

func TestSecurityHeadersContentTypeOptions(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	if got := w.Result().Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("expected X-Content-Type-Options 'nosniff', got %q", got)
	}
}

func TestSecurityHeadersFrameOptions(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	if got := w.Result().Header.Get("X-Frame-Options"); got != "DENY" {
		t.Errorf("expected X-Frame-Options 'DENY', got %q", got)
	}
}

func TestSecurityHeadersReferrerPolicy(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	if got := w.Result().Header.Get("Referrer-Policy"); got != "strict-origin-when-cross-origin" {
		t.Errorf("expected Referrer-Policy 'strict-origin-when-cross-origin', got %q", got)
	}
}

func TestSecurityHeadersPermissionsPolicy(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	SecurityHeaders()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})).ServeHTTP(w, r)

	got := w.Result().Header.Get("Permissions-Policy")
	if got == "" {
		t.Fatal("expected Permissions-Policy header to be set")
	}
	if !containsStr(got, "geolocation=()") {
		t.Errorf("expected Permissions-Policy to restrict geolocation, got %q", got)
	}
}
