package core

import (
	"context"
	"io/fs"
	"net/http"
	"sync"
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

// orderPlugin appends its tag to a shared log inside its middleware, so tests
// can assert the exact execution order of plugin middleware in the stack.
type orderPlugin struct {
	tag string
	log *[]string
	mu  *sync.Mutex
}

func (p *orderPlugin) Name() string { return "order-" + p.tag }
func (p *orderPlugin) RegisterRoutes() {
	Register("GET", "/plugin-order", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
func (p *orderPlugin) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				p.mu.Lock()
				*p.log = append(*p.log, p.tag)
				p.mu.Unlock()
				next.ServeHTTP(w, r)
			})
		},
	}
}
func (p *orderPlugin) Assets() fs.FS                     { return fstest.MapFS{} }
func (p *orderPlugin) OnStart(ctx context.Context) error { return nil }
func (p *orderPlugin) OnShutdown(ctx context.Context) error {
	return nil
}

var _ Plugin = (*orderPlugin)(nil)

// nilMiddlewarePlugin returns a slice containing a nil middleware func to
// verify the stack tolerates a nil entry without panicking.
type nilMiddlewarePlugin struct{}

func (p *nilMiddlewarePlugin) Name() string { return "nil-mw-plugin" }
func (p *nilMiddlewarePlugin) RegisterRoutes() {
	Register("GET", "/plugin-nil-mw", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
func (p *nilMiddlewarePlugin) Middlewares() []func(http.Handler) http.Handler {
	return []func(http.Handler) http.Handler{nil}
}
func (p *nilMiddlewarePlugin) Assets() fs.FS                     { return fstest.MapFS{} }
func (p *nilMiddlewarePlugin) OnStart(ctx context.Context) error { return nil }
func (p *nilMiddlewarePlugin) OnShutdown(ctx context.Context) error {
	return nil
}

var _ Plugin = (*nilMiddlewarePlugin)(nil)
