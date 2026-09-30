package dreefile

import (
	"path/filepath"
	"strings"
	"testing"
)

func relPaths(files []routeFileAt) []string {
	out := make([]string, 0, len(files))
	for _, f := range files {
		key := f.dirRel + "/" + f.name
		out = append(out, key)
	}
	return out
}

func TestCollectRouteFilesLocalOverridesGlobal(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"www/routes/+page.dreego":                "<body>global</body>",
		"www/routes/about/+page.dreego":          "<body>global about</body>",
		"www/routes/only-global/+page.dreego":    "<body>only global</body>",
		"www/app/routes/+page.dreego":            "<body>local</body>",
		"www/app/routes/only-local/+page.dreego": "<body>only local</body>",
	})
	files, err := collectRouteFiles(filepath.Join(root, "www", "app"), filepath.Join(root, "www"))
	if err != nil {
		t.Fatal(err)
	}
	got := relPaths(files)
	seen := map[string]int{}
	for _, key := range got {
		seen[key]++
	}
	for _, want := range []string{"/+page.dreego", "about/+page.dreego", "only-global/+page.dreego", "only-local/+page.dreego"} {
		if seen[want] != 1 {
			t.Errorf("expected exactly one entry for %q, got %d in %v", want, seen[want], got)
		}
	}
	for _, f := range files {
		if f.dirRel == "" && f.name == "+page.dreego" && !strings.Contains(f.path, "app") {
			t.Errorf("local +page must shadow the global one, got %q", f.path)
		}
	}
}

func TestCollectStaticFilesLocalWins(t *testing.T) {
	root := writeTestProject(t, map[string]string{
		"www/static/shared.css":      ".shared{}",
		"www/static/favicon.ico":     "global",
		"www/app/static/favicon.ico": "local",
		"www/app/static/app.css":     ".app{}",
	})
	files, err := collectStaticFiles(filepath.Join(root, "www", "app"), filepath.Join(root, "www"))
	if err != nil {
		t.Fatal(err)
	}
	byRel := map[string]string{}
	for _, f := range files {
		byRel[f.rel] = f.path
	}
	if byRel["shared.css"] != filepath.Join(root, "www", "static", "shared.css") {
		t.Errorf("global static must be inherited: %v", byRel)
	}
	if byRel["app.css"] != filepath.Join(root, "www", "app", "static", "app.css") {
		t.Errorf("local static must be included: %v", byRel)
	}
	if !strings.Contains(byRel["favicon.ico"], filepath.Join("app", "static")) {
		t.Errorf("local favicon must win, got %q", byRel["favicon.ico"])
	}
}
