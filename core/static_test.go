package core

import (
	"net/http/httptest"
	"testing"
)

func TestMimeByExt(t *testing.T) {
	cases := map[string]string{
		".css":   "text/css; charset=utf-8",
		".js":    "application/javascript; charset=utf-8",
		".svg":   "image/svg+xml",
		".png":   "image/png",
		".ico":   "image/x-icon",
		".jpg":   "image/jpeg",
		".jpeg":  "image/jpeg",
		".html":  "text/html; charset=utf-8",
		".json":  "application/json; charset=utf-8",
		".woff2": "font/woff2",
		".woff":  "font/woff",
	}
	for ext, want := range cases {
		if got := MimeByExt(ext); got != want {
			t.Errorf("MimeByExt(%q) = %q, want %q", ext, got, want)
		}
	}
}

func TestMimeByExtDefault(t *testing.T) {
	for _, ext := range []string{".txt", ".md", ".bin", "", ".unknown"} {
		if got := MimeByExt(ext); got != "application/octet-stream" {
			t.Errorf("MimeByExt(%q) = %q, want application/octet-stream", ext, got)
		}
	}
}

func TestMimeByExtCaseInsensitive(t *testing.T) {
	for _, ext := range []string{".CSS", ".Js", ".PNG", ".HTML", ".Woff2"} {
		if got := MimeByExt(ext); got == "application/octet-stream" {
			t.Errorf("MimeByExt(%q) returned default, expected case-insensitive match", ext)
		}
	}
}

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

func TestStaticPattern(t *testing.T) {
	cases := map[string]string{
		"css/style.css": "/css/style.css",
		"app.js":        "/app.js",
		"img/logo.png":  "/img/logo.png",
	}
	for in, want := range cases {
		if got := staticPattern(in); got != want {
			t.Errorf("staticPattern(%q) = %q, want %q", in, got, want)
		}
	}
}
