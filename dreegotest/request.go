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
	Header     http.Header
	Cookies    []*http.Cookie
}

// Location returns the Location header value of the response.
func (r *Response) Location() string {
	if r == nil || r.Header == nil {
		return ""
	}
	return r.Header.Get("Location")
}

type AppClient struct {
	handler http.Handler
}

func NewApp(app *dreego.App) *AppClient {
	return &AppClient{handler: app.Handler()}
}

func (c *AppClient) Get(t *testing.T, path string) *Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	return c.serve(req)
}

func (c *AppClient) PostForm(t *testing.T, path string, form url.Values) *Response {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return c.serve(req)
}

func (c *AppClient) serve(req *http.Request) *Response {
	rec := httptest.NewRecorder()
	c.handler.ServeHTTP(rec, req)
	return &Response{
		StatusCode: rec.Code,
		Body:       rec.Body.String(),
		Header:     rec.Header(),
		Cookies:    rec.Result().Cookies(),
	}
}
