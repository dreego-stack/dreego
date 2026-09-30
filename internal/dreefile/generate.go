package dreefile

import (
	"fmt"
	"maps"
	"path/filepath"
	"strings"
	"time"

	"github.com/dreego-stack/dreego/internal/dreefile/codegen"
	transpileri18n "github.com/dreego-stack/dreego/internal/dreefile/i18n"
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
	var rootDirs []string
	for _, root := range roots {
		if len(root.apps) == 0 {
			return genPlan{}, genStats{}, fmt.Errorf("website root %s contains no apps: add an app directory with a routes/ tree (for example %s/www/routes/)", root.dir, root.dir)
		}
		rootDirs = append(rootDirs, root.dir)
		rootFiles, rootStats, err := buildRootPlan(root, module)
		if err != nil {
			return genPlan{}, genStats{}, err
		}
		maps.Copy(files, rootFiles)
		stats.routes += rootStats.routes
		stats.components += rootStats.components
		stats.static += rootStats.static
	}
	return genPlan{files: files, roots: rootDirs}, stats, nil
}

// buildRootPlan generates the shared packages and every app package of one
// website root. Shared layouts/ and components/ live directly in the root; each
// app is a subdirectory with its own routes/ and static/ and becomes its own Go
// package exporting `var App dreego.Registrar`.
func buildRootPlan(root websiteRoot, module string) (map[string]string, genStats, error) {
	gen := NewGenerator()
	gen.Module = module
	gen.RootRel = relToRoot(".", root.dir)
	gen.Pkg = sanitizePkgName(root.name)

	files := map[string]string{}

	allLayouts, layoutIndex, err := discoverLayouts(root.dir)
	if err != nil {
		return nil, genStats{}, err
	}
	gen.LayoutPkg = "layouts"
	gen.LayoutImportPath = module + "/" + relToRoot(".", filepath.Join(root.dir, "layouts"))

	rootSettings, err := loadSettings(root.dir)
	if err != nil {
		return nil, genStats{}, err
	}
	if err := validateUrlRules(rootSettings); err != nil {
		return nil, genStats{}, err
	}
	if rootSettings != nil && rootSettings.I18n.Enabled {
		rootCatalogs, err := transpileri18n.Load(filepath.Join(root.dir, "locales"), rootSettings.I18n.Locales, rootSettings.I18n.DefaultLocale)
		if err != nil {
			return nil, genStats{}, fmt.Errorf("i18n catalogs: %w", err)
		}
		gen.MessageArguments = transpileri18n.ArgumentKinds(rootCatalogs)
	}

	sharedComps, _, err := scanComponents(gen, root.dir)
	if err != nil {
		return nil, genStats{}, err
	}
	moduleComps, err := loadModuleComponentsAcrossApps(gen, root.dir, root.apps)
	if err != nil {
		return nil, genStats{}, err
	}
	modulePkgs, err := generateComponentSources(gen, moduleComps)
	if err != nil {
		return nil, genStats{}, err
	}
	sharedComps = mergeComponentPackages(sharedComps, modulePkgs)
	if err := collectComponentAliases(gen, root.dir, root.apps); err != nil {
		return nil, genStats{}, err
	}
	if err := validateComponentAliases(gen); err != nil {
		return nil, genStats{}, err
	}
	compCount, err := writeComponentPackages(files, gen, sharedComps)
	if err != nil {
		return nil, genStats{}, err
	}

	rootLayoutPkg := layoutPackagesForScope(gen, allLayouts, "root")
	if len(rootLayoutPkg.entries) > 0 {
		layoutPkg := gen.Pkg
		gen.Pkg = "layouts"
		layoutSrcs, err := generateLayouts(gen, "layouts", rootLayoutPkg.entries)
		if err != nil {
			return nil, genStats{}, err
		}
		if len(layoutSrcs) > 0 {
			imports := gen.Imports["layouts"]
			importLine := buildImportLine(imports, "layouts")
			stdImports := stdImportsFor(gen, "layouts", strings.Join(layoutSrcs, ""))
			layoutOut := fmt.Sprintf("package layouts\n\nimport (\n\t%s\n\t%s\n\n\tdreego \"github.com/dreego-stack/dreego/core\"\n)\n\n", stdImports, importLine)
			layoutOut += strings.Join(layoutSrcs, "")
			if layoutNeedsHeadHelpers(layoutSrcs) {
				layoutOut += headMergeHelpers()
			}
			files[filepath.Join(root.dir, "layouts", "dree.go")] = layoutOut
		}
		gen.Pkg = layoutPkg
	}

	var stats genStats
	stats.components = compCount

	if err := validateRootI18nUses(gen, root, rootSettings); err != nil {
		return nil, genStats{}, err
	}

	for _, app := range root.apps {
		appGen := NewGenerator()
		appGen.Module = module
		appGen.RootRel = gen.RootRel
		appGen.LayoutPkg = "layouts"
		appGen.LayoutImportPath = gen.LayoutImportPath
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

	return files, stats, nil
}

// validateRootI18nUses checks message expressions collected from the shared
// root components and layouts against the root i18n catalog.
func validateRootI18nUses(gen *Generator, root websiteRoot, settings *Settings) error {
	if len(gen.MessageUses) == 0 {
		return nil
	}
	if settings == nil || !settings.I18n.Enabled {
		return fmt.Errorf("i18n templates: enable i18n in %s before using message expressions", configFileName)
	}
	catalogs, err := transpileri18n.Load(filepath.Join(root.dir, "locales"), settings.I18n.Locales, settings.I18n.DefaultLocale)
	if err != nil {
		return fmt.Errorf("i18n catalogs: %w", err)
	}
	if err := transpileri18n.ValidateUses(catalogs, toI18nUses(gen.MessageUses)); err != nil {
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

