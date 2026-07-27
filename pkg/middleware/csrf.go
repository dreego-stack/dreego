package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"codeberg.org/dreego/dreego/pkg/session"
)

func CSRF(store session.Store) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, _ := store.Get(r, "csrf_token")
			if token == "" {
				token = generateCSRFToken()
				store.Set(w, r, "csrf_token", token, &session.Options{
					HttpOnly: true,
					Path:     "/",
				})
			}
			http.SetCookie(w, &http.Cookie{
				Name:     "csrf_token",
				Value:    token,
				Path:     "/",
				HttpOnly: false,
				SameSite: http.SameSiteStrictMode,
			})

			if isUnsafeMethod(r.Method) {
				clientToken := r.Header.Get("X-CSRF-Token")
				if clientToken == "" {
					r.ParseForm()
					clientToken = r.FormValue("csrf_token")
				}
				if clientToken != token {
					http.Error(w, "invalid csrf token", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func isUnsafeMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return true
	}
	return false
}
