package transpiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type componentSource struct {
	path      string
	raw       string
	file      *File
	def       *ComponentDef
	scopeHash string
	pkgDir    string
}

func scanComponents(gen *Generator, root string) (map[string][]string, map[string]string, error) {
	components, err := loadComponents(gen, root)
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
	gen.Pkg = "components"
	for _, component := range components {
		gen.Src = component.raw
		src, err := GenerateComponent(gen, component.file, component.scopeHash)
		if err != nil {
			return nil, nil, fmt.Errorf("error generating component %s: %w", component.path, err)
		}
		sourcesByPkg[component.pkgDir] = append(sourcesByPkg[component.pkgDir], src)
	}
	return sourcesByPkg, pathsByName, nil
}

func loadComponents(gen *Generator, root string) ([]componentSource, error) {
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
		pkg := "components"
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
				pkg = sanitizePkgName(valid[len(valid)-1])
			}
		}
		component.pkgDir = pkgDir
		gen.RegisterCompPkg(component.def.Name, pkg, relToRoot(root, pkgDir))
		components = append(components, component)
		return nil
	})
	return components, err
}

func loadComponent(path string) (componentSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return componentSource{}, fmt.Errorf("error reading component %s: %w", path, err)
	}
	raw := string(data)
	def, _, body := ParseHeader(raw)
	if def == nil || def.Name == "" {
		return componentSource{}, nil
	}
	tokens, err := Lex(body)
	if err != nil {
		return componentSource{}, fmt.Errorf("error lexing component %s: %w", path, err)
	}
	file, err := NewParser(tokens).Parse()
	if err != nil {
		return componentSource{}, fmt.Errorf("error parsing component %s: %w", path, err)
	}
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
	if file.Body != nil {
		setNodeSource(file.Body.Nodes, path, bodyOffset)
		def.Slots = mergeUnique(def.Slots, collectSlotNames(file.Body.Nodes))
		def.HasDefaultSlot = hasDefaultSlot(file.Body.Nodes)
		def.HasNamedSlot = hasNamedSlot(file.Body.Nodes) || len(def.Slots) > 0
	}
	if len(file.Server) == 0 {
		file.Server = []ServerSection{{Method: ""}}
	}
}
