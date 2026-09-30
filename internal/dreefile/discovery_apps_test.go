package dreefile

import (
	"os"
	"path/filepath"
	"testing"
)

func mustMkdir(t *testing.T, paths ...string) {
	t.Helper()
	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			t.Fatal(err)
		}
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	mustMkdir(t, filepath.Dir(path))
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestScanAppsFindsOnlyRouteDirs(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "www", "routes"), filepath.Join(root, "blog", "routes"), filepath.Join(root, "notapp"))
	apps, err := scanApps(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 2 {
		t.Fatalf("scanApps found %d apps, want 2: %+v", len(apps), apps)
	}
	if apps[0].name != "blog" || apps[1].name != "www" {
		t.Fatalf("apps not sorted by name: %+v", apps)
	}
	for _, a := range apps {
		if a.routes != filepath.Join(root, a.name, "routes") {
			t.Errorf("app %s routes = %q", a.name, a.routes)
		}
		if a.pkg != a.name {
			t.Errorf("app %s pkg = %q", a.name, a.pkg)
		}
	}
}

func TestScanAppsSkipsDotAndVendor(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t,
		filepath.Join(root, ".hidden", "routes"),
		filepath.Join(root, "vendor", "routes"),
		filepath.Join(root, "node_modules", "routes"),
		filepath.Join(root, "www", "routes"),
	)
	apps, err := scanApps(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(apps) != 1 || apps[0].name != "www" {
		t.Fatalf("scanApps = %+v, want only www", apps)
	}
}

func TestScanAppsSanitizesPkgName(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "my-app", "routes"), filepath.Join(root, "2cool", "routes"))
	apps, err := scanApps(root)
	if err != nil {
		t.Fatal(err)
	}
	byName := map[string]string{}
	for _, a := range apps {
		byName[a.name] = a.pkg
	}
	if byName["my-app"] != "my_app" {
		t.Errorf("pkg for my-app = %q, want my_app", byName["my-app"])
	}
	if byName["2cool"] != "pkg2cool" {
		t.Errorf("pkg for 2cool = %q, want pkg2cool", byName["2cool"])
	}
}

func TestScanAppsSetsStaticAndConfigPaths(t *testing.T) {
	root := t.TempDir()
	mustMkdir(t, filepath.Join(root, "www", "routes"))
	apps, err := scanApps(root)
	if err != nil || len(apps) != 1 {
		t.Fatalf("scanApps = %+v, %v", apps, err)
	}
	if apps[0].static != filepath.Join(root, "www", "static") {
		t.Errorf("static = %q", apps[0].static)
	}
	if apps[0].config != filepath.Join(root, "www", configFileName) {
		t.Errorf("config = %q", apps[0].config)
	}
}

func TestFindWebsiteRootsRejectsLegacyRoot(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, configFileName), "{}")
	mustWrite(t, filepath.Join(dir, "routes", "+page.dreego"), "<body>x</body>")
	runFindRoots(t, dir, func(t *testing.T, err error) {
		if err == nil {
			t.Fatal("expected legacy layout error")
		}
	})
}

func TestFindWebsiteRootsAcceptsAppLayout(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, configFileName), "{}")
	mustWrite(t, filepath.Join(dir, "www", "routes", "+page.dreego"), "<body>x</body>")
	runFindRoots(t, dir, func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestFindWebsiteRootsRejectsLegacyRootWithSecondRoot(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, configFileName), "{}")
	mustWrite(t, filepath.Join(dir, "routes", "+page.dreego"), "<body>x</body>")
	mustWrite(t, filepath.Join(dir, "other", configFileName), "{}")
	mustWrite(t, filepath.Join(dir, "other", "www", "routes", "+page.dreego"), "<body>y</body>")
	runFindRoots(t, dir, func(t *testing.T, err error) {
		if err == nil {
			t.Fatal("expected legacy layout error even with a second website root")
		}
	})
}

func TestFindWebsiteRootsNoRoots(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "some", "file.txt"), "no website here")
	runFindRoots(t, dir, func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestFindWebsiteRootsSkipsNestedAppRoutes(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, configFileName), "{}")
	mustWrite(t, filepath.Join(dir, "www", "routes", "+page.dreego"), "<body>x</body>")
	mustWrite(t, filepath.Join(dir, "www", "sub", "routes", "+page.dreego"), "<body>sub</body>")
	runFindRoots(t, dir, func(t *testing.T, err error) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

// runFindRoots runs findWebsiteRoots from dir as the working directory.
func runFindRoots(t *testing.T, dir string, assert func(*testing.T, error)) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	_, err = findWebsiteRoots()
	assert(t, err)
}
