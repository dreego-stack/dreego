package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// A <server type="custom"> route writes its own Content-Type. The generated GET
// handler must not overwrite that header with the HTML default after render
// returns. With gzip compression the response headers are still buffered when
// the handler runs, so the clobber is observable on the wire.
func TestBugCustomRouteContentTypeNotClobbered(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/calendar/de.ics.dreego": `<server type="custom">
    c.W.Header().Set("Content-Type", "text/calendar; charset=utf-8")
    c.W.WriteHeader(200)
</server>
<body>BEGIN:VCALENDAR
END:VCALENDAR</body>`,
	})
	code, body, headers := c.Request(t, "GET", "/calendar/de.ics", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "BEGIN:VCALENDAR")
	got := headers.Get("Content-Type")
	if !strings.HasPrefix(got, "text/calendar") {
		t.Fatalf("custom route Content-Type = %q, want text/calendar prefix", got)
	}
}

// c.Write sets the Content-Type directly; the generated handler must preserve
// it for custom routes.
func TestBugCustomRouteWriteContentTypeNotClobbered(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/calendar/ru.ics.dreego": `<server type="custom">
    c.Write(200, "text/calendar", []byte("BEGIN:VCALENDAR\r\nEND:VCALENDAR\r\n"))
</server>`,
	})
	code, body, headers := c.Request(t, "GET", "/calendar/ru.ics", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "BEGIN:VCALENDAR")
	got := headers.Get("Content-Type")
	if !strings.HasPrefix(got, "text/calendar") {
		t.Fatalf("custom route Content-Type = %q, want text/calendar prefix", got)
	}
	if strings.Count(got, "charset") > 1 {
		t.Fatalf("custom route Content-Type has duplicate charset: %q", got)
	}
}

// A plain HTML route must still get the HTML default Content-Type.
func TestBugHTMLRouteContentTypeDefault(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego": `<server>
    msg := "hello"
</server>
<body><h1>{{ msg }}</h1></body>`,
	})
	code, _, headers := c.Request(t, "GET", "/", "", nil)
	dreegotest.MustStatus(t, code, 200)
	if got := headers.Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Fatalf("HTML route Content-Type = %q, want text/html prefix", got)
	}
}
