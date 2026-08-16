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

type Generator struct {
	defs map[string]*ComponentDef
	src  string
}

func NewGenerator() *Generator {
	return &Generator{defs: map[string]*ComponentDef{}}
}

func (g *Generator) registerDef(name string, def *ComponentDef) {
	g.defs[name] = def
}

func (g *Generator) lookupDef(name string) *ComponentDef {
	return g.defs[name]
}

type generator = Generator

func Run(force bool) error {
	start := time.Now()
	var found int

	gen := NewGenerator()

	layout, err := findLayout()
	if err != nil {
		return err
	}
	var allSources []string
	var allRegistrations []string

	var settings *Settings
	var genDir string

	routePatterns := map[string]bool{}
	routeSources := map[string]string{}

	_, compSrcs, err := scanComponents(gen)
	if err != nil {
		return err
	}

	err = filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
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
		if seg := doubleBracketSegment(path); seg != "" {
			return fmt.Errorf("optional segment %q in %s is not supported; define each route explicitly", seg, path)
		}
		routePatterns["GET"+" "+pattern] = true
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

			raw := string(data)
			_, imports, body := ParseHeader(raw)
			bodyOffset := len(raw) - len(body)

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
			file.SourceContent = raw
			if file.Template != nil {
				setNodeSource(file.Template.Nodes, fpath, bodyOffset)
				file.FormActions = scanFormActions(file.Template.Nodes)
			}

			if len(file.Go) == 0 {
				file.Go = []GoSection{{Method: method}}
			}
			for i := range file.Go {
				if !file.Go[i].MethodExplicit {
					file.Go[i].Method = method
				}
			}

			scopeHash := hex.EncodeToString(h[:])[:12]
			pkgName := filepath.Base(path)

			gen.src = raw
			if base == "404" || base == "500" {
				errCode := 404
				if base == "500" {
					errCode = 500
				}
				catchPattern := errorCatchPattern(pattern)
				src, reg, err := GenerateErrorHandler(gen, file, pkgName, errCode, catchPattern, scopeHash)
				if err != nil {
					return fmt.Errorf("error generating error page %s: %w", fpath, err)
				}
				allSources = append(allSources, src)
				allRegistrations = append(allRegistrations, reg)
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
			src, reg, err := GenerateMethodHandler(gen, file, layout, pkgName, pageName, pattern, scopeHash)
			if err != nil {
				return fmt.Errorf("error generating %s: %w", fpath, err)
			}
			allSources = append(allSources, src)
			allRegistrations = append(allRegistrations, reg)
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

	staticSrc, staticCount, err := generateStaticAssets(routePatterns)
	if err != nil {
		return fmt.Errorf("static assets: %w", err)
	}

	settings = loadSettings(genDir)

	src := strings.Join(allSources, "")
	compSrc := strings.Join(compSrcs, "")

	routeImportLine := routeImports(src)

	compImports := []string{}
	if strings.Contains(compSrc, "fmt.") {
		compImports = append(compImports, "\"fmt\"")
	}
	compImports = append(compImports, "\"strings\"")
	compImportLine := strings.Join(compImports, "\n\t")

	var configReg strings.Builder
	if settings != nil {
		configReg.WriteString(registrationStatement(fmt.Sprintf("app.SetLogging(%t)", settings.Logging.Enabled)))
		for _, rd := range settings.Redirects {
			configReg.WriteString(registrationStatement(fmt.Sprintf("app.RegisterRedirect(%q, %q, %d)", rd.From, rd.To, rd.Status)))
		}
		for _, rw := range settings.Rewrites {
			configReg.WriteString(registrationStatement(fmt.Sprintf("app.RegisterRewrite(%q, %q)", rw.From, rw.To)))
		}
	}
	configReg.WriteString(staticSrc)

	routesOut := fmt.Sprintf("package gen\n\nimport (\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", routeImportLine)
	routesOut += src
	routesOut += "func Register(app *dreego.App) error {\n"
	routesOut += configReg.String()
	routesOut += strings.Join(allRegistrations, "")
	routesOut += "\treturn nil\n}\n"

	if !isUpToDate(filepath.Join(genDir, "routes.go"), routesOut) {
		if err := os.WriteFile(filepath.Join(genDir, "routes.go"), []byte(routesOut), 0644); err != nil {
			return fmt.Errorf("error writing gen/routes.go: %w", err)
		}
	}

	if compSrc != "" {
		compOut := fmt.Sprintf("package gen\n\nimport (\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", compImportLine)
		compOut += compSrc
		if !isUpToDate(filepath.Join(genDir, "components.go"), compOut) {
			if err := os.WriteFile(filepath.Join(genDir, "components.go"), []byte(compOut), 0644); err != nil {
				return fmt.Errorf("error writing gen/components.go: %w", err)
			}
		}
	}

	if err := touchDreeGo(genDir); err != nil {
		return fmt.Errorf("error writing gen/dree.go: %w", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Found %d routes + %d components + %d static\n", found, len(compSrcs), staticCount)
	fmt.Printf("Generated gen/routes.go + gen/components.go + gen/dree.go\n")
	fmt.Printf("in %.2fms\n", float64(elapsed.Microseconds())/1000.0)
	return nil
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

func touchDreeGo(genDir string) error {
	dreeGo := "package gen\n"
	if !isUpToDate(filepath.Join(genDir, "dree.go"), dreeGo) {
		if err := os.WriteFile(filepath.Join(genDir, "dree.go"), []byte(dreeGo), 0644); err != nil {
			return err
		}
	}
	return nil
}

func isUpToDate(path, content string) bool {
	existing, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return string(existing) == content
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

// fileRegisteredMethods returns the HTTP methods GenerateMethodHandler will
// register for a route file. It mirrors the codegen decision: a form action
// registers GET plus POST only when the action has both a form struct and a
// handler function; otherwise exactly one method is registered (the first
// non-GET method, or GET).
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

// doubleBracketSegment returns the name of the first double-bracket segment in
// a route directory path, or "" if the path contains none. Double brackets
// were the historical optional-segment syntax; optional segments are not
// supported, so generation rejects them with a source-aware diagnostic.
func doubleBracketSegment(path string) string {
	rel := filepath.ToSlash(strings.TrimPrefix(path, "./"))
	if idx := strings.Index(rel, "routes/"); idx >= 0 {
		rel = rel[idx+len("routes/"):]
	} else if strings.HasSuffix(rel, "/routes") || rel == "routes" {
		rel = ""
	}
	for _, seg := range strings.Split(rel, "/") {
		if strings.HasPrefix(seg, "[[") && strings.HasSuffix(seg, "]]") {
			return seg
		}
	}
	return ""
}

func findLayout() (*File, error) {
	var layout *File
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || layout != nil {
			return nil
		}
		if filepath.Base(path) == "default.dreego" || filepath.Base(path) == "layout.dreego" {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading layout %s: %w", path, err)
			}
			tokens, err := Lex(string(data))
			if err != nil {
				return fmt.Errorf("error lexing layout %s: %w", path, err)
			}
			p := NewParser(tokens)
			f, err := p.Parse()
			if err != nil {
				return fmt.Errorf("error parsing layout %s: %w", path, err)
			}
			if f != nil {
				f.SourceContent = string(data)
				layout = f
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return layout, nil
}

func isLayoutDir(path string) bool {
	return strings.HasSuffix(path, "/layouts") || strings.Contains(path, "/layouts/") ||
		strings.HasSuffix(path, "/components") || strings.Contains(path, "/components/")
}

func generateStaticAssets(routePatterns map[string]bool) (src string, count int, err error) {
	staticDir := filepath.Join("dreego", "static")
	if _, e := os.Stat(staticDir); os.IsNotExist(e) {
		return "", 0, nil
	}

	var buf strings.Builder

	err = filepath.WalkDir(staticDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil || d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(staticDir, path)
		urlPath := "/" + filepath.ToSlash(rel)

		methodPattern := "GET" + " " + urlPath
		if routePatterns[methodPattern] {
			return fmt.Errorf("static file %q conflicts with existing route %q", path, urlPath)
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}

		ext := filepath.Ext(path)
		mime := MimeByExt(ext)

		content := []byte(data)
		buf.WriteString(registrationStatement(fmt.Sprintf("app.RegisterStatic(%q, %q, %#v)", urlPath, mime, content)))
		count++
		return nil
	})

	if err != nil {
		return "", 0, err
	}

	return buf.String(), count, nil
}
