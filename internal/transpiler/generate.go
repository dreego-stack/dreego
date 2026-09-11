package transpiler

import (
	"fmt"
	"maps"
	"path/filepath"
	"sort"
	"strings"
	"time"

	transpileri18n "github.com/dreego-stack/dreego/internal/transpiler/i18n"
	luainput "github.com/dreego-stack/dreego/internal/transpiler/js/lua"
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
		maps.Copy(files, rootFiles)
		stats.routes += rootStats.routes
		stats.components += rootStats.components
		stats.static += rootStats.static
	}
	return genPlan{files: files, roots: roots}, stats, nil
}

func buildRootPlan(root, module string) (map[string]string, genStats, error) {
	settings, err := loadSettings(root)
	if err != nil {
		return nil, genStats{}, err
	}
	gen := NewGenerator()
	gen.Module = module
	gen.RootRel = relToRoot(".", root)
	gen.Pkg = sanitizePkgName(filepath.Base(root))
	var catalogs transpileri18n.Set
	var generatedI18n string
	if settings != nil && settings.I18n.Enabled {
		catalogs, err = transpileri18n.Load(filepath.Join(root, "locales"), settings.I18n.Locales, settings.I18n.DefaultLocale)
		if err != nil {
			return nil, genStats{}, fmt.Errorf("i18n catalogs: %w", err)
		}
		gen.MessageArguments = transpileri18n.ArgumentKinds(catalogs)
		generatedI18n = transpileri18n.GoConfig(catalogs, settings.I18n.Detection, settings.I18n.URLStrategy, settings.I18n.Domains, settings.I18n.Fallbacks)
	}

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

	pluginSrc, pluginCount, err := generatePluginClientAssets(".", settings, routePatterns)
	if err != nil {
		return nil, genStats{}, err
	}
	staticSrc += pluginSrc
	staticCount += pluginCount
	files := map[string]string{}

	for _, rd := range routeDirs {
		imports := gen.Imports[rd.pkg]
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
			imports := gen.Imports[pkg]
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
		imports := gen.Imports["layouts"]
		importLine := buildImportLine(imports, "layouts")
		stdImports := stdImportsFor(strings.Join(layoutSrcs, ""))
		layoutOut := fmt.Sprintf("package layouts\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", stdImports, importLine)
		layoutOut += strings.Join(layoutSrcs, "")
		files[filepath.Join(layoutDir, "dree.go")] = layoutOut
	}

	if settings != nil && settings.I18n.Enabled {
		uses := make([]transpileri18n.Use, 0, len(gen.MessageUses))
		for _, use := range gen.MessageUses {
			uses = append(uses, transpileri18n.Use{Key: use.Key, Arguments: use.Arguments})
		}
		if err := transpileri18n.ValidateUses(catalogs, uses); err != nil {
			return nil, genStats{}, fmt.Errorf("i18n templates: %w", err)
		}
	} else if len(gen.MessageUses) > 0 {
		return nil, genStats{}, fmt.Errorf("i18n templates: enable i18n in %s before using message expressions", configFileName)
	}

	if runtime := luainput.BundleFeatures(gen.Lua); runtime != "" {
		path := "/_dreego/lua.js"
		if routePatterns["GET "+path] {
			return nil, genStats{}, fmt.Errorf("generated Lua runtime conflicts with route %q", path)
		}
		staticSrc += registrationStatement(fmt.Sprintf("app.RegisterStatic(%q, %q, []byte(%q))", path, "text/javascript; charset=utf-8", runtime))
		staticCount++
	}

	rootOut := buildRootFile(root, module, routeDirs, staticSrc, settings, generatedI18n)
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

func buildRootFile(root, module string, routeDirs []routeDir, staticSrc string, settings *Settings, i18nConfig ...string) string {
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
		if len(i18nConfig) > 0 && i18nConfig[0] != "" {
			b.WriteString(registrationStatement("app.SetI18n(" + i18nConfig[0] + ")"))
		}
	}
	b.WriteString(staticSrc)
	for _, call := range regCalls {
		b.WriteString(call)
	}
	b.WriteString("\treturn nil\n}\n")
	return b.String()
}
