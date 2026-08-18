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
	plan, stats, err := buildPlan(force)
	if err != nil {
		return err
	}
	if err := applyPlan(plan); err != nil {
		return err
	}
	elapsed := time.Since(start)
	fmt.Printf("Found %d routes + %d components + %d static\n", stats.routes, stats.components, stats.static)
	fmt.Printf("Generated gen/routes.go + gen/components.go + gen/dree.go\n")
	fmt.Printf("in %.2fms\n", float64(elapsed.Microseconds())/1000.0)
	return nil
}

func RunCheck() error {
	plan, _, err := buildPlan(false)
	if err != nil {
		return err
	}
	disk, err := readDiskFiles(plan.genDir)
	if err != nil {
		return err
	}
	diffs := plan.diff(disk)
	if len(diffs) == 0 {
		fmt.Println("generated code is up-to-date")
		return nil
	}
	return fmt.Errorf("generated code is out of date:\n%s", diffReport(diffs))
}

type genStats struct {
	routes     int
	components int
	static     int
}

func buildPlan(_ bool) (genPlan, genStats, error) {
	gen := NewGenerator()

	layouts, err := discoverLayouts()
	if err != nil {
		return genPlan{}, genStats{}, err
	}

	_, compSrcs, err := scanComponents(gen)
	if err != nil {
		return genPlan{}, genStats{}, err
	}

	var allSources []string
	var allRegistrations []string
	routePatterns := map[string]bool{}
	routeSources := map[string]string{}
	var found int
	var genDir string

	err = filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
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
		if isDreegoLayoutsDir(path) || isDreegoComponentsDir(path) {
			return nil
		}
		if isGeneratedDir(path) {
			return filepath.SkipDir
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
		if !isDreegoRoutesDir(path) {
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
		layout := resolveLayoutForRoute(path, layouts)

		for _, fpath := range dreegoFiles {
			data, err := os.ReadFile(fpath)
			if err != nil {
				return fmt.Errorf("error reading %s: %w", fpath, err)
			}
			h := sha256.Sum256(data)
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

			scopeHash := hex.EncodeToString(h[:])[:12]
			pkgName := filepath.Base(path)

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
					routePatterns[catchKey] = true
				}
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
		return genPlan{}, genStats{}, err
	}

	if genDir == "" {
		genDir = findGenDirFallback()
	}

	staticSrc, staticCount, err := generateStaticAssets(routePatterns)
	if err != nil {
		return genPlan{}, genStats{}, fmt.Errorf("static assets: %w", err)
	}

	settings := loadSettings(genDir)
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

	files := map[string]string{
		filepath.Join(genDir, "routes.go"): routesOut,
		filepath.Join(genDir, "dree.go"):   "package gen\n",
	}
	if compSrc != "" {
		compOut := fmt.Sprintf("package gen\n\nimport (\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", compImportLine)
		compOut += compSrc
		files[filepath.Join(genDir, "components.go")] = compOut
	}

	return genPlan{files: files, genDir: genDir}, genStats{routes: found, components: len(compSrcs), static: staticCount}, nil
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
