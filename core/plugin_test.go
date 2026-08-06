package core

import (
	"context"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
)

// testPlugin is a minimal full implementation of the frozen v1 Plugin
// contract. Used for the compile-time interface-satisfaction check.
type testPlugin struct {
	name string
}

func (p *testPlugin) Name() string                                   { return p.name }
func (p *testPlugin) RegisterRoutes()                                {}
func (p *testPlugin) Middlewares() []func(http.Handler) http.Handler { return nil }
func (p *testPlugin) Assets() fs.FS                                  { return fstest.MapFS{} }
func (p *testPlugin) OnStart(ctx context.Context) error              { return nil }
func (p *testPlugin) OnShutdown(ctx context.Context) error           { return nil }

// Compile-time guarantee: testPlugin satisfies the frozen Plugin contract.
var _ Plugin = (*testPlugin)(nil)

// routePlugin registers a route via core.Register inside RegisterRoutes, the
// way an external plugin module would (importing core and calling Register).
type routePlugin struct{}

func (p *routePlugin) Name() string { return "route-plugin" }
func (p *routePlugin) RegisterRoutes() {
	Register("GET", "/plugin-route", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("plugin-route-ok"))
	})
}
func (p *routePlugin) Middlewares() []func(http.Handler) http.Handler { return nil }
func (p *routePlugin) Assets() fs.FS                                  { return fstest.MapFS{} }
func (p *routePlugin) OnStart(ctx context.Context) error              { return nil }
func (p *routePlugin) OnShutdown(ctx context.Context) error           { return nil }

var _ Plugin = (*routePlugin)(nil)

// headerPlugin appends a response header via middleware, proving that
// UsePlugin collects and injects middleware into the Build() stack.
type headerPlugin struct{}

func (p *headerPlugin) Name() string { return "header-plugin" }
func (p *headerPlugin) RegisterRoutes() {
	Register("GET", "/plugin-middleware", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
}
func (p *headerPlugin) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Plugin-Middleware", "active")
				next.ServeHTTP(w, r)
			})
		},
	}
}
func (p *headerPlugin) Assets() fs.FS                     { return fstest.MapFS{} }
func (p *headerPlugin) OnStart(ctx context.Context) error { return nil }
func (p *headerPlugin) OnShutdown(ctx context.Context) error {
	return nil
}

var _ Plugin = (*headerPlugin)(nil)

// lifecyclePlugin records OnStart/OnShutdown invocations so the tests can
// assert that UsePlugin registers them and StartPlugins/ShutdownPlugins call
// them exactly once each.
type lifecyclePlugin struct {
	started  bool
	shutdown bool
	startErr error
	shutErr  error
}

func (p *lifecyclePlugin) Name() string { return "lifecycle-plugin" }
func (p *lifecyclePlugin) RegisterRoutes() {
	Register("GET", "/plugin-lifecycle", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
func (p *lifecyclePlugin) Middlewares() []func(http.Handler) http.Handler { return nil }
func (p *lifecyclePlugin) Assets() fs.FS                                  { return fstest.MapFS{} }
func (p *lifecyclePlugin) OnStart(ctx context.Context) error {
	p.started = true
	return p.startErr
}
func (p *lifecyclePlugin) OnShutdown(ctx context.Context) error {
	p.shutdown = true
	return p.shutErr
}

var _ Plugin = (*lifecyclePlugin)(nil)

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
