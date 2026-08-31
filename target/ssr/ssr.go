// Package ssr is the explicit server-side rendering host. It owns net/http
// startup, server timeouts and graceful shutdown. The App stays target-neutral
// and may be passed to other hosts or builders.
package ssr

import (
	"context"

	dreego "github.com/dreego-stack/dreego/core"
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
