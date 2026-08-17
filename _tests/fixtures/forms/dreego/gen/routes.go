package gen

import (
	"fmt"
	"net/http"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
)

type EntryForm struct {
        Name    string `form:"name" validate:"required"`
        Message string `form:"message" validate:"required"`
    }

    var entries []string

    func AddEntry(c dreego.Context, form EntryForm) error {
        entries = append(entries, form.Name+": "+form.Message)
        return c.Redirect("/entries", 303)
    }

func renderIndex(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

b.WriteString(`<title>Guestbook</title>`)
	b.WriteString("<div data-scope=\"dcc7dffa4cd2\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Guestbook`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Post a message with the g-action form handler.`)
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<form g-action="AddEntry" method="post">`)
	b.WriteString(`
        `)
	b.WriteString(`<input type="hidden" name="csrf_token" value="` + dreego.SafeAttr(fmt.Sprintf("%v", c.CSRFToken())) + `">`)
	b.WriteString(`
        `)
	b.WriteString(`<label>`)
	b.WriteString(`Name
            `)
	b.WriteString(`<input name="name" type="text" value="` + dreego.SafeAttr(fmt.Sprintf("%v", c.Old("name"))) + `">`)
	b.WriteString(`
        `)
	b.WriteString(`</label>`)
	b.WriteString(`
        `)
	if c.Errors("name") != "" {
		b.WriteString(`<p class="error">`)
		b.WriteString(dreego.SafeText(fmt.Sprintf("%v", c.Errors("name"))))
		b.WriteString(`</p>`)
	}
	b.WriteString(`
        `)
	b.WriteString(`<label>`)
	b.WriteString(`Message
            `)
	b.WriteString(`<input name="message" type="text" value="` + dreego.SafeAttr(fmt.Sprintf("%v", c.Old("message"))) + `">`)
	b.WriteString(`
        `)
	b.WriteString(`</label>`)
	b.WriteString(`
        `)
	if c.Errors("message") != "" {
		b.WriteString(`<p class="error">`)
		b.WriteString(dreego.SafeText(fmt.Sprintf("%v", c.Errors("message"))))
		b.WriteString(`</p>`)
	}
	b.WriteString(`
        `)
	b.WriteString(`<button type="submit">`)
	b.WriteString(`Post`)
	b.WriteString(`</button>`)
	b.WriteString(`
    `)
	b.WriteString(`</form>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`<a href="/entries">`)
	b.WriteString(`View entries`)
	b.WriteString(`</a>`)
	b.WriteString(`</p>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleIndexGet(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderIndex(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}
func HandleIndexPost(app *dreego.App, w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)

	var form EntryForm
	if err := dreego.BindForm(r, &form); err != nil {
		c.Set("error__form", err.Error())
		html, _ := renderIndex(c)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
		return
	}

	errs := app.ValidateForm(form)
	if len(errs) > 0 {
		dreego.SaveErrors(c, errs)
		dreego.SaveOld(c, form)
		html, _ := renderIndex(c)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
		return
	}

	if err := AddEntry(c, form); err != nil {
		if err == dreego.ErrRedirect {
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.URL.Path, 303)
}


func renderCounter(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

	count := c.SessionVal("count")

b.WriteString(`<title>Counter</title>`)
	b.WriteString("<div data-scope=\"466f6563f3d8\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Counter`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`Count: `)
	b.WriteString(dreego.SafeText(fmt.Sprintf("%v", len(count))))
	b.WriteString(`</p>`)
	b.WriteString(`
    `)
	b.WriteString(`<form method="post">`)
	b.WriteString(`
        `)
	b.WriteString(`<input type="hidden" name="csrf_token" value="` + dreego.SafeAttr(fmt.Sprintf("%v", c.CSRFToken())) + `">`)
	b.WriteString(`
        `)
	b.WriteString(`<button type="submit">`)
	b.WriteString(`Increment`)
	b.WriteString(`</button>`)
	b.WriteString(`
    `)
	b.WriteString(`</form>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleCounter(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderCounter(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func renderCounterPOST(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

	n := c.SessionVal("count")
	c.SetSessionVal("count", n+"x")

	b.WriteString("<div data-scope=\"c3e18a00dd5d\">")
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`incremented`)
	b.WriteString(`</p>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleCounterPOST(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderCounterPOST(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(html))
}

func renderEntries(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

	empty := len(entries) == 0

b.WriteString(`<title>Entries</title>`)
	b.WriteString("<div data-scope=\"748ada11387f\">")
	b.WriteString(`
    `)
	b.WriteString(`<h1>`)
	b.WriteString(`Entries`)
	b.WriteString(`</h1>`)
	b.WriteString(`
    `)
	if empty {
		b.WriteString(`
        `)
		b.WriteString(`<p>`)
		b.WriteString(`No entries yet`)
		b.WriteString(`</p>`)
		b.WriteString(`
    `)
	} else {
		b.WriteString(`
        `)
		b.WriteString(`<ul>`)
		b.WriteString(`
        `)
			for i, entry := range entries {
				loop := dreego.EachLoop{Index: i, First: i == 0, Last: i == len(entries)-1, Even: i%2 == 0, Odd: i%2 != 0}
				_ = loop
				b.WriteString(`
            `)
				b.WriteString(`<li>`)
				b.WriteString(dreego.SafeText(fmt.Sprintf("%v", entry)))
				b.WriteString(`</li>`)
				b.WriteString(`
        `)
			}
		b.WriteString(`
        `)
		b.WriteString(`</ul>`)
		b.WriteString(`
    `)
	}
	b.WriteString(`
    `)
	b.WriteString(`<p>`)
	b.WriteString(`<a href="/">`)
	b.WriteString(`Post a message`)
	b.WriteString(`</a>`)
	b.WriteString(`</p>`)
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleEntries(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderEntries(c)
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
	if err := app.Register("GET", "/{$}", HandleIndexGet); err != nil {
		return err
	}
	if err := app.Register("POST", "/{$}", func(w http.ResponseWriter, r *http.Request) { HandleIndexPost(app, w, r) }); err != nil {
		return err
	}
	if err := app.Register("GET", "/counter", HandleCounter); err != nil {
		return err
	}
	if err := app.Register("POST", "/counter", HandleCounterPOST); err != nil {
		return err
	}
	if err := app.Register("GET", "/entries", HandleEntries); err != nil {
		return err
	}
	return nil
}
