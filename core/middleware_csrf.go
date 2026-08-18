package core

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net/http"
	"os"
)

func CSRF(store Store) func(http.Handler) http.Handler {
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := store.Get(r, "csrf_token")
			if err != nil {
				logger.Error("csrf token read failed", "error", err)
			}
			if token == "" {
				token = generateCSRFToken()
				if err := store.Set(w, r, "csrf_token", token, nil); err != nil {
					logger.Error("csrf token write failed", "error", err)
					token = ""
				}
			}
			if token != "" {
				secure := isSecureForCSRF(r, store)
				http.SetCookie(w, &http.Cookie{
					Name:     "csrf_token",
					Value:    token,
					Path:     "/",
					HttpOnly: false,
					Secure:   secure,
					SameSite: http.SameSiteStrictMode,
				})
			}

			if isUnsafeMethod(r.Method) {
				clientToken := r.Header.Get("X-CSRF-Token")
				if clientToken == "" {
					r.ParseForm()
					clientToken = r.FormValue("csrf_token")
				}
				if subtle.ConstantTimeCompare([]byte(clientToken), []byte(token)) != 1 {
					http.Error(w, "invalid csrf token", http.StatusForbidden)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func isSecureForCSRF(r *http.Request, store Store) bool {
	if r.TLS != nil {
		return true
	}
	if cs, ok := store.(*CookieStore); ok {
		return isTLS(r, cs.TrustedProxies())
	}
	return false
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic("csrf: failed to read random bytes: " + err.Error())
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

func isUnsafeMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
		return true
	}
	return false
}
