package transpiler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func Run(force bool) error {
	start := time.Now()
	plan, stats, err := buildPlan(force)
	if err != nil {
		return err
	}
	if err := applyPlan(plan, force); err != nil {
		return err
	}
	elapsed := time.Since(start)
	fmt.Printf("Found %d routes + %d components + %d static\n", stats.routes, stats.components, stats.static)
	fmt.Printf("Generated %d dree.go files\n", len(plan.files))
	fmt.Printf("in %.2fms\n", float64(elapsed.Microseconds())/1000.0)
	return nil
}

func RunCheck() error {
	plan, _, err := buildPlan(false)
	if err != nil {
		return err
	}
	disk, err := readDiskFiles(plan.roots)
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

func buildPlan(force bool) (genPlan, genStats, error) {
	module := modulePath()
	roots, err := findWebsiteRoots()
	if err != nil {
		return genPlan{}, genStats{}, err
	}
	if len(roots) == 0 {
		return genPlan{}, genStats{}, fmt.Errorf("no website found: create a %s file in the website root directory", configFileName)
	}

	files := map[string]string{}
	var stats genStats
	for _, root := range roots {
		rootFiles, rootStats, err := buildRootPlan(root, module)
		if err != nil {
			return genPlan{}, genStats{}, err
		}
		for p, c := range rootFiles {
			files[p] = c
		}
		stats.routes += rootStats.routes
		stats.components += rootStats.components
		stats.static += rootStats.static
	}
	return genPlan{files: files, roots: roots}, stats, nil
}

func buildRootPlan(root, module string) (map[string]string, genStats, error) {
	gen := NewGenerator()
	gen.module = module
	gen.rootRel = relToRoot(".", root)
	gen.pkg = sanitizePkgName(filepath.Base(root))

	layouts, err := discoverLayouts(root)
	if err != nil {
		return nil, genStats{}, err
	}

	compSrcs, compPkgs, err := scanComponents(gen, root)
	if err != nil {
		return nil, genStats{}, err
	}

	routeDirs, routePatterns, routeCount, err := scanRoutes(gen, root, layouts)
	if err != nil {
		return nil, genStats{}, err
	}

	staticSrc, staticCount, err := generateStaticAssets(root, routePatterns)
	if err != nil {
		return nil, genStats{}, fmt.Errorf("static assets: %w", err)
	}

	settings := loadSettings(root)
	files := map[string]string{}

	for _, rd := range routeDirs {
		imports := gen.imports[rd.pkg]
		importLine := buildImportLine(imports, rd.pkg)
		stdImports := stdImportsFor(rd.src)
		coreImport := "dreego \"github.com/dreego-stack/dreego/core\""
		if strings.Contains(rd.src, "ssr.") {
			coreImport += "\n\tssr \"github.com/dreego-stack/dreego/core/ssr\""
		}
		out := fmt.Sprintf("package %s\n\nimport (\n\t%s\n\n\t%s\n)\n\n", rd.pkg, importLine, coreImport)
		if stdImports != "" {
			out = fmt.Sprintf("package %s\n\nimport (\n\t%s\n\t%s\n\n\t%s\n)\n\n", rd.pkg, stdImports, importLine, coreImport)
		}
		out += rd.src
		out += "func Register(app *dreego.App) error {\n"
		out += strings.Join(rd.regs, "")
		out += "\treturn nil\n}\n"
		files[filepath.Join(rd.dir, "dree.go")] = out
	}

	if len(compSrcs) > 0 {
		for pkgDir, srcs := range compSrcs {
			rel := relToRoot(root, pkgDir)
			pkg := sanitizePkgName(filepath.Base(pkgDir))
			imports := gen.imports[pkg]
			importLine := buildImportLine(imports, pkg)
			stdImports := stdImportsFor(strings.Join(srcs, ""))
			compOut := fmt.Sprintf("package %s\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", pkg, stdImports, importLine)
			compOut += strings.Join(srcs, "")
			files[filepath.Join(pkgDir, "dree.go")] = compOut
			_ = rel
		}
	}

	layoutSrcs, err := generateLayouts(gen, root, layouts)
	if err != nil {
		return nil, genStats{}, err
	}
	if len(layoutSrcs) > 0 {
		layoutDir := filepath.Join(root, "layouts")
		imports := gen.imports["layouts"]
		importLine := buildImportLine(imports, "layouts")
		stdImports := stdImportsFor(strings.Join(layoutSrcs, ""))
		layoutOut := fmt.Sprintf("package layouts\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", stdImports, importLine)
		layoutOut += strings.Join(layoutSrcs, "")
		layoutOut += headMergeHelpers()
		files[filepath.Join(layoutDir, "dree.go")] = layoutOut
	}

	rootOut := buildRootFile(root, module, routeDirs, staticSrc, settings)
	files[filepath.Join(root, "dree.go")] = rootOut

	return files, genStats{routes: routeCount, components: len(compPkgs), static: staticCount}, nil
}

func buildImportLine(imports map[string]string, selfPkg string) string {
	if len(imports) == 0 {
		return ""
	}
	var lines []string
	for alias, path := range imports {
		if alias == selfPkg {
			lines = append(lines, fmt.Sprintf("%q", path))
		} else {
			lines = append(lines, fmt.Sprintf("%s %q", alias, path))
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n\t")
}

func stdImportsFor(src string) string {
	var imports []string
	if strings.Contains(src, "strings.") {
		imports = append(imports, "\"strings\"")
	}
	if strings.Contains(src, "http.") {
		imports = append(imports, "\"net/http\"")
	}
	if strings.Contains(src, "fmt.") {
		imports = append(imports, "\"fmt\"")
	}
	return strings.Join(imports, "\n\t")
}

func buildRootFile(root, module string, routeDirs []routeDir, staticSrc string, settings *Settings) string {
	pkg := sanitizePkgName(filepath.Base(root))
	var imports []string
	var regCalls []string
	for _, rd := range routeDirs {
		imports = append(imports, fmt.Sprintf("%s %q", rd.pkg, module+"/"+relToRoot(".", root)+"/"+relToRoot(root, rd.dir)))
		regCalls = append(regCalls, fmt.Sprintf("\tif err := %s.Register(app); err != nil {\n\t\treturn err\n\t}\n", rd.pkg))
	}
	importLine := strings.Join(imports, "\n\t")

	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", pkg))
	if len(imports) > 0 {
		b.WriteString("import (\n\t" + importLine + "\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n")
	} else {
		b.WriteString("import (\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n")
	}
	b.WriteString("func Register(app *dreego.App) error {\n")
	if settings != nil {
		b.WriteString(registrationStatement(fmt.Sprintf("app.SetLogging(%t)", settings.Logging.Enabled)))
		for _, rd := range settings.Redirects {
			b.WriteString(registrationStatement(fmt.Sprintf("app.RegisterRedirect(%q, %q, %d)", rd.From, rd.To, rd.Status)))
		}
		for _, rw := range settings.Rewrites {
			b.WriteString(registrationStatement(fmt.Sprintf("app.RegisterRewrite(%q, %q)", rw.From, rw.To)))
		}
	}
	b.WriteString(staticSrc)
	for _, call := range regCalls {
		b.WriteString(call)
	}
	b.WriteString("\treturn nil\n}\n")
	return b.String()
}

func modulePath() string {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if rest, ok := strings.CutPrefix(line, "module "); ok {
			return strings.TrimSpace(rest)
		}
	}
	return ""
}

func loadSettings(root string) *Settings {
	settingsPath := filepath.Join(root, configFileName)
	s, err := LoadConfig(settingsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			slog.Warn("dreego: "+configFileName+" is invalid; using defaults",
				"path", settingsPath, "error", err)
		}
		return nil
	}
	return s
}

func isUpToDate(path, content string) bool {
	existing, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return string(existing) == content
}

func hashOf(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])[:12]
}
