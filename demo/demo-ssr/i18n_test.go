package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"demo/www"
	dreego "github.com/dreego-stack/dreego/core"
)

func TestPublicDemoSwitchesLanguageWithoutChangingRoute(t *testing.T) {
	app := dreego.New()
	if err := configure(app); err != nil {
		t.Fatal(err)
	}
	if err := www.Register(app); err != nil {
		t.Fatal(err)
	}
	if err := registerLocaleSelection(app); err != nil {
		t.Fatal(err)
	}

	germanRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	germanRequest.Header.Set("Accept-Language", "de-DE,de;q=0.9")
	germanResponse := httptest.NewRecorder()
	app.ServeHTTP(germanResponse, germanRequest)
	if germanResponse.Code != http.StatusOK {
		t.Fatalf("German response status = %d", germanResponse.Code)
	}
	if body := germanResponse.Body.String(); !strings.Contains(body, "Kleine Apps. Echte Routen. Reines Go.") ||
		!strings.Contains(body, `<html lang="de">`) ||
		!strings.Contains(body, `aria-label="Sprache auswählen"`) ||
		!strings.Contains(body, "4 Anwendungen") {
		t.Fatalf("German page is not localized: %s", body)
	}
	if germanResponse.Header().Get("Content-Language") != "de" {
		t.Fatalf("German Content-Language = %q", germanResponse.Header().Get("Content-Language"))
	}

	form := url.Values{"locale": {"en"}, "return": {"/"}}
	selectionRequest := httptest.NewRequest(http.MethodPost, "/locale", strings.NewReader(form.Encode()))
	selectionRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	selectionResponse := httptest.NewRecorder()
	app.ServeHTTP(selectionResponse, selectionRequest)
	if selectionResponse.Code != http.StatusSeeOther || selectionResponse.Header().Get("Location") != "/" {
		t.Fatalf("selection response = %d, location = %q", selectionResponse.Code, selectionResponse.Header().Get("Location"))
	}

	englishRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	englishRequest.Header.Set("Accept-Language", "de-DE,de;q=0.9")
	for _, cookie := range selectionResponse.Result().Cookies() {
		englishRequest.AddCookie(cookie)
	}
	englishResponse := httptest.NewRecorder()
	app.ServeHTTP(englishResponse, englishRequest)
	if body := englishResponse.Body.String(); !strings.Contains(body, "Small apps. Real routes. Plain Go.") || !strings.Contains(body, `<html lang="en">`) {
		t.Fatalf("English page is not localized after selection: %s", body)
	}
}
