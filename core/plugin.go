package core

import (
	"context"
	"io/fs"
	"net/http"
)

// Plugin is the frozen v1 contract for external plugin modules. Plugins
// import core and satisfy this interface; core never imports a plugin.
type Plugin interface {
	Name() string
	RegisterRoutes() // plugin calls dreego.Register(...) internally
	Middlewares() []func(http.Handler) http.Handler
	Assets() fs.FS
	OnStart(ctx context.Context) error
	OnShutdown(ctx context.Context) error
}

var plugins []Plugin

var pluginMiddlewares []func(http.Handler) http.Handler

var pluginAssets []fs.FS

// UsePlugin registers a plugin's routes, middleware, assets and lifecycle
// hooks with the core runtime.
func UsePlugin(p Plugin) {
	p.RegisterRoutes()
	pluginMiddlewares = append(pluginMiddlewares, p.Middlewares()...)
	pluginAssets = append(pluginAssets, p.Assets())
	plugins = append(plugins, p)
}

// StartPlugins calls OnStart on every registered plugin and returns the first
// error encountered.
func StartPlugins(ctx context.Context) error {
	for _, p := range plugins {
		if err := p.OnStart(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownPlugins calls OnShutdown on every registered plugin and returns the
// first error encountered.
func ShutdownPlugins(ctx context.Context) error {
	for _, p := range plugins {
		if err := p.OnShutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}
