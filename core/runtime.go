package core

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var routes []route

var redirects []redirectRule
var rewrites []rewriteRule

var loggingEnabled = true

var csrfEnabled = true

var errorHandlers = map[int]http.HandlerFunc{}

var sessionStore Store

var builtHandler http.Handler

func Reset() {
	builtHandler = nil
	plugins = nil
	pluginMiddlewares = nil
	pluginAssets = nil
}

type route struct {
	method  string
	pattern string
	handler http.HandlerFunc
}

type redirectRule struct {
	from   string
	to     string
	status int
}

type rewriteRule struct {
	from string
	to   string
}

func Register(method, pattern string, handler http.HandlerFunc) {
	for i, r := range routes {
		if r.method == method && r.pattern == pattern {
			routes[i].handler = handler
			return
		}
	}
	routes = append(routes, route{method, pattern, handler})
}

func RegisterRedirect(from, to string, status int) {
	redirects = append(redirects, redirectRule{from: from, to: to, status: status})
}

func RegisterRewrite(from, to string) {
	rewrites = append(rewrites, rewriteRule{from: from, to: to})
}

func Build() {
	if builtHandler != nil {
		return
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler())
	mux.HandleFunc("GET /ready", readyHandler())

	for _, r := range routes {
		if r.method != "" {
			mux.HandleFunc(r.method+" "+r.pattern, r.handler)
		} else {
			mux.HandleFunc(r.pattern, r.handler)
		}
	}

	var h http.Handler = mux
	for i := len(pluginMiddlewares) - 1; i >= 0; i-- {
		if pluginMiddlewares[i] == nil {
			continue
		}
		h = pluginMiddlewares[i](h)
	}
	h = redirectRewriteMiddleware(h)
	if sessionStore != nil && csrfEnabled {
		h = CSRF(sessionStore)(h)
	}
	if sessionStore != nil {
		h = sessionMiddleware(sessionStore)(h)
	}
	if loggingEnabled {
		h = RequestLogging()(h)
	}
	h = RequestID()(h)
	h = Compress()(h)
	h = SecurityHeaders()(h)
	h = Recovery(errorHandlers[500])(h)
	builtHandler = h
}

func ServeMux() http.Handler {
	Build()
	return builtHandler
}

func redirectRewriteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, rw := range rewrites {
			if matchRewrite(rw, r.URL.Path) {
				r.URL.Path = strings.Replace(r.URL.Path, strings.TrimSuffix(rw.from, "/*"), strings.TrimSuffix(rw.to, "/*"), 1)
			}
		}

		for _, rd := range redirects {
			if r.URL.Path == strings.TrimSuffix(rd.from, "/*") || rd.from == r.URL.Path {
				http.Redirect(w, r, rd.to, rd.status)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func matchRewrite(rw rewriteRule, path string) bool {
	return strings.HasPrefix(path, strings.TrimSuffix(rw.from, "/*"))
}

func sessionMiddleware(store Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, WithStore(r, store))
		})
	}
}

func SetLogging(enabled bool) {
	loggingEnabled = enabled
}

func SetCSRF(enabled bool) {
	csrfEnabled = enabled
}

func SetErrorHandler(code int, handler http.HandlerFunc) {
	errorHandlers[code] = handler
}

func SetSessionStore(s Store) {
	sessionStore = s
}

func SessionStore() Store {
	return sessionStore
}

func Listen(addr string) error {
	Build()
	srv := &http.Server{
		Addr:    addr,
		Handler: ServeMux(),
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
