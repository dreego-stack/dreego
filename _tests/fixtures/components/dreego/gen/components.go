package gen

import (
	"fmt"
	"strings"

	dreego "github.com/dreego-stack/dreego/core"
)

func PageShell(title string) dreego.Component {
	return dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {
		var b strings.Builder

		b.WriteString("<div data-scope=\"81d31ce60019\">")
		b.WriteString(`
    `)
		b.WriteString(`<header>`)
		b.WriteString(`
        `)
		b.WriteString(`<h1>`)
		b.WriteString(dreego.SafeText(fmt.Sprintf("%v", title)))
		b.WriteString(`</h1>`)
		b.WriteString(`
        `)
		b.WriteString(`<nav>`)
		b.WriteString(`<a href="/">`)
		b.WriteString(`Shop`)
		b.WriteString(`</a>`)
		b.WriteString(`</nav>`)
		b.WriteString(`
    `)
		b.WriteString(`</header>`)
		b.WriteString(`
    `)
		b.WriteString(`<main>`)
		b.WriteString(ctx.Get("slot"))
		b.WriteString(`</main>`)
		b.WriteString(`
`)
		b.WriteString("</div>")

		return b.String(), nil
	})
}

func ProductCard(name string, price string, inStock bool) dreego.Component {
	return dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {
		var b strings.Builder

		b.WriteString("<div data-scope=\"4477f389c981\">")
		b.WriteString(`
    `)
		b.WriteString(`<article class="product-card">`)
		b.WriteString(`
        `)
		b.WriteString(`<h2>`)
		b.WriteString(dreego.SafeText(fmt.Sprintf("%v", name)))
		b.WriteString(`</h2>`)
		b.WriteString(`
        `)
		b.WriteString(`<p class="price">`)
		b.WriteString(dreego.SafeText(fmt.Sprintf("%v", price)))
		b.WriteString(`</p>`)
		b.WriteString(`
        `)
		if inStock {
		b.WriteString(`
            `)
		b.WriteString(`<p class="badge">`)
		b.WriteString(`In stock`)
		b.WriteString(`</p>`)
		b.WriteString(`
        `)
	} else {
		b.WriteString(`
            `)
		b.WriteString(`<p class="badge">`)
		b.WriteString(`Sold out`)
		b.WriteString(`</p>`)
		b.WriteString(`
        `)
	}
		b.WriteString(`
    `)
		b.WriteString(`</article>`)
		b.WriteString(`
`)
		b.WriteString("</div>")
		b.WriteString("<style>")
		b.WriteString(`[data-scope=4477f389c981] .product-card{ border: 1px solid #e2e8f0; padding: 1rem; border-radius: 8px; }[data-scope=4477f389c981] .badge{ font-weight: bold; }`)
		b.WriteString("</style>")

		return b.String(), nil
	})
}

