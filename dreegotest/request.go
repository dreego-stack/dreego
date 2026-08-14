package dreegotest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	dreego "github.com/dreego-stack/dreego/core"
)

type Response struct {
	StatusCode int
	Body       string
}

func Get(t *testing.T, path string) *Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	return serve(req)
}

func PostForm(t *testing.T, path string, form url.Values) *Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return serve(req)
}

func serve(req *http.Request) *Response {
	rec := httptest.NewRecorder()
	testApp.Handler().ServeHTTP(rec, req)
	return &Response{StatusCode: rec.Code, Body: rec.Body.String()}
}

var testApp = dreego.New()

func Register(method, pattern string, handler http.HandlerFunc) {
	testApp.Register(method, pattern, handler)
}

func Reset() {
	testApp = dreego.New()
}
