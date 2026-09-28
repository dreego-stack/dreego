package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestReferenceForms(t *testing.T) {
	t.Parallel()
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/dreego.config.json": `{"logging":{"enabled":false},"redirects":[],"rewrites":[]}`,
		"store/store.go":         "package store\n\nimport \"sync\"\n\nvar (\n\tmu      sync.Mutex\n\tentries []string\n)\n\nfunc Add(entry string) {\n\tmu.Lock()\n\tdefer mu.Unlock()\n\tentries = append(entries, entry)\n}\n\nfunc All() []string {\n\tmu.Lock()\n\tdefer mu.Unlock()\n\treturn append([]string(nil), entries...)\n}\n",
		"www/routes/+page.dreego": `GOIMPORT { store "t/store" }

<head>
    <title>Guestbook</title>
</head>

<server>
    type EntryForm struct {
        Name    string ` + "`form:\"name\" validate:\"required\"`" + `
        Message string ` + "`form:\"message\" validate:\"required\"`" + `
    }

    func AddEntry(c dreego.Context, form EntryForm) error {
        store.Add(form.Name + ": " + form.Message)
        return c.Redirect("/entries", 303)
    }
</server>

<body>
    <h1>Guestbook</h1>
    <form g-action="AddEntry" method="post">
        <input type="hidden" name="csrf_token" value="{{ c.CSRFToken() }}">
        <label>Name <input name="name" type="text" value="{{ c.Old("name") }}"></label>
        <label>Message <input name="message" type="text" value="{{ c.Old("message") }}"></label>
        <button type="submit">Post</button>
    </form>
    <p><a href="/entries">View entries</a></p>
</body>`,
		"www/routes/entries/+page.dreego": `GOIMPORT { store "t/store" }

<head>
    <title>Entries</title>
</head>

<server>
    entries := store.All()
    empty := len(entries) == 0
</server>

<body>
    <h1>Entries</h1>
    {#if empty}
        <p>No entries yet</p>
    {#else}
        <ul>{#each entries as entry}<li>{{ entry }}</li>{/each}</ul>
    {/if}
    <p><a href="/">Post a message</a></p>
</body>`,
		"www/routes/counter/+page.dreego": `<head>
    <title>Counter</title>
</head>

<server>
    count := c.SessionVal("count")
</server>

<body>
    <h1>Counter</h1>
    <p>Count: {{ len(count) }}</p>
    <form method="post">
        <input type="hidden" name="csrf_token" value="{{ c.CSRFToken() }}">
        <button type="submit">Increment</button>
    </form>
</body>

<server method="post">
    n := c.SessionVal("count")
    c.SetSessionVal("count", n+"x")
</server>
<body method="post"><p>incremented</p></body>`,
	}, "store := dreego.NewCookieStore([]byte(\"reference-apps-secret-key-32-bytes!\")); if err := app.SetSessionStore(store); err != nil { panic(err) }; ")

	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Guestbook") {
		t.Fatalf("guestbook page missing heading: %s", body)
	}
	code, body = c.Get(t, "/entries")
	if code != 200 {
		t.Fatalf("GET /entries = %d, want 200", code)
	}
	if !strings.Contains(body, "No entries yet") {
		t.Fatalf("empty entries page missing message: %s", body)
	}
	token := c.Cookie("csrf_token")
	if token == "" {
		t.Fatalf("no csrf_token cookie issued")
	}
	code, _, headers := c.Request(t, "POST", "/", "name=Ada&message=Hello+world&csrf_token="+token, map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 303 {
		t.Fatalf("POST / = %d, want 303 redirect", code)
	}
	if !strings.Contains(headers.Get("Location"), "/entries") {
		t.Fatalf("POST / redirect location = %q, want /entries", headers.Get("Location"))
	}
	code, body = c.Get(t, "/entries")
	if code != 200 {
		t.Fatalf("GET /entries after POST = %d, want 200", code)
	}
	if !strings.Contains(body, "Ada") || !strings.Contains(body, "Hello world") {
		t.Fatalf("entry missing after POST: %s", body)
	}
	code, body = c.Get(t, "/counter")
	if code != 200 {
		t.Fatalf("GET /counter = %d, want 200", code)
	}
	if !strings.Contains(body, "Count: 0") {
		t.Fatalf("counter page missing initial count: %s", body)
	}
	code, _, _ = c.Request(t, "POST", "/counter", "csrf_token="+c.Cookie("csrf_token"), map[string]string{"Content-Type": "application/x-www-form-urlencoded"})
	if code != 200 {
		t.Fatalf("POST /counter = %d, want 200", code)
	}
	code, body = c.Get(t, "/counter")
	if code != 200 {
		t.Fatalf("GET /counter after POST = %d, want 200", code)
	}
	if !strings.Contains(body, "Count: 1") {
		t.Fatalf("counter did not increment: %s", body)
	}
}

