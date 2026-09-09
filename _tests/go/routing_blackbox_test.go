package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestRoutingBlackboxCatchall(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/blog/[...catchall]/+page.dreego": `<server>p := c.Param("catchall")</server>
<body><p>blog:{{ p }}</p></body>`,
	})
	code, body := c.Get(t, "/blog/hello/world")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "blog:hello/world")
}

func TestRoutingBlackboxCatchallRoot(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/[...path]/+page.dreego": `<server>p := c.Param("path")</server>
<body><p>root:{{ p }}</p></body>`,
	})
	code, body := c.Get(t, "/a/b/c")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "root:a/b/c")
}

func TestRoutingBlackboxGroup(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/(admin)/dashboard/+page.dreego": `<body><p>admin dashboard</p></body>`,
	})
	code, body := c.Get(t, "/dashboard")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "admin dashboard")
	code, _ = c.Get(t, "/admin/dashboard")
	if code == 200 {
		t.Fatal("group segment must not appear in the URL")
	}
}

func TestRoutingBlackboxDynamic(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/users/[id]/+page.dreego": `<server>id := c.Param("id")</server>
<body><p>user:{{ id }}</p></body>`,
	})
	code, body := c.Get(t, "/users/42")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "user:42")
}

func TestRoutingBlackboxNamedFilesAreURLSegments(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego":   `<body><p>home route</p></body>`,
		"www/routes/page.dreego":    `<body><p>page route</p></body>`,
		"www/routes/profile.dreego": `<body><p>profile route</p></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 || !strings.Contains(body, "home route") {
		t.Fatalf("GET / = %d %q, want 200 with home route", code, body)
	}
	for path, want := range map[string]string{"/page": "page route", "/profile": "profile route"} {
		code, body = c.Get(t, path)
		if code != 200 || !strings.Contains(body, want) {
			t.Fatalf("GET %s = %d %q, want 200 with %s", path, code, body, want)
		}
	}
}

func TestRoutingBlackboxNested(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/about/+page.dreego":       `<body><p>about</p></body>`,
		"www/routes/users/about/+page.dreego": `<body><p>users about</p></body>`,
	})
	code, body := c.Get(t, "/about")
	if code != 200 || !strings.Contains(body, "about") {
		t.Fatalf("GET /about = %d %q, want 200", code, body)
	}
	code, body = c.Get(t, "/users/about")
	if code != 200 || !strings.Contains(body, "users about") {
		t.Fatalf("GET /users/about = %d %q, want 200", code, body)
	}
}

func TestRoutingBlackboxFlatRoutes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego":            `<body><p>home</p></body>`,
		"www/routes/about.dreego":            `<body><p>about flat</p></body>`,
		"www/routes/users/[id]/+page.dreego": `<server>id := c.Param("id")</server><body><p>user:{{ id }}</p></body>`,
	})
	code, body := c.Get(t, "/about")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "about flat")
	code, body = c.Get(t, "/users/42")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "user:42")
}

func TestRoutingBlackboxMethodAttr(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/+page.dreego": `<server method="post">msg := "posted"</server>
<body method="post"><p>{{ msg }}</p></body>`,
	})
	code, body, _ := c.Request(t, "POST", "/", "", nil)
	if code != 200 || !strings.Contains(body, "posted") {
		t.Fatalf("POST / = %d %q, want 200 with posted body", code, body)
	}
	code, _ = c.Get(t, "/")
	if code == 200 {
		t.Fatal("GET / must not match a route registered only for POST")
	}
}

func TestRoutingBlackboxMethodSections(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/about.dreego": `<server>msg := "get"</server>
<body><p>{{ msg }}</p></body>
<server method="post">msg := "post"</server>
<body method="post"><p>{{ msg }}</p></body>`,
	})
	code, body := c.Get(t, "/about")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "get")
	if strings.Contains(body, "post") {
		t.Fatal("GET must not render POST template")
	}
	code, body, _ = c.Request(t, "POST", "/about", "", nil)
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "post")
	if strings.Contains(body, "get") {
		t.Fatal("POST must not render GET template")
	}
}
