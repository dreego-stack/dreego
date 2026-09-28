package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestReferenceHello(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json": `{"logging":{"enabled":false},"redirects":[],"rewrites":[]}`,
		"www/routes/+page.dreego": `<head>
    <title>Hello</title>
    <meta name="description" content="Minimal Dreego app">
</head>

<server>
    message := "Hello from Dreego"
</server>

<body>
    <h1>{{ message }}</h1>
    <p>This is the minimal Dreego reference app.</p>
    <nav><a href="/about">About</a> <a href="/users/1">User 1</a></nav>
</body>`,
		"www/routes/404.dreego": `<head>
    <title>Not found</title>
</head>

<body>
    <h1>Page not found</h1>
    <p>The URL you requested does not exist.</p>
    <a href="/">Back home</a>
</body>`,
		"www/routes/about/+page.dreego": `<head>
    <title>About</title>
</head>

<body>
    <h1>About this app</h1>
    <p>One route file per URL, one method per file.</p>
    <a href="/">Back home</a>
</body>`,
		"www/routes/users/[id]/+page.dreego": `<head>
    <title>User</title>
</head>

<server>
    id := c.Param("id")
</server>

<body>
    <h1>User {{ id }}</h1>
    <p>Dynamic segments use [brackets] in the file name.</p>
    <a href="/">Back home</a>
</body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "Hello from Dreego") {
		t.Fatalf("home page missing greeting: %s", body)
	}
	if !strings.Contains(body, "<title>Hello</title>") {
		t.Fatalf("home page missing title: %s", body)
	}
	code, body = c.Get(t, "/about")
	if code != 200 {
		t.Fatalf("GET /about = %d, want 200", code)
	}
	if !strings.Contains(body, "About this app") {
		t.Fatalf("about page missing content: %s", body)
	}
	code, body = c.Get(t, "/users/42")
	if code != 200 {
		t.Fatalf("GET /users/42 = %d, want 200", code)
	}
	if !strings.Contains(body, "User 42") {
		t.Fatalf("dynamic route missing param: %s", body)
	}
	code, _ = c.Get(t, "/missing")
	if code != 404 {
		t.Fatalf("GET /missing = %d, want 404", code)
	}
}

