package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSafeReturnPath(t *testing.T) {
	cases := []struct {
		candidate string
		want      string
	}{
		{"/account?tab=profile", "/account?tab=profile"},
		{"https://evil.example", "/"},
		{"//evil.example", "/"},
		{"locale", "/"},
		{"/locale", "/"},
	}
	for _, tc := range cases {
		if got := safeReturnPath(tc.candidate, "/locale", "/"); got != tc.want {
			t.Errorf("safeReturnPath(%q) = %q, want %q", tc.candidate, got, tc.want)
		}
	}
}

func TestLocaleSelectionHandlerWorksWithoutJavaScript(t *testing.T) {
	app := New()
	if err := app.SetI18n(I18nConfig{DefaultLocale: "de", Locales: []LocaleCatalog{
		{Locale: "de", Messages: map[string]LocalizedMessage{}},
		{Locale: "en", Messages: map[string]LocalizedMessage{}},
	}}); err != nil {
		t.Fatal(err)
	}
	if err := app.Register(http.MethodPost, "/locale", LocaleSelectionHandler(LocaleSelectionOptions{})); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/locale", strings.NewReader("locale=en&return=%2Faccount"))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()
	app.ServeHTTP(response, request)
	if response.Code != http.StatusSeeOther || response.Header().Get("Location") != "/account" {
		t.Fatalf("status = %d, location = %q", response.Code, response.Header().Get("Location"))
	}
	cookies := response.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Value != "en" {
		t.Fatalf("cookies = %+v", cookies)
	}
}
