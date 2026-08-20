package core

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func benchRenderPage(c *SSRContext) (string, error) {
	var b strings.Builder
	b.WriteString(`<div data-scope="abc">`)
	b.WriteString(`<h1>`)
	b.WriteString(html.EscapeString(fmt.Sprintf("%v", "Benchmark")))
	b.WriteString(`</h1>`)
	show := true
	if show {
		b.WriteString(`<p>visible</p>`)
	}
	items := []string{"a", "b", "c"}
	for i, item := range items {
		loop := EachLoop{Index: i, First: i == 0, Last: i == len(items)-1, Even: i%2 == 0, Odd: i%2 != 0}
		_ = loop
		b.WriteString(`<p>`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", item)))
		b.WriteString(`</p>`)
	}
	cardHTML, err := benchCard("Hello").Render(c)
	if err != nil {
		return "", err
	}
	b.WriteString(cardHTML)
	b.WriteString(`</div>`)
	return b.String(), nil
}

func benchCard(title string) Component {
	return ComponentFunc(func(ctx *SSRContext) (string, error) {
		var b strings.Builder
		b.WriteString(`<div data-scope="def">`)
		b.WriteString(`<h2>`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", title)))
		b.WriteString(`</h2>`)
		b.WriteString(`</div>`)
		return b.String(), nil
	})
}

func benchApp() *App {
	app := New()
	app.SetLogging(false)
	app.Register("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		c := NewSSR(w, r)
		html, err := benchRenderPage(c)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	})
	app.Register("GET", "/api", func(w http.ResponseWriter, r *http.Request) {
		NewSSR(w, r).JSON(200, map[string]string{"status": "ok"})
	})
	app.Build()
	return app
}

func BenchmarkRequestPage(b *testing.B) {
	h := benchApp().Handler()
	req := httptest.NewRequest("GET", "/", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	}
}

func BenchmarkRequestJSON(b *testing.B) {
	h := benchApp().Handler()
	req := httptest.NewRequest("GET", "/api", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	}
}

func BenchmarkRequestSimple(b *testing.B) {
	app := New()
	app.SetLogging(false)
	app.Register("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	app.Build()
	h := app.Handler()
	req := httptest.NewRequest("GET", "/", nil)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
	}
}
