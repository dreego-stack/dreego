package core

import (
	"net/http"
	"strings"

)

var routes []route

var redirects []redirectRule
var rewrites []rewriteRule

var loggingEnabled = true

var csrfEnabled = true

var errorHandlers = map[int]http.HandlerFunc{}

var sessionStore Store

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
	routes = append(routes, route{method, pattern, handler})
}

func RegisterRedirect(from, to string, status int) {
	redirects = append(redirects, redirectRule{from: from, to: to, status: status})
}

func RegisterRewrite(from, to string) {
	rewrites = append(rewrites, rewriteRule{from: from, to: to})
}

func ServeMux() http.Handler {
	mux := http.NewServeMux()
	for _, r := range routes {
		if r.method != "" {
			mux.HandleFunc(r.method+" "+r.pattern, r.handler)
		} else {
			mux.HandleFunc(r.pattern, r.handler)
		}
	}

	var h http.Handler = mux
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
	h = Recovery(errorHandlers[500])(h)
	return h
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
	return http.ListenAndServe(addr, ServeMux())
}
