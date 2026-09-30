package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const configFileName = "dreego.config.json"

func isSkippedDir(name string) bool {
	switch name {
	case "vendor", "node_modules", ".git", ".worktrees", ".tmp":
		return true
	}
	return strings.HasPrefix(name, ".") && name != "." && name != ".."
}

func isWebsiteRoot(path string) bool {
	info, err := os.Stat(filepath.Join(path, configFileName))
	return err == nil && !info.IsDir()
}

func hasDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// appEntry is one app subdirectory of a website root. An app is a directory that
// contains a routes/ tree; each app generates its own Go package.
type appEntry struct {
	name   string
	dir    string
	pkg    string
	routes string
	static string
	config string
}

// websiteRoot is a directory holding dreego.config.json with one or more apps.
type websiteRoot struct {
	dir    string
	name   string
	config string
	apps   []appEntry
}

// rootOwnedDirs are the directories a website root owns itself. A root's
// routes/, static/, layouts/, and components/ are shared defaults for every app;
// a reserved child name is never treated as an app directory. locales/ holds
// i18n catalogs and is likewise not an app.
func rootOwnedDirs(name string) bool {
	switch name {
	case "routes", "static", "layouts", "components", "locales":
		return true
	}
	return false
}

// findWebsiteRoots locates every website root in the current tree. A website
// root is a directory with dreego.config.json; it may itself hold shared
// routes/, static/, layouts/, and components/ as defaults for its apps. An app
// is an immediate subdirectory that carries any of those trees or its own
// dreego.config.json; a config directory that is an immediate child of another
// website root is an app, not a new root.
func findWebsiteRoots() ([]websiteRoot, error) {
	var configDirs []string
	if isWebsiteRoot(".") {
		configDirs = append(configDirs, ".")
	}
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !d.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if base == "." {
			return nil
		}
		if strings.HasPrefix(base, ".") {
			return filepath.SkipDir
		}
		if isSkippedDir(base) {
			return filepath.SkipDir
		}
		if isWebsiteRoot(path) {
			configDirs = append(configDirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(configDirs)

	isConfigDir := map[string]bool{}
	for _, dir := range configDirs {
		isConfigDir[filepath.Clean(dir)] = true
	}

	var roots []websiteRoot
	for _, dir := range configDirs {
		clean := filepath.Clean(dir)
		parent := filepath.Clean(filepath.Dir(clean))
		if parent != clean && isConfigDir[parent] {
			continue
		}
		apps, err := scanApps(dir)
		if err != nil {
			return nil, err
		}
		roots = append(roots, websiteRoot{
			dir:    dir,
			name:   filepath.Base(dir),
			config: filepath.Join(dir, configFileName),
			apps:   apps,
		})
	}
	return roots, nil
}

// scanApps returns the app subdirectories of a website root, sorted by name. An
// app carries a routes/, static/, layouts/, or components/ tree, or its own
// dreego.config.json; everything else is a root-owned directory.
func scanApps(root string) ([]appEntry, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("error reading website root %s: %w", root, err)
	}
	var apps []appEntry
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || isSkippedDir(e.Name()) || rootOwnedDirs(e.Name()) {
			continue
		}
		dir := filepath.Join(root, e.Name())
		if !isAppDir(dir) {
			continue
		}
		apps = append(apps, appEntry{
			name:   e.Name(),
			dir:    dir,
			pkg:    sanitizePkgName(e.Name()),
			routes: filepath.Join(dir, "routes"),
			static: filepath.Join(dir, "static"),
			config: filepath.Join(dir, configFileName),
		})
	}
	sort.Slice(apps, func(i, j int) bool { return apps[i].name < apps[j].name })
	return apps, nil
}

// isAppDir reports whether dir is an app of a website root: it carries at least
// one of the app-owned trees, or its own dreego.config.json override.
func isAppDir(dir string) bool {
	for _, sub := range []string{"routes", "static", "layouts", "components"} {
		if hasDir(filepath.Join(dir, sub)) {
			return true
		}
	}
	return isWebsiteRoot(dir)
}

func isRoutesDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "routes" || strings.HasPrefix(rel, "routes/")
}

func isComponentsDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "components" || strings.HasPrefix(rel, "components/")
}

func isLayoutsDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "layouts" || strings.HasPrefix(rel, "layouts/") || strings.HasSuffix(rel, "/layouts")
}

func isStaticDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "static" || strings.HasPrefix(rel, "static/")
}

func relToRoot(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	return filepath.ToSlash(rel)
}

func routeDirRel(root, path string) string {
	rel := relToRoot(root, path)
	if rel == "routes" {
		return ""
	}
	return strings.TrimPrefix(rel, "routes/")
}

func sanitizePkgName(name string) string {
	var b strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			b.WriteRune(r)
		case r == '-', r == '.', r == ' ':
			b.WriteRune('_')
		}
	}
	s := b.String()
	if s == "" {
		return "app"
	}
	if s[0] >= '0' && s[0] <= '9' {
		s = "pkg" + s
	}
	return s
}
