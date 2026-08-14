package core

import (
	"context"
	"io/fs"
	"net/http"
	"testing/fstest"
)

type multiRoutePlugin struct{}

func (p *multiRoutePlugin) Name() string { return "multi-route-plugin" }
func (p *multiRoutePlugin) RegisterRoutes(app *App) {
	app.Register("GET", "/plugin/multi", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("multi-get"))
	})
	app.Register("POST", "/plugin/multi", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("multi-post"))
	})
	app.Register("GET", "/plugin/multi/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("multi-get-" + r.PathValue("id")))
	})
}
func (p *multiRoutePlugin) Middlewares() []func(http.Handler) http.Handler { return nil }
func (p *multiRoutePlugin) Assets() fs.FS                                  { return fstest.MapFS{} }
func (p *multiRoutePlugin) OnStart(ctx context.Context) error              { return nil }
func (p *multiRoutePlugin) OnShutdown(ctx context.Context) error           { return nil }

var _ Plugin = (*multiRoutePlugin)(nil)

type overlapRoutePlugin struct {
	path   string
	marker string
}

func (p *overlapRoutePlugin) Name() string { return "overlap-" + p.marker }
func (p *overlapRoutePlugin) RegisterRoutes(app *App) {
	app.Register("GET", p.path, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(p.marker))
	})
}
func (p *overlapRoutePlugin) Middlewares() []func(http.Handler) http.Handler { return nil }
func (p *overlapRoutePlugin) Assets() fs.FS                                  { return fstest.MapFS{} }
func (p *overlapRoutePlugin) OnStart(ctx context.Context) error              { return nil }
func (p *overlapRoutePlugin) OnShutdown(ctx context.Context) error           { return nil }

var _ Plugin = (*overlapRoutePlugin)(nil)
