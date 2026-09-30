package dreefile

import (
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
)

// layoutFuncName names a layout renderer after its scope and file so several
// layouts can coexist in the single generated layouts package. The default and
// legacy files keep the plain scope-based name; a named layout file adds its
// own PascalCased base name so an explicit LAYOUT path can select it.
func layoutFuncName(scopeKey, fileName string) string {
	base := "Layout"
	switch fileName {
	case "default.dreego":
		base = "Default"
	case "layout.dreego":
		base = "Layout"
	default:
		stem := strings.TrimSuffix(fileName, ".dreego")
		base = "Layout" + gogen.ToPascalCase(stem)
	}
	if scopeKey == "" || scopeKey == "root" {
		return base
	}
	return base + gogen.ToPascalCase(scopeKey)
}

// resolveRouteLayout picks the layout for a route. An explicit LAYOUT path wins
// over the directory cascade; resolveLayoutForRoute then tries the route-local
// scopes (deepest first), the app layout, and the shared root layout.
func resolveRouteLayout(file *File, appName, routeRel string, layouts, index map[string]*layoutEntry) (*layoutEntry, error) {
	if file != nil && file.Layout != "" {
		return resolveExplicitLayout(file, index)
	}
	return resolveLayoutForRoute(appName, routeRel, layouts, index)
}

func resolveExplicitLayout(file *File, index map[string]*layoutEntry) (*layoutEntry, error) {
	entry, ok := index[normaliseLayoutPath(file.Layout)]
	if !ok {
		return nil, layoutChainMissingError(&layoutEntry{file: file, source: file.SourcePath}, file.Layout)
	}
	return resolveLayoutEntry(entry, index)
}

// resolveLayoutForRoute resolves the layout for one route of an app. It tries
// the app's route-local scopes (deepest first), then a shared route-local scope
// under the root routes/ tree, then the app layout, then the shared root layout.
func resolveLayoutForRoute(appName, routeRel string, layouts, index map[string]*layoutEntry) (*layoutEntry, error) {
	routeRel = strings.TrimPrefix(routeRel, "/")
	scopes := cascadeScopes(routeRel)
	for _, scope := range scopes {
		key := appName
		if scope != "" {
			key = appName + "/" + scope
		}
		if e := pickLayout(layouts, key); e != nil {
			return resolveLayoutEntry(e, index)
		}
	}
	for _, scope := range scopes {
		if scope == "" {
			continue
		}
		if e := pickLayout(layouts, "root/"+scope); e != nil {
			return resolveLayoutEntry(e, index)
		}
	}
	if e := pickLayout(layouts, "root"); e != nil {
		return resolveLayoutEntry(e, index)
	}
	return nil, nil
}

func pickLayout(layouts map[string]*layoutEntry, scopeKey string) *layoutEntry {
	for _, name := range []string{"default.dreego", "layout.dreego"} {
		if e, ok := layouts[scopeKey+":"+name]; ok {
			return e
		}
	}
	return nil
}

func resolveLayoutEntry(e *layoutEntry, index map[string]*layoutEntry) (*layoutEntry, error) {
	if e.file == nil || e.file.Layout == "" {
		return e, nil
	}
	chain, err := resolveLayoutChain(e, index)
	if err != nil {
		return nil, err
	}
	if len(chain) > 1 {
		return nil, layoutChainNestedError(e, chain)
	}
	return chain[len(chain)-1], nil
}

func cascadeScopes(routeRel string) []string {
	if routeRel == "" {
		return []string{""}
	}
	parts := strings.Split(routeRel, "/")
	scopes := []string{""}
	cur := ""
	for _, p := range parts {
		if p == "" {
			continue
		}
		if cur == "" {
			cur = p
		} else {
			cur = cur + "/" + p
		}
		scopes = append(scopes, cur)
	}
	for i, j := 0, len(scopes)-1; i < j; i, j = i+1, j-1 {
		scopes[i], scopes[j] = scopes[j], scopes[i]
	}
	return scopes
}
