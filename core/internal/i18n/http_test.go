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
