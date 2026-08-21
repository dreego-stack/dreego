package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type layoutEntry struct {
	rel    string
	source string
	file   *File
	name   string
}

func discoverLayouts(root string) (map[string]*layoutEntry, error) {
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
			tokens, lexErr := Lex(string(data))
			if lexErr != nil {
				return fmt.Errorf("error lexing layout %s: %w", full, lexErr)
			}
			f, parseErr := NewParser(tokens).Parse()
			if parseErr != nil {
				return fmt.Errorf("error parsing layout %s: %w", full, parseErr)
			}
			if f != nil {
				f.SourceContent = string(data)
				rel := layoutScopeRel(root, path)
				funcName := "Layout"
				if name == "default.dreego" {
					funcName = "Default"
				}
				entries[rel+":"+name] = &layoutEntry{rel: rel, source: full, file: f, name: funcName}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := detectAmbiguousLayouts(entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func detectAmbiguousLayouts(entries map[string]*layoutEntry) error {
	byScope := map[string][]*layoutEntry{}
	for _, e := range entries {
		byScope[e.rel] = append(byScope[e.rel], e)
	}
	var ambiguous []string
	for scope, list := range byScope {
		if len(list) > 1 {
			var files []string
			for _, e := range list {
				files = append(files, e.source)
			}
			sort.Strings(files)
			ambiguous = append(ambiguous, fmt.Sprintf("ambiguous layout in %s: %s", scope, strings.Join(files, ", ")))
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

func resolveLayoutForRoute(routeRel string, layouts map[string]*layoutEntry) *layoutEntry {
	routeRel = strings.TrimPrefix(routeRel, "/")
	scopes := cascadeScopes(routeRel)
	for _, scope := range scopes {
		for _, name := range []string{"default.dreego", "layout.dreego"} {
			if e, ok := layouts[scope+":"+name]; ok {
				return e
			}
		}
	}
	return nil
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

func generateLayouts(gen *Generator, root string, layouts map[string]*layoutEntry) ([]string, error) {
	var srcs []string
	scopes := map[string]bool{}
	for _, e := range layouts {
		scopes[e.rel] = true
	}
	var scopeList []string
	for s := range scopes {
		scopeList = append(scopeList, s)
	}
	sort.Strings(scopeList)

	for _, scope := range scopeList {
		for _, name := range []string{"default.dreego", "layout.dreego"} {
			e, ok := layouts[scope+":"+name]
			if !ok {
				continue
			}
			funcName := "Layout"
			if name == "default.dreego" {
				funcName = "Default"
			}
			src, err := GenerateLayout(gen, e.file, funcName)
			if err != nil {
				return nil, err
			}
			srcs = append(srcs, src)
		}
	}
	return srcs, nil
}
