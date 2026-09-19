package server

import (
	"net/http"
	"testing"
)

func TestStateChangingRouteDetectsUnsafeMethods(t *testing.T) {
	cases := []struct {
		name    string
		method  string
		pattern string
		want    bool
	}{
		{"get", http.MethodGet, "/", false},
		{"head", http.MethodHead, "/", false},
		{"options", http.MethodOptions, "/", false},
		{"post", http.MethodPost, "/notes", true},
		{"put", http.MethodPut, "/notes", true},
		{"patch", http.MethodPatch, "/notes", true},
		{"delete", http.MethodDelete, "/notes", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			app := New()
			if err := app.Register(tc.method, tc.pattern, func(http.ResponseWriter, *http.Request) {}); err != nil {
				t.Fatalf("register: %v", err)
			}
			if got := app.stateChangingRoute(); got != tc.want {
				t.Errorf("stateChangingRoute() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStateChangingRouteEmpty(t *testing.T) {
	app := New()
	if app.stateChangingRoute() {
		t.Fatal("an app with no routes must not report a state-changing route")
	}
}
