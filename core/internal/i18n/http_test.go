package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNegotiatorResolutionOrder(t *testing.T) {
	config := Config{
		DefaultLocale: "de",
		Locales:       []LocaleCatalog{{Locale: "de"}, {Locale: "en"}, {Locale: "ar"}},
		Detection:     []string{"account", "cookie", "browser", "custom", "default"},
		Account:       func(*http.Request) string { return "" },
		Resolvers:     []Resolver{func(*http.Request) string { return "ar" }},
	}
	negotiator, err := NewNegotiator(config, nil)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Language", "en-GB,en;q=0.8")
	request.AddCookie(&http.Cookie{Name: "dreego_locale", Value: "fr"})
	if got := negotiator.Resolve(request); got != "en" {
		t.Fatalf("Resolve = %q", got)
	}
}

func TestNegotiatorCookieWinsBrowser(t *testing.T) {
	config := Config{DefaultLocale: "de", Locales: []LocaleCatalog{{Locale: "de"}, {Locale: "en"}}}
	negotiator, err := NewNegotiator(config, nil)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Language", "de")
	request.AddCookie(&http.Cookie{Name: "dreego_locale", Value: "en"})
	if got := negotiator.Resolve(request); got != "en" {
		t.Fatalf("Resolve = %q", got)
	}
}

func TestNegotiatorPrefixWinsAndIsRemovedBeforeRouting(t *testing.T) {
	config := Config{DefaultLocale: "de", URLStrategy: "prefix", Locales: []LocaleCatalog{{Locale: "de"}, {Locale: "en"}}}
	negotiator, err := NewNegotiator(config, nil)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/en/account", nil)
	var path, locale string
	handler := negotiator.Middleware(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		locale = Locale(r.Context())
	}))
	handler.ServeHTTP(httptest.NewRecorder(), request)
	if path != "/account" || locale != "en" {
		t.Fatalf("path = %q, locale = %q", path, locale)
	}
}

func TestNegotiatorDomainWins(t *testing.T) {
	config := Config{DefaultLocale: "de", URLStrategy: "domain", Domains: map[string]string{"en": "example.com"}, Locales: []LocaleCatalog{{Locale: "de"}, {Locale: "en"}}}
	negotiator, err := NewNegotiator(config, nil)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "https://example.com/account", nil)
	if got := negotiator.Resolve(request); got != "en" {
		t.Fatalf("Resolve = %q", got)
	}
}

func TestSetLocaleRejectsUnsupportedChoiceAndWritesSecureCookie(t *testing.T) {
	config := Config{DefaultLocale: "de", Locales: []LocaleCatalog{{Locale: "de"}, {Locale: "en"}}}
	negotiator, err := NewNegotiator(config, nil)
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "https://example.com/locale", nil)
	response := httptest.NewRecorder()
	negotiator.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := SetLocale(w, r, "fr"); err == nil {
			t.Error("expected unsupported locale error")
		}
		if err := SetLocale(w, r, "en"); err != nil {
			t.Error(err)
		}
	})).ServeHTTP(response, request)
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Value != "en" || !cookies[0].HttpOnly || !cookies[0].Secure {
		t.Fatalf("cookies = %+v", cookies)
	}
}
