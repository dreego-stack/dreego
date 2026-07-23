package main

import (
	"fmt"
	"net/http"
	"strings"

	"codeberg.org/dreego/dreego/pkg/context"
)

const head_index = `<title>Dreego</title>`

const style_b148a717 = `body { font-family: system-ui; max-width: 800px; margin: 2rem auto; }
    h1 { color: #333; }
    a { color: #06f; }`

func renderIndex(ctx *context.Context) (string, error) {
	var b strings.Builder

	message := "Hello from Dreego!"

	b.WriteString(head_index)
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(fmt.Sprintf("%v", message))
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
	b.WriteString("<style>")
	b.WriteString(style_b148a717)
	b.WriteString("</style>")

	return b.String(), nil
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	ctx := &context.Context{W: w, R: r}
	html, err := renderIndex(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func HandleIndexCSS() string {
	return style_b148a717
}
