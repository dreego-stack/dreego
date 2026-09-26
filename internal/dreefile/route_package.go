package dreefile

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type routePkg struct {
	dir       string
	rel       string
	pkg       string
	key       string
	src       strings.Builder
	regs      []string
	needsHead bool
	children  []*routePkg
	alias     string
}

func (p *routePkg) isRoot(root string) bool {
	return filepath.Clean(p.dir) == filepath.Clean(filepath.Join(root, "routes"))
}

// routePackageDir maps a route directory to the nearest ancestor that is a
// valid Go package directory. Segment names such as "[id]" or "(group)" cannot
// be import paths, so they merge into their closest valid parent.
func routePackageDir(root, dir string) string {
	rootRoutes := filepath.Join(root, "routes")
	for filepath.Clean(dir) != filepath.Clean(rootRoutes) {
		if validRoutePackagePath(relToRoot(rootRoutes, dir)) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return rootRoutes
}

func validRoutePackagePath(rel string) bool {
	if rel == "" || rel == "." {
		return true
	}
	for seg := range strings.SplitSeq(rel, "/") {
		if !validRoutePackageSegment(seg) {
			return false
		}
	}
	return true
}

func validRoutePackageSegment(seg string) bool {
	if seg == "" || seg == "." || seg == ".." {
		return false
	}
	hasAlnum := false
	for _, r := range seg {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z':
			hasAlnum = true
		case r >= '0' && r <= '9':
			hasAlnum = true
		case r == '-', r == '_', r == '.':
		default:
			return false
		}
	}
	return hasAlnum && sanitizePkgName(seg) != ""
}

func routeImportPath(gen *Generator, rel string) string {
	path := gen.Module + "/" + gen.RootRel + "/routes"
	if rel != "" {
		path += "/" + rel
	}
	return path
}

type routeImport struct {
	alias string
	path  string
}

// buildRoutePackageFile renders one route folder's dree.go. For the top-level
// routes folder it also imports every sub-package and calls its Register, so
// the website entry point only has to call routes.Register(app).
func buildRoutePackageFile(gen *Generator, p *routePkg) string {
	src := p.src.String()
	stdImports := stdImportsFor(gen, p.key, src)
	lines := []routeImport{}
	used := map[string]bool{"dreego": true, "ssr": true, p.pkg: true}
	for _, line := range strings.Split(stdImports, "\n") {
		used[baseImportName(strings.Trim(strings.TrimSpace(line), `"`))] = true
	}
	for alias, path := range gen.Imports[p.key] {
		lines = append(lines, routeImport{alias: alias, path: path})
		used[alias] = true
	}
	children := append([]*routePkg(nil), p.children...)
	sort.Slice(children, func(i, j int) bool { return children[i].rel < children[j].rel })
	for _, child := range children {
		alias := child.pkg
		if used[alias] {
			alias = uniqueImportAlias(alias, used)
		}
		used[alias] = true
		child.alias = alias
		lines = append(lines, routeImport{alias: alias, path: routeImportPath(gen, child.rel)})
	}
	sort.Slice(lines, func(i, j int) bool { return lines[i].path < lines[j].path })

	importLines := make([]string, 0, len(lines))
	for _, l := range lines {
		if l.alias == p.pkg {
			importLines = append(importLines, strconv.Quote(l.path))
			continue
		}
		importLines = append(importLines, fmt.Sprintf("%s %s", l.alias, strconv.Quote(l.path)))
	}
	importLine := strings.Join(importLines, "\n\t")

	coreImport := "dreego \"github.com/dreego-stack/dreego/core\""
	if strings.Contains(src, "ssr.") {
		coreImport += "\n\tssr \"github.com/dreego-stack/dreego/adapter/ssr\""
	}

	var b strings.Builder
	fmt.Fprintf(&b, "package %s\n\nimport (\n", p.pkg)
	if stdImports != "" {
		b.WriteString("\t" + stdImports + "\n")
	}
	if importLine != "" {
		b.WriteString("\t" + importLine + "\n\n")
	}
	b.WriteString("\t" + coreImport + "\n)\n\n")
	b.WriteString(src)

	b.WriteString("func Register(app *dreego.App) error {\n")
	b.WriteString(strings.Join(p.regs, ""))
	for _, child := range children {
		b.WriteString(registrationStatement(fmt.Sprintf("%s.Register(app)", child.alias)))
	}
	b.WriteString("\treturn nil\n}\n")
	return b.String()
}

func uniqueImportAlias(base string, used map[string]bool) string {
	for i := 2; ; i++ {
		candidate := base + strconv.Itoa(i)
		if !used[candidate] {
			return candidate
		}
	}
}

func baseImportName(path string) string {
	if path == "" {
		return ""
	}
	slash := strings.LastIndex(path, "/")
	name := path[slash+1:]
	if dot := strings.Index(name, "."); dot > 0 {
		name = name[:dot]
	}
	return name
}
