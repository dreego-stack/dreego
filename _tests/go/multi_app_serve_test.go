package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// Each app serves on its own port and responds independently.

func TestMultiAppIndependentPorts(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/www/routes/+page.dreego":  `<body><h1>www home</h1></body>`,
		"site/blog/routes/+page.dreego": `<body><h1>blog home</h1></body>`,
	}, "site/www", "site/blog")

	code, body := m.App(t, "www").Get(t, "/")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "www home")

	code, body = m.App(t, "blog").Get(t, "/")
	dreegotest.MustStatus(t, code, 200)
	dreegotest.MustContainBody(t, body, "blog home")
}

func TestMultiAppRoutesDoNotLeak(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":               `{"logging":{"enabled":false}}`,
		"site/www/routes/+page.dreego":          `<body><p>www root</p></body>`,
		"site/www/routes/only-www/+page.dreego": `<body><p>www only</p></body>`,
		"site/blog/routes/+page.dreego":         `<body><p>blog root</p></body>`,
	}, "site/www", "site/blog")

	code, _ := m.App(t, "blog").Get(t, "/only-www")
	if code != 404 {
		t.Fatalf("blog must not see www's route: status = %d, want 404", code)
	}
	code, body := m.App(t, "www").Get(t, "/only-www")
	if code != 200 || !strings.Contains(body, "www only") {
		t.Fatalf("www route missing: %d %s", code, body)
	}
}

func TestMultiAppSharedComponentRendersOnBothApps(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/components/Badge.dreego":  "DREEFILE component (label string)\n<body><span class=\"badge\">{{ label }}</span></body>",
		"site/www/routes/+page.dreego":  "COMPONENT \"site/components\" IMPORT { Badge }\n<body><@Badge label=\"w\"/></body>",
		"site/blog/routes/+page.dreego": "COMPONENT \"site/components\" IMPORT { Badge }\n<body><@Badge label=\"b\"/></body>",
	}, "site/www", "site/blog")

	for name, want := range map[string]string{"www": "w", "blog": "b"} {
		code, body := m.App(t, name).Get(t, "/")
		dreegotest.MustStatus(t, code, 200)
		if !strings.Contains(body, `<span class="badge">`+want+`</span>`) {
			t.Fatalf("app %s missing shared component output %q: %s", name, want, body)
		}
	}
}

func TestMultiAppSharedLayoutAppliesToBoth(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/layouts/default.dreego":   `<body><html><body><nav id="shared-nav">Shared</nav>{#slot}</body></html></body>`,
		"site/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"site/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
	}, "site/www", "site/blog")

	for _, name := range []string{"www", "blog"} {
		code, body := m.App(t, name).Get(t, "/")
		dreegotest.MustStatus(t, code, 200)
		dreegotest.MustContainBody(t, body, `id="shared-nav"`)
	}
}

func TestMultiAppAppLocalLayoutOverridesShared(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":          `{"logging":{"enabled":false}}`,
		"site/layouts/default.dreego":      `<body><html><body><nav id="shared">S</nav>{#slot}</body></html></body>`,
		"site/blog/layouts/default.dreego": `<body><html><body><nav id="blog-local">B</nav>{#slot}</body></html></body>`,
		"site/www/routes/+page.dreego":     `<body><p>www</p></body>`,
		"site/blog/routes/+page.dreego":    `<body><p>blog</p></body>`,
	}, "site/www", "site/blog")

	_, wwwBody := m.App(t, "www").Get(t, "/")
	dreegotest.MustContainBody(t, wwwBody, `id="shared"`)
	dreegotest.MustNotContainBody(t, wwwBody, `id="blog-local"`)

	_, blogBody := m.App(t, "blog").Get(t, "/")
	dreegotest.MustContainBody(t, blogBody, `id="blog-local"`)
	dreegotest.MustNotContainBody(t, blogBody, `id="shared"`)
}

func TestMultiAppAppConfigOverridesRootLogging(t *testing.T) {
	t.Parallel()
	// Root disables logging; blog re-enables it. This only asserts the app
	// still serves (the logging flag is exercised by the generated file test).
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/blog/dreego.config.json":  `{"logging":{"enabled":true}}`,
		"site/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"site/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
	}, "site/www", "site/blog")

	code, _ := m.App(t, "www").Get(t, "/")
	dreegotest.MustStatus(t, code, 200)
	code, _ = m.App(t, "blog").Get(t, "/")
	dreegotest.MustStatus(t, code, 200)
}

func TestMultiAppOwnStaticAssets(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"site/www/static/only.css":      `.www-only { color: red; }`,
		"site/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
		"site/blog/static/blog.css":     `.blog-only { color: blue; }`,
	}, "site/www", "site/blog")

	code, body := m.App(t, "www").Get(t, "/only.css")
	code2, _ := m.App(t, "www").Get(t, "/blog.css")
	if code != 200 || !strings.Contains(body, "www-only") {
		t.Fatalf("www static asset: %d %s", code, body)
	}
	if code2 != 404 {
		t.Fatalf("blog static asset must not be visible on www: %d", code2)
	}
}

func TestMultiApp404IsPerApp(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"site/www/routes/404.dreego":    `<body><p>WWW NOT FOUND</p></body>`,
		"site/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
		"site/blog/routes/404.dreego":   `<body><p>BLOG NOT FOUND</p></body>`,
	}, "site/www", "site/blog")

	code, body := m.App(t, "www").Get(t, "/nope")
	if code != 404 || !strings.Contains(body, "WWW NOT FOUND") {
		t.Fatalf("www 404: %d %s", code, body)
	}
	code, body = m.App(t, "blog").Get(t, "/nope")
	if code != 404 || !strings.Contains(body, "BLOG NOT FOUND") {
		t.Fatalf("blog 404: %d %s", code, body)
	}
}

func TestMultiAppRedirectIsPerApp(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"site/blog/dreego.config.json":  `{"redirects":[{"from":"/old","to":"/new","status":301}]}`,
		"site/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
	}, "site/www", "site/blog")

	code, _, h := m.App(t, "blog").Request(t, "GET", "/old", "", nil)
	if code != 301 || h.Get("Location") != "/new" {
		t.Fatalf("blog redirect: status=%d location=%q", code, h.Get("Location"))
	}
	code, _, _ = m.App(t, "www").Request(t, "GET", "/old", "", nil)
	if code != 404 {
		t.Fatalf("www must not inherit blog's redirect: status=%d", code)
	}
}

func TestMultiAppRootRedirectInherited(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false},"redirects":[{"from":"/home","to":"/","status":302}]}`,
		"site/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"site/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
	}, "site/www", "site/blog")

	for _, name := range []string{"www", "blog"} {
		code, _, h := m.App(t, name).Request(t, "GET", "/home", "", nil)
		if code != 302 || h.Get("Location") != "/" {
			t.Fatalf("app %s root redirect: status=%d location=%q", name, code, h.Get("Location"))
		}
	}
}

func TestMultiAppHeadTitleIsolated(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"site/www/routes/+page.dreego":  `<head><title>WWW</title></head><body><p>www</p></body>`,
		"site/blog/routes/+page.dreego": `<head><title>BLOG</title></head><body><p>blog</p></body>`,
	}, "site/www", "site/blog")

	_, wwwBody := m.App(t, "www").Get(t, "/")
	dreegotest.MustContainBody(t, wwwBody, "<title>WWW</title>")
	dreegotest.MustNotContainBody(t, wwwBody, "BLOG")

	_, blogBody := m.App(t, "blog").Get(t, "/")
	dreegotest.MustContainBody(t, blogBody, "<title>BLOG</title>")
	dreegotest.MustNotContainBody(t, blogBody, "WWW")
}

func TestMultiAppDynamicRoutesPerApp(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":                    `{"logging":{"enabled":false}}`,
		"site/www/routes/users/[id]/+page.dreego":    `<server>id := c.Param("id")</server><body><p>www user {{ id }}</p></body>`,
		"site/blog/routes/posts/[slug]/+page.dreego": `<server>slug := c.Param("slug")</server><body><p>blog post {{ slug }}</p></body>`,
	}, "site/www", "site/blog")

	code, body := m.App(t, "www").Get(t, "/users/42")
	if code != 200 || !strings.Contains(body, "www user 42") {
		t.Fatalf("www dynamic: %d %s", code, body)
	}
	code, body = m.App(t, "blog").Get(t, "/posts/hello")
	if code != 200 || !strings.Contains(body, "blog post hello") {
		t.Fatalf("blog dynamic: %d %s", code, body)
	}
	code, _ = m.App(t, "blog").Get(t, "/users/42")
	if code != 404 {
		t.Fatalf("blog must not serve www's dynamic route: %d", code)
	}
}
