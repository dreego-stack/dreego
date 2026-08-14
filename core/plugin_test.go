package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

func TestPluginInterfaceSatisfaction(t *testing.T) {
	var p Plugin = &testPlugin{name: "satisfaction"}
	if p.Name() != "satisfaction" {
		t.Fatalf("expected plugin name to be preserved through the interface, got %q", p.Name())
	}
}

func TestUsePluginRegistersRoutes(t *testing.T) {
	app := New()
	app.UsePlugin(&routePlugin{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-route", nil)
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("plugin-registered route returned %d, want %d", rr.Code, http.StatusOK)
	}
	if body := rr.Body.String(); body != "plugin-route-ok" {
		t.Errorf("plugin-registered route body = %q, want %q", body, "plugin-route-ok")
	}
}

func TestUsePluginCollectsMiddleware(t *testing.T) {
	app := New()
	app.UsePlugin(&headerPlugin{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-middleware", nil)
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("X-Plugin-Middleware"); got != "active" {
		t.Errorf("plugin middleware header = %q, want %q", got, "active")
	}
}

func TestUsePluginLifecycle(t *testing.T) {
	app := New()

	first := &lifecyclePlugin{}
	second := &lifecyclePlugin{}
	app.UsePlugin(first)
	app.UsePlugin(second)

	if err := app.StartPlugins(context.Background()); err != nil {
		t.Fatalf("StartPlugins returned error: %v", err)
	}
	if !first.started || !second.started {
		t.Error("StartPlugins must call OnStart on every registered plugin")
	}

	if err := app.ShutdownPlugins(context.Background()); err != nil {
		t.Fatalf("ShutdownPlugins returned error: %v", err)
	}
	if !first.shutdown || !second.shutdown {
		t.Error("ShutdownPlugins must call OnShutdown on every registered plugin")
	}
}

func TestUsePluginLifecycleErrors(t *testing.T) {
	app := New()
	app.UsePlugin(&lifecyclePlugin{startErr: context.Canceled})
	if err := app.StartPlugins(context.Background()); err == nil {
		t.Error("StartPlugins must propagate an OnStart error")
	}

	app2 := New()
	app2.UsePlugin(&lifecyclePlugin{shutErr: context.DeadlineExceeded})
	if err := app2.ShutdownPlugins(context.Background()); err == nil {
		t.Error("ShutdownPlugins must propagate an OnShutdown error")
	}
}

func TestPluginMiddlewareFIFOOrder(t *testing.T) {
	app := New()

	var log []string
	var mu sync.Mutex
	app.UsePlugin(&orderPlugin{tag: "A", log: &log, mu: &mu})
	app.UsePlugin(&orderPlugin{tag: "B", log: &log, mu: &mu})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-order", nil)
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}
	want := []string{"A", "B"}
	if !reflect.DeepEqual(log, want) {
		t.Errorf("plugin middleware execution order = %v, want %v (FIFO: first registered runs first)", log, want)
	}
}

func TestPluginMiddlewareOrderFixatedOnFirstBuild(t *testing.T) {
	app := New()

	var log []string
	var mu sync.Mutex
	app.UsePlugin(&orderPlugin{tag: "A", log: &log, mu: &mu})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-order", nil)
	app.Handler().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}

	app.UsePlugin(&orderPlugin{tag: "B", log: &log, mu: &mu})

	log = nil
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/plugin-order", nil)
	app.Handler().ServeHTTP(rr, req)

	want := []string{"A"}
	if !reflect.DeepEqual(log, want) {
		t.Errorf("after first Build(), middleware order = %v, want %v (order fixated on first Build)", log, want)
	}
}

func TestPluginMiddlewareNilEntryStable(t *testing.T) {
	app := New()

	app.UsePlugin(&nilMiddlewarePlugin{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-nil-mw", nil)
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}
}
