package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestRoutingBlackboxCatchall(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/blog/[...catchall]/get.dreego": `<server>p := c.Param("catchall")</server>
<body><p>blog:{{ p }}</p></body>`,
	})
	code, body := c.Get(t, "/blog/hello/world")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "blog:hello/world")
}

func TestRoutingBlackboxCatchallRoot(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/[...path]/get.dreego": `<server>p := c.Param("path")</server>
<body><p>root:{{ p }}</p></body>`,
	})
	code, body := c.Get(t, "/a/b/c")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "root:a/b/c")
}

func TestRoutingBlackboxGroup(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/(admin)/dashboard/get.dreego": `<body><p>admin dashboard</p></body>`,
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
		"www/routes/users/[id]/get.dreego": `<server>id := c.Param("id")</server>
<body><p>user:{{ id }}</p></body>`,
	})
	code, body := c.Get(t, "/users/42")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "user:42")
}

func TestRoutingBlackboxMethods(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/get.dreego":    `<body><p>get route</p></body>`,
		"www/routes/post.dreego":   `<body><p>post route</p></body>`,
		"www/routes/put.dreego":    `<body><p>put route</p></body>`,
		"www/routes/delete.dreego": `<body><p>delete route</p></body>`,
	})
	code, body := c.Get(t, "/")
	if code != 200 || !strings.Contains(body, "get route") {
		t.Fatalf("GET / = %d %q, want 200 with get route", code, body)
	}
	code, body, _ = c.Request(t, "POST", "/", "", nil)
	if code != 200 || !strings.Contains(body, "post route") {
		t.Fatalf("POST / = %d %q, want 200 with post route", code, body)
	}
	code, body, _ = c.Request(t, "PUT", "/", "", nil)
	if code != 200 || !strings.Contains(body, "put route") {
		t.Fatalf("PUT / = %d %q, want 200 with put route", code, body)
	}
	code, body, _ = c.Request(t, "DELETE", "/", "", nil)
	if code != 200 || !strings.Contains(body, "delete route") {
		t.Fatalf("DELETE / = %d %q, want 200 with delete route", code, body)
	}
}

func TestRoutingBlackboxNested(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/routes/about/get.dreego":       `<body><p>about</p></body>`,
		"www/routes/users/about/get.dreego": `<body><p>users about</p></body>`,
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
		"www/routes/get.dreego": `<server method="post">msg := "posted"</server>
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
