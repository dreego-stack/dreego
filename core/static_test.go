package core

import (
	"net/http/httptest"
	"testing"
)

func TestRegisterStaticServesContent(t *testing.T) {
	app := New()
	before := len(app.routes)
	app.RegisterStatic("/assets/app.js", "application/javascript; charset=utf-8", []byte("console.log(1)"))

	if len(app.routes) != before+1 {
		t.Fatalf("expected %d routes, got %d", before+1, len(app.routes))
	}
	r := app.routes[len(app.routes)-1]
	if r.method != "GET" {
		t.Errorf("expected method GET, got %q", r.method)
	}
	if r.pattern != "/assets/app.js" {
		t.Errorf("expected pattern /assets/app.js, got %q", r.pattern)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/assets/app.js", nil)
	r.handler(w, req)

	if ct := w.Header().Get("Content-Type"); ct != "application/javascript; charset=utf-8" {
		t.Errorf("expected Content-Type application/javascript; charset=utf-8, got %q", ct)
	}
	if body := w.Body.String(); body != "console.log(1)" {
		t.Errorf("expected body console.log(1), got %q", body)
	}
}
