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

func routeDirRel(path string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if idx := strings.Index(rel, "dreego/routes/"); idx >= 0 {
		return rel[idx+len("dreego/routes/"):]
	}
	if idx := strings.Index(rel, "routes/"); idx >= 0 {
		return rel[idx+len("routes/"):]
	}
	if strings.HasSuffix(rel, "/routes") || rel == "routes" {
		return ""
	}
	return rel
}

func buildPageName(path string) string {
	rel := routeDirRel(path)
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

func buildPattern(path string) string {
	rel := routeDirRel(path)
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
	for _, g := range file.Go {
		if g.Method != "GET" {
			return []string{g.Method}
		}
	}
	return []string{"GET"}
}

func doubleBracketSegment(path string) string {
	rel := routeDirRel(path)
	for _, seg := range strings.Split(rel, "/") {
		if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
			return seg
		}
	}
	return ""
}

func detectGenDir(firstRoutePath string) string {
	idx := strings.Index(firstRoutePath, "dreego/routes")
	if idx < 0 {
		idx = strings.Index(firstRoutePath, "dreego"+string(filepath.Separator)+"routes")
	}
	if idx >= 0 {
		return firstRoutePath[:idx] + "dreego/gen"
	}
	return "dreego/gen"
}

func findGenDirFallback() string {
	return "dreego/gen"
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

func routeImports(src string) string {
	imports := []string{}
	if strings.Contains(src, "fmt.") {
		imports = append(imports, "\"fmt\"")
	}
	if src != "" {
		imports = append(imports, "\"net/http\"", "\"strings\"")
	}
	return strings.Join(imports, "\n\t")
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
