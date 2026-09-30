package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// Generated-file contracts for the one-root-many-apps model.

func TestGeneratedAppHasRegistrarVar(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/app/routes/+page.dreego": `<body><p>hi</p></body>`,
	})
	out := gen["www/app/dree.go"]
	if !strings.Contains(out, "var App dreego.Registrar = registerApp") {
		t.Fatalf("app file must export a Registrar: %s", out)
	}
	if !strings.Contains(out, "package app") {
		t.Fatalf("app file must use the app package name: %s", out)
	}
}

func TestGeneratedRoutePackageRegisterStillPresent(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/app/routes/+page.dreego": `<body><p>hi</p></body>`,
	})
	out := gen["www/app/routes/dree.go"]
	if !strings.Contains(out, "func Register(app *dreego.App) error") {
		t.Fatalf("route package must keep Register(app): %s", out)
	}
	if !strings.Contains(out, "package routes") {
		t.Fatalf("route package name changed: %s", out)
	}
}

func TestGeneratedAppRegistersRoutePackage(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/app/routes/+page.dreego": `<body><p>hi</p></body>`,
	})
	out := gen["www/app/dree.go"]
	if !strings.Contains(out, "routes.Register(app)") {
		t.Fatalf("app file must register its route package: %s", out)
	}
}

func TestGeneratedSetLoggingOnlyWhenPresent(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name         string
		config       string
		wantContains bool
	}{
		{"logging present false", `{"logging":{"enabled":false}}`, true},
		{"logging present true", `{"logging":{"enabled":true}}`, true},
		{"logging absent", `{"redirects":[]}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := dreegotest.ProjectDir(t, map[string]string{
				"www/dreego.config.json":      tc.config,
				"www/app/routes/+page.dreego": `<body><p>hi</p></body>`,
			})
			if out, err := dreegotest.RunCLI(t, dir, "generate"); err != nil {
				t.Fatalf("generate: %v\n%s", err, out)
			}
			data, err := os.ReadFile(filepath.Join(dir, "www", "app", "dree.go"))
			if err != nil {
				t.Fatal(err)
			}
			has := strings.Contains(string(data), "app.SetLogging(")
			if has != tc.wantContains {
				t.Fatalf("SetLogging emitted=%v, want %v:\n%s", has, tc.wantContains, data)
			}
		})
	}
}

func TestGeneratedTwoAppsProduceTwoPackages(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"www/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"www/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
	})
	if _, ok := gen["www/www/dree.go"]; !ok {
		t.Fatalf("missing www app dree.go; got %v", keys(gen))
	}
	if _, ok := gen["www/blog/dree.go"]; !ok {
		t.Fatalf("missing blog app dree.go; got %v", keys(gen))
	}
}

func TestGeneratedSharedLayoutPackage(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/layouts/default.dreego":  `<body><html><body>{#slot}</body></html></body>`,
		"www/app/routes/+page.dreego": `<body><p>hi</p></body>`,
	})
	out, ok := gen["www/layouts/dree.go"]
	if !ok {
		t.Fatalf("missing shared layout package; got %v", keys(gen))
	}
	if !strings.Contains(out, "package layouts") || !strings.Contains(out, "func Default(") {
		t.Fatalf("shared layout package wrong:\n%s", out)
	}
}

func TestGeneratedSharedComponentPackage(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/components/Card.dreego":  "DREEFILE component ()\n<body><b>c</b></body>",
		"www/app/routes/+page.dreego": "COMPONENT \"www/components\" IMPORT { Card }\n<body><@Card/></body>",
	})
	out, ok := gen["www/components/dree.go"]
	if !ok {
		t.Fatalf("missing shared component package; got %v", keys(gen))
	}
	if !strings.Contains(out, "package components") || !strings.Contains(out, "func Card(") {
		t.Fatalf("shared component package wrong:\n%s", out)
	}
}

func TestGeneratedAppLocalComponentPackageName(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":           `{"logging":{"enabled":false}}`,
		"www/components/Badge.dreego":      "DREEFILE component ()\n<body><b>shared</b></body>",
		"www/blog/components/Local.dreego": "DREEFILE component ()\n<body><b>local</b></body>",
		"www/blog/routes/+page.dreego":     "COMPONENT \"www/components\" IMPORT { Badge }\nCOMPONENT \"www/blog/components\" IMPORT { Local }\n<body><@Badge/><@Local/></body>",
	})
	local, ok := gen["www/blog/components/dree.go"]
	if !ok {
		t.Fatalf("missing app-local component package; got %v", keys(gen))
	}
	if !strings.Contains(local, "package blog_components") {
		t.Fatalf("app-local component package must be uniquely named:\n%s", local)
	}
}

func TestGeneratedRootRedirectOnBothApps(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":       `{"logging":{"enabled":false},"redirects":[{"from":"/old","to":"/new","status":301}]}`,
		"www/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"www/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
	})
	for _, app := range []string{"www", "blog"} {
		out := gen["www/"+app+"/dree.go"]
		if !strings.Contains(out, `app.RegisterRedirect("/old", "/new", 301)`) {
			t.Fatalf("app %s missing inherited redirect:\n%s", app, out)
		}
	}
}

func TestGeneratedAppRedirectOnlyInThatApp(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":       `{"logging":{"enabled":false}}`,
		"www/blog/dreego.config.json":  `{"redirects":[{"from":"/b","to":"/c","status":302}]}`,
		"www/www/routes/+page.dreego":  `<body><p>www</p></body>`,
		"www/blog/routes/+page.dreego": `<body><p>blog</p></body>`,
	})
	if strings.Contains(gen["www/www/dree.go"], `app.RegisterRedirect(`) {
		t.Fatalf("www must not inherit blog's redirect:\n%s", gen["www/www/dree.go"])
	}
	if !strings.Contains(gen["www/blog/dree.go"], `app.RegisterRedirect("/b", "/c", 302)`) {
		t.Fatalf("blog missing its redirect:\n%s", gen["www/blog/dree.go"])
	}
}

func TestGeneratedAppStaticRegistration(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/app/routes/+page.dreego": `<body><p>hi</p></body>`,
		"www/app/static/style.css":    `.x{color:red}`,
	})
	out := gen["www/app/dree.go"]
	if !strings.Contains(out, `app.RegisterStatic("/style.css"`) {
		t.Fatalf("app static asset not registered in the app file:\n%s", out)
	}
}

func TestGeneratedNoCollectorRootFile(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":      `{"logging":{"enabled":false}}`,
		"www/app/routes/+page.dreego": `<body><p>hi</p></body>`,
	})
	if _, ok := gen["www/dree.go"]; ok {
		t.Fatalf("no collector root dree.go must be generated:\n%s", gen["www/dree.go"])
	}
}

func keys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
