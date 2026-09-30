package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

func routeLocalLayoutFixture() map[string]string {
	return map[string]string{
		"www/layouts/default.dreego":                      "<body><html lang=\"en\"><body><nav>Root nav</nav>{#slot}</body></html></body>",
		"www/app/routes/registrierung/layouts/default.dreego": "<body><html lang=\"de\"><body class=\"auth\"><nav>Registrierung nav</nav>{#slot}</body></html></body>",
		"www/app/routes/+page.dreego":                         "<body><p>Home page</p></body>",
		"www/app/routes/registrierung/+page.dreego":           "<body><p>Registrierung page</p></body>",
	}
}

func TestBugRouteLocalLayoutBuilds(t *testing.T) {
	t.Parallel()
	dreegotest.MustBuild(t, routeLocalLayoutFixture())
}

func TestBugRouteLocalLayoutGeneratedNames(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, routeLocalLayoutFixture())

	layouts := gen["www/layouts/dree.go"]
	if strings.Count(layouts, "func Default(") != 1 {
		t.Fatalf("root layout must emit exactly one func Default, got:\n%s", layouts)
	}
	dreegotest.MustContain(t, layouts, "func DefaultAppRegistrierung(")

	routes := gen["www/app/routes/dree.go"]
	dreegotest.MustContain(t, routes, "layouts.Default(c, pageContent, head)")

	routeLocal := gen["www/app/routes/registrierung/dree.go"]
	dreegotest.MustContain(t, routeLocal, "layouts.DefaultAppRegistrierung(c, pageContent, head)")
	dreegotest.MustNotContain(t, routes, "/registrierung/layouts")
	dreegotest.MustNotContain(t, routeLocal, "/registrierung/layouts")
}

func TestBugRouteLocalLayoutCascade(t *testing.T) {
	t.Parallel()
	c := dreegotest.Serve(t, routeLocalLayoutFixture())

	code, body := c.Get(t, "/")
	if code != 200 {
		t.Fatalf("GET / status = %d, want 200", code)
	}
	for _, want := range []string{"Root nav", "Home page"} {
		if !strings.Contains(body, want) {
			t.Fatalf("root route missing %q, got:\n%s", want, body)
		}
	}

	code, body = c.Get(t, "/registrierung")
	if code != 200 {
		t.Fatalf("GET /registrierung status = %d, want 200", code)
	}
	for _, want := range []string{"Registrierung nav", "Registrierung page", `class="auth"`} {
		if !strings.Contains(body, want) {
			t.Fatalf("route-local layout missing %q, got:\n%s", want, body)
		}
	}
}
