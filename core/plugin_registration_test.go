package core

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type featureOptions struct {
	Path string
}

func registerFeature(app *App, options featureOptions, log *[]string) error {
	if err := app.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			*log = append(*log, "feature")
			next.ServeHTTP(w, r)
		})
	}); err != nil {
		return err
	}
	return app.Register(http.MethodGet, options.Path, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("registered"))
	})
}

func TestAppBoundRegistrationFunction(t *testing.T) {
	app := New()
	var log []string
	if err := registerFeature(app, featureOptions{Path: "/feature"}, &log); err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	app.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/feature", nil))
	if recorder.Body.String() != "registered" {
		t.Fatalf("body = %q, want registered", recorder.Body.String())
	}
	if !reflect.DeepEqual(log, []string{"feature"}) {
		t.Fatalf("middleware log = %v, want [feature]", log)
	}
}

func TestAppBoundRegistrationPropagatesLateError(t *testing.T) {
	app := New()
	app.Build()
	if err := registerFeature(app, featureOptions{Path: "/late"}, new([]string)); err == nil {
		t.Fatal("registration after Build must fail")
	}
}
