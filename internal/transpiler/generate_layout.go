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
}

func discoverLayouts() (map[string]*layoutEntry, error) {
	entries := map[string]*layoutEntry{}
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
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
		if !isDreegoLayoutsDir(path) {
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
				rel := layoutScopeRel(path)
				entries[rel+":"+name] = &layoutEntry{rel: rel, source: full, file: f}
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

func layoutScopeRel(layoutsDir string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(layoutsDir, "./"))
	idx := strings.Index(rel, "dreego/routes/")
	if idx < 0 {
		return ""
	}
	rest := rel[idx+len("dreego/routes/"):]
	return strings.TrimSuffix(rest, "/layouts")
}

func resolveLayoutForRoute(routePath string, layouts map[string]*layoutEntry) *File {
	routeRel := routeDirRel(routePath)
	routeRel = strings.TrimPrefix(routeRel, "/")
	scopes := cascadeScopes(routeRel)
	for _, scope := range scopes {
		for _, name := range []string{"default.dreego", "layout.dreego"} {
			if e, ok := layouts[scope+":"+name]; ok {
				return e.file
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
