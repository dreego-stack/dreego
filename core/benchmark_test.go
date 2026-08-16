package core

import (
	"fmt"
	"html"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const benchPageSrc = `<head>
    <title>Benchmark</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>

<go>
    title := "Benchmark"
    show := true
    items := []string{"a", "b", "c"}
</go>

<div class="page">
    <h1>{{ title }}</h1>
    {#if show}
        <p>visible</p>
    {/if}
    {#each items as item}
        <p>{{ item }}</p>
    {/each}
    <@Card title="Hello" />
</div>`

const benchComponentSrc = `Component Card (title string)

<div>
    <h2>{{ title }}</h2>
</div>`

func transpilePage(src string) (string, error) {
	_, imports, body := ParseHeader(src)
	tokens, err := Lex(body)
	if err != nil {
		return "", err
	}
	p := NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return "", err
	}
	file.Imports = imports
	file.SourceContent = src
	if len(file.Go) == 0 {
		file.Go = []GoSection{{Method: "GET"}}
	}
	for i := range file.Go {
		if !file.Go[i].MethodExplicit {
			file.Go[i].Method = "GET"
		}
	}
	gen := NewGenerator()
	out, _, err := GenerateMethodHandler(gen, file, nil, "routes", "index", "/{$}", "abc")
	return out, err
}

func BenchmarkGeneratePage(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := transpilePage(benchPageSrc); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGenerateComponent(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		comp, _, body := ParseHeader(benchComponentSrc)
		tokens, err := Lex(body)
		if err != nil {
			b.Fatal(err)
		}
		p := NewParser(tokens)
		file, err := p.Parse()
		if err != nil {
			b.Fatal(err)
		}
		file.Component = comp
		file.SourceContent = benchComponentSrc
		if len(file.Go) == 0 {
			file.Go = []GoSection{{Method: ""}}
		}
		gen := NewGenerator()
		if _, err := GenerateComponent(gen, file, "abc"); err != nil {
			b.Fatal(err)
		}
	}
}

func benchRenderPage(c *SSRContext) (string, error) {
	var b strings.Builder
	b.WriteString(`<div data-scope="abc">`)
	b.WriteString(`<h1>`)
	b.WriteString(html.EscapeString(fmt.Sprintf("%v", "Benchmark")))
	b.WriteString(`</h1>`)
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
