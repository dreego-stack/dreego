package core

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPluginRegistersMultipleRoutes verifies that a plugin can programmatically
// register several routes — GET, POST and a dynamic pattern — via core.Register
// inside RegisterRoutes, and that all of them are reachable through ServeMux()
// after UsePlugin + Build().
//
// This is the pragmatic core of route-hooks.1: because RegisterRoutes() calls
// core.Register at runtime, plugin routes land in the `routes` slice
// automatically and need no new codegen behavior.
func TestPluginRegistersMultipleRoutes(t *testing.T) {
	Reset()
	UsePlugin(&multiRoutePlugin{})

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
		ServeMux().ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("%s %s returned %d, want %d", c.method, c.path, rr.Code, http.StatusOK)
			continue
		}
		if got := rr.Body.String(); got != c.want {
			t.Errorf("%s %s body = %q, want %q", c.method, c.path, got, c.want)
		}
	}
}

// TestPluginRoutesSurviveReset documents that the `routes` slice intentionally
// outlives Reset() (matching the existing TestResetClearsCache semantics). A
// plugin route registered before Reset() is still served afterwards.
func TestPluginRoutesSurviveReset(t *testing.T) {
	Reset()
	UsePlugin(&multiRoutePlugin{})
	Reset()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/plugin/multi", nil)
	ServeMux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("plugin route after Reset() returned %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Body.String(); got != "multi-get" {
		t.Errorf("plugin route after Reset() body = %q, want %q", got, "multi-get")
	}
}

// TestPluginRouteLastWinsOverridesAppRoute pins down the route-overlap
// semantics: Register() is idempotent per method+pattern, so the LAST
// registration replaces the handler (last-wins) instead of appending a
// duplicate. When an app route is registered first and a plugin overrides the
// same pattern, the plugin handler wins.
func TestPluginRouteLastWinsOverridesAppRoute(t *testing.T) {
	Reset()

	Register("GET", "/shared", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("app"))
	})
	UsePlugin(&overlapRoutePlugin{path: "/shared", marker: "plugin"})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shared", nil)
	ServeMux().ServeHTTP(rr, req)

	if got := rr.Body.String(); got != "plugin" {
		t.Errorf("after plugin override, body = %q, want %q (last registration wins)", got, "plugin")
	}
}

// TestPluginRouteLastWinsAppWinsOverPlugin pins the reverse of the last-wins
// semantics: an app route registered after a plugin route replaces the plugin
// handler for the same method+pattern.
func TestPluginRouteLastWinsAppWinsOverPlugin(t *testing.T) {
	Reset()

	UsePlugin(&overlapRoutePlugin{path: "/shared2", marker: "plugin"})
	Register("GET", "/shared2", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("app"))
	})

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/shared2", nil)
	ServeMux().ServeHTTP(rr, req)

	if got := rr.Body.String(); got != "app" {
		t.Errorf("after app override, body = %q, want %q (last registration wins)", got, "app")
	}
}

// TestGenerateRouterRendersRouteInfo documents how "gen/dree.go collects
// plugin routes" is realized: GenerateRouter is a pure codegen function that
// emits registration code for whatever []RouteInfo it is handed. There is no
// plugin discovery at codegen time — the tooling must collect plugin RouteInfo
// entries (here the same ones a plugin's RegisterRoutes would use) and pass
// them in. This test pins that GenerateRouter renders such entries correctly.
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
