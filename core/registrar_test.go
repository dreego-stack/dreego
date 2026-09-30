package core

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewAppliesRegistrarsInOrder(t *testing.T) {
	var order []string
	first := func(app *App) error {
		order = append(order, "first")
		return app.Register(http.MethodGet, "/first", func(http.ResponseWriter, *http.Request) {})
	}
	second := func(app *App) error {
		order = append(order, "second")
		return app.Register(http.MethodGet, "/second", func(http.ResponseWriter, *http.Request) {})
	}

	app := New(first, second)

	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Fatalf("registrar order = %v, want [first second]", order)
	}
	if app == nil {
		t.Fatal("New returned nil App")
	}
}

func TestNewWithoutRegistrarsReturnsApp(t *testing.T) {
	if app := New(); app == nil {
		t.Fatal("New() returned nil App")
	}
}

func TestNewPanicsOnRegistrarError(t *testing.T) {
	boom := errors.New("registrar failed")
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("New did not panic on registrar error")
		}
		err, ok := recovered.(error)
		if !ok || !errors.Is(err, boom) {
			t.Fatalf("recovered = %v, want %v", recovered, boom)
		}
	}()

	New(func(*App) error { return boom })
}

func TestNewSkipsNilRegistrar(t *testing.T) {
	if app := New(nil); app == nil {
		t.Fatal("New(nil) returned nil App")
	}
}
