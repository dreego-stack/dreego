package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type routeDir = routePkg

func scanRoutes(gen *Generator, root string, layouts, layoutIndex map[string]*layoutEntry) ([]*routePkg, map[string]bool, int, error) {
	rootRoutes := filepath.Join(root, "routes")
	pkgs := map[string]*routePkg{}
	routePatterns := map[string]bool{}
	found := 0
	routeSources := map[string]string{}
	declSources := map[string]map[string]string{}

	getPkg := func(dir string) *routePkg {
		pkgDir := routePackageDir(root, dir)
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

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if !d.IsDir() {
			return nil
		}
		if !isRoutesDir(root, path) {
			return nil
		}
		entries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", path, err)
		}
		var dreegoFiles []string
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".dreego") {
				dreegoFiles = append(dreegoFiles, filepath.Join(path, e.Name()))
			}
		}
		if len(dreegoFiles) == 0 {
			return nil
		}

		p := getPkg(path)
		gen.Pkg = p.pkg
		gen.ImportKey = p.key
		gen.ImportKeySet = true
		decls := declSources[p.key]
		if decls == nil {
			decls = map[string]string{}
			declSources[p.key] = decls
		}

		var src strings.Builder
		var regs []string

		for _, fpath := range dreegoFiles {
			rel := routeFileRel(root, path, filepath.Base(fpath))
			pattern := buildPattern(rel)
			if seg := doubleBracketSegment(rel); seg != "" {
				return fmt.Errorf("optional segment %q in %s is not supported; define each route explicitly", seg, fpath)
			}
			pageName := buildPageName(rel)
			layout, err := resolveLayoutForRoute(rel, layouts, layoutIndex)
			if err != nil {
				return err
			}
			data, err := os.ReadFile(fpath)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", fpath, err)
			}
			baseName := strings.TrimSuffix(filepath.Base(fpath), ".dreego")
			method := "GET"

			file, raw, perr := parseRouteFile(gen, fpath, data)
			if perr != nil {
				return perr
			}
			if err := registerGoImports(gen, p.key, fpath, file.GoImports); err != nil {
				return err
			}
			for _, name := range hoistedDeclarationNames(file) {
				if prev, dup := decls[name]; dup {
					return serverDeclarationConflict(name, prev, fpath)
				}
				decls[name] = fpath
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
						return fmt.Errorf("duplicate catch-all %s: %s and %s", catchPattern, prev, fpath)
					}
					routeSources[catchKey] = fpath
				}
				s, reg, err := GenerateErrorHandler(gen, file, p.pkg, errCode, catchPattern, scopeHash)
				if err != nil {
					return fmt.Errorf("error generating error page %s: %w", fpath, err)
				}
				src.WriteString(s)
				regs = append(regs, reg)
				continue
			}

			for _, m := range fileRegisteredMethods(file) {
				key := m + " " + pattern
				if prev, dup := routeSources[key]; dup {
					return fmt.Errorf("duplicate route %s %s: %s and %s", m, pattern, prev, fpath)
				}
				routeSources[key] = fpath
				routePatterns[key] = true
			}

			s, reg, err := GenerateMethodHandler(gen, file, layout, p.pkg, pageName, pattern, scopeHash)
			if err != nil {
				return fmt.Errorf("error generating %s: %w", fpath, err)
			}
			src.WriteString(s)
			regs = append(regs, reg)
			if layout != nil {
				p.needsHead = true
			}
		}

		p.src.WriteString(src.String())
		p.regs = append(p.regs, regs...)
		found += len(regs)
		return nil
	})
	if err != nil {
		return nil, nil, 0, err
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
