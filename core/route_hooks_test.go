package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPluginRegistersMultipleRoutes(t *testing.T) {
	app := New()
	if err := registerMultiRoutes(app); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		method string
		path   string
		want   string
	}{
		{"GET", "/plugin/multi", "multi-get"},
		{"POST", "/plugin/multi", "multi-post"},
		{"GET", "/plugin/multi/42", "multi-get-42"},
	}
	for _, c := range cases {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(c.method, c.path, nil)
		app.Handler().ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("%s %s returned %d, want %d", c.method, c.path, rr.Code, http.StatusOK)
			continue
		}
		if got := rr.Body.String(); got != c.want {
			t.Errorf("%s %s body = %q, want %q", c.method, c.path, got, c.want)
		}
	}
}

func TestPluginRoutesInNewApp(t *testing.T) {
	app := New()
	if err := registerMultiRoutes(app); err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin/multi", nil)
	app.Handler().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("plugin route returned %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Body.String(); got != "multi-get" {
		t.Errorf("plugin route body = %q, want %q", got, "multi-get")
	}
}

func TestRegistrationFunctionCannotOverrideAppRoute(t *testing.T) {
	app := New()

	app.Register("GET", "/shared", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("app"))
	})
	if err := registerOverlapRoute(app, "/shared", "plugin"); err == nil {
		t.Fatal("duplicate plugin route must fail")
	}
}

func TestAppCannotOverrideRegistrationFunctionRoute(t *testing.T) {
	app := New()

	if err := registerOverlapRoute(app, "/shared2", "plugin"); err != nil {
		t.Fatal(err)
	}
	if err := app.Register("GET", "/shared2", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("app"))
	}); err == nil {
		t.Fatal("duplicate app route must fail")
	}
}

func registerMultiRoutes(app *App) error {
	routes := []struct {
		method  string
		pattern string
		handler http.HandlerFunc
	}{
		{"GET", "/plugin/multi", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("multi-get")) }},
		{"POST", "/plugin/multi", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("multi-post")) }},
		{"GET", "/plugin/multi/{id}", func(w http.ResponseWriter, r *http.Request) { _, _ = w.Write([]byte("multi-get-" + r.PathValue("id"))) }},
	}
	for _, route := range routes {
		if err := app.Register(route.method, route.pattern, route.handler); err != nil {
			return err
		}
	}
	return nil
}

func registerOverlapRoute(app *App, path, marker string) error {
	return app.Register("GET", path, func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(marker))
	})
}

func TestGenerateRouterRendersRouteInfo(t *testing.T) {
	routes := []RouteInfo{
		{HandlerName: "HandleAdmin", RoutePath: "/admin", Method: "GET"},
		{HandlerName: "HandleAuthLogin", RoutePath: "/api/auth/login", Method: "POST"},
	}
	got := GenerateRouter(routes)

	for _, want := range []string{
		`mux.HandleFunc("GET /admin", HandleAdmin)`,
		`mux.HandleFunc("POST /api/auth/login", HandleAuthLogin)`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("GenerateRouter output missing %q\n--- got ---\n%s", want, got)
		}
	}
}
