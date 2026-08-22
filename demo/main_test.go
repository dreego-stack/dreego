package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHostRouterSelectsSaaSByHost(t *testing.T) {
	public := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("public")) })
	product := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("saas")) })
	h := hostRouter(public, product)
	for _, tc := range []struct {
		host string
		want string
	}{
		{host: "localhost:8080", want: "public"},
		{host: "saas.localhost:8080", want: "saas"},
		{host: "saas.example.test:8080", want: "saas"},
		{host: "not-saas.localhost:8080", want: "public"},
	} {
		t.Run(tc.host, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://"+tc.host+"/", nil)
			r.Host = tc.host
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if got := w.Body.String(); got != tc.want {
				t.Fatalf("body = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHostRouterPreservesPathAndMethod(t *testing.T) {
	product := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/billing" {
			t.Errorf("request = %s %s, want POST /billing", r.Method, r.URL.Path)
		}
	})
	hostRouter(http.NotFoundHandler(), product).ServeHTTP(httptest.NewRecorder(), func() *http.Request {
		r := httptest.NewRequest(http.MethodPost, "http://saas.localhost/billing", nil)
		r.Host = "saas.localhost"
		return r
	}())
}
