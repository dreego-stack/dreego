package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

)

var methodExt = map[string]string{
	"get":    "GET",
	"post":   "POST",
	"put":    "PUT",
	"delete": "DELETE",
}

func Run(force bool) error {
	start := time.Now()
	var found int

	layout := findLayout()
	var allSources []string

	var settings *Settings
	var genDir string

	_, compSrcs := scanComponents()

	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
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
		if isLayoutDir(path) {
			return nil
		}
		if base == "gen" && strings.Contains(path, "dreego/gen") {
			return filepath.SkipDir
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			return nil
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
		if !isRouteDir(path) {
			return nil
		}

		if genDir == "" {
			genDir = detectGenDir(path)
		}

		found++
		pattern := buildPattern(path)
		pageName := buildPageName(path)

		for _, fpath := range dreegoFiles {
			data, _ := os.ReadFile(fpath)
			h := sha256.Sum256(data)
			base := strings.TrimSuffix(filepath.Base(fpath), ".dreego")

			method := "GET"
			for prefix, m := range methodExt {
				if base == prefix || strings.HasPrefix(base, prefix+"-") {
					method = m
					break
				}
			}

			_, imports, body := ParseHeader(string(data))

			tokens, err := Lex(body)
			if err != nil {
				return fmt.Errorf("error lexing %s: %w", fpath, err)
			}

			p := NewParser(tokens)
			file, err := p.Parse()
			if err != nil {
				return fmt.Errorf("error parsing %s: %w", fpath, err)
			}
			file.Imports = imports

			if len(file.Go) == 0 {
				file.Go = []GoSection{{Method: method}}
			}
			for i := range file.Go {
				file.Go[i].Method = method
			}

			scopeHash := hex.EncodeToString(h[:])[:12]
			pkgName := filepath.Base(path)

			if base == "404" || base == "500" {
				errCode := 404
				if base == "500" {
					errCode = 500
				}
				catchPattern := errorCatchPattern(pattern)
				src, err := GenerateErrorHandler(file, pkgName, errCode, catchPattern, scopeHash)
				if err != nil {
					return fmt.Errorf("error generating error page %s: %w", fpath, err)
				}
				allSources = append(allSources, src)
				continue
			}

			for _, g := range file.Go {
				src, err := GenerateMethodHandler(file, layout, pkgName, pageName, pattern, g, scopeHash)
				if err != nil {
					return fmt.Errorf("error generating %s: %w", fpath, err)
				}
				allSources = append(allSources, src)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	if genDir == "" {
		genDir = findGenDirFallback()
	}

	if err := os.MkdirAll(genDir, 0755); err != nil {
		return err
	}

	src := strings.Join(allSources, "")
	compSrc := strings.Join(compSrcs, "")

	var imports []string
	if strings.Contains(src+compSrc, "fmt.") {
		imports = append(imports, "\"fmt\"")
	}
	if strings.Contains(src+compSrc, "html.EscapeString") {
		imports = append(imports, "\"html\"")
	}
	imports = append(imports, "\"net/http\"")
	imports = append(imports, "\"strings\"")
	importLine := strings.Join(imports, "\n\t")

	routesOut := fmt.Sprintf("package gen\n\nimport (\n\t%s\n\n\tcore \"codeberg.org/dreego/dreego/dreego-core\"\n)\n\n", importLine)
	routesOut += compSrc + src

	if err := os.WriteFile(filepath.Join(genDir, "routes.go"), []byte(routesOut), 0644); err != nil {
		return fmt.Errorf("error writing gen/routes.go: %w", err)
	}

	settings = loadSettings(genDir)
	if err := writeDreeGo(genDir, settings); err != nil {
		return fmt.Errorf("error writing gen/dree.go: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Found %d routes + %d components\n", found, len(compSrcs))
	fmt.Printf("Generated gen/routes.go + gen/dree.go\n")
	fmt.Printf("in %dns\n", elapsed.Nanoseconds())
	return nil
}

func writeDreeGo(genDir string, settings *Settings) error {
	var buf strings.Builder
	buf.WriteString("package gen\n\n")
	if settings != nil && (settings.Logging.Enabled || len(settings.Redirects) > 0 || len(settings.Rewrites) > 0) {
		buf.WriteString("import (\n")
		buf.WriteString("\tcore \"codeberg.org/dreego/dreego/dreego-core\"\n")
		buf.WriteString(")\n\n")
	}
	buf.WriteString("func init() {\n")
	if settings != nil {
		buf.WriteString(fmt.Sprintf("\tcore.SetLogging(%t)\n", settings.Logging.Enabled))
		for _, rd := range settings.Redirects {
			buf.WriteString(fmt.Sprintf("\tcore.RegisterRedirect(\"%s\", \"%s\", %d)\n", rd.From, rd.To, rd.Status))
		}
		for _, rw := range settings.Rewrites {
			buf.WriteString(fmt.Sprintf("\tcore.RegisterRewrite(\"%s\", \"%s\")\n", rw.From, rw.To))
		}
	}
	buf.WriteString("}\n")
	return os.WriteFile(filepath.Join(genDir, "dree.go"), []byte(buf.String()), 0644)
}

func loadSettings(genDir string) *Settings {
	settingsPath := filepath.Join(genDir, "..", "config.json")
	s, err := LoadConfig(settingsPath)
	if err != nil {
		return nil
	}
	return s
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

func isRouteDir(path string) bool {
	return strings.Contains(path, "routes/") || strings.HasSuffix(path, "/routes") || filepath.Base(path) == "routes"
}

func buildPageName(path string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if idx := strings.Index(rel, "routes/"); idx >= 0 {
		rel = rel[idx+len("routes/"):]
	} else if strings.HasSuffix(rel, "/routes") || rel == "routes" {
		rel = ""
	}
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
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if idx := strings.Index(rel, "routes/"); idx >= 0 {
		rel = rel[idx+len("routes/"):]
	} else if strings.HasSuffix(rel, "/routes") || rel == "routes" {
		rel = ""
	}
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
	if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
		return seg[1 : len(seg)-1]
	}
	if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
		return seg[1 : len(seg)-1]
	}
	if strings.HasPrefix(seg, "(") && strings.HasSuffix(seg, ")") {
		return seg[1 : len(seg)-1]
	}
	return seg
}

func patternSegment(seg string) string {
	if strings.HasPrefix(seg, "[") && strings.HasSuffix(seg, "]") {
		return "{" + seg[1:len(seg)-1] + "}"
	}
	if strings.HasPrefix(seg, "_") && strings.HasSuffix(seg, "_") {
		return "{" + seg[1:len(seg)-1] + "}"
	}
	if strings.HasPrefix(seg, "(") && strings.HasSuffix(seg, ")") {
		return ""
	}
	return seg
}

func findLayout() *File {
	var layout *File
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || layout != nil {
			return nil
		}
		if filepath.Base(path) == "default.dreego" || filepath.Base(path) == "layout.dreego" {
			data, _ := os.ReadFile(path)
			tokens, _ := Lex(string(data))
			p := NewParser(tokens)
			f, _ := p.Parse()
			if f != nil {
				layout = f
			}
		}
		return nil
	})
	return layout
}

func isLayoutDir(path string) bool {
	return strings.HasSuffix(path, "/layouts") || strings.Contains(path, "/layouts/") ||
		strings.HasSuffix(path, "/components") || strings.Contains(path, "/components/")
}

func scanComponents() (genDir string, sources []string) {
	filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".dreego") {
			return nil
		}
		if !strings.Contains(path, "/components/") && !strings.HasSuffix(filepath.Dir(path), "/components") {
			return nil
		}

		data, _ := os.ReadFile(path)
		raw := string(data)

		comp, _, body := ParseHeader(raw)
		if comp == nil || comp.Name == "" {
			return nil
		}

		h := sha256.Sum256(data)
		scopeHash := hex.EncodeToString(h[:])[:12]

		if genDir == "" {
			genDir = detectGenDir(path)
		}

		tokens, _ := Lex(body)
		if tokens == nil {
			return nil
		}

		p := NewParser(tokens)
		file, err := p.Parse()
		if err != nil {
			return nil
		}
		file.Component = comp

		if len(file.Go) == 0 {
			file.Go = []GoSection{{Method: ""}}
		}

		src, err := GenerateComponent(file, scopeHash)
		if err != nil {
			return nil
		}
		sources = append(sources, src)
		return nil
	})
	return
}
