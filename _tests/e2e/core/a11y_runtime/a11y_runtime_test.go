package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// The scaffolded accessibility shell (skip link, main landmark, nav label) must
// survive in the rendered HTML.
func TestAccessibilityRuntimeAttrs(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json": `{"logging":{"enabled":false},"redirects":[],"rewrites":[]}`,
		"www/components/PageShell.dreego": `DREEFILE component (title string)

<body>
    <a href="#main" class="skip-link">skip to content</a>
    <header>
        <h1>{{ title }}</h1>
        <nav aria-label="Primary"><a href="/">Shop</a></nav>
    </header>
    <main id="main">{#slot}</main>
</body>`,
		"www/routes/+page.dreego": `<head>
    <title>Home</title>
</head>

<body>
    <@PageShell title="Home"><p>Page body</p></@PageShell>
</body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, `id="main"`) {
		t.Errorf("rendered HTML missing skip-link target id=\"main\": %s", body)
	}
	if !strings.Contains(body, "skip to content") {
		t.Errorf("rendered HTML missing skip-link text: %s", body)
	}
	if !strings.Contains(body, `aria-label="Primary"`) {
		t.Errorf("rendered HTML missing nav aria-label=\"Primary\": %s", body)
	}
}
