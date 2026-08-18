package core

import (
	"crypto/sha256"
	"encoding/hex"
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
}

func scanComponents(gen *generator) (string, []string, error) {
	components, genDir, err := loadComponents()
	if err != nil {
		return "", nil, err
	}
	pathsByName := map[string]string{}
	for _, component := range components {
		if previous, exists := pathsByName[component.def.Name]; exists {
			return "", nil, fmt.Errorf("duplicate component %s: %s and %s", component.def.Name, previous, component.path)
		}
		pathsByName[component.def.Name] = component.path
		gen.registerDef(component.def.Name, component.def)
	}

	sources := make([]string, 0, len(components))
	for _, component := range components {
		gen.src = component.raw
		src, err := GenerateComponent(gen, component.file, component.scopeHash)
		if err != nil {
			return "", nil, fmt.Errorf("error generating component %s: %w", component.path, err)
		}
		sources = append(sources, src)
	}
	return genDir, sources, nil
}

func loadComponents() ([]componentSource, string, error) {
	var components []componentSource
	var genDir string
	err := filepath.WalkDir(".", func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("error walking %s: %w", path, walkErr)
		}
		if d.IsDir() || !strings.HasSuffix(path, ".dreego") {
			return nil
		}
		if !isInDreegoRoot(path) {
			return nil
		}
		dir := filepath.Dir(path)
		if !isDreegoComponentsDir(dir) {
			return nil
		}
		component, err := loadComponent(path)
		if err != nil {
			return err
		}
		if component.def == nil {
			return nil
		}
		if genDir == "" {
			genDir = detectGenDir(path)
		}
		components = append(components, component)
		return nil
	})
	return components, genDir, err
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
	hash := sha256.Sum256(data)
	return componentSource{
		path:      path,
		raw:       raw,
		file:      file,
		def:       def,
		scopeHash: hex.EncodeToString(hash[:])[:12],
	}, nil
}

func prepareComponentFile(file *File, def *ComponentDef, path, raw string, bodyOffset int) {
	file.Component = def
	file.SourceContent = raw
	if file.Template != nil {
		setNodeSource(file.Template.Nodes, path, bodyOffset)
		def.Slots = mergeUnique(def.Slots, collectSlotNames(file.Template.Nodes))
		def.HasDefaultSlot = hasDefaultSlot(file.Template.Nodes)
		def.HasNamedSlot = hasNamedSlot(file.Template.Nodes) || len(def.Slots) > 0
	}
	if len(file.Go) == 0 {
		file.Go = []GoSection{{Method: ""}}
	}
}
