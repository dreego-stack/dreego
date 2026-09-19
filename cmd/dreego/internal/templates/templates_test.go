package templates

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
)

func TestListSorted(t *testing.T) {
	metas := List()
	got := make([]string, 0, len(metas))
	for _, meta := range metas {
		got = append(got, meta.Name)
	}
	want := []string{"web-app", "web-minimal"}
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("List() = %v, want %v", got, want)
	}
}

func TestEveryDiskTemplateRegistered(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	registered := map[string]bool{}
	for _, meta := range List() {
		registered[meta.Name] = true
	}
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() || strings.HasPrefix(name, "_") {
			continue
		}
		if _, err := os.Stat(filepath.Join(name, metaFile)); err != nil {
			continue
		}
		if !registered[name] {
			t.Errorf("template dir %q has a %s but is not registered by List()", name, metaFile)
		}
	}
}

func TestInstallSkipsMetaStripsTmplAndSubstitutes(t *testing.T) {
	mem := fstest.MapFS{
		"_common/main.go.tmpl":      {Data: []byte("common §$name$§\n")},
		"_common/logo.txt":          {Data: []byte("shared §$name$§ logo")},
		"_common/.DS_Store":         {Data: []byte("junk")},
		"web-minimal/template.json": {Data: []byte(`{"name":"web-minimal"}`)},
		"web-minimal/main.go.tmpl":  {Data: []byte("template §$name$§ www\n")},
		"web-minimal/readme.md":     {Data: []byte("welcome to §$name$§")},
	}
	target := t.TempDir()
	if err := install(mem, target, "example.com/acme/app", "web-minimal"); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(target, metaFile)); !os.IsNotExist(err) {
		t.Errorf("%s must not be copied into the target", metaFile)
	}
	if _, err := os.Stat(filepath.Join(target, ".DS_Store")); !os.IsNotExist(err) {
		t.Errorf(".DS_Store must not be copied into the target")
	}
	if _, err := os.Stat(filepath.Join(target, "main.go.tmpl")); !os.IsNotExist(err) {
		t.Errorf(".tmpl must be stripped from destination names")
	}

	main, err := os.ReadFile(filepath.Join(target, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(main) != "template example.com/acme/app www\n" {
		t.Errorf("overlay must win, got %q", main)
	}
	logo, err := os.ReadFile(filepath.Join(target, "logo.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(logo) != "shared example.com/acme/app logo" {
		t.Errorf("_common-only file not copied correctly, got %q", logo)
	}
	if _, err := os.Stat(filepath.Join(target, "readme.md")); err != nil {
		t.Errorf("template file not copied: %v", err)
	}

	err = filepath.WalkDir(target, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), nameToken) {
			t.Errorf("%s still contains %s", p, nameToken)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestInstallSkipsGeneratedOutput(t *testing.T) {
	mem := fstest.MapFS{
		"_common/www/dreego.config.json":     {Data: []byte("{}")},
		"web-minimal/template.json":          {Data: []byte(`{"name":"web-minimal"}`)},
		"web-minimal/www/dree.go":            {Data: []byte("package www")},
		"web-minimal/www/register_dreego.go": {Data: []byte("package www")},
		"web-minimal/www/handle_index.go":    {Data: []byte("package www")},
		"web-minimal/www/keep.go":            {Data: []byte("package www")},
	}
	target := t.TempDir()
	if err := install(mem, target, "example.com/acme/app", "web-minimal"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"dree.go", "register_dreego.go", "handle_index.go"} {
		if _, err := os.Stat(filepath.Join(target, "www", name)); !os.IsNotExist(err) {
			t.Errorf("generated file %s must not be copied into the target", name)
		}
	}
	if _, err := os.Stat(filepath.Join(target, "www", "keep.go")); err != nil {
		t.Errorf("handwritten file must be copied: %v", err)
	}
}

func TestInstallRealTemplate(t *testing.T) {
	target := t.TempDir()
	if err := Install(target, "example.com/acme/app", DefaultName); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"main.go", "Taskfile.yml", "Dockerfile", "docker-compose.yml", ".gitignore", filepath.Join("www", "dreego.config.json")} {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Errorf("expected %s in target: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(target, metaFile)); !os.IsNotExist(err) {
		t.Errorf("%s must not be copied into the target", metaFile)
	}
}

func TestInstallUnknownTemplate(t *testing.T) {
	if err := Install(t.TempDir(), "example.com/acme/app", "does-not-exist"); err == nil {
		t.Fatal("Install must fail for an unknown template")
	}
}

func TestMetaOfValid(t *testing.T) {
	meta, err := MetaOf("web-minimal")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "web-minimal" || meta.Title != "Minimal web app" || meta.Type != "web" || meta.Adapter != "ssr" {
		t.Errorf("unexpected meta: %+v", meta)
	}
}

func TestMetaOfWebApp(t *testing.T) {
	meta, err := MetaOf("web-app")
	if err != nil {
		t.Fatal(err)
	}
	if meta.Name != "web-app" || meta.Title != "Full web application" || meta.Type != "web" || meta.Adapter != "ssr" {
		t.Errorf("unexpected meta: %+v", meta)
	}
	if len(meta.ExtraRequires) != 0 {
		t.Errorf("web-app extraRequires = %v, want empty", meta.ExtraRequires)
	}
}

func TestInstallWebApp(t *testing.T) {
	target := t.TempDir()
	if err := Install(target, "example.com/acme/app", "web-app"); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{
		"main.go",
		"Taskfile.yml",
		filepath.Join("www", "dreego.config.json"),
		filepath.Join("www", "layouts", "default.dreego"),
		filepath.Join("www", "components", "Nav.dreego"),
		filepath.Join("www", "components", "Card.dreego"),
		filepath.Join("www", "routes", "+page.dreego"),
		filepath.Join("www", "routes", "notes_store.go"),
		filepath.Join("www", "routes", "dashboard", "+page.dreego"),
	} {
		if _, err := os.Stat(filepath.Join(target, name)); err != nil {
			t.Errorf("expected %s in target: %v", name, err)
		}
	}
	if _, err := os.Stat(filepath.Join(target, metaFile)); !os.IsNotExist(err) {
		t.Errorf("%s must not be copied into the target", metaFile)
	}
	err := filepath.WalkDir(target, func(p string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		if strings.Contains(string(data), nameToken) {
			t.Errorf("%s still contains %s", p, nameToken)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestMetaOfUnknown(t *testing.T) {
	if _, err := MetaOf("does-not-exist"); err == nil {
		t.Fatal("MetaOf must fail for an unknown name")
	}
	if Exists("does-not-exist") {
		t.Fatal("Exists must be false for an unknown name")
	}
}

func TestMetaOfUnderscore(t *testing.T) {
	if _, err := MetaOf("_common"); err == nil {
		t.Fatal("MetaOf must fail for a _-prefixed layer")
	}
	if Exists("_common") {
		t.Fatal("Exists must be false for a _-prefixed layer")
	}
}

func TestExistsValid(t *testing.T) {
	if !Exists(DefaultName) {
		t.Fatal("default template name must exist")
	}
}

func TestMetaOfRejectsInvalidNames(t *testing.T) {
	for _, name := range []string{
		"does-not-exist",
		"_common",
		"_hidden",
		"",
		"../x",
		"/abs",
		"a/../b",
		`a\b`,
		".hidden",
	} {
		if _, err := MetaOf(name); err == nil {
			t.Errorf("MetaOf(%q) must fail", name)
		}
		if Exists(name) {
			t.Errorf("Exists(%q) must be false", name)
		}
	}
}
