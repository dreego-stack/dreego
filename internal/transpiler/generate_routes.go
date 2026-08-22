package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var methodExt = map[string]string{
	"get":    "GET",
	"post":   "POST",
	"put":    "PUT",
	"delete": "DELETE",
}

type routeDir struct {
	dir  string
	pkg  string
	src  string
	regs []string
}

func scanRoutes(gen *Generator, root string, layouts map[string]*layoutEntry) ([]routeDir, map[string]bool, int, error) {
	rd := &routeDir{dir: filepath.Join(root, "routes"), pkg: "routes"}
	routePatterns := map[string]bool{}
	found := 0
	routeSources := map[string]string{}

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

		gen.pkg = "routes"

		var src strings.Builder
		var regs []string
		needsHeadHelpers := false

		for _, fpath := range dreegoFiles {
			rel := routeFileRel(root, path, filepath.Base(fpath))
			pattern := buildPattern(rel)
			if seg := doubleBracketSegment(rel); seg != "" {
				return fmt.Errorf("optional segment %q in %s is not supported; define each route explicitly", seg, fpath)
			}
			pageName := buildPageName(rel)
			layout := resolveLayoutForRoute(rel, layouts)
			data, err := os.ReadFile(fpath)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", fpath, err)
			}
			baseName := strings.TrimSuffix(filepath.Base(fpath), ".dreego")
			method := methodForFile(baseName)

			file, raw, perr := parseRouteFile(gen, fpath, data)
			if perr != nil {
				return perr
			}

			if len(file.Go) == 0 {
				file.Go = []GoSection{{Method: method}}
			}
			for i := range file.Go {
				if !file.Go[i].MethodExplicit {
					file.Go[i].Method = method
				}
			}
			for i := range file.Templates {
				if !file.Templates[i].MethodExplicit {
					file.Templates[i].Method = method
				}
			}

			scopeHash := hashOf(data)
			gen.src = raw

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
				s, reg, err := GenerateErrorHandler(gen, file, "routes", errCode, catchPattern, scopeHash)
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

			s, reg, err := GenerateMethodHandler(gen, file, layout, "routes", pageName, pattern, scopeHash)
			if err != nil {
				return fmt.Errorf("error generating %s: %w", fpath, err)
			}
			src.WriteString(s)
			regs = append(regs, reg)
			if layout != nil {
				needsHeadHelpers = true
			}
		}

		if needsHeadHelpers {
			src.WriteString(headMergeHelpers())
		}

		rd.src += src.String()
		rd.regs = append(rd.regs, regs...)
		found += len(regs)
		return nil
	})
	if err != nil {
		return nil, nil, 0, err
	}

	if found == 0 {
		return nil, routePatterns, 0, nil
	}
	return []routeDir{*rd}, routePatterns, found, nil
}

func routeFileRel(root, dir, name string) string {
	rel := routeDirRel(root, dir)
	base := strings.TrimSuffix(name, ".dreego")
	if base == "+page" || base == "index" || base == "404" || base == "500" || isLegacyMethodFile(base) {
		return rel
	}
	if rel == "" {
		return base
	}
	return filepath.ToSlash(filepath.Join(rel, base))
}

func isLegacyMethodFile(base string) bool {
	for prefix := range methodExt {
		if base == prefix || strings.HasPrefix(base, prefix+"-") {
			return true
		}
	}
	return false
}

func buildPageName(rel string) string {
	parts := []string{}
	for _, seg := range strings.Split(rel, "/") {
		if seg == "" {
			continue
		}
		parts = append(parts, cleanSegment(seg))
	}
	if len(parts) == 0 {
		return "index"
	}
	return strings.Join(parts, "_")
}

func buildPattern(rel string) string {
	if rel == "" || rel == "." {
		return "/{$}"
	}
	segments := []string{}
	for _, seg := range strings.Split(rel, "/") {
		if seg == "" {
			continue
		}
		s := patternSegment(seg)
		if s != "" {
			segments = append(segments, s)
		}
	}
	if len(segments) == 0 {
		return "/{$}"
	}
	return "/" + strings.Join(segments, "/")
}

func errorCatchPattern(dirPattern string) string {
	if strings.HasSuffix(dirPattern, "/{$}") {
		return strings.TrimSuffix(dirPattern, "/{$}") + "/{p...}"
	}
	return dirPattern + "/{p...}"
}

func cleanSegment(seg string) string {
	for {
		if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
			return ""
		}
		if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
			seg = seg[1 : len(seg)-1]
			continue
		}
		if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
			seg = seg[1 : len(seg)-1]
			continue
		}
		if strings.HasPrefix(seg, "(") && strings.HasSuffix(seg, ")") {
			seg = seg[1 : len(seg)-1]
			continue
		}
		return seg
	}
}

func patternSegment(seg string) string {
	if strings.HasPrefix(seg, "(") && strings.HasSuffix(seg, ")") {
		return ""
	}
	if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
		return ""
	}
	wrapped := false
	for {
		if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
			seg = seg[1 : len(seg)-1]
			wrapped = true
			continue
		}
		if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
			seg = seg[1 : len(seg)-1]
			wrapped = true
			continue
		}
		break
	}
	if !wrapped {
		return seg
	}
	if seg == "" {
		return ""
	}
	if strings.HasPrefix(seg, "...") {
		return "{" + strings.TrimPrefix(seg, "...") + "...}"
	}
	return "{" + seg + "}"
}

func fileRegisteredMethods(file *File) []string {
	if len(file.FormActions) > 0 {
		action := file.FormActions[0]
		if findFormStruct(file.Go, action) != "" && findFormHandler(file.Go, action) {
			return []string{"GET", "POST"}
		}
		return []string{"GET"}
	}
	seen := map[string]bool{}
	methods := []string{}
	add := func(method string) {
		if method == "" {
			method = "GET"
		}
		if !seen[method] {
			seen[method] = true
			methods = append(methods, method)
		}
	}
	for _, g := range file.Go {
		add(g.Method)
	}
	for _, t := range file.Templates {
		add(t.Method)
	}
	if len(methods) == 0 {
		add("GET")
	}
	return methods
}

func doubleBracketSegment(rel string) string {
	for _, seg := range strings.Split(rel, "/") {
		if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
			return seg
		}
	}
	return ""
}

func methodForFile(base string) string {
	method := "GET"
	for prefix, m := range methodExt {
		if base == prefix || strings.HasPrefix(base, prefix+"-") {
			method = m
			break
		}
	}
	return method
}

func parseRouteFile(gen *Generator, fpath string, data []byte) (*File, string, error) {
	raw := string(data)
	_, imports, body := ParseHeader(raw)
	tokens, err := Lex(body)
	if err != nil {
		return nil, "", fmt.Errorf("error lexing %s: %w", fpath, err)
	}
	p := NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return nil, "", fmt.Errorf("error parsing %s: %w", fpath, err)
	}
	file.Imports = imports
	file.SourceContent = raw
	bodyOffset := len(raw) - len(body)
	if file.Template != nil {
		setNodeSource(file.Template.Nodes, fpath, bodyOffset)
		setSourceText(file.Template.Nodes, raw)
		file.FormActions = scanFormActions(file.Template.Nodes)
		for _, d := range a11yDiagnostics(file.Template.Nodes) {
			fmt.Fprintf(os.Stderr, "warning: %s\n", d)
		}
	}
	return file, raw, nil
}
