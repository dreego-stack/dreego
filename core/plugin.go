package core

import (
	"context"
	"io/fs"
	"net/http"
)

type Plugin interface {
	Name() string
	RegisterRoutes(app *App)
	Middlewares() []func(http.Handler) http.Handler
	Assets() fs.FS
	OnStart(ctx context.Context) error
	OnShutdown(ctx context.Context) error
}
