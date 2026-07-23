package main

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/dreego/dreego/pkg/context"
)

const head_about = `<meta name="description" content="About page">`

const script_about = `console.log("about page loaded");`

const style_d3c26d9d = `h1 { color: blue; }`

func renderAbout(ctx *context.Context) (string, error) {
	var b strings.Builder

	name := "Welt"
	items := []string{"Apfel", "Banane", "Kirsche"}
	showDetails := true

	b.WriteString(head_about)
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Hallo `)
	b.WriteString(fmt.Sprintf("%v", name))
	b.WriteString(`!`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	if showDetails {
		b.WriteString(`
        `)
		b.WriteString(`<p>`)
		b.WriteString(`Willkommen auf der About-Seite.`)
		b.WriteString(`</p>`)
		b.WriteString(`
    `)
	}
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Meine Fruchte:`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<ul>`)
	b.WriteString(`
        `)
	for _, item := range items {
		b.WriteString(`
            `)
		b.WriteString(`<li>`)
		b.WriteString(fmt.Sprintf("%v", item))
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
	b.WriteString(`Home`)
	b.WriteString(`</a>`)
	b.WriteString(`
`)
	b.WriteString("<script>")
	b.WriteString(script_about)
	b.WriteString("</script>")
	b.WriteString("<style>")
	b.WriteString(style_d3c26d9d)
	b.WriteString("</style>")

	return b.String(), nil
}

func HandleAbout(w http.ResponseWriter, r *http.Request) {
	ctx := &context.Context{W: w, R: r}
	html, err := renderAbout(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func HandleAboutCSS() string {
	return style_d3c26d9d
}

func HandleAboutJS() string {
	return script_about
}
