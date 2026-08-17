package gen

import (
	"fmt"
	"net/http"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
)

func renderErrorroutes404(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

b.WriteString(`<title>Not found</title>`)
	b.WriteString("<div data-scope=\"9c85579232d7\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Page not found`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`The URL you requested does not exist.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Back home`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleErrorroutes404(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderErrorroutes404(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(404)
	w.Write([]byte(html))
}


func renderIndex(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

	message := "Hello from Dreego"

b.WriteString(`<title>Hello</title>
    <meta name="description" content="Minimal Dreego app">`)
	b.WriteString("<div data-scope=\"cfb89299fe5a\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(dreego.SafeText(fmt.Sprintf("%v", message)))
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`This is the minimal Dreego reference app.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<nav>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/about">`)
	b.WriteString(`About`)
	b.WriteString(`</a>`)
	b.WriteString(`
        `)
	b.WriteString(`<a href="/users/1">`)
	b.WriteString(`User 1`)
	b.WriteString(`</a>`)
	b.WriteString(`
    `)
	b.WriteString(`</nav>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderIndex(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func renderAbout(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

b.WriteString(`<title>About</title>`)
	b.WriteString("<div data-scope=\"736160029a0d\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`About this app`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`One route file per URL, one method per file.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Back home`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleAbout(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderAbout(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func renderUsersId(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

	id := c.Param("id")

b.WriteString(`<title>User</title>`)
	b.WriteString("<div data-scope=\"0cd0244e06d0\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`User `)
	b.WriteString(dreego.SafeText(fmt.Sprintf("%v", id)))
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Dynamic segments use [brackets] in the file name.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Back home`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleUsersId(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderUsersId(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
func Register(app *dreego.App) error {
	if err := app.SetLogging(false); err != nil {
		return err
	}
	if err := app.Register("", "/{p...}", HandleErrorroutes404); err != nil {
		return err
	}
	if err := app.Register("GET", "/{$}", HandleIndex); err != nil {
		return err
	}
	if err := app.Register("GET", "/about", HandleAbout); err != nil {
		return err
	}
	if err := app.Register("GET", "/users/{id}", HandleUsersId); err != nil {
		return err
	}
	return nil
}
