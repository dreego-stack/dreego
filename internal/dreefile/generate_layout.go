package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
)

type layoutEntry struct {
	scopeKey string
	appName  string
	source   string
	dir      string
	file     *File
	name     string
}

// layoutPackage is one generated Go package holding layouts of a single scope.
type layoutPackage struct {
	pkg        string
	importPath string
	dir        string
	entries    []*layoutEntry
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
		if filepath.Base(path) != "layouts" {
			return nil
		}
		scopeKey, ok := layoutScopeKey(root, path)
		if !ok {
			return nil
		}
		appName := ""
		if scopeKey != "root" {
			appName = strings.SplitN(scopeKey, "/", 2)[0]
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
			if name != "default.dreego" && name != "layout.dreego" {
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
			funcName := "Layout"
			if name == "default.dreego" {
				funcName = "Default"
			}
			entries[scopeKey+":"+name] = &layoutEntry{
				scopeKey: scopeKey,
				appName:  appName,
				source:   full,
				dir:      path,
				file:     f,
				name:     funcName,
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
func layoutScopeKey(root, dir string) (string, bool) {
	rel := relToRoot(root, dir)
	if rel == "layouts" {
		return "root", true
	}
	parts := strings.Split(rel, "/")
	if len(parts) == 2 && parts[1] == "layouts" {
		return parts[0], true
	}
	if len(parts) >= 4 && parts[1] == "routes" && parts[len(parts)-1] == "layouts" {
		sub := strings.Join(parts[2:len(parts)-1], "/")
		return parts[0] + "/" + sub, true
	}
	return "", false
}

func detectAmbiguousLayouts(entries map[string]*layoutEntry) error {
	var ambiguous []string
	byDir := map[string][]string{}
	for _, e := range entries {
		if e == nil {
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

// resolveLayoutForRoute resolves the layout for one route of an app. It tries
// route-local scopes (deepest first), then the app scope, then the shared root
// scope.
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
	if appName != "root" {
		if e := pickLayout(layouts, "root"); e != nil {
			return resolveLayoutEntry(e, index)
		}
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
