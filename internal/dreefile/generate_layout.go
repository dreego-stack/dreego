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
	if err := detectLayoutNameCollisions(entries); err != nil {
		return nil, nil, err
	}
	return entries, buildLayoutIndex(root, entries), nil
}

// detectLayoutNameCollisions rejects two layouts that would generate the same
// Go renderer name in the single layouts package.
func detectLayoutNameCollisions(entries map[string]*layoutEntry) error {
	byName := map[string]string{}
	var collisions []string
	for _, e := range entries {
		if e == nil {
			continue
		}
		if prev, ok := byName[e.name]; ok {
			collisions = append(collisions, fmt.Sprintf("layout name collision: %s and %s both render as %s", prev, e.source, e.name))
			continue
		}
		byName[e.name] = e.source
	}
	if len(collisions) == 0 {
		return nil
	}
	sort.Strings(collisions)
	return fmt.Errorf("%s", strings.Join(collisions, "; "))
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
