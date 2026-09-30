package dreefile

import (
	"path/filepath"
	"sort"
	"strings"
)

// layoutPackagesForApp groups the layouts an app resolves against into their Go
// packages (route-local scopes, app scope, and the shared root scope).
func layoutPackagesForApp(gen *Generator, all map[string]*layoutEntry, appName string) []layoutPackage {
	order := []string{}
	seen := map[string]bool{}
	for _, e := range all {
		if e == nil {
			continue
		}
		if e.scopeKey != "root" && e.appName != appName {
			continue
		}
		if !seen[e.scopeKey] {
			seen[e.scopeKey] = true
			order = append(order, e.scopeKey)
		}
	}
	sort.Strings(order)
	var pkgs []layoutPackage
	for _, scopeKey := range order {
		pkg := layoutPackagesForScope(gen, all, scopeKey)
		if len(pkg.entries) > 0 {
			pkgs = append(pkgs, pkg)
		}
	}
	return pkgs
}

func layoutPackagesForScope(gen *Generator, all map[string]*layoutEntry, scopeKey string) layoutPackage {
	pkg := layoutPackage{pkg: layoutPkgName(scopeKey)}
	for _, name := range []string{"default.dreego", "layout.dreego"} {
		e, ok := all[scopeKey+":"+name]
		if !ok {
			continue
		}
		pkg.dir = e.dir
		pkg.entries = append(pkg.entries, e)
	}
	if pkg.dir != "" {
		pkg.importPath = gen.Module + "/" + relToRoot(".", pkg.dir)
	}
	return pkg
}

// layoutPkgName maps a scope key to a Go package name: "root" -> "layouts", and
// otherwise a unique name derived from the scope (for example "blog_layouts").
func layoutPkgName(scopeKey string) string {
	if scopeKey == "" || scopeKey == "root" {
		return "layouts"
	}
	return sanitizePkgName(strings.ReplaceAll(scopeKey, "/", "_")) + "_layouts"
}

// layoutPkg returns the Go package name of an entry's layout package.
func layoutPkg(e *layoutEntry) string {
	return layoutPkgName(e.scopeKey)
}

// layoutImportPath returns the module-qualified import path of an entry's
// layout package.
func layoutImportPath(module, root string, e *layoutEntry) string {
	if module == "" {
		return relToRoot(".", e.dir)
	}
	return module + "/" + relToRoot(".", e.dir)
}

// generateLayouts emits the Go source for one layout package. The caller writes
// the surrounding package clause and file.
func generateLayouts(gen *Generator, pkgName string, entries []*layoutEntry) ([]string, error) {
	var srcs []string
	for _, name := range []string{"default.dreego", "layout.dreego"} {
		for _, e := range entries {
			if filepath.Base(e.source) != name {
				continue
			}
			gen.Src = e.file.SourceContent
			if err := registerGoImports(gen, pkgName, e.source, e.file.GoImports); err != nil {
				return nil, err
			}
			src, err := GenerateLayout(gen, e.file, e.name)
			if err != nil {
				return nil, err
			}
			srcs = append(srcs, src)
		}
	}
	return srcs, nil
}
