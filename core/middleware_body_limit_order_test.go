package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAppBodyLimitRunsBeforeCSRFFormParsing(t *testing.T) {
	app := New()
	store := NewCookieStore([]byte("01234567890123456789012345678901"))
	if err := app.SetSessionStore(store); err != nil {
		t.Fatal(err)
	}
	if err := app.Use(MaxBodyReader(32)); err != nil {
		t.Fatal(err)
	}
	if err := app.Register(http.MethodPost, "/submit", func(http.ResponseWriter, *http.Request) {
		t.Fatal("handler must not run for an oversized body")
	}); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/submit", strings.NewReader("csrf_token="+strings.Repeat("x", 64)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusRequestEntityTooLarge)
	}
}
