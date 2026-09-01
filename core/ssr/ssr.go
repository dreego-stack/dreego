// Package ssr is the explicit server-side rendering host. It owns net/http
// startup, server timeouts and graceful shutdown. The App stays target-neutral
// and may be passed to other hosts or builders.
package ssr

import (
	"context"
	"net/http"

	dreego "github.com/dreego-stack/dreego/core"
	"github.com/dreego-stack/dreego/core/internal/middleware"
	"github.com/dreego-stack/dreego/core/internal/server"
	"github.com/dreego-stack/dreego/core/internal/session"
	"github.com/dreego-stack/dreego/core/internal/validate"
)

// Listen starts the App's HTTP server on addr and blocks until it shuts down.
func Listen(app *dreego.App, addr string) error {
	return app.Listen(addr)
}

// Shutdown gracefully stops the App's HTTP server, waiting for in-flight
// requests to drain within the context deadline.
func Shutdown(app *dreego.App, ctx context.Context) error {
	return app.Shutdown(ctx)
}

// DefaultAddr returns the default listen address used by generated mains.
func DefaultAddr() string {
	return ":8080"
}

var ErrServerRunning = server.ErrServerRunning

type ServerConfig = server.ServerConfig

func DefaultServerConfig() ServerConfig {
	return server.DefaultServerConfig()
}

func RequestLogging() func(http.Handler) http.Handler {
	return middleware.RequestLogging()
}

func MaxBodyReader(max int64) func(http.Handler) http.Handler {
	return middleware.MaxBodyReader(max)
}

func Compress() func(http.Handler) http.Handler {
	return middleware.Compress()
}

func CSRF(store dreego.Store) func(http.Handler) http.Handler {
	return middleware.CSRF(store)
}

func Recovery(errorHandler http.HandlerFunc) func(http.Handler) http.Handler {
	return middleware.Recovery(errorHandler)
}

func RequestID() func(http.Handler) http.Handler {
	return middleware.RequestID()
}

func RequestIDFromCtx(ctx context.Context) string {
	return middleware.RequestIDFromCtx(ctx)
}

func WithStore(r *http.Request, s dreego.Store) *http.Request {
	return session.WithStore(r, s)
}

func BindForm(r *http.Request, target any) error {
	return validate.BindForm(r, target)
}
