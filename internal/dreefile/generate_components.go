package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
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

// componentPackage is one generated Go package holding component sources.
type componentPackage struct {
	dir  string
	pkg  string
	srcs []string
}

// componentPkgName returns a unique Go package name for a component directory
// relative to the website root. The shared root components/ keeps the plain
// name "components"; an app-local or nested directory is prefixed with its path
// so two component packages never collide in one generated file.
func componentPkgName(baseRoot, pkgDir string) string {
	rel := relToRoot(baseRoot, pkgDir)
	if rel == "components" {
		return "components"
	}
	return sanitizePkgName(strings.ReplaceAll(rel, "/", "_"))
}

func scanComponents(gen *Generator, root string) ([]componentPackage, map[string]string, error) {
	return scanComponentsWithBase(gen, root, root, true)
}

// scanComponentsWithBase scans component sources under `root`. When
// includeModuleComponents is true it also loads components imported from other Go
// modules and places them in the components package of `baseRoot`, so the whole
// website shares one generated location for a module component. App passes set
// it false so the shared root package is not overwritten.
func scanComponentsWithBase(gen *Generator, root, baseRoot string, includeModuleComponents bool) ([]componentPackage, map[string]string, error) {
	components, err := loadComponents(gen, root, baseRoot, includeModuleComponents)
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

	byDir := map[string]*componentPackage{}
	var order []string
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
		if byDir[component.pkgDir] == nil {
			byDir[component.pkgDir] = &componentPackage{dir: component.pkgDir, pkg: component.pkg}
			order = append(order, component.pkgDir)
		}
		byDir[component.pkgDir].srcs = append(byDir[component.pkgDir].srcs, src)
	}
	sort.Strings(order)
	pkgs := make([]componentPackage, 0, len(order))
	for _, dir := range order {
		pkgs = append(pkgs, *byDir[dir])
	}
	return pkgs, pathsByName, nil
}

// generateComponentSources emits the Go source for pre-loaded component sources
// and groups them into their output packages.
func generateComponentSources(gen *Generator, components []componentSource) ([]componentPackage, error) {
	byDir := map[string]*componentPackage{}
	var order []string
	for _, component := range components {
		gen.Src = component.raw
		gen.Pkg = component.pkg
		if err := registerGoImports(gen, component.pkg, component.path, component.file.GoImports); err != nil {
			return nil, err
		}
		src, err := GenerateComponent(gen, component.file, component.scopeHash)
		if err != nil {
			return nil, fmt.Errorf("error generating component %s: %w", component.path, err)
		}
		if byDir[component.pkgDir] == nil {
			byDir[component.pkgDir] = &componentPackage{dir: component.pkgDir, pkg: component.pkg}
			order = append(order, component.pkgDir)
		}
		byDir[component.pkgDir].srcs = append(byDir[component.pkgDir].srcs, src)
	}
	sort.Strings(order)
	pkgs := make([]componentPackage, 0, len(order))
	for _, dir := range order {
		pkgs = append(pkgs, *byDir[dir])
	}
	return pkgs, nil
}
func loadComponents(gen *Generator, root, baseRoot string, includeModuleComponents bool) ([]componentSource, error) {
	var components []componentSource
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".dreego") {
			return nil
		}
		dir := filepath.Dir(path)
		if !isComponentsDir(root, dir) {
			return nil
		}
		component, err := loadComponent(path)
		if err != nil {
			return err
		}
		if component.def == nil {
			return nil
		}
		rel := relToRoot(root, dir)
		pkgDir := dir
		if rel != "components" {
			segments := strings.Split(strings.TrimPrefix(rel, "components/"), "/")
			valid := []string{}
			for _, seg := range segments {
				if seg == "" {
					continue
				}
				if sanitizePkgName(seg) != seg {
					break
				}
				valid = append(valid, seg)
			}
			if len(valid) > 0 {
				pkgDir = filepath.Join(append([]string{root, "components"}, valid...)...)
			}
		}
		component.pkgDir = pkgDir
		component.pkg = componentPkgName(baseRoot, pkgDir)
		gen.RegisterCompPkg(component.def.Name, component.pkg, gen.Module+"/"+relToRoot(".", pkgDir))
		components = append(components, component)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !includeModuleComponents {
		return components, nil
	}
	moduleComponents, err := loadModuleComponents(gen, root, baseRoot)
	if err != nil {
		return nil, err
	}
	return append(components, moduleComponents...), nil
}

// loadModuleComponentsAcrossApps loads components imported from other Go modules
// for every app route tree under the website root, placing them in the shared
// root components package so the root pass owns them exactly once.
func loadModuleComponentsAcrossApps(gen *Generator, root string, apps []appEntry) ([]componentSource, error) {
	seen := map[string]bool{}
	var components []componentSource
	for _, app := range apps {
		paths, err := importedComponentPaths(app.dir)
		if err != nil {
			return nil, err
		}
		for _, importPath := range paths {
			if seen[importPath] {
				continue
			}
			seen[importPath] = true
			loaded, err := loadModuleComponent(gen, root, importPath)
			if err != nil {
				return nil, err
			}
			components = append(components, loaded...)
		}
	}
	return components, nil
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
