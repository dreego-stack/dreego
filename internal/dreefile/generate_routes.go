package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type routeDir = routePkg

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
	localRoutes := filepath.Join(appRoot, "routes")
	globalRoutes := ""
	if websiteRoot != "" {
		globalRoutes = filepath.Join(websiteRoot, "routes")
	}
	seen := map[string]bool{}
	var out []routeFileAt
	for _, src := range []string{localRoutes, globalRoutes} {
		if !hasDir(src) {
			continue
		}
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
		sort.Slice(found, func(i, j int) bool {
			if found[i].dirRel != found[j].dirRel {
				return found[i].dirRel < found[j].dirRel
			}
			return found[i].name < found[j].name
		})
		for _, f := range found {
			key := f.dirRel + "\x00" + f.name
			if seen[key] {
				continue
			}
			seen[key] = true
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].dirRel != out[j].dirRel {
			return out[i].dirRel < out[j].dirRel
		}
		return out[i].name < out[j].name
	})
	return out, nil
}

// mergedRouteProfiles reads PROFILE directives from the shared root routes and
// the app routes. The app value wins for a folder; a genuine mismatch fails.
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
			if prev, ok := merged[rel]; ok && prev != name {
				return nil, fmt.Errorf("conflicting PROFILE in route folder %q: %q and %q", displayRouteFolder(rel), prev, name)
			}
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

// scanRoutes builds the app's route packages from its own routes/ tree plus the
// website root's shared routes/, where a local file overrides the global file
// with the same relative path.
func scanRoutes(gen *Generator, root, appName string, layouts, layoutIndex map[string]*layoutEntry, websiteRoot ...string) ([]*routePkg, map[string]bool, int, error) {
	shared := ""
	if len(websiteRoot) > 0 {
		shared = websiteRoot[0]
	}
	files, err := collectRouteFiles(root, shared)
	if err != nil {
		return nil, nil, 0, err
	}
	profiles, err := mergedRouteProfiles(root, shared)
	if err != nil {
		return nil, nil, 0, err
	}

	rootRoutes := filepath.Join(root, "routes")
	pkgs := map[string]*routePkg{}
	routePatterns := map[string]bool{}
	found := 0
	routeSources := map[string]string{}
	declSources := map[string]map[string]string{}

	getPkg := func(dirRel string) *routePkg {
		logical := filepath.Join(rootRoutes, filepath.FromSlash(dirRel))
		pkgDir := routePackageDir(root, logical)
		rel := routeDirRel(root, pkgDir)
		if p, ok := pkgs[rel]; ok {
			return p
		}
		pkg := "routes"
		if rel != "" {
			pkg = sanitizePkgName(filepath.Base(pkgDir))
		}
		p := &routePkg{dir: pkgDir, rel: rel, pkg: pkg, key: rel}
		pkgs[rel] = p
		return p
	}

	appliedProfiles := map[string]bool{}
	index := 0
	for index < len(files) {
		dirRel := files[index].dirRel
		var group []routeFileAt
		for index < len(files) && files[index].dirRel == dirRel {
			group = append(group, files[index])
			index++
		}

		p := getPkg(dirRel)
		gen.Pkg = p.pkg
		gen.ImportKey = p.key
		gen.ImportKeySet = true
		decls := declSources[p.key]
		if decls == nil {
			decls = map[string]string{}
			declSources[p.key] = decls
		}

		folderProfile := resolveRouteProfile(profiles, dirRel)
		var src strings.Builder
		var regs []string
		var profilePatterns []string

		for _, fpath := range group {
			rel := routeFileRelFromDir(dirRel, fpath.name)
			pattern := buildPattern(rel)
			if seg := doubleBracketSegment(rel); seg != "" {
				return nil, nil, 0, fmt.Errorf("optional segment %q in %s is not supported; define each route explicitly", seg, fpath.path)
			}
			pageName := buildPageName(rel)
			data, err := os.ReadFile(fpath.path)
			if err != nil {
				return nil, nil, 0, fmt.Errorf("error reading %s: %w", fpath.path, err)
			}
			baseName := strings.TrimSuffix(fpath.name, ".dreego")
			method := "GET"

			file, raw, perr := parseRouteFile(gen, fpath.path, data)
			if perr != nil {
				return nil, nil, 0, perr
			}
			if err := registerGoImports(gen, p.key, fpath.path, file.GoImports); err != nil {
				return nil, nil, 0, err
			}
			for _, name := range hoistedDeclarationNames(file) {
				if prev, dup := decls[name]; dup {
					return nil, nil, 0, serverDeclarationConflict(name, prev, fpath.path)
				}
				decls[name] = fpath.path
			}

			layout, err := resolveRouteLayout(file, appName, rel, layouts, layoutIndex)
			if err != nil {
				return nil, nil, 0, err
			}

			if len(file.Server) == 0 {
				file.Server = []ServerSection{{Method: method}}
			}
			for i := range file.Server {
				if !file.Server[i].MethodExplicit {
					file.Server[i].Method = method
				}
			}
			for i := range file.Bodies {
				if !file.Bodies[i].MethodExplicit {
					file.Bodies[i].Method = method
				}
			}

			scopeHash := hashOf(data)
			gen.Src = raw

			if baseName == "404" || baseName == "500" {
				errCode := 404
				if baseName == "500" {
					errCode = 500
				}
				catchPattern := errorCatchPattern(pattern)
				if errCode == 404 {
					catchKey := "GET" + " " + catchPattern
					if prev, dup := routeSources[catchKey]; dup {
						return nil, nil, 0, fmt.Errorf("duplicate catch-all %s: %s and %s", catchPattern, prev, fpath.path)
					}
					routeSources[catchKey] = fpath.path
				}
				s, reg, err := GenerateErrorHandler(gen, file, p.pkg, errCode, catchPattern, scopeHash)
				if err != nil {
					return nil, nil, 0, fmt.Errorf("error generating error page %s: %w", fpath.path, err)
				}
				src.WriteString(s)
				regs = append(regs, reg)
				if folderProfile != "" && !appliedProfiles[pattern] {
					appliedProfiles[pattern] = true
					profilePatterns = append(profilePatterns, pattern)
				}
				continue
			}

			for _, m := range fileRegisteredMethods(file) {
				key := m + " " + pattern
				if prev, dup := routeSources[key]; dup {
					return nil, nil, 0, fmt.Errorf("duplicate route %s %s: %s and %s", m, pattern, prev, fpath.path)
				}
				routeSources[key] = fpath.path
				routePatterns[key] = true
			}

			s, reg, err := GenerateMethodHandler(gen, file, layout, p.pkg, pageName, pattern, scopeHash)
			if err != nil {
				return nil, nil, 0, fmt.Errorf("error generating %s: %w", fpath.path, err)
			}
			src.WriteString(s)
			regs = append(regs, reg)
			if folderProfile != "" && !appliedProfiles[pattern] {
				appliedProfiles[pattern] = true
				profilePatterns = append(profilePatterns, pattern)
			}
			if layout != nil {
				p.needsHead = true
			}
		}

		p.src.WriteString(src.String())
		for _, pattern := range profilePatterns {
			p.regs = append(p.regs, registrationStatement(fmt.Sprintf("app.ApplyProfile(%q, %q)", pattern, folderProfile)))
		}
		p.regs = append(p.regs, regs...)
		found += len(regs)
	}

	if found == 0 {
		return nil, routePatterns, 0, nil
	}

	if _, ok := pkgs[""]; !ok {
		pkgs[""] = &routePkg{dir: rootRoutes, pkg: "routes", key: ""}
	}

	for rel, p := range pkgs {
		if rel == "" {
			continue
		}
		if parent := parentRoutePkg(pkgs, rel); parent != nil {
			parent.children = append(parent.children, p)
		}
	}

	list := make([]*routePkg, 0, len(pkgs))
	for _, p := range pkgs {
		if p.needsHead {
			p.src.WriteString(headMergeHelpers())
		}
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].rel < list[j].rel })
	return list, routePatterns, found, nil
}

func parentRoutePkg(pkgs map[string]*routePkg, rel string) *routePkg {
	for {
		idx := strings.LastIndex(rel, "/")
		if idx < 0 {
			return pkgs[""]
		}
		rel = rel[:idx]
		if p, ok := pkgs[rel]; ok {
			return p
		}
	}
}

func routeFileRel(root, dir, name string) string {
	rel := routeDirRel(root, dir)
	base := strings.TrimSuffix(name, ".dreego")
	if base == "+page" || base == "index" || base == "404" || base == "500" {
		return rel
	}
	if rel == "" {
		return base
	}
	return filepath.ToSlash(filepath.Join(rel, base))
}
