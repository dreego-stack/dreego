package gen

import (
	"net/http"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
)


func renderIndex(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

b.WriteString(`<title>Plugin demo</title>`)
	b.WriteString("<div data-scope=\"bd5740b4d68e\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Plugin demo`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`This app registers a local plugin package before the generated routes.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<ul>`)
	b.WriteString(`
        `)
	b.WriteString(`<li>`)
	b.WriteString(`<a href="/plugin/hello">`)
	b.WriteString(`Plugin hello`)
	b.WriteString(`</a>`)
	b.WriteString(`</li>`)
	b.WriteString(`
        `)
	b.WriteString(`<li>`)
	b.WriteString(`<a href="/plugin/hello/42">`)
	b.WriteString(`Plugin hello 42`)
	b.WriteString(`</a>`)
	b.WriteString(`</li>`)
	b.WriteString(`
        `)
	b.WriteString(`<li>`)
	b.WriteString(`<a href="/plugin/health">`)
	b.WriteString(`Plugin health`)
	b.WriteString(`</a>`)
	b.WriteString(`</li>`)
	b.WriteString(`
    `)
	b.WriteString(`</ul>`)
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
func Register(app *dreego.App) error {
	if err := app.SetLogging(false); err != nil {
		return err
	}
	if err := app.Register("GET", "/{$}", HandleIndex); err != nil {
		return err
	}
	return nil
}
