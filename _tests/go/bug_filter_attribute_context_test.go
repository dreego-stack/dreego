package tests

import (
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// Feedback 3.4: |raw in an attribute must behave like |raw in text context.
// The generated code must compile (no "undefined: raw") and the raw value must
// bypass the scheme allowlist, so a server-built webcal link survives.
func TestBugFilterRawInAttributeContext(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego": `<server>webcal := "webcal://example.com/feed.ics"</server>
<body><a href="{{ webcal|raw }}">subscribe</a></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `href="webcal://example.com/feed.ics"`)
	dreegotest.MustNotContainBody(t, body, `href="#"`)
}

// Every filter must produce identical output in text and attribute contexts.
func TestBugFilterContextParity(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego": `<server>v := "MiXeD"</server>
<body>
<p id="t">{{ v|upper }}</p>
<a id="a" title="{{ v|upper }}" href="{{ v|upper }}">x</a>
</body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `<p id="t">MIXED</p>`)
	dreegotest.MustContainBody(t, body, `title="MIXED"`)
	dreegotest.MustContainBody(t, body, `href="MIXED"`)
}

// A server-built webcal URL is replaced with "#" by SafeURL unless |raw is
// used; this guards the documented escape hatch in every attribute context.
func TestBugAttributeRawCompilesWithUnknownFilter(t *testing.T) {
	t.Parallel()
	dreegotest.MustFailWith(t, `<body><a href="{{ v|nosuchfilter }}">x</a></body>`, "unknown filter 'nosuchfilter'")
}

// Feedback 3.4 was reported inside a component: the component code path has its
// own attribute scanner and must apply |raw the same way the route path does.
func TestBugFilterRawInComponentAttribute(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Subscribe.dreego": `DREEFILE component (feed string)
<body><a href="{{ feed|raw }}">component subscribe</a></body>`,
		"www/routes/+page.dreego": `<server>feed := "webcal://example.com/c.ics"</server>
<body><@Subscribe feed={feed}/></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `href="webcal://example.com/c.ics"`)
	dreegotest.MustNotContainBody(t, body, `href="#"`)
}

// Filter parity inside a component: raw and upper behave like the route path.
func TestBugComponentFilterContextParity(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/components/Badge.dreego": `DREEFILE component (label string, url string)
<body><span title="{{ label|upper }}">{{ label|upper }}</span><a href="{{ url|raw }}">go</a></body>`,
		"www/routes/+page.dreego": `<server>label := "ready"; url := "webcal://example.com/c.ics"</server>
<body><@Badge label={label} url={url}/></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `<span title="READY">READY</span>`)
	dreegotest.MustContainBody(t, body, `href="webcal://example.com/c.ics"`)
}
