package gen

import (
	"net/http"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
)


func renderIndex(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

	type Product struct {
	Name    string
	Price   string
	InStock bool
	}
	products := []Product{
	{Name: "Dreego Mug", Price: "$12", InStock: true},
	{Name: "Dreego Tee", Price: "$24", InStock: false},
	}

b.WriteString(`<title>Shop</title>`)
	b.WriteString("<div data-scope=\"c4e2c4969c6d\">")
	b.WriteString(`
    `)
	{
		previousSlot1 := c.Data("slot")
		var slotBuilder1 strings.Builder
			slotBuilder1.WriteString(`
        `)
			slotBuilder1.WriteString(`<div class="grid">`)
			slotBuilder1.WriteString(`
        `)
				for i, product := range products {
					loop := dreego.EachLoop{Index: i, First: i == 0, Last: i == len(products)-1, Even: i%2 == 0, Odd: i%2 != 0}
					_ = loop
					slotBuilder1.WriteString(`
            `)
					slotBuilder1.WriteString(func() string { h, _ := ProductCard(product.Name, product.Price, product.InStock).Render(c); return h }())
					slotBuilder1.WriteString(`
        `)
				}
			slotBuilder1.WriteString(`
        `)
			slotBuilder1.WriteString(`</div>`)
			slotBuilder1.WriteString(`
    `)
		c.Set("slot", slotBuilder1.String())
		html, err := PageShell("Welcome to the shop").Render(c)
		if previousSlot1 == nil { c.Delete("slot") } else { c.Set("slot", previousSlot1) }
		if err != nil { return "", err }
		b.WriteString(html)
	}
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

func renderProductsId(c *dreego.SSRContext) (string, error) {
	var b strings.Builder

	type Product struct {
	Name    string
	Price   string
	InStock bool
	}
	products := []Product{
	{Name: "Dreego Mug", Price: "$12", InStock: true},
	{Name: "Dreego Tee", Price: "$24", InStock: false},
	}
	id := c.Param("id")
	var product Product
	if id == "1" {
	product = products[0]
	} else {
	product = products[1]
	}

b.WriteString(`<title>Product</title>`)
	b.WriteString("<div data-scope=\"12b2cce79d99\">")
	b.WriteString(`
    `)
	{
		previousSlot1 := c.Data("slot")
		var slotBuilder1 strings.Builder
			slotBuilder1.WriteString(`
        `)
			slotBuilder1.WriteString(func() string { h, _ := ProductCard(product.Name, product.Price, product.InStock).Render(c); return h }())
			slotBuilder1.WriteString(`
        `)
			slotBuilder1.WriteString(`<p>`)
			slotBuilder1.WriteString(`<a href="/">`)
			slotBuilder1.WriteString(`Back to shop`)
			slotBuilder1.WriteString(`</a>`)
			slotBuilder1.WriteString(`</p>`)
			slotBuilder1.WriteString(`
    `)
		c.Set("slot", slotBuilder1.String())
		html, err := PageShell("Product detail").Render(c)
		if previousSlot1 == nil { c.Delete("slot") } else { c.Set("slot", previousSlot1) }
		if err != nil { return "", err }
		b.WriteString(html)
	}
	b.WriteString(`
`)
	b.WriteString("</div>")

	return b.String(), nil
}

func HandleProductsId(w http.ResponseWriter, r *http.Request) {
	c := dreego.NewSSR(w, r)
	html, err := renderProductsId(c)
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
	if err := app.Register("GET", "/products/{id}", HandleProductsId); err != nil {
		return err
	}
	return nil
}
