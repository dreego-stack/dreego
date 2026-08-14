package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPluginRegistersMultipleRoutes(t *testing.T) {
	app := New()
	app.UsePlugin(&multiRoutePlugin{})

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
	app.UsePlugin(&multiRoutePlugin{})

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

func TestPluginRouteLastWinsOverridesAppRoute(t *testing.T) {
	app := New()

	app.Register("GET", "/shared", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("app"))
	})
	app.UsePlugin(&overlapRoutePlugin{path: "/shared", marker: "plugin"})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shared", nil)
	app.Handler().ServeHTTP(rr, req)

	if got := rr.Body.String(); got != "plugin" {
		t.Errorf("after plugin override, body = %q, want %q (last registration wins)", got, "plugin")
	}
}

func TestPluginRouteLastWinsAppWinsOverPlugin(t *testing.T) {
	app := New()

	app.UsePlugin(&overlapRoutePlugin{path: "/shared2", marker: "plugin"})
	app.Register("GET", "/shared2", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("app"))
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shared2", nil)
	app.Handler().ServeHTTP(rr, req)

	if got := rr.Body.String(); got != "app" {
		t.Errorf("after app override, body = %q, want %q (last registration wins)", got, "app")
	}
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
