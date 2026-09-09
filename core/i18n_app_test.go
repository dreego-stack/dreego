package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAppLocalizesRequestFromBrowserPreference(t *testing.T) {
	app := New()
	err := app.SetI18n(I18nConfig{
		DefaultLocale: "de",
		Locales: []LocaleCatalog{
			{Locale: "de", Messages: map[string]LocalizedMessage{"hello": {Value: MessageText("Hallo")}}},
			{Locale: "en", Messages: map[string]LocalizedMessage{"hello": {Value: MessageText("Hello")}}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := app.Register(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(Message(r.Context(), "hello")))
	}); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("Accept-Language", "en-GB")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Body.String() != "Hello" {
		t.Fatalf("body = %q", response.Body.String())
	}
	if response.Header().Get("Content-Language") != "en" || !strings.Contains(response.Header().Get("Vary"), "Accept-Language") {
		t.Fatalf("localization headers = %+v", response.Header())
	}
}

func TestAppCanReplaceAndDisableLocalizer(t *testing.T) {
	app := New()
	localizer, err := NewLocalizer(I18nConfig{DefaultLocale: "en", Locales: []LocaleCatalog{{Locale: "en", Messages: map[string]LocalizedMessage{}}}})
	if err != nil {
		t.Fatal(err)
	}
	config := I18nConfig{DefaultLocale: "en", Locales: []LocaleCatalog{{Locale: "en"}}}
	if err := app.SetLocalizer(config, localizer); err != nil {
		t.Fatal(err)
	}
	if err := app.DisableI18n(); err != nil {
		t.Fatal(err)
	}
}

func TestAppCustomLocaleResolver(t *testing.T) {
	app := New()
	config := I18nConfig{DefaultLocale: "de", Detection: []string{"custom", "default"}, Locales: []LocaleCatalog{
		{Locale: "de", Messages: map[string]LocalizedMessage{}},
		{Locale: "ar", Messages: map[string]LocalizedMessage{}},
	}}
	if err := app.SetI18n(config); err != nil {
		t.Fatal(err)
	}
	if err := app.RegisterLocaleResolver(func(*http.Request) string { return "ar" }); err != nil {
		t.Fatal(err)
	}
	if err := app.Register(http.MethodGet, "/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(Locale(r.Context())))
	}); err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	app.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/", nil))
	if response.Body.String() != "ar" || response.Header().Get("Vary") != "*" {
		t.Fatalf("body = %q, headers = %+v", response.Body.String(), response.Header())
	}
}
