package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestReferenceComponents(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json": `{"logging":{"enabled":false},"redirects":[],"rewrites":[]}`,
		"catalog/catalog.go":     "package catalog\n\ntype Product struct {\n\tName    string\n\tPrice   string\n\tInStock bool\n}\n\nfunc Products() []Product {\n\treturn []Product{\n\t\t{Name: \"Dreego Mug\", Price: \"$12\", InStock: true},\n\t\t{Name: \"Dreego Tee\", Price: \"$24\", InStock: false},\n\t}\n}\n",
		"www/components/PageShell.dreego": `DREEFILE component (title string)

<body>
    <a href="#main" class="skip-link">skip to content</a>
    <header>
        <h1>{{ title }}</h1>
        <nav aria-label="Primary"><a href="/">Shop</a></nav>
    </header>
    <main id="main">{#slot}</main>
</body>

<style>
.skip-link { position: absolute; left: -9999px; }
.skip-link:focus { position: fixed; top: 0.5rem; left: 0.5rem; z-index: 100; padding: 0.5rem 1rem; background: #000; color: #fff; }
</style>`,
		"www/components/ProductCard.dreego": `DREEFILE component (name string, price string, inStock bool)

<body>
    <article class="product-card">
        <h2>{{ name }}</h2>
        <p class="price">{{ price }}</p>
        {#if inStock}<p class="badge">In stock</p>{#else}<p class="badge">Sold out</p>{/if}
    </article>
</body>

<style>
.product-card { border: 1px solid #e2e8f0; padding: 1rem; border-radius: 8px; }
.badge { font-weight: bold; }
</style>`,
		"www/routes/+page.dreego": `GOIMPORT { catalog "t/catalog" }

<head>
    <title>Shop</title>
</head>

<server>
    products := catalog.Products()
</server>

<body>
    <@PageShell title="Welcome to the shop">
        <div class="grid">
        {#each products as product}
            <@ProductCard name={product.Name} price={product.Price} inStock={product.InStock}/>
        {/each}
        </div>
    </@PageShell>
</body>`,
		"www/routes/products/[id]/+page.dreego": `GOIMPORT { catalog "t/catalog" }

<head>
    <title>Product</title>
</head>

<server>
    products := catalog.Products()
    id := c.Param("id")
    var product catalog.Product
    if id == "1" {
        product = products[0]
    } else {
        product = products[1]
    }
</server>

<body>
    <@PageShell title="Product detail">
        <@ProductCard name={product.Name} price={product.Price} inStock={product.InStock}/>
        <p><a href="/">Back to shop</a></p>
    </@PageShell>
</body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Welcome to the shop") {
		t.Fatalf("shop page missing heading: %s", body)
	}
	if !strings.Contains(body, "Dreego Mug") {
		t.Fatalf("product card missing name: %s", body)
	}
	if !strings.Contains(body, "data-scope=") {
		t.Fatalf("component scoped style missing data-scope: %s", body)
	}
	if !strings.Contains(body, "In stock") {
		t.Fatalf("in-stock badge missing: %s", body)
	}
	if !strings.Contains(body, "Sold out") {
		t.Fatalf("sold-out badge missing: %s", body)
	}
	code, body = c.Get(t, "/products/1")
	if code != 200 {
		t.Fatalf("GET /products/1 = %d, want 200", code)
	}
	if !strings.Contains(body, "Dreego Mug") {
		t.Fatalf("product detail missing name: %s", body)
	}
	code, body = c.Get(t, "/products/2")
	if code != 200 {
		t.Fatalf("GET /products/2 = %d, want 200", code)
	}
	if !strings.Contains(body, "Dreego Tee") {
		t.Fatalf("product detail missing name: %s", body)
	}
}

