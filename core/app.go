package core

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type App struct {
	mu                sync.Mutex
	routes            []route
	redirects         []redirectRule
	rewrites          []rewriteRule
	loggingEnabled    bool
	csrfEnabled       bool
	errorHandlers     map[int]http.HandlerFunc
	sessionStore      Store
	builtHandler      http.Handler
	plugins           []Plugin
	pluginMiddlewares []func(http.Handler) http.Handler
	pluginAssets      []fs.FS
	ready             atomic.Bool
	cspHeader         string
	built             bool
}

func New() *App {
	a := &App{
		loggingEnabled: true,
		csrfEnabled:    true,
		errorHandlers:  map[int]http.HandlerFunc{},
		cspHeader:      defaultCSP,
	}
	a.ready.Store(true)
	return a
}

func (a *App) Register(method, pattern string, handler http.HandlerFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, r := range a.routes {
		if r.method == method && r.pattern == pattern {
			a.routes[i].handler = handler
			return
		}
	}
	a.routes = append(a.routes, route{method, pattern, handler})
}

func (a *App) RegisterRedirect(from, to string, status int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.redirects = append(a.redirects, redirectRule{from: from, to: to, status: status})
}

func (a *App) RegisterRewrite(from, to string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.rewrites = append(a.rewrites, rewriteRule{from: from, to: to})
}

func (a *App) RegisterStatic(path, mime string, content []byte) {
	data := make([]byte, len(content))
	copy(data, content)
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", mime)
		w.Write(data)
	}
	a.Register("GET", path, handler)
}

func (a *App) UsePlugin(p Plugin) {
	p.RegisterRoutes(a)
	a.pluginMiddlewares = append(a.pluginMiddlewares, p.Middlewares()...)
	a.pluginAssets = append(a.pluginAssets, p.Assets())
	a.plugins = append(a.plugins, p)
}

func (a *App) StartPlugins(ctx context.Context) error {
	for _, p := range a.plugins {
		if err := p.OnStart(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) ShutdownPlugins(ctx context.Context) error {
	for _, p := range a.plugins {
		if err := p.OnShutdown(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) SetLogging(enabled bool) {
	a.loggingEnabled = enabled
}

func (a *App) SetCSRF(enabled bool) {
	a.csrfEnabled = enabled
}

func (a *App) SetErrorHandler(code int, handler http.HandlerFunc) {
	a.errorHandlers[code] = handler
}

func (a *App) SetSessionStore(s Store) {
	a.sessionStore = s
}

func (a *App) SessionStore() Store {
	return a.sessionStore
}

func (a *App) SetCSP(value string) {
	a.cspHeader = value
}

func (a *App) SetReady(r bool) {
	a.ready.Store(r)
}

func (a *App) RegisterRule(name string, fn func(string) string) {
	registerCustomRule(name, fn)
}

func (a *App) Build() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.built {
		return
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", a.healthHandler())
	mux.HandleFunc("GET /ready", a.readyHandler())

	for _, r := range a.routes {
		if r.method != "" {
			mux.HandleFunc(r.method+" "+r.pattern, r.handler)
		} else {
			mux.HandleFunc(r.pattern, r.handler)
		}
	}

	var h http.Handler = mux
	for i := len(a.pluginMiddlewares) - 1; i >= 0; i-- {
		if a.pluginMiddlewares[i] == nil {
			continue
		}
		h = a.pluginMiddlewares[i](h)
	}
	h = a.redirectRewriteMiddleware(h)
	if a.sessionStore != nil && a.csrfEnabled {
		h = CSRF(a.sessionStore)(h)
	}
	if a.sessionStore != nil {
		h = a.sessionMiddleware(h)
	}
	if a.loggingEnabled {
		h = RequestLogging()(h)
	}
	h = RequestID()(h)
	h = Compress()(h)
	h = a.securityHeadersMiddleware(h)
	h = Recovery(a.errorHandlers[500])(h)
	a.builtHandler = h
	a.built = true
}

func (a *App) Handler() http.Handler {
	a.Build()
	return a.builtHandler
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Handler().ServeHTTP(w, r)
}

func (a *App) Listen(addr string) error {
	a.Build()
	srv := &http.Server{
		Addr:    addr,
		Handler: a.Handler(),
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (a *App) redirectRewriteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rw := range a.rewrites {
			if matchRewrite(rw, r.URL.Path) {
				r.URL.Path = strings.Replace(r.URL.Path, strings.TrimSuffix(rw.from, "/*"), strings.TrimSuffix(rw.to, "/*"), 1)
			}
		}

		for _, rd := range a.redirects {
			if target, ok := matchRedirect(rd, r.URL.Path); ok {
				http.Redirect(w, r, target, rd.status)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (a *App) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, WithStore(r, a.sessionStore))
	})
}

func (a *App) healthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func (a *App) readyHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if a.ready.Load() {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ready"))
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
	}
}

func (a *App) securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := a.cspHeader
		if csp == "" {
			csp = "default-src 'self'"
		}
		w.Header().Set("Content-Security-Policy", csp)
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		next.ServeHTTP(w, r)
	})
}
