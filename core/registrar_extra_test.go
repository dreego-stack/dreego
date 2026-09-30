package core

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegistrarTypeIsFuncAppError(t *testing.T) {
	var r Registrar = func(app *App) error {
		return app.Register(http.MethodGet, "/x", func(http.ResponseWriter, *http.Request) {})
	}
	if r == nil {
		t.Fatal("Registrar must be assignable from a func literal")
	}
	if r(New()) != nil {
		t.Fatal("a valid registrar must return nil")
	}
}

func TestNewAppliesRegistrarRoutes(t *testing.T) {
	app := New(func(a *App) error {
		return a.Register(http.MethodGet, "/ping", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		})
	})
	rec := httptest.NewRecorder()
	app.Handler().ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ping", nil))
	if rec.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
}

func TestNewStopsAtFirstRegistrarError(t *testing.T) {
	second := false
	boom := errors.New("stop")
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic")
		}
		if second {
			t.Fatal("New must not call registrars after the first error")
		}
	}()
	New(
		func(*App) error { return boom },
		func(*App) error { second = true; return nil },
	)
}

func TestNewAppliesMultipleRegistrarsAcrossRoutes(t *testing.T) {
	app := New(
		func(a *App) error {
			return a.Register(http.MethodGet, "/a", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("A")) })
		},
		func(a *App) error {
			return a.Register(http.MethodGet, "/b", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("B")) })
		},
	)
	h := app.Handler()
	for path, want := range map[string]string{"/a": "A", "/b": "B"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
		if rec.Code != 200 || strings.TrimSpace(rec.Body.String()) != want {
			t.Fatalf("GET %s = %d %q, want 200 %q", path, rec.Code, rec.Body.String(), want)
		}
	}
}

func TestNewPanicsWithOriginalError(t *testing.T) {
	sentinel := errors.New("sentinel")
	defer func() {
		if got := recover(); !errors.Is(got.(error), sentinel) {
			t.Fatalf("recovered %v, want sentinel", got)
		}
	}()
	New(func(*App) error { return sentinel })
}

func TestNewDuplicateRouteRegistrarPanics(t *testing.T) {
	register := func(a *App) error {
		return a.Register(http.MethodGet, "/dup", func(http.ResponseWriter, *http.Request) {})
	}
	defer func() {
		if recover() == nil {
			t.Fatal("a duplicate route registrar must panic")
		}
	}()
	New(register, register)
}

func TestNewNilRegistrarsOnly(t *testing.T) {
	app := New(nil, nil, nil)
	if app == nil {
		t.Fatal("New with only nils must return an App")
	}
}
