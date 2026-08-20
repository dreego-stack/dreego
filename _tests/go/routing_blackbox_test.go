package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func TestRoutingBlackboxCatchall(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/blog/[...catchall]/get.dreego": `<go>p := c.Param("catchall")</go>
<div><p>blog:{{ p }}</p></div>`,
	})
	code, body := c.Get(t, "/blog/hello/world")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "blog:hello/world")
}

func TestRoutingBlackboxCatchallRoot(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/[...path]/get.dreego": `<go>p := c.Param("path")</go>
<div><p>root:{{ p }}</p></div>`,
	})
	code, body := c.Get(t, "/a/b/c")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "root:a/b/c")
}

func TestRoutingBlackboxGroup(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/(admin)/dashboard/get.dreego": `<div><p>admin dashboard</p></div>`,
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
		"dreego/routes/users/[id]/get.dreego": `<go>id := c.Param("id")</go>
<div><p>user:{{ id }}</p></div>`,
	})
	code, body := c.Get(t, "/users/42")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "user:42")
}

func TestRoutingBlackboxMethods(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego":    `<div><p>get route</p></div>`,
		"dreego/routes/post.dreego":   `<div><p>post route</p></div>`,
		"dreego/routes/put.dreego":    `<div><p>put route</p></div>`,
		"dreego/routes/delete.dreego": `<div><p>delete route</p></div>`,
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
		"dreego/routes/about/get.dreego":       `<div><p>about</p></div>`,
		"dreego/routes/users/about/get.dreego": `<div><p>users about</p></div>`,
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

func TestRoutingBlackboxMethodAttr(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"dreego/routes/get.dreego": `<go method="post">msg := "posted"</go>
<div><p>{{ msg }}</p></div>`,
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
