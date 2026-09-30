package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
	"github.com/dreego-stack/dreego/internal/gomod"
)

type componentSource struct {
	path      string
	raw       string
	file      *File
	def       *ComponentDef
	scopeHash string
	pkgDir    string
	pkg       string
}

func scanComponents(gen *Generator, root string) (map[string][]string, map[string]string, error) {
	return scanComponentsAt(gen, root, root, true)
}

// scanComponentsAt scans components under scopeRoot and writes generated files
// under baseRoot's component packages. When includeModule is true it also loads
// components imported from other Go modules (placed at baseRoot); app passes set
// it false so a root-owned package is never overwritten.
func scanComponentsAt(gen *Generator, scopeRoot, baseRoot string, includeModule bool) (map[string][]string, map[string]string, error) {
	components, err := loadComponents(gen, scopeRoot, baseRoot, includeModule)
	if err != nil {
		return nil, nil, err
	}
	pathsByName := map[string]string{}
	for _, component := range components {
		if previous, exists := pathsByName[component.def.Name]; exists {
			return nil, nil, fmt.Errorf("duplicate component %s: %s and %s", component.def.Name, previous, component.path)
		}
		pathsByName[component.def.Name] = component.path
		gen.RegisterDef(component.def.Name, component.def)
	}

	sourcesByPkg := map[string][]string{}
	for _, component := range components {
		gen.Src = component.raw
		gen.Pkg = component.pkg
		if err := registerGoImports(gen, component.pkg, component.path, component.file.GoImports); err != nil {
			return nil, nil, err
		}
		src, err := GenerateComponent(gen, component.file, component.scopeHash)
		if err != nil {
			return nil, nil, fmt.Errorf("error generating component %s: %w", component.path, err)
		}
		sourcesByPkg[component.pkgDir] = append(sourcesByPkg[component.pkgDir], src)
	}
	return sourcesByPkg, pathsByName, nil
}

// loadModuleComponentsAcrossApps loads components imported from other Go modules
// for every app route tree and generates their Go source, placing them in the
// shared root components package so the root pass owns them exactly once.
func loadModuleComponentsAcrossApps(gen *Generator, root string, apps []appEntry) (map[string][]string, error) {
	seen := map[string]bool{}
	var all []componentSource
	for _, app := range apps {
		paths, err := importedComponentPaths(app.dir)
		if err != nil {
			return nil, err
		}
		var fresh []string
		for _, importPath := range paths {
			if seen[importPath] {
				continue
			}
			seen[importPath] = true
			fresh = append(fresh, importPath)
		}
		if len(fresh) == 0 {
			continue
		}
		loaded, err := loadModuleComponents(gen, root, root, fresh)
		if err != nil {
			return nil, err
		}
		all = append(all, loaded...)
	}
	byDir := map[string][]string{}
	for _, component := range all {
		gen.Src = component.raw
		gen.Pkg = component.pkg
		if err := registerGoImports(gen, component.pkg, component.path, component.file.GoImports); err != nil {
			return nil, err
		}
		src, err := GenerateComponent(gen, component.file, component.scopeHash)
		if err != nil {
			return nil, fmt.Errorf("error generating component %s: %w", component.path, err)
		}
		byDir[component.pkgDir] = append(byDir[component.pkgDir], src)
	}
	return byDir, nil
}

func loadComponents(gen *Generator, scopeRoot, baseRoot string, includeModule bool) ([]componentSource, error) {	var components []componentSource
	err := filepath.WalkDir(scopeRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".dreego") {
			return nil
		}
		dir := filepath.Dir(path)
		if !isComponentsDir(scopeRoot, dir) {
			return nil
		}
		component, err := loadComponent(path)
		if err != nil {
			return err
		}
		if component.def == nil {
			return nil
		}
		pkgDir := dir
		component.pkgDir = pkgDir
		component.pkg = componentPkgName(baseRoot, pkgDir)
		gen.RegisterCompPkg(component.def.Name, component.pkg, gen.Module+"/"+relToRoot(".", pkgDir))
		components = append(components, component)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !includeModule {
		return components, nil
	}
	moduleComponents, err := loadModuleComponents(gen, scopeRoot, baseRoot, nil)
	if err != nil {
		return nil, err
	}
	return append(components, moduleComponents...), nil
}

// componentPkgName returns a unique Go package name for a component directory.
// The shared root components/ keeps "components"; an app-local or nested
// directory is prefixed with its path so package names never collide.
func componentPkgName(baseRoot, pkgDir string) string {
	rel := relToRoot(baseRoot, pkgDir)
	if rel == "components" {
		return "components"
	}
	return sanitizePkgName(strings.ReplaceAll(rel, "/", "_"))
}

func loadModuleComponents(gen *Generator, routeRoot, baseRoot string, paths []string) ([]componentSource, error) {
	if paths == nil {
		var err error
		paths, err = importedComponentPaths(routeRoot)
		if err != nil {
			return nil, err
		}
	}
	var components []componentSource
	for _, importPath := range paths {
		sourceDir, packageName, err := resolveComponentImport(importPath)
		if err != nil {
			return nil, err
		}
		entries, err := os.ReadDir(sourceDir)
		if err != nil {
			return nil, fmt.Errorf("read imported components %s: %w", importPath, err)
		}
		pkgDir := filepath.Join(baseRoot, "components")
		if packageName != "components" {
			pkgDir = filepath.Join(pkgDir, packageName)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".dreego") {
				continue
			}
			component, err := loadComponent(filepath.Join(sourceDir, entry.Name()))
			if err != nil {
				return nil, err
			}
			if component.def == nil {
				continue
			}
			component.pkgDir = pkgDir
			component.pkg = componentPkgName(baseRoot, pkgDir)
			gen.RegisterDef(component.def.Name, component.def)
			gen.RegisterCompPkg(component.def.Name, component.pkg, gen.Module+"/"+relToRoot(".", pkgDir))
			components = append(components, component)
		}
	}
	return components, nil
}

func importedComponentPaths(root string) ([]string, error) {
	seen := map[string]bool{}
	var paths []string
	routesDir := filepath.Join(root, "routes")
	if _, err := os.Stat(routesDir); os.IsNotExist(err) {
		return paths, nil
	}
	err := filepath.WalkDir(routesDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() || !strings.HasSuffix(path, ".dreego") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, imports, _ := ParseHeader(string(data))
		for _, imp := range imports {
			if !strings.HasSuffix(imp.Path, ".dreego") && strings.Contains(imp.Path, ".") && !seen[imp.Path] {
				seen[imp.Path] = true
				paths = append(paths, imp.Path)
			}
		}
		return nil
	})
	return paths, err
}

func resolveComponentImport(importPath string) (string, string, error) {
	mod, err := gomod.Read("go.mod")
	if err != nil {
		return "", "", fmt.Errorf("read go.mod for component import %q: %w", importPath, err)
	}
	best := ""
	for modulePath := range mod.Requires {
		if strings.HasPrefix(importPath, modulePath+"/") && len(modulePath) > len(best) {
			best = modulePath
		}
	}
	if best == "" {
		return "", "", fmt.Errorf("component import %q is not provided by a required Go module", importPath)
	}
	moduleDir, err := moduleDirectory(best)
	if err != nil {
		return "", "", err
	}
	relative := strings.TrimPrefix(importPath, best+"/")
	packageName := sanitizePkgName(filepath.Base(relative))
	return filepath.Join(moduleDir, filepath.FromSlash(relative)), packageName, nil
}

func loadComponent(path string) (componentSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return componentSource{}, fmt.Errorf("error reading component %s: %w", path, err)
	}
	raw := string(data)
	header, body, err := ParseFileHeaderStrict(raw)
	if err != nil {
		return componentSource{}, fmt.Errorf("%s:%w", path, err)
	}
	if header.Kind != FileKindComponent {
		return componentSource{}, nil
	}
	name := header.Name
	if header.Name == "" {
		name, err = componentNameFromPath(path)
		if err != nil {
			return componentSource{}, fmt.Errorf("%s: %w", path, err)
		}
	}
	if name == "" {
		return componentSource{}, nil
	}
	def := &ComponentDef{Name: name, Props: header.Props, Slots: header.Slots}
	tokens, err := Lex(body)
	if err != nil {
		return componentSource{}, fmt.Errorf("error lexing component %s: %w", path, err)
	}
	file, err := NewParserConcatServer(tokens).Parse()
	if err != nil {
		return componentSource{}, fmt.Errorf("error parsing component %s: %w", path, err)
	}
	file.Imports = header.Imports
	file.Kind = header.Kind
	file.Layout = header.Layout
	file.GoImports = header.GoImports
	prepareComponentFile(file, def, path, raw, len(raw)-len(body))
	return componentSource{
		path:      path,
		raw:       raw,
		file:      file,
		def:       def,
		scopeHash: hashOf(data),
	}, nil
}

func prepareComponentFile(file *File, def *ComponentDef, path, raw string, bodyOffset int) {
	file.Component = def
	file.SourceContent = raw
	file.SourcePath = path
	if file.Client != nil {
		file.Client.Pos += bodyOffset
	}
	if file.Body != nil {
		gogen.SetNodeSource(file.Body.Nodes, path, bodyOffset)
		def.Slots = mergeUnique(def.Slots, collectSlotNames(file.Body.Nodes))
		def.HasDefaultSlot = hasDefaultSlot(file.Body.Nodes)
		def.HasNamedSlot = hasNamedSlot(file.Body.Nodes) || len(def.Slots) > 0
	}
	if len(file.Server) == 0 {
		file.Server = []ServerSection{{Method: ""}}
	}
}
