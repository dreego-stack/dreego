package tests

import (
	"strings"
	"testing"

	"github.com/dreego-stack/dreego/dreegotest"
)

// A layout may reference a component from the app-local components/ tree of the
// app that uses it, not only a component from the shared root components/. The
// layout is generated in the root pass, so app-local components must be
// registered before layouts are generated.
func TestBugLayoutUsesAppLocalComponent(t *testing.T) {
	t.Parallel()
	gen := dreegotest.Build(t, map[string]string{
		"www/dreego.config.json":           `{"logging":{"enabled":false}}`,
		"www/app/components/Header.dreego": "DREEFILE component ()\n<body><header>app header</header></body>",
		"www/app/layouts/default.dreego":   "COMPONENT \"www/app/components\" IMPORT { Header }\n<body><html><body><@Header/>{#slot}</body></html></body>",
		"www/app/routes/+page.dreego":      "<body><p>home</p></body>",
	})

	layout, ok := gen["www/layouts/dree.go"]
	if !ok {
		t.Fatalf("layout package not generated; got %v", genKeys(gen))
	}
	if !strings.Contains(layout, "Header()") {
		t.Fatalf("layout must call the app-local Header component:\n%s", layout)
	}

	component, ok := gen["www/app/components/dree.go"]
	if !ok {
		t.Fatalf("app-local component package not generated; got %v", genKeys(gen))
	}
	if !strings.Contains(component, "func Header(") {
		t.Fatalf("app-local Header must be generated:\n%s", component)
	}
}
