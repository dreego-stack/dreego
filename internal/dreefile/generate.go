package dreefile

import (
	"fmt"
	"maps"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/dreego-stack/dreego/internal/dreefile/codegen"
	transpileri18n "github.com/dreego-stack/dreego/internal/dreefile/i18n"
	luainput "github.com/dreego-stack/dreego/internal/dreefile/sections/client/lua"
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
	if err := validateGeneratedMarkers(disk); err != nil {
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
	var rootsList []string
	for _, root := range roots {
		if len(root.apps) == 0 {
			return genPlan{}, genStats{}, fmt.Errorf("website root %s contains no apps: add an app directory with a routes/, static/, layouts/, components/ tree or a dreego.config.json (for example %s/www/routes/)", root.dir, root.dir)
		}
		rootFiles, rootStats, err := buildRootPlan(root, module)
		if err != nil {
			return genPlan{}, genStats{}, err
		}
		maps.Copy(files, rootFiles)
		stats.routes += rootStats.routes
		stats.components += rootStats.components
		stats.static += rootStats.static
		rootsList = append(rootsList, root.dir)
	}
	return genPlan{files: files, roots: rootsList}, stats, nil
}

// buildRootPlan generates the shared packages (layouts, components) of one
// website root and then every app package.
func buildRootPlan(root websiteRoot, module string) (map[string]string, genStats, error) {
	settings, err := loadSettings(root.dir)
	if err != nil {
		return nil, genStats{}, err
	}
	if err := validateUrlRules(settings); err != nil {
		return nil, genStats{}, err
	}

	gen := NewGenerator()
	gen.Module = module
	gen.Requires = moduleRequires()
	gen.RootRel = relToRoot(".", root.dir)
	gen.Pkg = sanitizePkgName(root.name)

	var catalogs transpileri18n.Set
	if settings != nil && settings.I18n.Enabled {
		catalogs, err = transpileri18n.Load(filepath.Join(root.dir, "locales"), settings.I18n.Locales, settings.I18n.DefaultLocale)
		if err != nil {
			return nil, genStats{}, fmt.Errorf("i18n catalogs: %w", err)
		}
		gen.MessageArguments = transpileri18n.ArgumentKinds(catalogs)
	}

	allLayouts, layoutIndex, err := discoverLayouts(root.dir)
	if err != nil {
		return nil, genStats{}, err
	}

	files := map[string]string{}

	if err := collectComponentAliases(gen, root.dir, root.apps); err != nil {
		return nil, genStats{}, err
	}
	compSrcs, compPkgs, err := scanComponents(gen, root.dir)
	if err != nil {
		return nil, genStats{}, err
	}
	moduleComps, err := loadModuleComponentsAcrossApps(gen, root.dir, root.apps)
	if err != nil {
		return nil, genStats{}, err
	}
	mergeComponentsByPkg(compSrcs, moduleComps)
	if err := validateComponentAliases(gen); err != nil {
		return nil, genStats{}, err
	}
	for pkgDir, srcs := range compSrcs {
		pkg := componentPkgName(root.dir, pkgDir)
		imports := gen.Imports[pkg]
		importLine := buildImportLine(imports, pkg)
		stdImports, err := stdImportsFor(gen, pkg, strings.Join(srcs, ""))
		if err != nil {
			return nil, genStats{}, err
		}
		compOut := fmt.Sprintf("package %s\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", pkg, stdImports, importLine)
		compOut += strings.Join(srcs, "")
		files[filepath.Join(pkgDir, "dree.go")] = withGeneratedMarker(relToRoot(".", pkgDir), compOut)
	}

	// A layout may reference a component that lives in the app-local
	// components/ tree of the app using it. Layouts are generated here in the
	// root pass, before the per-app component scan, so register every app's
	// local components for layout resolution. The root-only set is restored
	// afterwards so each app pass stays isolated (an app must not see a
	// sibling's components).
	rootDefs := maps.Clone(gen.Defs)
	rootCompPkgs := maps.Clone(gen.CompPkgs)
	rootCompPaths := maps.Clone(gen.CompPaths)
	if err := registerAppComponents(gen, root.dir, root.apps); err != nil {
		return nil, genStats{}, err
	}

	layoutPkg := gen.Pkg
	gen.Pkg = "layouts"
	layoutSrcs, err := generateLayouts(gen, root.dir, allLayouts)
	if err != nil {
		return nil, genStats{}, err
	}
	gen.Defs = rootDefs
	gen.CompPkgs = rootCompPkgs
	gen.CompPaths = rootCompPaths
	if len(layoutSrcs) > 0 {
		layoutDir := filepath.Join(root.dir, "layouts")
		imports := gen.Imports["layouts"]
		importLine := buildImportLine(imports, "layouts")
		stdImports, err := stdImportsFor(gen, "layouts", strings.Join(layoutSrcs, ""))
		if err != nil {
			return nil, genStats{}, err
		}
		layoutOut := fmt.Sprintf("package layouts\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", stdImports, importLine)
		layoutOut += strings.Join(layoutSrcs, "")
		if layoutNeedsHeadHelpers(layoutSrcs) {
			layoutOut += headMergeHelpers()
		}
		files[filepath.Join(layoutDir, "dree.go")] = withGeneratedMarker(relToRoot(".", layoutDir), layoutOut)
	}
	gen.Pkg = layoutPkg

	usesBefore := len(gen.MessageUses)

	var stats genStats
	stats.components = len(compPkgs)

	for _, app := range root.apps {
		appGen := NewGenerator()
		appGen.Module = module
		appGen.Requires = gen.Requires
		appGen.RootRel = gen.RootRel
		appGen.AppRel = relToRoot(".", app.dir)
		appGen.Defs = maps.Clone(gen.Defs)
		appGen.CompPkgs = maps.Clone(gen.CompPkgs)
		appGen.CompPaths = maps.Clone(gen.CompPaths)
		appGen.CompAliases = maps.Clone(gen.CompAliases)
		appGen.Lua = maps.Clone(gen.Lua)
		appGen.MessageArguments = gen.MessageArguments
		appStats, err := buildAppPlan(files, appGen, root, app, allLayouts, layoutIndex)
		if err != nil {
			return nil, genStats{}, err
		}
		stats.routes += appStats.routes
		stats.components += appStats.components
		stats.static += appStats.static
	}

	sharedUses := gen.MessageUses[:usesBefore]
	if err := validateRootI18nUses(gen, root, settings, sharedUses); err != nil {
		return nil, genStats{}, err
	}

	return files, stats, nil
}

// registerAppComponents registers the app-local components of every app into
// the shared generator so layout generation can resolve them. The returned
// sources are discarded: each app pass regenerates its own component files.
func registerAppComponents(gen *Generator, rootDir string, apps []appEntry) error {
	for _, app := range apps {
		if _, _, err := scanComponentsAt(gen, app.dir, rootDir, false); err != nil {
			return err
		}
	}
	return nil
}

// mergeComponentsByPkg merges extra package sources into a base map that is
// keyed by output directory, so shared and module components share one file.
func mergeComponentsByPkg(base, extra map[string][]string) {
	for dir, srcs := range extra {
		base[dir] = append(base[dir], srcs...)
	}
}

// buildAppPlan generates one app package: its routes, static assets, and the
// app entry file exporting a Registrar.
func buildAppPlan(files map[string]string, gen *Generator, root websiteRoot, app appEntry, allLayouts, layoutIndex map[string]*layoutEntry) (genStats, error) {
	settings, err := loadAppSettings(root.dir, app.dir)
	if err != nil {
		return genStats{}, err
	}
	if err := validateUrlRules(settings); err != nil {
		return genStats{}, err
	}

	var catalogs transpileri18n.Set
	generatedI18n := ""
	if settings != nil && settings.I18n.Enabled {
		catalogs, err = transpileri18n.Load(filepath.Join(root.dir, "locales"), settings.I18n.Locales, settings.I18n.DefaultLocale)
		if err != nil {
			return genStats{}, fmt.Errorf("i18n catalogs: %w", err)
		}
		gen.MessageArguments = transpileri18n.ArgumentKinds(catalogs)
		generatedI18n = transpileri18n.GoConfig(catalogs, settings.I18n.Detection, settings.I18n.URLStrategy, settings.I18n.Domains, settings.I18n.Fallbacks)
	}

	usesBefore := len(gen.MessageUses)
	gen.Pkg = app.pkg

	appComps, _, err := scanComponentsAt(gen, app.dir, root.dir, false)
	if err != nil {
		return genStats{}, err
	}
	if err := collectComponentAliases(gen, app.dir, nil); err != nil {
		return genStats{}, err
	}
	if err := validateComponentAliases(gen); err != nil {
		return genStats{}, err
	}
	appCompCount := 0
	for pkgDir, srcs := range appComps {
		pkg := componentPkgName(root.dir, pkgDir)
		imports := gen.Imports[pkg]
		importLine := buildImportLine(imports, pkg)
		stdImports, err := stdImportsFor(gen, pkg, strings.Join(srcs, ""))
		if err != nil {
			return genStats{}, err
		}
		compOut := fmt.Sprintf("package %s\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", pkg, stdImports, importLine)
		compOut += strings.Join(srcs, "")
		files[filepath.Join(pkgDir, "dree.go")] = withGeneratedMarker(relToRoot(".", pkgDir), compOut)
		appCompCount++
	}

	routeDirs, routePatterns, routeCount, err := scanRoutes(gen, app.dir, app.name, allLayouts, layoutIndex, root.dir)
	if err != nil {
		return genStats{}, err
	}
	appUses := gen.MessageUses[usesBefore:]

	if settings != nil && settings.I18n.Enabled {
		if err := transpileri18n.ValidateUses(catalogs, toI18nUses(appUses)); err != nil {
			return genStats{}, fmt.Errorf("i18n templates: %w", err)
		}
	} else if len(appUses) > 0 {
		return genStats{}, fmt.Errorf("i18n templates: enable i18n in %s before using message expressions", configFileName)
	}

	staticSrc, staticCount, err := generateStaticAssets(app.dir, root.dir, routePatterns)
	if err != nil {
		return genStats{}, fmt.Errorf("static assets: %w", err)
	}
	pluginSrc, pluginCount, err := generatePluginClientAssets(".", settings, routePatterns)
	if err != nil {
		return genStats{}, err
	}
	staticSrc += pluginSrc
	staticCount += pluginCount

	if runtime := luainput.BundleFeatures(gen.Lua); runtime != "" {
		path := "/_dreego/lua.js"
		if routePatterns["GET "+path] {
			return genStats{}, fmt.Errorf("generated Lua runtime conflicts with route %q", path)
		}
		staticSrc += registrationStatement(fmt.Sprintf("app.RegisterStatic(%q, %q, []byte(%q))", path, "text/javascript; charset=utf-8", runtime))
		staticCount++
	}

	for _, rd := range routeDirs {
		out, err := buildRoutePackageFile(gen, rd)
		if err != nil {
			return genStats{}, err
		}
		files[filepath.Join(rd.dir, "dree.go")] = out
	}

	files[filepath.Join(app.dir, "dree.go")] = withGeneratedMarker(relToRoot(".", app.dir), buildAppFile(app, gen, routeDirs, staticSrc, settings, generatedI18n))

	return genStats{routes: routeCount, components: appCompCount, static: staticCount}, nil
}

// validateRootI18nUses checks message expressions collected from the shared root
// layouts and components against the root i18n catalog.
func validateRootI18nUses(gen *Generator, root websiteRoot, settings *Settings, uses []codegen.MessageUse) error {
	if len(uses) == 0 {
		return nil
	}
	if settings == nil || !settings.I18n.Enabled {
		return fmt.Errorf("i18n templates: enable i18n in %s before using message expressions", configFileName)
	}
	catalogs, err := transpileri18n.Load(filepath.Join(root.dir, "locales"), settings.I18n.Locales, settings.I18n.DefaultLocale)
	if err != nil {
		return fmt.Errorf("i18n catalogs: %w", err)
	}
	if err := transpileri18n.ValidateUses(catalogs, toI18nUses(uses)); err != nil {
		return fmt.Errorf("i18n templates: %w", err)
	}
	return nil
}

func toI18nUses(uses []codegen.MessageUse) []transpileri18n.Use {
	out := make([]transpileri18n.Use, 0, len(uses))
	for _, use := range uses {
		out = append(out, transpileri18n.Use{Key: use.Key, Arguments: use.Arguments})
	}
	return out
}

func layoutNeedsHeadHelpers(srcs []string) bool {
	return strings.Contains(strings.Join(srcs, ""), "dedupeLayoutHead(")
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

// buildAppFile generates the app package entry point. It exports a Registrar so
// a caller can write `dreego.New(www.App)`.
func buildAppFile(app appEntry, gen *Generator, routeDirs []*routePkg, staticSrc string, settings *Settings, i18nConfig ...string) string {
	var imports []string
	var regCalls []string
	for _, rd := range routeDirs {
		if rd.rel != "" {
			continue
		}
		path := gen.Module + "/" + gen.AppRel + "/" + relToRoot(app.dir, rd.dir)
		imports = append(imports, fmt.Sprintf("routes %q", path))
		regCalls = append(regCalls, "\tif err := routes.Register(app); err != nil {\n\t\treturn err\n\t}\n")
	}
	importLine := strings.Join(imports, "\n\t")

	var b strings.Builder
	b.WriteString(fmt.Sprintf("package %s\n\n", app.pkg))
	if len(imports) > 0 {
		b.WriteString("import (\n\t" + importLine + "\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n")
	} else {
		b.WriteString("import (\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n")
	}
	b.WriteString("// App registers this app on the given application.\n")
	b.WriteString("var App dreego.Registrar = registerApp\n\n")
	b.WriteString("func registerApp(app *dreego.App) error {\n")
	if settings != nil {
		if settings.Logging.Present {
			b.WriteString(registrationStatement(fmt.Sprintf("app.SetLogging(%t)", settings.Logging.Enabled)))
		}
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
