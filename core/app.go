package core

import (
	"context"
	"net/http"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

type App struct {
	mu             sync.RWMutex
	routes         []route
	redirects      []redirectRule
	rewrites       []rewriteRule
	loggingEnabled bool
	csrfEnabled    bool
	errorHandlers  map[int]http.HandlerFunc
	sessionStore   Store
	builtHandler   http.Handler
	middlewares    []func(http.Handler) http.Handler
	customRules    map[string]validatorFunc
	ready          atomic.Bool
	cspHeader      string
	built          bool
	buildDone      chan struct{}
	server         *http.Server
	serverConfig   ServerConfig
	shutdownDone   chan error
}

func New() *App {
	a := &App{
		loggingEnabled: true,
		csrfEnabled:    true,
		errorHandlers:  map[int]http.HandlerFunc{},
		customRules:    map[string]validatorFunc{},
		cspHeader:      defaultCSP,
		buildDone:      make(chan struct{}),
	}
	a.ready.Store(true)
	return a
}

func (a *App) SetReady(r bool) {
	a.ready.Store(r)
}

func (a *App) Build() {
	a.mu.Lock()
	if a.built {
		a.mu.Unlock()
		return
	}
	a.built = true
	buildDone := a.buildDone
	a.mu.Unlock()

	completed := false
	defer func() {
		if completed {
			return
		}
		a.mu.Lock()
		a.built = false
		close(buildDone)
		a.buildDone = make(chan struct{})
		a.mu.Unlock()
	}()

	if v, ok := a.sessionStore.(storeValidator); ok {
		if err := v.Validate(); err != nil {
			panic(err)
		}
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
	for i := len(a.middlewares) - 1; i >= 0; i-- {
		if a.middlewares[i] == nil {
			continue
		}
		h = a.middlewares[i](h)
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
	a.mu.Lock()
	a.builtHandler = h
	close(buildDone)
	a.mu.Unlock()
	completed = true
}

func (a *App) Handler() http.Handler {
	for {
		a.Build()
		a.mu.RLock()
		h := a.builtHandler
		built := a.built
		buildDone := a.buildDone
		a.mu.RUnlock()
		if h != nil {
			return h
		}
		if !built {
			continue
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			<-buildDone
			a.Handler().ServeHTTP(w, r)
		})
	}
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.Handler().ServeHTTP(w, r)
}

func (a *App) Listen(addr string) error {
	a.Build()
	cfg := a.serverConfig.withDefaults()
	srv := &http.Server{
		Addr:              addr,
		Handler:           a.Handler(),
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		MaxHeaderBytes:    cfg.MaxHeaderBytes,
	}

	a.mu.Lock()
	a.server = srv
	shutdownDone := make(chan error, 1)
	a.shutdownDone = shutdownDone
	a.mu.Unlock()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() { serverErr <- srv.ListenAndServe() }()

	clearState := func() {
		a.mu.Lock()
		a.server = nil
		a.shutdownDone = nil
		a.mu.Unlock()
	}

	select {
	case <-sigCtx.Done():
		shutCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		err := srv.Shutdown(shutCtx)
		if serr := <-serverErr; serr != nil && serr != http.ErrServerClosed {
			clearState()
			return serr
		}
		clearState()
		return err
	case err := <-serverErr:
		if err == http.ErrServerClosed {
			rerr := <-shutdownDone
			clearState()
			return rerr
		}
		clearState()
		return err
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	a.mu.Lock()
	srv := a.server
	done := a.shutdownDone
	a.mu.Unlock()
	if srv == nil {
		return nil
	}
	err := srv.Shutdown(ctx)
	if done != nil {
		select {
		case done <- err:
		default:
		}
	}
	return err
}

func (a *App) redirectRewriteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rw := range a.rewrites {
			if target, ok := applyRewrite(rw, r.URL.Path); ok {
				r2 := r.Clone(r.Context())
				r2.URL.Path = target
				r2.URL.RawPath = ""
				r = r2
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
