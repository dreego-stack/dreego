package core

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestConcurrentRouteRegistrationNoRace(t *testing.T) {
	app := New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = app.Register(http.MethodGet, "/r"+itoa(n), func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
		}(i)
	}
	wg.Wait()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}

func TestConcurrentAppConfigNoRace(t *testing.T) {
	app := New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(4)
		go func(i int) {
			defer wg.Done()
			_ = app.SetLogging(i%2 == 0)
		}(i)
		go func(i int) {
			defer wg.Done()
			_ = app.SetCSRF(i%2 == 0)
		}(i)
		go func(i int) {
			defer wg.Done()
			_ = app.SetCSP("default-src 'self'")
		}(i)
		go func(i int) {
			defer wg.Done()
			_ = app.Use(func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})
			})
		}(i)
	}
	wg.Wait()
}

func TestConcurrentSessionStoreNoRace(t *testing.T) {
	store := NewCookieStore(testSecret)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			_ = store.Set(w, r, "k"+itoa(i), "v"+itoa(i), nil)
		}(i)
		go func(i int) {
			defer wg.Done()
			r := httptest.NewRequest("GET", "/", nil)
			_, _ = store.Get(r, "k"+itoa(i))
		}(i)
	}
	wg.Wait()
}

func TestConcurrentCookiePolicyNoRace(t *testing.T) {
	store := NewCookieStore(testSecret)
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func(i int) {
			defer wg.Done()
			store.SetCookiePolicy(CookiePolicy{
				SameSite: http.SameSiteLaxMode,
				Secure:   i%2 == 0,
				HttpOnly: true,
				Path:     "/",
			})
		}(i)
		go func(i int) {
			defer wg.Done()
			store.SetTrustedProxies([]string{"10.0.0." + itoa(i%9)})
		}(i)
		go func(i int) {
			defer wg.Done()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			_ = store.Set(w, r, "k"+itoa(i), "v", nil)
			_ = store.TrustedProxies()
		}(i)
	}
	wg.Wait()
}

func TestConcurrentCSRFMiddlewareNoRace(t *testing.T) {
	store := NewCookieStore(testSecret)
	mw := CSRF(store)
	handler := mw(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			handler.ServeHTTP(w, r)
		}()
	}
	wg.Wait()
}

func TestConcurrentReadyHandlerNoRace(t *testing.T) {
	app := New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			app.SetReady(true)
		}()
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/ready", nil)
			app.readyHandler().ServeHTTP(w, r)
		}()
	}
	wg.Wait()
}
