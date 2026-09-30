package dreefile

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	transpileri18n "github.com/dreego-stack/dreego/internal/dreefile/i18n"
	luainput "github.com/dreego-stack/dreego/internal/dreefile/sections/client/lua"
)

// buildAppPlan generates one app package: its components, its app-local
// layouts, its routes, its static assets, and the app entry file.
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

	appComps, _, err := scanComponentsWithBase(gen, app.dir, root.dir, false)
	if err != nil {
		return genStats{}, err
	}
	if err := collectComponentAliases(gen, app.dir, nil); err != nil {
		return genStats{}, err
	}
	if err := validateComponentAliases(gen); err != nil {
		return genStats{}, err
	}
	compCount, err := writeComponentPackages(files, gen, appComps)
	if err != nil {
		return genStats{}, err
	}

	if err := writeAppLocalLayouts(files, gen, app.name, allLayouts); err != nil {
		return genStats{}, err
	}

	gen.Pkg = app.pkg
	routeDirs, routePatterns, routeCount, err := scanRoutes(gen, app.dir, app.name, allLayouts, layoutIndex)
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

	staticSrc, staticCount, err := generateStaticAssets(app.dir, routePatterns)
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
		files[filepath.Join(rd.dir, "dree.go")] = buildRouteFile(gen, rd)
	}
	files[filepath.Join(app.dir, "dree.go")] = buildAppFile(app, gen, routeDirs, staticSrc, settings, generatedI18n)

	return genStats{routes: routeCount, components: compCount, static: staticCount}, nil
}

// writeAppLocalLayouts emits the Go packages for an app's own layouts (app-local
// and route-local). The shared root layouts package is generated once in the
// root pass.
func writeAppLocalLayouts(files map[string]string, gen *Generator, appName string, allLayouts map[string]*layoutEntry) error {
	for _, lp := range layoutPackagesForApp(gen, allLayouts, appName) {
		if lp.pkg == "layouts" {
			continue
		}
		layoutPkg := gen.Pkg
		gen.Pkg = lp.pkg
		layoutSrcs, err := generateLayouts(gen, lp.pkg, lp.entries)
		if err != nil {
			return err
		}
		if len(layoutSrcs) > 0 {
			imports := gen.Imports[lp.pkg]
			importLine := buildImportLine(imports, lp.pkg)
			stdImports := stdImportsFor(gen, lp.pkg, strings.Join(layoutSrcs, ""))
			layoutOut := fmt.Sprintf("package %s\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", lp.pkg, stdImports, importLine)
			layoutOut += strings.Join(layoutSrcs, "")
			if layoutNeedsHeadHelpers(layoutSrcs) {
				layoutOut += headMergeHelpers()
			}
			files[filepath.Join(lp.dir, "dree.go")] = layoutOut
		}
		gen.Pkg = layoutPkg
	}
	return nil
}

// mergeComponentPackages combines component packages that share an output
// directory (the shared root components and module-imported components both land
// there) into a single package entry so they are written to one file.
func mergeComponentPackages(base, extra []componentPackage) []componentPackage {
	byDir := map[string]*componentPackage{}
	var order []string
	for _, cp := range append(append([]componentPackage{}, base...), extra...) {
		if existing, ok := byDir[cp.dir]; ok {
			existing.srcs = append(existing.srcs, cp.srcs...)
			continue
		}
		copy := cp
		byDir[cp.dir] = &copy
		order = append(order, cp.dir)
	}
	sort.Strings(order)
	out := make([]componentPackage, 0, len(order))
	for _, dir := range order {
		out = append(out, *byDir[dir])
	}
	return out
}

func writeComponentPackages(files map[string]string, gen *Generator, compPkgs []componentPackage) (int, error) {
	if len(compPkgs) == 0 {
		return 0, nil
	}
	count := 0
	for _, cp := range compPkgs {
		imports := gen.Imports[cp.pkg]
		importLine := buildImportLine(imports, cp.pkg)
		stdImports := stdImportsFor(gen, cp.pkg, strings.Join(cp.srcs, ""))
		compOut := fmt.Sprintf("package %s\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", cp.pkg, stdImports, importLine)
		compOut += strings.Join(cp.srcs, "")
		files[filepath.Join(cp.dir, "dree.go")] = compOut
		count++
	}
	return count, nil
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

// buildRouteFile wraps one route directory's generated sources in a package with
// a Register function.
func buildRouteFile(gen *Generator, rd routeDir) string {
	imports := gen.Imports[rd.pkg]
	importLine := buildImportLine(imports, rd.pkg)
	stdImports := stdImportsFor(gen, rd.pkg, rd.src)
	coreImport := "dreego \"github.com/dreego-stack/dreego/core\""
	if strings.Contains(rd.src, "ssr.") {
		coreImport += "\n\tssr \"github.com/dreego-stack/dreego/adapter/ssr\""
	}
	out := fmt.Sprintf("package %s\n\nimport (\n\t%s\n\n\t%s\n)\n\n", rd.pkg, importLine, coreImport)
	if stdImports != "" {
		out = fmt.Sprintf("package %s\n\nimport (\n\t%s\n\t%s\n\n\t%s\n)\n\n", rd.pkg, stdImports, importLine, coreImport)
	}
	out += rd.src
	out += "func Register(app *dreego.App) error {\n"
	out += strings.Join(rd.regs, "")
	out += "\treturn nil\n}\n"
	return out
}

// buildAppFile generates the app package entry point. It exports an App
// Registrar so a caller can write `dreego.New(www.App)`.
func buildAppFile(app appEntry, gen *Generator, routeDirs []routeDir, staticSrc string, settings *Settings, i18nConfig ...string) string {
	var imports []string
	var regCalls []string
	for _, rd := range routeDirs {
		imports = append(imports, fmt.Sprintf("%s %q", rd.pkg, gen.Module+"/"+relToRoot(".", rd.dir)))
		regCalls = append(regCalls, fmt.Sprintf("\tif err := %s.Register(app); err != nil {\n\t\treturn err\n\t}\n", rd.pkg))
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
