package templates

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed all:_common all:web-minimal
var fsys embed.FS

const DefaultName = "web-minimal"

const (
	commonLayer = "_common"
	metaFile    = "template.json"
	nameToken   = "§$name$§"
)

type Meta struct {
	Name          string   `json:"name"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Type          string   `json:"type"`
	Adapter       string   `json:"adapter"`
	ExtraRequires []string `json:"extraRequires"`
}

func List() []Meta {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil
	}
	metas := make([]Meta, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		meta, err := metaOf(fsys, entry.Name())
		if err != nil {
			continue
		}
		metas = append(metas, meta)
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Name < metas[j].Name })
	return metas
}

func Exists(name string) bool {
	_, err := MetaOf(name)
	return err == nil
}

func MetaOf(name string) (Meta, error) {
	return metaOf(fsys, name)
}

func metaOf(fsys fs.FS, name string) (Meta, error) {
	var meta Meta
	if !validName(name) {
		return meta, fmt.Errorf("unknown template %q", name)
	}
	data, err := fs.ReadFile(fsys, path.Join(name, metaFile))
	if err != nil {
		return meta, fmt.Errorf("template %q: %w", name, err)
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return meta, fmt.Errorf("template %q: %w", name, err)
	}
	return meta, nil
}

func Install(target, moduleName, name string) error {
	return install(fsys, target, moduleName, name)
}

func install(fsys fs.FS, target, moduleName, name string) error {
	if _, err := metaOf(fsys, name); err != nil {
		return err
	}
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	for _, layer := range []string{commonLayer, name} {
		if err := installLayer(fsys, layer, target, moduleName); err != nil {
			return err
		}
	}
	return nil
}

func installLayer(fsys fs.FS, layer, target, moduleName string) error {
	return fs.WalkDir(fsys, layer, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(layer, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		if d.IsDir() {
			return os.MkdirAll(filepath.Join(target, rel), 0755)
		}
		if skipFile(rel) {
			return nil
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		content := strings.ReplaceAll(string(data), nameToken, moduleName)
		dest := filepath.Join(target, strings.TrimSuffix(rel, ".tmpl"))
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		return os.WriteFile(dest, []byte(content), 0644)
	})
}

func skipFile(rel string) bool {
	if rel == metaFile || filepath.Base(rel) == ".DS_Store" {
		return true
	}
	base := filepath.Base(rel)
	if base == "dree.go" || strings.HasSuffix(base, "_dreego.go") {
		return true
	}
	matched, err := filepath.Match("handle_*.go", base)
	return err == nil && matched
}

func validName(name string) bool {
	if name == "" || strings.HasPrefix(name, "_") || strings.HasPrefix(name, ".") {
		return false
	}
	return !strings.ContainsAny(name, `/\`)
}
