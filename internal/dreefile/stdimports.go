package dreefile

import (
	"fmt"
	"go/build"
	"path"
	"sort"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

func isStdlibImport(importPath string) bool {
	pkg, err := build.Default.Import(importPath, ".", build.FindOnly)
	if err == nil && pkg.Goroot {
		return true
	}
	if build.Default.GOROOT == "" {
		return isDotlessImport(importPath)
	}
	return false
}

func isDotlessImport(importPath string) bool {
	first, _, _ := strings.Cut(importPath, "/")
	return first != "" && !strings.Contains(first, ".")
}

func importBaseName(importPath string) string {
	base := path.Base(importPath)
	if isModuleVersionSegment(base) {
		if parent := path.Base(path.Dir(importPath)); parent != "." && parent != "/" {
			return parent
		}
	}
	return base
}

// isModuleVersionSegment reports whether seg is a semantic import version
// suffix such as "v2". Such a segment is not a usable Go identifier, so the
// alias derives from the parent path instead.
func isModuleVersionSegment(seg string) bool {
	if len(seg) < 2 || seg[0] != 'v' {
		return false
	}
	for i := 1; i < len(seg); i++ {
		if seg[i] < '0' || seg[i] > '9' {
			return false
		}
	}
	return true
}

func importSpecName(imp ir.GoImport) string {
	if imp.Alias != "" {
		return imp.Alias
	}
	return importBaseName(imp.Path)
}

func resolvesGoImport(gen *Generator, importPath string) bool {
	if gen.Module != "" && (importPath == gen.Module || strings.HasPrefix(importPath, gen.Module+"/")) {
		return true
	}
	for modulePath := range gen.Requires {
		if importPath == modulePath || strings.HasPrefix(importPath, modulePath+"/") {
			return true
		}
	}
	return isStdlibImport(importPath)
}

func registerGoImports(gen *Generator, pkg, source string, imports []ir.GoImport) error {
	for _, imp := range imports {
		if imp.Path == "" {
			return fmt.Errorf("%s: GOIMPORT contains an empty package path", source)
		}
		if !resolvesGoImport(gen, imp.Path) {
			return fmt.Errorf("%s: GOIMPORT %q is not in go.mod; run 'go get %s' to add it", source, imp.Path, imp.Path)
		}
		gen.AddGoImport(pkg, imp)
	}
	return nil
}

type goImportSpec struct {
	Name string
	Path string
}

func stdImportsFor(gen *Generator, pkg, src string) (string, error) {
	specs, err := collectGoImportSpecs(gen, pkg, src)
	if err != nil {
		return "", err
	}
	lines := make([]string, 0, len(specs))
	for _, spec := range specs {
		if spec.Name == importBaseName(spec.Path) {
			lines = append(lines, fmt.Sprintf("%q", spec.Path))
			continue
		}
		lines = append(lines, fmt.Sprintf("%s %q", spec.Name, spec.Path))
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n\t"), nil
}

func collectGoImportSpecs(gen *Generator, pkg, src string) ([]goImportSpec, error) {
	byPath := map[string]bool{}
	var specs []goImportSpec
	add := func(name, importPath string) {
		if byPath[importPath] {
			return
		}
		byPath[importPath] = true
		specs = append(specs, goImportSpec{Name: name, Path: importPath})
	}
	for _, imp := range gen.GoImports[pkg] {
		add(importSpecName(imp), imp.Path)
	}
	for _, detected := range []struct{ name, path string }{
		{"strings", "strings"},
		{"http", "net/http"},
		{"fmt", "fmt"},
	} {
		if containsQualifiedIdent(src, detected.name) {
			add(detected.name, detected.path)
		}
	}
	seenName := map[string]string{}
	var collisions []string
	for _, spec := range specs {
		if previous, ok := seenName[spec.Name]; ok && previous != spec.Path {
			collisions = append(collisions, fmt.Sprintf("GOIMPORT %q and %q share the name %q in package %s; add an alias", previous, spec.Path, spec.Name, pkg))
			continue
		}
		seenName[spec.Name] = spec.Path
	}
	if len(collisions) > 0 {
		sort.Strings(collisions)
		return nil, fmt.Errorf("%s", strings.Join(collisions, "; "))
	}
	return specs, nil
}

func containsQualifiedIdent(src, name string) bool {
	needle := name + "."
	for idx := 0; ; {
		at := strings.Index(src[idx:], needle)
		if at < 0 {
			return false
		}
		at += idx
		if at == 0 || !isIdentByte(src[at-1]) {
			return true
		}
		idx = at + 1
	}
}

func isIdentByte(b byte) bool {
	return b == '_' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}
