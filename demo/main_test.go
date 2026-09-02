package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHostRouterSelectsSaaSByHost(t *testing.T) {
	public := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("public")) })
	product := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("saas")) })
	blog := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("blog")) })
	h := hostRouter(public, product, blog)
	for _, tc := range []struct {
		host string
		want string
	}{
		{host: "localhost:8080", want: "public"},
		{host: "saas.localhost:8080", want: "saas"},
		{host: "saas.example.test:8080", want: "saas"},
		{host: "blog.localhost:8080", want: "blog"},
		{host: "blog.example.test:8080", want: "blog"},
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

func TestHostRouterSelectsByPath(t *testing.T) {
	public := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("public")) })
	product := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("saas")) })
	blog := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("blog")) })
	h := hostRouter(public, product, blog)
	for _, tc := range []struct {
		path string
		want string
	}{
		{path: "/", want: "public"},
		{path: "/blog/", want: "blog"},
		{path: "/blog/posts/hello-dreego", want: "blog"},
		{path: "/saas/", want: "saas"},
		{path: "/saas/dashboard", want: "saas"},
		{path: "/bloggy", want: "public"},
		{path: "/saasify", want: "public"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://localhost:8080"+tc.path, nil)
			r.Host = "localhost:8080"
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if got := w.Body.String(); got != tc.want {
				t.Fatalf("body = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHostRouterRedirectsWithoutTrailingSlash(t *testing.T) {
	public := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("public")) })
	product := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("saas")) })
	blog := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("blog")) })
	h := hostRouter(public, product, blog)
	for _, tc := range []struct {
		path string
		want string
	}{
		{path: "/blog", want: "/blog/"},
		{path: "/saas", want: "/saas/"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "http://localhost:8080"+tc.path, nil)
			r.Host = "localhost:8080"
			w := httptest.NewRecorder()
			h.ServeHTTP(w, r)
			if w.Code != http.StatusMovedPermanently {
				t.Fatalf("code = %d, want %d", w.Code, http.StatusMovedPermanently)
			}
			if got := w.Header().Get("Location"); got != tc.want {
				t.Fatalf("Location = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestHostRouterStripsPathPrefix(t *testing.T) {
	blog := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/posts/hello-dreego" {
			t.Errorf("path = %q, want /posts/hello-dreego", r.URL.Path)
		}
	})
	hostRouter(http.NotFoundHandler(), http.NotFoundHandler(), blog).ServeHTTP(httptest.NewRecorder(), func() *http.Request {
		r := httptest.NewRequest(http.MethodGet, "http://localhost:8080/blog/posts/hello-dreego", nil)
		r.Host = "localhost:8080"
		return r
	}())
}

func TestHostRouterPreservesPathAndMethod(t *testing.T) {
	product := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/billing" {
			t.Errorf("request = %s %s, want POST /billing", r.Method, r.URL.Path)
		}
	})
	hostRouter(http.NotFoundHandler(), product, http.NotFoundHandler()).ServeHTTP(httptest.NewRecorder(), func() *http.Request {
		r := httptest.NewRequest(http.MethodPost, "http://saas.localhost/billing", nil)
		r.Host = "saas.localhost"
		return r
	}())
}
