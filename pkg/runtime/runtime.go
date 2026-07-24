package runtime

import (
	"net/http"
	"strings"

	"codeberg.org/dreego/dreego/pkg/middleware"
)

var routes []route

var redirects []redirectRule
var rewrites []rewriteRule

var loggingEnabled = true

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
		mux.HandleFunc(r.method+" "+r.pattern, r.handler)
	}

	var h http.Handler = mux
	h = redirectRewriteMiddleware(h)
	if loggingEnabled {
		h = middleware.RequestLogging()(h)
	}
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

func SetLogging(enabled bool) {
	loggingEnabled = enabled
}

func Listen(addr string) error {
	return http.ListenAndServe(addr, ServeMux())
}
