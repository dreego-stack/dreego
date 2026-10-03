package dreefile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// A layout may reference a component from the app-local components/ tree of the
// app that uses it, not only a component from the shared root components/. The
// layout is generated in the root pass, before the per-app component scan, so
// app-local components must be registered before layouts are generated.
func TestLayoutUsesAppLocalComponent(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"www/dreego.config.json":           "{}",
		"www/app/components/Header.dreego": "DREEFILE component ()\n<body><header>app header</header></body>",
		"www/app/layouts/default.dreego":   "COMPONENT \"www/app/components\" IMPORT { Header }\n<body><html><body><@Header/>{#slot}</body></html></body>",
		"www/app/routes/+page.dreego":      "<body><p>home</p></body>",
	})

	if err := runInDir(t, dir); err != nil {
		t.Fatalf("Run: %v", err)
	}

	layout, err := os.ReadFile(filepath.Join(dir, "www", "layouts", "dree.go"))
	if err != nil {
		t.Fatalf("read layout output: %v", err)
	}
	if !strings.Contains(string(layout), "Header()") {
		t.Fatalf("layout must call the app-local Header component:\n%s", layout)
	}

	component, err := os.ReadFile(filepath.Join(dir, "www", "app", "components", "dree.go"))
	if err != nil {
		t.Fatalf("read app-local component output: %v", err)
	}
	if !strings.Contains(string(component), "func Header(") {
		t.Fatalf("app-local Header must be generated:\n%s", component)
	}
}
