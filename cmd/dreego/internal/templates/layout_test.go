package templates

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Layout contracts for the one-root-many-apps scaffold (a `dreego/` root with a
// `www` app).

func TestInstallWebMinimalLayout(t *testing.T) {
	target := t.TempDir()
	if err := Install(target, "example.com/acme/app", "web-minimal"); err != nil {
		t.Fatal(err)
	}
	want := []string{
		"main.go",
		filepath.Join("dreego", "dreego.config.json"),
		filepath.Join("dreego", "layouts", "default.dreego"),
		filepath.Join("dreego", "www", "routes", "+page.dreego"),
	}
	for _, name := range want {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Errorf("web-minimal missing %s: %v", name, err)
		}
	}
	forbidden := []string{
		filepath.Join("www", "routes"),
		filepath.Join("www", "dreego.config.json"),
		filepath.Join("dreego", "routes"),
	}
	for _, name := range forbidden {
		if _, err := os.Stat(filepath.Join(target, name)); !os.IsNotExist(err) {
			t.Errorf("web-minimal must not create legacy path %s", name)
		}
	}
}

func TestInstallWebAppLayout(t *testing.T) {
	target := t.TempDir()
	if err := Install(target, "example.com/acme/app", "web-app"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		filepath.Join("dreego", "dreego.config.json"),
		filepath.Join("dreego", "layouts", "default.dreego"),
		filepath.Join("dreego", "components", "Nav.dreego"),
		filepath.Join("dreego", "components", "Card.dreego"),
		filepath.Join("dreego", "www", "routes", "+page.dreego"),
		filepath.Join("dreego", "www", "routes", "dashboard", "+page.dreego"),
	} {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Errorf("web-app missing %s: %v", name, err)
		}
	}
}

func TestCommonMainImportsAppRegistrar(t *testing.T) {
	target := t.TempDir()
	if err := Install(target, "example.com/acme/myapp", "web-minimal"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(target, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, `"example.com/acme/myapp/dreego/www"`) {
		t.Errorf("main.go must import the app package under the module path:\n%s", content)
	}
	if !strings.Contains(content, "dreego.New(www.App)") {
		t.Errorf("main.go must call dreego.New(www.App):\n%s", content)
	}
	if strings.Contains(content, nameToken) {
		t.Errorf("main.go still contains the name placeholder:\n%s", content)
	}
}

func TestInstalledAppLayoutHasNoLegacyRootRoutes(t *testing.T) {
	for _, name := range []string{"web-minimal", "web-app"} {
		t.Run(name, func(t *testing.T) {
			target := t.TempDir()
			if err := Install(target, "example.com/acme/app", name); err != nil {
				t.Fatal(err)
			}
			rootRoutes := filepath.Join(target, "dreego", "routes")
			if _, err := os.Stat(rootRoutes); !os.IsNotExist(err) {
				t.Fatalf("%s must not install a root-level dreego/routes (legacy layout)", name)
			}
		})
	}
}

func TestInstallSkipsGeneratedFilesInNewLayout(t *testing.T) {
	target := t.TempDir()
	if err := Install(target, "example.com/acme/app", "web-minimal"); err != nil {
		t.Fatal(err)
	}
	err := filepath.WalkDir(target, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return walkErr
		}
		base := filepath.Base(p)
		if base == "dree.go" || base == "template.json" {
			t.Errorf("scaffold must not contain generated/meta file %s", p)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstalledWebMinimalConfigHasLogging(t *testing.T) {
	target := t.TempDir()
	if err := Install(target, "example.com/acme/app", "web-minimal"); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(target, "dreego", "dreego.config.json"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"logging"`) {
		t.Fatalf("scaffolded config missing logging: %s", data)
	}
}

func TestListDeterministic(t *testing.T) {
	first := List()
	second := List()
	if len(first) != len(second) {
		t.Fatalf("List length changed: %d vs %d", len(first), len(second))
	}
	for i := range first {
		if first[i].Name != second[i].Name {
			t.Fatalf("List order not deterministic: %v vs %v", first[i].Name, second[i].Name)
		}
	}
}
