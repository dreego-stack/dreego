package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// Global defaults: root routes/, static/, layouts/, and components/ are shared
// by every app; an app-local file with the same relative path wins.

func globalDefaultsFixture() map[string]string {
	return map[string]string{
		"www/dreego.config.json":               `{"logging":{"enabled":false}}`,
		"www/layouts/default.dreego":           `<body><html><body><nav id="root-nav">Root</nav>{#slot}</body></html></body>`,
		"www/routes/+page.dreego":              `<body><p>global home</p></body>`,
		"www/routes/about/+page.dreego":        `<body><p>global about</p></body>`,
		"www/routes/only-global/+page.dreego":  `<body><p>only global</p></body>`,
		"www/app/routes/+page.dreego":          `<body><p>app home</p></body>`,
		"www/app/routes/only-app/+page.dreego": `<body><p>only app</p></body>`,
	}
}

func TestGlobalRoutesInheritedAndLocalOverrides(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, globalDefaultsFixture())

	code, body := c.Get(t, "/")
	if code != 200 || !strings.Contains(body, "app home") {
		t.Fatalf("local +page must override the global one: %d %s", code, body)
	}
	if strings.Contains(body, "global home") {
		t.Fatalf("global +page leaked despite a local override: %s", body)
	}

	code, body = c.Get(t, "/about")
	if code != 200 || !strings.Contains(body, "global about") {
		t.Fatalf("global-only route must be inherited: %d %s", code, body)
	}

	code, body = c.Get(t, "/only-global")
	if code != 200 || !strings.Contains(body, "only global") {
		t.Fatalf("global-only route must be inherited: %d %s", code, body)
	}

	code, body = c.Get(t, "/only-app")
	if code != 200 || !strings.Contains(body, "only app") {
		t.Fatalf("app-only route must be served: %d %s", code, body)
	}
}

func TestGlobalLayoutAppliesToInheritedRoutes(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, globalDefaultsFixture())

	_, body := c.Get(t, "/about")
	dreegotest.MustContainBody(t, body, `id="root-nav"`)
}

func TestAppLocalDefaultLayoutWins(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json":         `{"logging":{"enabled":false}}`,
		"www/layouts/default.dreego":     `<body><html><body><nav id="root-nav">Root</nav>{#slot}</body></html></body>`,
		"www/app/layouts/default.dreego": `<body><html><body><nav id="app-nav">App</nav>{#slot}</body></html></body>`,
		"www/app/routes/+page.dreego":    `<body><p>home</p></body>`,
	})
	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `id="app-nav"`)
	dreegotest.MustNotContainBody(t, body, `id="root-nav"`)
}

func TestGlobalStaticInheritedAndLocalWins(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/static/shared.css":       `.shared { color: red; }`,
		"www/static/favicon.ico":      `root-favicon`,
		"www/app/routes/+page.dreego": `<body><p>home</p></body>`,
		"www/app/static/favicon.ico":  `local-favicon`,
		"www/app/static/local.css":    `.local { color: blue; }`,
	})

	code, body := c.Get(t, "/shared.css")
	if code != 200 || !strings.Contains(body, ".shared") {
		t.Fatalf("global static asset must be inherited: %d %s", code, body)
	}

	code, body = c.Get(t, "/favicon.ico")
	if code != 200 || body != "local-favicon" {
		t.Fatalf("local static file must win: %d %q", code, body)
	}

	code, body = c.Get(t, "/local.css")
	if code != 200 || !strings.Contains(body, ".local") {
		t.Fatalf("app-local static asset must be served: %d %s", code, body)
	}
}

func TestMinimalAppInheritsAllGlobalTrees(t *testing.T) {
	t.Parallel()
	// The smallest app carries only its own dreego.config.json and inherits the
	// global routes, static and layouts.
	c := dreegotest.ServeSetup(t, map[string]string{
		"www/dreego.config.json":     `{"logging":{"enabled":false}}`,
		"www/routes/+page.dreego":    `<body><p>shared home</p></body>`,
		"www/static/global.css":      `.global {}`,
		"www/layouts/default.dreego": `<body><html><body><nav id="root-nav">Root</nav>{#slot}</body></html></body>`,
		"www/www/dreego.config.json": `{}`,
	}, "")

	code, body := c.Get(t, "/")
	if code != 200 || !strings.Contains(body, "shared home") {
		t.Fatalf("minimal app must serve the global route: %d %s", code, body)
	}
	dreegotest.MustContainBody(t, body, `id="root-nav"`)

	code, body = c.Get(t, "/global.css")
	if code != 200 || !strings.Contains(body, ".global") {
		t.Fatalf("minimal app must serve the global static: %d %s", code, body)
	}
}

func TestAppWithoutRoutesInheritsGlobal(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"site/routes/+page.dreego":     `<body><p>shared home</p></body>`,
		"site/www/routes/+page.dreego": `<body><p>www home</p></body>`,
		"site/blog/dreego.config.json": `{"logging":{"enabled":false}}`,
		"site/blog/static/blog.css":    `.blog {}`,
	}, "site/www", "site/blog")

	_, body := m.App(t, "blog").Get(t, "/")
	dreegotest.MustContainBody(t, body, "shared home")

	code, _ := m.App(t, "blog").Get(t, "/blog.css")
	if code != 200 {
		t.Fatalf("app without routes must still serve its own static: %d", code)
	}
}

func TestGlobalRouteLocalLayoutApplies(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json":                 `{"logging":{"enabled":false}}`,
		"www/layouts/default.dreego":             `<body><html><body><nav id="root-nav">Root</nav>{#slot}</body></html></body>`,
		"www/routes/blog/layouts/default.dreego": `<body><html><body><nav id="blog-nav">Blog</nav>{#slot}</body></html></body>`,
		"www/routes/blog/+page.dreego":           `<body><p>blog index</p></body>`,
		"www/routes/+page.dreego":                `<body><p>home</p></body>`,
		"www/app/routes/blog/+page.dreego":       `<body><p>app blog index</p></body>`,
	})

	_, body := c.Get(t, "/blog")
	dreegotest.MustContainBody(t, body, `id="blog-nav"`)
	dreegotest.MustNotContainBody(t, body, `id="root-nav"`)

	_, body = c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `id="root-nav"`)
}

func TestGlobalRouteLocalLayoutFromSharedTree(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, map[string]string{
		"www/dreego.config.json":                  `{"logging":{"enabled":false}}`,
		"www/layouts/default.dreego":              `<body><html><body>{#slot}</body></html></body>`,
		"www/routes/admin/layouts/default.dreego": `<body><html><body><nav id="admin-nav">Admin</nav>{#slot}</body></html></body>`,
		"www/routes/admin/+page.dreego":           `<body><p>admin</p></body>`,
		"www/app/routes/admin/+page.dreego":       `<body><p>app admin</p></body>`,
	})
}

func TestGlobalRoutesAndLayoutsVisibleToEveryApp(t *testing.T) {
	t.Parallel()
	m := dreegotest.ServeApps(t, map[string]string{
		"site/dreego.config.json":              `{"logging":{"enabled":false}}`,
		"site/layouts/default.dreego":          `<body><html><body><nav id="shared-nav">Shared</nav>{#slot}</body></html></body>`,
		"site/routes/global-page/+page.dreego": `<body><p>global page</p></body>`,
		"site/www/routes/+page.dreego":         `<body><p>www</p></body>`,
		"site/www/routes/own/+page.dreego":     `<body><p>www own</p></body>`,
		"site/blog/routes/+page.dreego":        `<body><p>blog</p></body>`,
	}, "site/www", "site/blog")

	for _, name := range []string{"www", "blog"} {
		code, body := m.App(t, name).Get(t, "/global-page")
		if code != 200 || !strings.Contains(body, "global page") {
			t.Fatalf("app %s missing global route: %d %s", name, code, body)
		}
		dreegotest.MustContainBody(t, body, `id="shared-nav"`)
	}

	code, _ := m.App(t, "blog").Get(t, "/own")
	if code != 404 {
		t.Fatalf("blog must not see www's local route: %d", code)
	}
}

func TestExplicitLayoutOnRouteSelectsLayout(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/layouts/default.dreego":  `<body><html><body><nav id="default-nav">Default</nav>{#slot}</body></html></body>`,
		"www/layouts/admin.dreego":    `<body><html><body><nav id="admin-nav">Admin</nav>{#slot}</body></html></body>`,
		"www/app/routes/+page.dreego": "LAYOUT \"www/layouts/admin.dreego\"\n\n<body><p>home</p></body>",
	})

	_, body := c.Get(t, "/")
	dreegotest.MustContainBody(t, body, `id="admin-nav"`)
	dreegotest.MustNotContainBody(t, body, `id="default-nav"`)
}

func TestExplicitLayoutMissingTargetFailsGenerate(t *testing.T) {
	t.Parallel()
	out, err := dreegotest.RunCLI(t, dreegotest.ProjectDir(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/layouts/default.dreego":  `<body><html><body>{#slot}</body></html></body>`,
		"www/app/routes/+page.dreego": "LAYOUT \"www/layouts/missing.dreego\"\n\n<body><p>home</p></body>",
	}), "generate")
	if err == nil {
		t.Fatalf("generate accepted a missing explicit LAYOUT target:\n%s", out)
	}
	if !strings.Contains(out, "www/layouts/missing.dreego") {
		t.Fatalf("diagnostic must name the missing layout, got:\n%s", out)
	}
}
