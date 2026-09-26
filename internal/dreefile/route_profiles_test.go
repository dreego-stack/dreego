package dreefile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func allRouteRegs(dirs []*routePkg) string {
	var b strings.Builder
	var collect func(p *routePkg)
	collect = func(p *routePkg) {
		for _, r := range p.regs {
			b.WriteString(r)
		}
		for _, c := range p.children {
			collect(c)
		}
	}
	for _, d := range dirs {
		collect(d)
	}
	return b.String()
}

func TestResolveRouteProfileNearestAncestorWins(t *testing.T) {
	profiles := map[string]string{
		"":             "root",
		"admin":        "outer",
		"admin/secret": "inner",
	}
	cases := map[string]string{
		"":                  "root",
		"blog":              "root",
		"admin":             "outer",
		"admin/secret":      "inner",
		"admin/secret/x":    "inner",
		"admin/public":      "outer",
		"admin/public/team": "outer",
	}
	for routeDir, want := range cases {
		if got := resolveRouteProfile(profiles, routeDir); got != want {
			t.Errorf("resolveRouteProfile(%q) = %q, want %q", routeDir, got, want)
		}
	}
}

func TestResolveRouteProfileEmptyWithoutAnyProfile(t *testing.T) {
	if got := resolveRouteProfile(map[string]string{}, "admin/secret"); got != "" {
		t.Fatalf("expected empty profile, got %q", got)
	}
}

func TestDiscoverRouteProfilesReadsRouteFolders(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/+page.dreego":                `<body>home</body>`,
		"routes/(hooks)/github/+page.dreego": "PROFILE \"hooks\"\n\n<body>hook</body>",
		"routes/admin/+page.dreego":          "PROFILE \"admin\"\n\n<body>admin</body>",
		"routes/admin/settings/+page.dreego": "PROFILE \"settings\"\n\n<body>settings</body>",
		"routes/public/+page.dreego":         "<body>public</body>",
	})
	profiles, err := discoverRouteProfiles(root)
	if err != nil {
		t.Fatalf("discoverRouteProfiles: %v", err)
	}
	if profiles["(hooks)/github"] != "hooks" {
		t.Errorf("group folder profile = %q, want hooks", profiles["(hooks)/github"])
	}
	if profiles["admin"] != "admin" {
		t.Errorf("admin profile = %q, want admin", profiles["admin"])
	}
	if profiles["admin/settings"] != "settings" {
		t.Errorf("settings profile = %q, want settings", profiles["admin/settings"])
	}
	if _, ok := profiles[""]; ok {
		t.Errorf("root without PROFILE must not be recorded")
	}
}

func TestDiscoverRouteProfilesRejectsConflictingProfileInFolder(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/hooks/+page.dreego": "PROFILE \"a\"\n\n<body>one</body>",
		"routes/hooks/other.dreego": "PROFILE \"b\"\n\n<body>two</body>",
	})
	_, err := discoverRouteProfiles(root)
	if err == nil {
		t.Fatal("expected conflicting PROFILE error")
	}
	for _, want := range []string{"hooks", `"a"`, `"b"`} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error must contain %q, got: %v", want, err)
		}
	}
}

func TestScanRoutesEmitsApplyProfileForGroupFolder(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/(hooks)/hooks.dreego":        "PROFILE \"hooks\"\n\n<body>group root</body>",
		"routes/(hooks)/github/+page.dreego": "<body>descendant inherits the group profile</body>",
		"routes/(hooks)/gitlab/+page.dreego": "<body>descendant inherits the group profile</body>",
		"routes/marketing/+page.dreego":      "<body>plain</body>",
	})
	dirs, _, _, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err != nil {
		t.Fatalf("scanRoutes: %v", err)
	}
	regs := allRouteRegs(dirs)
	for _, want := range []string{
		`app.ApplyProfile("/hooks", "hooks")`,
		`app.ApplyProfile("/github", "hooks")`,
		`app.ApplyProfile("/gitlab", "hooks")`,
	} {
		if !strings.Contains(regs, want) {
			t.Errorf("missing %s in: %s", want, regs)
		}
	}
	if strings.Contains(regs, `ApplyProfile("/marketing"`) {
		t.Errorf("unprofiled route must not receive ApplyProfile: %s", regs)
	}
}

func TestScanRoutesNearestProfileWins(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/admin/+page.dreego":             "PROFILE \"outer\"\n\n<body>admin</body>",
		"routes/admin/secret/+page.dreego":      "PROFILE \"inner\"\n\n<body>secret</body>",
		"routes/admin/secret/deep/+page.dreego": "<body>deep inherits inner</body>",
	})
	dirs, _, _, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err != nil {
		t.Fatalf("scanRoutes: %v", err)
	}
	regs := allRouteRegs(dirs)
	if !strings.Contains(regs, `app.ApplyProfile("/admin", "outer")`) {
		t.Errorf("missing outer profile for /admin: %s", regs)
	}
	if !strings.Contains(regs, `app.ApplyProfile("/admin/secret", "inner")`) {
		t.Errorf("missing inner profile for /admin/secret: %s", regs)
	}
	if !strings.Contains(regs, `app.ApplyProfile("/admin/secret/deep", "inner")`) {
		t.Errorf("nearest ancestor must win for nested folder: %s", regs)
	}
	if strings.Contains(regs, `app.ApplyProfile("/admin/secret", "outer")`) {
		t.Errorf("outer profile must not leak into the inner folder: %s", regs)
	}
}

func TestScanRoutesNoProfileEmitsNothing(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/+page.dreego":       "<body>home</body>",
		"routes/about/+page.dreego": "<body>about</body>",
	})
	dirs, _, _, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err != nil {
		t.Fatalf("scanRoutes: %v", err)
	}
	if strings.Contains(allRouteRegs(dirs), "ApplyProfile") {
		t.Fatalf("no PROFILE must emit no ApplyProfile: %s", allRouteRegs(dirs))
	}
}

func TestScanRoutesRootProfileAppliesToDescendants(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"routes/+page.dreego":       "PROFILE \"app\"\n\n<body>home</body>",
		"routes/about/+page.dreego": "<body>about</body>",
		"routes/admin/+page.dreego": "<body>admin inherits root</body>",
	})
	dirs, _, _, err := scanRoutes(NewGenerator(), root, map[string]*layoutEntry{}, map[string]*layoutEntry{})
	if err != nil {
		t.Fatalf("scanRoutes: %v", err)
	}
	regs := allRouteRegs(dirs)
	for _, want := range []string{
		`app.ApplyProfile("/{$}", "app")`,
		`app.ApplyProfile("/about", "app")`,
		`app.ApplyProfile("/admin", "app")`,
	} {
		if !strings.Contains(regs, want) {
			t.Errorf("missing %s in: %s", want, regs)
		}
	}
}

func TestParseRouteFileExposesProfile(t *testing.T) {
	file, _, err := parseRouteFile(NewGenerator(), "www/routes/hooks/+page.dreego", []byte("PROFILE \"hooks\"\n\n<body>x</body>"))
	if err != nil {
		t.Fatalf("parseRouteFile: %v", err)
	}
	if file.Profile != "hooks" {
		t.Fatalf("expected profile hooks, got %q", file.Profile)
	}
}

func TestRunGeneratesApplyProfileRegistrations(t *testing.T) {
	dir := writeTestProject(t, map[string]string{
		"go.mod":                                 "module example.com/site\n\ngo 1.27\n",
		"www/dreego.config.json":                 "{}",
		"www/routes/(hooks)/hooks.dreego":        "PROFILE \"hooks\"\n\n<body>webhook</body>",
		"www/routes/(hooks)/github/+page.dreego": "<body>github</body>",
		"www/routes/marketing/+page.dreego":      "<body>marketing</body>",
	})
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	if err := Run(false); err != nil {
		t.Fatalf("Run: %v", err)
	}
	root := filepath.Join(dir, "www", "routes", "dree.go")
	data, err := os.ReadFile(root)
	if err != nil {
		t.Fatal(err)
	}
	out := string(data)
	for _, want := range []string{
		`app.ApplyProfile("/hooks", "hooks")`,
		`app.ApplyProfile("/github", "hooks")`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("generated root dree.go missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, `ApplyProfile("/marketing"`) {
		t.Errorf("unprofiled route must not be bound:\n%s", out)
	}
}
