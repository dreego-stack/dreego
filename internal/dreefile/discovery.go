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

// findWebsiteRoots locates every website root in the current tree. A website
// root is a directory with dreego.config.json that does not itself contain a
// routes/ directory (a root-level routes/ is the legacy single-site layout and
// is rejected). Apps are the immediate subdirectories containing routes/.
func findWebsiteRoots() ([]websiteRoot, error) {
	var configDirs []string
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

	var roots []websiteRoot
	for _, dir := range configDirs {
		if hasDir(filepath.Join(dir, "routes")) {
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

	for _, dir := range configDirs {
		if !hasDir(filepath.Join(dir, "routes")) {
			continue
		}
		if !isAppOfAnyRoot(dir, roots) {
			return nil, legacyLayoutError(dir)
		}
	}
	if isWebsiteRoot(".") && hasDir("routes") && !isAppOfAnyRoot(".", roots) {
		return nil, legacyLayoutError(".")
	}
	return roots, nil
}

func legacyLayoutError(dir string) error {
	example := filepath.Join(dir, "www", "routes")
	if filepath.Clean(dir) == "." {
		example = filepath.Join("www", "routes")
	}
	return fmt.Errorf("legacy website layout detected at %s: a website root no longer contains routes/ directly; create an app subdirectory (for example %s) instead", dir, example)
}

func isAppOfAnyRoot(dir string, roots []websiteRoot) bool {
	parent := filepath.Dir(dir)
	for _, root := range roots {
		if filepath.Clean(root.dir) == filepath.Clean(parent) {
			return true
		}
	}
	return false
}

// scanApps returns the app subdirectories of a website root, sorted by name. An
// app requires a routes/ directory; static/ is optional.
func scanApps(root string) ([]appEntry, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, fmt.Errorf("error reading website root %s: %w", root, err)
	}
	var apps []appEntry
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || isSkippedDir(e.Name()) {
			continue
		}
		dir := filepath.Join(root, e.Name())
		if !hasDir(filepath.Join(dir, "routes")) {
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

func isComponentsDir(root, path string) bool {
	rel := relToRoot(root, path)
	return rel == "components" || strings.HasPrefix(rel, "components/")
}

func relToRoot(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}
	return filepath.ToSlash(rel)
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
