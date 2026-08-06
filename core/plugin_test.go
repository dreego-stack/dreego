package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

// Test 1: Interface satisfaction.
func TestPluginInterfaceSatisfaction(t *testing.T) {
	var p Plugin = &testPlugin{name: "satisfaction"}
	if p.Name() != "satisfaction" {
		t.Fatalf("expected plugin name to be preserved through the interface, got %q", p.Name())
	}
}

// Test 2: UsePlugin registers plugin routes.
func TestUsePluginRegistersRoutes(t *testing.T) {
	Reset()
	UsePlugin(&routePlugin{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-route", nil)
	ServeMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("plugin-registered route returned %d, want %d", rr.Code, http.StatusOK)
	}
	if body := rr.Body.String(); body != "plugin-route-ok" {
		t.Errorf("plugin-registered route body = %q, want %q", body, "plugin-route-ok")
	}
}

// Test 3: UsePlugin collects and injects middleware into the Build() stack.
func TestUsePluginCollectsMiddleware(t *testing.T) {
	Reset()
	UsePlugin(&headerPlugin{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-middleware", nil)
	ServeMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("X-Plugin-Middleware"); got != "active" {
		t.Errorf("plugin middleware header = %q, want %q", got, "active")
	}
}

// Test 4: UsePlugin registers lifecycle hooks; StartPlugins/ShutdownPlugins
// call OnStart/OnShutdown on every registered plugin.
func TestUsePluginLifecycle(t *testing.T) {
	Reset()

	first := &lifecyclePlugin{}
	second := &lifecyclePlugin{}
	UsePlugin(first)
	UsePlugin(second)

	if err := StartPlugins(context.Background()); err != nil {
		t.Fatalf("StartPlugins returned error: %v", err)
	}
	if !first.started || !second.started {
		t.Error("StartPlugins must call OnStart on every registered plugin")
	}

	if err := ShutdownPlugins(context.Background()); err != nil {
		t.Fatalf("ShutdownPlugins returned error: %v", err)
	}
	if !first.shutdown || !second.shutdown {
		t.Error("ShutdownPlugins must call OnShutdown on every registered plugin")
	}
}

// Test 4b: Lifecycle errors propagate to the caller.
func TestUsePluginLifecycleErrors(t *testing.T) {
	Reset()

	UsePlugin(&lifecyclePlugin{startErr: context.Canceled})
	if err := StartPlugins(context.Background()); err == nil {
		t.Error("StartPlugins must propagate an OnStart error")
	}

	Reset()
	UsePlugin(&lifecyclePlugin{shutErr: context.DeadlineExceeded})
	if err := ShutdownPlugins(context.Background()); err == nil {
		t.Error("ShutdownPlugins must propagate an OnShutdown error")
	}
}

// Test 5: Plugin middleware runs in FIFO order — the first registered plugin
// is the outermost middleware and runs first on request entry.
//
// Defined FIFO semantics: UsePlugin(pluginA); UsePlugin(pluginB) means on a
// request A runs first, then B, then the handler. This is the standard Go
// middleware convention (e.g. Chi). The current Build() loop
// `for _, mw := range pluginMiddlewares { h = mw(h) }` produces LIFO (B first),
// so this test is RED against the current implementation.
func TestPluginMiddlewareFIFOOrder(t *testing.T) {
	Reset()

	var log []string
	var mu sync.Mutex
	UsePlugin(&orderPlugin{tag: "A", log: &log, mu: &mu})
	UsePlugin(&orderPlugin{tag: "B", log: &log, mu: &mu})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-order", nil)
	ServeMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}
	want := []string{"A", "B"}
	if !reflect.DeepEqual(log, want) {
		t.Errorf("plugin middleware execution order = %v, want %v (FIFO: first registered runs first)", log, want)
	}
}

// Test 6: Middleware order is fixated on the first Build(). Registering a
// plugin after the handler is already built must not change the order.
func TestPluginMiddlewareOrderFixatedOnFirstBuild(t *testing.T) {
	Reset()

	var log []string
	var mu sync.Mutex
	UsePlugin(&orderPlugin{tag: "A", log: &log, mu: &mu})

	// First Build() fixates the stack.
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-order", nil)
	ServeMux().ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}

	// Registering B after the first Build() must not change the order.
	UsePlugin(&orderPlugin{tag: "B", log: &log, mu: &mu})

	log = nil
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/plugin-order", nil)
	ServeMux().ServeHTTP(rr, req)

	want := []string{"A"}
	if !reflect.DeepEqual(log, want) {
		t.Errorf("after first Build(), middleware order = %v, want %v (order fixated on first Build)", log, want)
	}
}

// Test 7: A nil plugin middleware entry must not panic; the stack stays stable
// and the route still serves.
func TestPluginMiddlewareNilEntryStable(t *testing.T) {
	Reset()

	UsePlugin(&nilMiddlewarePlugin{})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin-nil-mw", nil)
	ServeMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("route returned %d, want %d", rr.Code, http.StatusOK)
	}
}
