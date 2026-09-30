package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dreego-stack/dreego/internal/gomod"
)

func loadModuleComponents(gen *Generator, routeRoot, baseRoot string) ([]componentSource, error) {
	paths, err := importedComponentPaths(routeRoot)
	if err != nil {
		return nil, err
	}
	var components []componentSource
	for _, importPath := range paths {
		loaded, err := loadModuleComponent(gen, baseRoot, importPath)
		if err != nil {
			return nil, err
		}
		components = append(components, loaded...)
	}
	return components, nil
}

func loadModuleComponent(gen *Generator, baseRoot, importPath string) ([]componentSource, error) {
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
	var components []componentSource
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
