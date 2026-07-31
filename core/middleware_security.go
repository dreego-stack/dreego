package core

import "net/http"

var cspHeader = "default-src 'self'; script-src 'self' 'unsafe-inline' https://unpkg.com; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self'; base-uri 'self'; form-action 'self'"

func SecurityHeaders() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
			csp := cspHeader
			if csp == "" {
				csp = "default-src 'self'"
			}
			w.Header().Set("Content-Security-Policy", csp)
			next.ServeHTTP(w, r)
		})
	}
}

func SetCSP(value string) {
	cspHeader = value
}