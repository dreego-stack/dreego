package dreefile

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
)

type layoutEntry struct {
	scopeKey string
	rel      string
	source   string
	dir      string
	file     *File
	name     string
}

// discoverLayouts walks the website root and returns every layout keyed by
// "<scopeKey>:<filename>". The scopeKey is "root" for <root>/layouts, the app
// name for <root>/<app>/layouts, and "<app>/<sub>" for route-local layouts under
// <root>/<app>/routes/<sub>/layouts.
func discoverLayouts(root string) (map[string]*layoutEntry, map[string]*layoutEntry, error) {
	entries := map[string]*layoutEntry{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if !d.IsDir() {
			return nil
		}
		if !isLayoutsDir(root, path) {
			return nil
		}
		scopeKey, ok := layoutScopeKey(root, path)
		if !ok {
			return nil
		}
		dirEntries, err := os.ReadDir(path)
		if err != nil {
			return fmt.Errorf("error reading directory %s: %w", path, err)
		}
		for _, e := range dirEntries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if !strings.HasSuffix(name, ".dreego") {
				continue
			}
			full := filepath.Join(path, name)
			data, readErr := os.ReadFile(full)
			if readErr != nil {
				return fmt.Errorf("error reading layout %s: %w", full, readErr)
			}
			raw := string(data)
			header, body, headerErr := ParseFileHeaderStrict(raw)
			if headerErr != nil {
				return fmt.Errorf("%s:%w", full, headerErr)
			}
			tokens, lexErr := Lex(body)
			if lexErr != nil {
				return fmt.Errorf("error lexing layout %s: %w", full, lexErr)
			}
			f, parseErr := NewParser(tokens).Parse()
			if parseErr != nil {
				return fmt.Errorf("error parsing layout %s: %w", full, parseErr)
			}
			if f == nil {
				continue
			}
			f.Imports = header.Imports
			f.Kind = header.Kind
			f.Layout = header.Layout
			f.Profile = header.Profile
			f.GoImports = header.GoImports
			f.SourceContent = raw
			f.SourcePath = full
			bodyOffset := len(raw) - len(body)
			if f.Client != nil {
				f.Client.Pos += bodyOffset
			}
			if f.Body != nil {
				gogen.SetNodeSource(f.Body.Nodes, full, bodyOffset)
				gogen.SetSourceText(f.Body.Nodes, raw)
			}
			entries[scopeKey+":"+name] = &layoutEntry{
				scopeKey: scopeKey,
				rel:      layoutScopeRel(root, path),
				source:   full,
				dir:      path,
				file:     f,
				name:     layoutFuncName(scopeKey, name),
			}
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	if err := detectAmbiguousLayouts(entries); err != nil {
		return nil, nil, err
	}
	return entries, buildLayoutIndex(root, entries), nil
}

// layoutScopeKey classifies a layouts directory relative to the website root.
// It returns "root" for the shared root layouts, "<app>" for an app layout,
// "<app>/<sub>" for an app route-local layout, and "root/<sub>" for a
// route-local layout inside the shared root routes/ tree.
func layoutScopeKey(root, dir string) (string, bool) {
	rel := relToRoot(root, dir)
	if rel == "layouts" {
		return "root", true
	}
	parts := strings.Split(rel, "/")
	if parts[0] == "routes" {
		if len(parts) >= 2 && parts[len(parts)-1] == "layouts" {
			sub := strings.Join(parts[1:len(parts)-1], "/")
			if sub == "" {
				return "root", true
			}
			return "root/" + sub, true
		}
		return "", false
	}
	if len(parts) == 2 && parts[1] == "layouts" {
		return parts[0], true
	}
	if len(parts) >= 4 && parts[1] == "routes" && parts[len(parts)-1] == "layouts" {
		sub := strings.Join(parts[2:len(parts)-1], "/")
		return parts[0] + "/" + sub, true
	}
	return "", false
}

// detectAmbiguousLayouts rejects two default-ish layouts in one directory. A
// directory may hold several explicitly named layouts alongside one default.
func detectAmbiguousLayouts(entries map[string]*layoutEntry) error {
	var ambiguous []string
	byDir := map[string][]string{}
	for _, e := range entries {
		if e == nil {
			continue
		}
		name := filepath.Base(e.source)
		if name != "default.dreego" && name != "layout.dreego" {
			continue
		}
		byDir[e.dir] = append(byDir[e.dir], e.source)
	}
	for dir, files := range byDir {
		if len(files) > 1 {
			sort.Strings(files)
			ambiguous = append(ambiguous, fmt.Sprintf("ambiguous layout in %s: %s", dir, strings.Join(files, ", ")))
		}
	}
	if len(ambiguous) == 0 {
		return nil
	}
	sort.Strings(ambiguous)
	return fmt.Errorf("%s", strings.Join(ambiguous, "; "))
}

func layoutScopeRel(root, layoutsDir string) string {
	rel := relToRoot(root, layoutsDir)
	if rel == "layouts" {
		return ""
	}
	rel = strings.TrimSuffix(rel, "/layouts")
	rel = strings.TrimPrefix(rel, "routes/")
	return rel
}

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

// generateLayouts emits the Go source for every layout. The caller writes the
// surrounding package clause; all scopes share the single layouts package.
func generateLayouts(gen *Generator, root string, layouts map[string]*layoutEntry) ([]string, error) {
	var srcs []string
	gen.ImportKey = ""
	gen.ImportKeySet = false
	keys := slices.Sorted(maps.Keys(layouts))

	for _, key := range keys {
		e := layouts[key]
		if e == nil {
			continue
		}
		gen.Src = e.file.SourceContent
		if err := registerGoImports(gen, "layouts", e.source, e.file.GoImports); err != nil {
			return nil, err
		}
		src, err := GenerateLayout(gen, e.file, e.name)
		if err != nil {
			return nil, err
		}
		srcs = append(srcs, src)
	}
	return srcs, nil
}
