package gen

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"codeberg.org/dreego/dreego/pkg/context"
	"codeberg.org/dreego/dreego/pkg/runtime"
)

func renderErrorroutes404(c *context.SSRContext) (string, error) {
	var b strings.Builder

	b.WriteString(`<title>404 – Seite nicht gefunden</title>`)
	b.WriteString("<div data-scope=\"0410baf528e8\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`404`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Die Seite wurde nicht gefunden.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Zurück zur Startseite`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=0410baf528e8] h1 { color: #c00;}
[data-scope=0410baf528e8] a { color: #06f;}`)
	b.WriteString("</style>")

	return b.String(), nil
}

func HandleErrorroutes404(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderErrorroutes404(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(404)
	w.Write([]byte(html))
}

func init() {
	runtime.Register("", "/{p...}", HandleErrorroutes404)
}
func renderErrorroutes500(c *context.SSRContext) (string, error) {
	var b strings.Builder

	b.WriteString(`<title>500 – Interner Fehler</title>`)
	b.WriteString("<div data-scope=\"13d53c74b84a\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`500`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Ein interner Fehler ist aufgetreten.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Zurück zur Startseite`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=13d53c74b84a] h1 { color: #c00;}
[data-scope=13d53c74b84a] a { color: #06f;}`)
	b.WriteString("</style>")

	return b.String(), nil
}

func HandleErrorroutes500(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderErrorroutes500(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func init() {
	runtime.SetErrorHandler(500, HandleErrorroutes500)
}
func renderIndex(c *context.SSRContext) (string, error) {
	var b strings.Builder

	message := "Hello from Dreego!"

	b.WriteString("<div data-scope=\"961d4658a107\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(html.EscapeString(fmt.Sprintf("%v", message)))
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Welcome to the Dreego framework.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/about">`)
	b.WriteString(`About`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=961d4658a107] body { font-family: system-ui; max-width: 800px; margin: 2rem auto;}
[data-scope=961d4658a107] h1 { color: #333;}
[data-scope=961d4658a107] a { color: #06f;}`)
	b.WriteString("</style>")
	pageContent := b.String()
	b.Reset()
	c.Set("head", `<title>Dreego</title>`)
	c.Set("slot", pageContent)
	b.WriteString(`<!DOCTYPE html>`)
	b.WriteString(`
`)
	b.WriteString(`<html lang="de">`)
	b.WriteString(`
`)
	b.WriteString(`<head>`)
	b.WriteString(`
    `)
	b.WriteString(`<meta charset="utf-8">`)
	b.WriteString(`
    `)
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">`)
	b.WriteString(`
    `)
	b.WriteString(`<title>`)
	b.WriteString(`Dreego Demo`)
	b.WriteString(`</title>`)
	b.WriteString(`
    `)
	b.WriteString(c.Get("head"))
	b.WriteString(`
`)
	b.WriteString(`</head>`)
	b.WriteString(`
`)
	b.WriteString(`<body>`)
	b.WriteString(`
    `)
	b.WriteString(`<nav>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Home`)
	b.WriteString(`</a>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/about">`)
	b.WriteString(`About`)
	b.WriteString(`</a>`)
	b.WriteString(`
    `)
	b.WriteString(`</nav>`)
	b.WriteString(`
    `)
	b.WriteString(`<main>`)
	b.WriteString(`
        `)
	b.WriteString(c.Get("slot"))
	b.WriteString(`
    `)
	b.WriteString(`</main>`)
	b.WriteString(`
    `)
	b.WriteString(`<footer>`)
	b.WriteString(`
        `)
	b.WriteString(`<p>`)
	b.WriteString(`Powered by Dreego`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`</footer>`)
	b.WriteString(`
`)
	b.WriteString(`</body>`)
	b.WriteString(`
`)
	b.WriteString(`</html>`)
	b.WriteString(`
`)

	return b.String(), nil
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderIndex(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func init() {
	runtime.Register("GET", "/{$}", HandleIndex)
}
func renderGroupDemo(c *context.SSRContext) (string, error) {
	var b strings.Builder

	b.WriteString("<div data-scope=\"6223c6b49ea1\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Route Group: (group)`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Diese Seite liegt in `)
	b.WriteString(`<code>`)
	b.WriteString(`dreego/routes/(group)/demo/get.dreego`)
	b.WriteString(`</code>`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Die URL ist `)
	b.WriteString(`<strong>`)
	b.WriteString(`/demo`)
	b.WriteString(`</strong>`)
	b.WriteString(` — die Gruppe `)
	b.WriteString(`<code>`)
	b.WriteString(`(group)`)
	b.WriteString(`</code>`)
	b.WriteString(` erscheint nicht in der URL.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Zurück`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=6223c6b49ea1] code { background: #eee; padding: 0.2em 0.4em; border-radius: 3px;}`)
	b.WriteString("</style>")
	pageContent := b.String()
	b.Reset()
	c.Set("head", `<title>Route Group Demo</title>`)
	c.Set("slot", pageContent)
	b.WriteString(`<!DOCTYPE html>`)
	b.WriteString(`
`)
	b.WriteString(`<html lang="de">`)
	b.WriteString(`
`)
	b.WriteString(`<head>`)
	b.WriteString(`
    `)
	b.WriteString(`<meta charset="utf-8">`)
	b.WriteString(`
    `)
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">`)
	b.WriteString(`
    `)
	b.WriteString(`<title>`)
	b.WriteString(`Dreego Demo`)
	b.WriteString(`</title>`)
	b.WriteString(`
    `)
	b.WriteString(c.Get("head"))
	b.WriteString(`
`)
	b.WriteString(`</head>`)
	b.WriteString(`
`)
	b.WriteString(`<body>`)
	b.WriteString(`
    `)
	b.WriteString(`<nav>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Home`)
	b.WriteString(`</a>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/about">`)
	b.WriteString(`About`)
	b.WriteString(`</a>`)
	b.WriteString(`
    `)
	b.WriteString(`</nav>`)
	b.WriteString(`
    `)
	b.WriteString(`<main>`)
	b.WriteString(`
        `)
	b.WriteString(c.Get("slot"))
	b.WriteString(`
    `)
	b.WriteString(`</main>`)
	b.WriteString(`
    `)
	b.WriteString(`<footer>`)
	b.WriteString(`
        `)
	b.WriteString(`<p>`)
	b.WriteString(`Powered by Dreego`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`</footer>`)
	b.WriteString(`
`)
	b.WriteString(`</body>`)
	b.WriteString(`
`)
	b.WriteString(`</html>`)
	b.WriteString(`
`)

	return b.String(), nil
}

func HandleGroupDemo(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderGroupDemo(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func init() {
	runtime.Register("GET", "/demo", HandleGroupDemo)
}
func renderAbout(c *context.SSRContext) (string, error) {
	var b strings.Builder

	name := "Dreego"
	features := []string{"SSR-First", "File-based Routing", "Single Binary", "Layout-System", "5 Sektionen"}
	version := "0.0.1"

	b.WriteString("<div data-scope=\"80a846a20e2d\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Über `)
	b.WriteString(html.EscapeString(fmt.Sprintf("%v", name)))
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Version: `)
	b.WriteString(html.EscapeString(fmt.Sprintf("%v", version)))
	b.WriteString(`</p>`)
	b.WriteString(`

    `)
	b.WriteString(`<h2>`)
	b.WriteString(`Features`)
	b.WriteString(`</h2>`)
	b.WriteString(`
    `)
	b.WriteString(`<ul>`)
	b.WriteString(`
        `)
	for _, feature := range features {
		b.WriteString(`
            `)
		b.WriteString(`<li>`)
		b.WriteString(html.EscapeString(fmt.Sprintf("%v", feature)))
		b.WriteString(`</li>`)
		b.WriteString(`
        `)
	}
	b.WriteString(`
    `)
	b.WriteString(`</ul>`)
	b.WriteString(`

    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Zurück`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<script>")
	b.WriteString(`console.log("about page loaded");`)
	b.WriteString("</script>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=80a846a20e2d] h1 { color: #333;}
[data-scope=80a846a20e2d] h2 { color: #666; margin-top: 2rem;}
[data-scope=80a846a20e2d] li { margin: 0.5rem 0;}`)
	b.WriteString("</style>")
	pageContent := b.String()
	b.Reset()
	c.Set("head", `<meta name="description" content="About Dreego">`)
	c.Set("slot", pageContent)
	b.WriteString(`<!DOCTYPE html>`)
	b.WriteString(`
`)
	b.WriteString(`<html lang="de">`)
	b.WriteString(`
`)
	b.WriteString(`<head>`)
	b.WriteString(`
    `)
	b.WriteString(`<meta charset="utf-8">`)
	b.WriteString(`
    `)
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">`)
	b.WriteString(`
    `)
	b.WriteString(`<title>`)
	b.WriteString(`Dreego Demo`)
	b.WriteString(`</title>`)
	b.WriteString(`
    `)
	b.WriteString(c.Get("head"))
	b.WriteString(`
`)
	b.WriteString(`</head>`)
	b.WriteString(`
`)
	b.WriteString(`<body>`)
	b.WriteString(`
    `)
	b.WriteString(`<nav>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Home`)
	b.WriteString(`</a>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/about">`)
	b.WriteString(`About`)
	b.WriteString(`</a>`)
	b.WriteString(`
    `)
	b.WriteString(`</nav>`)
	b.WriteString(`
    `)
	b.WriteString(`<main>`)
	b.WriteString(`
        `)
	b.WriteString(c.Get("slot"))
	b.WriteString(`
    `)
	b.WriteString(`</main>`)
	b.WriteString(`
    `)
	b.WriteString(`<footer>`)
	b.WriteString(`
        `)
	b.WriteString(`<p>`)
	b.WriteString(`Powered by Dreego`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`</footer>`)
	b.WriteString(`
`)
	b.WriteString(`</body>`)
	b.WriteString(`
`)
	b.WriteString(`</html>`)
	b.WriteString(`
`)

	return b.String(), nil
}

func HandleAbout(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderAbout(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func init() {
	runtime.Register("GET", "/about", HandleAbout)
}
func renderHtmxdemo(c *context.SSRContext) (string, error) {
	var b strings.Builder

	b.WriteString("<div data-scope=\"ee3c813764ca\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`HTMX Demo`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<button hx-post="/htmxdemo" hx-target="#result">`)
	b.WriteString(`
        Klick mich
    `)
	b.WriteString(`</button>`)
	b.WriteString(`
    `)
	b.WriteString(`<div id="result">`)
	b.WriteString(`</div>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Zurück`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=ee3c813764ca] button { padding: 0.5rem 1rem; cursor: pointer; margin: 1rem 0;}
[data-scope=ee3c813764ca] #result { margin-top: 1rem; color: green;}`)
	b.WriteString("</style>")
	pageContent := b.String()
	b.Reset()
	c.Set("head", `<script src="https://unpkg.com/htmx.org@2.0.4"></script>
    <title>HTMX Demo</title>`)
	c.Set("slot", pageContent)
	b.WriteString(`<!DOCTYPE html>`)
	b.WriteString(`
`)
	b.WriteString(`<html lang="de">`)
	b.WriteString(`
`)
	b.WriteString(`<head>`)
	b.WriteString(`
    `)
	b.WriteString(`<meta charset="utf-8">`)
	b.WriteString(`
    `)
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">`)
	b.WriteString(`
    `)
	b.WriteString(`<title>`)
	b.WriteString(`Dreego Demo`)
	b.WriteString(`</title>`)
	b.WriteString(`
    `)
	b.WriteString(c.Get("head"))
	b.WriteString(`
`)
	b.WriteString(`</head>`)
	b.WriteString(`
`)
	b.WriteString(`<body>`)
	b.WriteString(`
    `)
	b.WriteString(`<nav>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Home`)
	b.WriteString(`</a>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/about">`)
	b.WriteString(`About`)
	b.WriteString(`</a>`)
	b.WriteString(`
    `)
	b.WriteString(`</nav>`)
	b.WriteString(`
    `)
	b.WriteString(`<main>`)
	b.WriteString(`
        `)
	b.WriteString(c.Get("slot"))
	b.WriteString(`
    `)
	b.WriteString(`</main>`)
	b.WriteString(`
    `)
	b.WriteString(`<footer>`)
	b.WriteString(`
        `)
	b.WriteString(`<p>`)
	b.WriteString(`Powered by Dreego`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`</footer>`)
	b.WriteString(`
`)
	b.WriteString(`</body>`)
	b.WriteString(`
`)
	b.WriteString(`</html>`)
	b.WriteString(`
`)

	return b.String(), nil
}

func HandleHtmxdemo(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderHtmxdemo(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func init() {
	runtime.Register("GET", "/htmxdemo", HandleHtmxdemo)
}
func renderHtmxdemoPOST(c *context.SSRContext) (string, error) {
	var b strings.Builder

	_ = fmt.Sprintf("")

	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`HTMX POST hat funktioniert!`)
	b.WriteString(`</p>`)
	b.WriteString(`
`)

	return b.String(), nil
}

func HandleHtmxdemoPOST(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderHtmxdemoPOST(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func init() {
	runtime.Register("POST", "/htmxdemo", HandleHtmxdemoPOST)
}
func renderErrorusers404(c *context.SSRContext) (string, error) {
	var b strings.Builder

	b.WriteString(`<title>404 – Benutzer nicht gefunden</title>`)
	b.WriteString("<div data-scope=\"9af6cfad8512\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`404`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Dieser Benutzer existiert nicht.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/users/1">`)
	b.WriteString(`Zum Beispiel-Benutzer`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=9af6cfad8512] h1 { color: #c00;}
[data-scope=9af6cfad8512] a { color: #06f;}`)
	b.WriteString("</style>")

	return b.String(), nil
}

func HandleErrorusers404(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderErrorusers404(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(404)
	w.Write([]byte(html))
}

func init() {
	runtime.Register("", "/users/{p...}", HandleErrorusers404)
}
func renderUsersId(c *context.SSRContext) (string, error) {
	var b strings.Builder

	userID := c.Param("id")

	b.WriteString("<div data-scope=\"11fde73a03bd\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Benutzer: `)
	b.WriteString(html.EscapeString(fmt.Sprintf("%v", userID)))
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Query ref: `)
	b.WriteString(html.EscapeString(fmt.Sprintf("%v", c.Query("ref"))))
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Zurück`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")
	b.WriteString("<style>")
	b.WriteString(`[data-scope=11fde73a03bd] h1 { color: #333;}`)
	b.WriteString("</style>")
	pageContent := b.String()
	b.Reset()
	c.Set("head", `<title>User {c.Param("id")}</title>`)
	c.Set("slot", pageContent)
	b.WriteString(`<!DOCTYPE html>`)
	b.WriteString(`
`)
	b.WriteString(`<html lang="de">`)
	b.WriteString(`
`)
	b.WriteString(`<head>`)
	b.WriteString(`
    `)
	b.WriteString(`<meta charset="utf-8">`)
	b.WriteString(`
    `)
	b.WriteString(`<meta name="viewport" content="width=device-width, initial-scale=1">`)
	b.WriteString(`
    `)
	b.WriteString(`<title>`)
	b.WriteString(`Dreego Demo`)
	b.WriteString(`</title>`)
	b.WriteString(`
    `)
	b.WriteString(c.Get("head"))
	b.WriteString(`
`)
	b.WriteString(`</head>`)
	b.WriteString(`
`)
	b.WriteString(`<body>`)
	b.WriteString(`
    `)
	b.WriteString(`<nav>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Home`)
	b.WriteString(`</a>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/about">`)
	b.WriteString(`About`)
	b.WriteString(`</a>`)
	b.WriteString(`
    `)
	b.WriteString(`</nav>`)
	b.WriteString(`
    `)
	b.WriteString(`<main>`)
	b.WriteString(`
        `)
	b.WriteString(c.Get("slot"))
	b.WriteString(`
    `)
	b.WriteString(`</main>`)
	b.WriteString(`
    `)
	b.WriteString(`<footer>`)
	b.WriteString(`
        `)
	b.WriteString(`<p>`)
	b.WriteString(`Powered by Dreego`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`</footer>`)
	b.WriteString(`
`)
	b.WriteString(`</body>`)
	b.WriteString(`
`)
	b.WriteString(`</html>`)
	b.WriteString(`
`)

	return b.String(), nil
}

func HandleUsersId(w http.ResponseWriter, r *http.Request) {
	c := context.NewSSR(w, r)
	html, err := renderUsersId(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func init() {
	runtime.Register("GET", "/users/{id}", HandleUsersId)
}
