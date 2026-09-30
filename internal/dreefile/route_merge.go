package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type routeFileAt struct {
	dirRel string
	name   string
	path   string
}

func relInRoutes(routesDir, path string) string {
	rel := relToRoot(routesDir, path)
	if rel == "." {
		return ""
	}
	return rel
}

// collectRouteFiles merges an app's routes/ tree with the website root's shared
// routes/ tree. A local file shadows a global file with the same relative path;
// a global file without a local counterpart is inherited. Local files win.
func collectRouteFiles(appRoot, websiteRoot string) ([]routeFileAt, error) {
	sources := []string{filepath.Join(appRoot, "routes")}
	if websiteRoot != "" {
		sources = append(sources, filepath.Join(websiteRoot, "routes"))
	}
	seen := map[string]bool{}
	var out []routeFileAt
	for _, src := range sources {
		if !hasDir(src) {
			continue
		}
		found, err := walkRouteFiles(src)
		if err != nil {
			return nil, err
		}
		for _, f := range found {
			key := f.dirRel + "\x00" + f.name
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, f)
		}
	}
	sortRouteFiles(out)
	return out, nil
}

func walkRouteFiles(src string) ([]routeFileAt, error) {
	var found []routeFileAt
	err := filepath.WalkDir(src, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if !d.IsDir() {
			return nil
		}
		if isLayoutsDir(src, path) {
			return filepath.SkipDir
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", path, err)
		}
		dirRel := relInRoutes(src, path)
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".dreego") {
				continue
			}
			found = append(found, routeFileAt{dirRel: dirRel, name: e.Name(), path: filepath.Join(path, e.Name())})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sortRouteFiles(found)
	return found, nil
}

func sortRouteFiles(files []routeFileAt) {
	sort.Slice(files, func(i, j int) bool {
		if files[i].dirRel != files[j].dirRel {
			return files[i].dirRel < files[j].dirRel
		}
		return files[i].name < files[j].name
	})
}

// mergedRouteProfiles reads PROFILE directives from the shared root routes and
// the app routes. The app value wins for a folder, matching the file override
// model; a conflicting PROFILE within one folder is still an error.
func mergedRouteProfiles(appRoot, websiteRoot string) (map[string]string, error) {
	merged := map[string]string{}
	for _, root := range []string{websiteRoot, appRoot} {
		if root == "" || !hasDir(filepath.Join(root, "routes")) {
			continue
		}
		profiles, err := discoverRouteProfiles(root)
		if err != nil {
			return nil, err
		}
		for rel, name := range profiles {
			merged[rel] = name
		}
	}
	return merged, nil
}

func routeFileRelFromDir(dirRel, name string) string {
	base := strings.TrimSuffix(name, ".dreego")
	if base == "+page" || base == "index" || base == "404" || base == "500" {
		return dirRel
	}
	if dirRel == "" {
		return base
	}
	return dirRel + "/" + base
}
