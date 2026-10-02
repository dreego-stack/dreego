package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// A static file whose extension is known to the Go standard library but absent
// from Dreego's explicit table (for example robots.txt, sitemap.xml, llms.txt)
// must be served with the standard Content-Type, not application/octet-stream.
func TestBugStaticFileServesStandardContentType(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/app/static/robots.txt":   "User-agent: *\nAllow: /\n",
		"www/app/static/sitemap.xml":  `<?xml version="1.0" encoding="UTF-8"?><urlset/>`,
		"www/app/routes/+page.dreego": `<body><p>home</p></body>`,
	})

	code, _, headers := c.Request(t, "GET", "/robots.txt", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustHeader(t, headers, "Content-Type", "text/plain; charset=utf-8")

	code, _, headers = c.Request(t, "GET", "/sitemap.xml", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustHeader(t, headers, "Content-Type", "text/xml; charset=utf-8")
}

// The explicit table still wins for extensions the host MIME database does not
// map consistently (fonts, icons).
func TestStaticFileServesPinnedContentType(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/app/static/font.woff2":   "FONT",
		"www/app/static/favicon.ico":  "ICON",
		"www/app/static/styles.css":   "body{}",
		"www/app/routes/+page.dreego": `<body><p>home</p></body>`,
	})

	for path, want := range map[string]string{
		"/font.woff2":  "font/woff2",
		"/favicon.ico": "image/x-icon",
		"/styles.css":  "text/css",
	} {
		code, _, headers := c.Request(t, "GET", path, "", nil)
		dreegotest.MustStatus(t, code, 200)
		if got := headers.Get("Content-Type"); !strings.HasPrefix(got, want) {
			t.Fatalf("%s Content-Type = %q, want prefix %q", path, got, want)
		}
	}
}

// An unknown extension still falls back to application/octet-stream.
func TestStaticFileUnknownExtensionFallsBack(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"www/app/static/data.notareal": "bytes",
		"www/app/routes/+page.dreego":  `<body><p>home</p></body>`,
	})

	code, _, headers := c.Request(t, "GET", "/data.notareal", "", nil)
	dreegotest.MustStatus(t, code, 200)
	if got := headers.Get("Content-Type"); got != "application/octet-stream" {
		t.Fatalf("unknown extension Content-Type = %q, want application/octet-stream", got)
	}
}
