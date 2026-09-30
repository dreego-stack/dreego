package dreefile

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type routeDir = routePkg

// scanRoutes builds the app's route packages from its own routes/ tree plus the
// website root's shared routes/, where a local file overrides the global file
// with the same relative path.
func scanRoutes(gen *Generator, root, appName string, layouts, layoutIndex map[string]*layoutEntry, websiteRoot ...string) ([]*routePkg, map[string]bool, int, error) {
	shared := ""
	if len(websiteRoot) > 0 {
		shared = websiteRoot[0]
	}
	files, err := collectRouteFiles(root, shared)
	if err != nil {
		return nil, nil, 0, err
	}
	profiles, err := mergedRouteProfiles(root, shared)
	if err != nil {
		return nil, nil, 0, err
	}

	s := &routeScan{
		gen:            gen,
		root:           root,
		appName:        appName,
		layouts:        layouts,
		layoutIndex:    layoutIndex,
		rootRoutes:     filepath.Join(root, "routes"),
		pkgs:           map[string]*routePkg{},
		routePatterns:  map[string]bool{},
		routeSources:   map[string]string{},
		declSources:    map[string]map[string]string{},
		appliedProfile: map[string]bool{},
	}
	if err := s.scan(files, profiles); err != nil {
		return nil, nil, 0, err
	}
	if s.found == 0 {
		return nil, s.routePatterns, 0, nil
	}

	if _, ok := s.pkgs[""]; !ok {
		s.pkgs[""] = &routePkg{dir: s.rootRoutes, pkg: "routes", key: ""}
	}
	for rel, p := range s.pkgs {
		if rel == "" {
			continue
		}
		if parent := parentRoutePkg(s.pkgs, rel); parent != nil {
			parent.children = append(parent.children, p)
		}
	}

	list := make([]*routePkg, 0, len(s.pkgs))
	for _, p := range s.pkgs {
		if p.needsHead {
			p.src.WriteString(headMergeHelpers())
		}
		list = append(list, p)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].rel < list[j].rel })
	return list, s.routePatterns, s.found, nil
}

type routeScan struct {
	gen            *Generator
	root           string
	appName        string
	layouts        map[string]*layoutEntry
	layoutIndex    map[string]*layoutEntry
	rootRoutes     string
	pkgs           map[string]*routePkg
	routePatterns  map[string]bool
	routeSources   map[string]string
	declSources    map[string]map[string]string
	appliedProfile map[string]bool
	found          int
}

func (s *routeScan) scan(files []routeFileAt, profiles map[string]string) error {
	for index := 0; index < len(files); {
		dirRel := files[index].dirRel
		var group []routeFileAt
		for index < len(files) && files[index].dirRel == dirRel {
			group = append(group, files[index])
			index++
		}
		if err := s.scanDir(dirRel, group, profiles); err != nil {
			return err
		}
	}
	return nil
}

func (s *routeScan) getPkg(dirRel string) *routePkg {
	logical := filepath.Join(s.rootRoutes, filepath.FromSlash(dirRel))
	pkgDir := routePackageDir(s.root, logical)
	rel := routeDirRel(s.root, pkgDir)
	if p, ok := s.pkgs[rel]; ok {
		return p
	}
	pkg := "routes"
	if rel != "" {
		pkg = sanitizePkgName(filepath.Base(pkgDir))
	}
	p := &routePkg{dir: pkgDir, rel: rel, pkg: pkg, key: rel}
	s.pkgs[rel] = p
	return p
}

func (s *routeScan) scanDir(dirRel string, group []routeFileAt, profiles map[string]string) error {
	p := s.getPkg(dirRel)
	s.gen.Pkg = p.pkg
	s.gen.ImportKey = p.key
	s.gen.ImportKeySet = true
	decls := s.declSources[p.key]
	if decls == nil {
		decls = map[string]string{}
		s.declSources[p.key] = decls
	}
	folderProfile := resolveRouteProfile(profiles, dirRel)

	var src strings.Builder
	var regs, profilePatterns []string
	for _, f := range group {
		chunk, reg, pattern, err := s.scanFile(p, decls, f, folderProfile)
		if err != nil {
			return err
		}
		src.WriteString(chunk)
		regs = append(regs, reg)
		if pattern != "" {
			profilePatterns = append(profilePatterns, pattern)
		}
	}

	p.src.WriteString(src.String())
	for _, pattern := range profilePatterns {
		p.regs = append(p.regs, registrationStatement(fmt.Sprintf("app.ApplyProfile(%q, %q)", pattern, folderProfile)))
	}
	p.regs = append(p.regs, regs...)
	s.found += len(regs)
	return nil
}

// scanFile generates one route file. It returns the rendered source, its
// registration statement, and a profile pattern when the folder profile first
// applies to this route pattern.
func (s *routeScan) scanFile(p *routePkg, decls map[string]string, f routeFileAt, folderProfile string) (string, string, string, error) {
	rel := routeFileRelFromDir(f.dirRel, f.name)
	pattern := buildPattern(rel)
	if seg := doubleBracketSegment(rel); seg != "" {
		return "", "", "", fmt.Errorf("optional segment %q in %s is not supported; define each route explicitly", seg, f.path)
	}
	pageName := buildPageName(rel)
	data, err := os.ReadFile(f.path)
	if err != nil {
		return "", "", "", fmt.Errorf("error reading %s: %w", f.path, err)
	}
	baseName := strings.TrimSuffix(f.name, ".dreego")

	file, raw, perr := parseRouteFile(s.gen, f.path, data)
	if perr != nil {
		return "", "", "", perr
	}
	if err := registerGoImports(s.gen, p.key, f.path, file.GoImports); err != nil {
		return "", "", "", err
	}
	for _, name := range hoistedDeclarationNames(file) {
		if prev, dup := decls[name]; dup {
			return "", "", "", serverDeclarationConflict(name, prev, f.path)
		}
		decls[name] = f.path
	}
	layout, err := resolveRouteLayout(file, s.appName, rel, s.layouts, s.layoutIndex)
	if err != nil {
		return "", "", "", err
	}
	normalizeMethodSections(file)

	scopeHash := hashOf(data)
	s.gen.Src = raw

	if baseName == "404" || baseName == "500" {
		src, reg, err := s.emitError(p, file, f, pattern, baseName, scopeHash)
		if err != nil {
			return "", "", "", err
		}
		return src, reg, s.profilePattern(pattern, folderProfile), nil
	}
	for _, m := range fileRegisteredMethods(file) {
		key := m + " " + pattern
		if prev, dup := s.routeSources[key]; dup {
			return "", "", "", fmt.Errorf("duplicate route %s %s: %s and %s", m, pattern, prev, f.path)
		}
		s.routeSources[key] = f.path
		s.routePatterns[key] = true
	}
	src, reg, err := GenerateMethodHandler(s.gen, file, layout, p.pkg, pageName, pattern, scopeHash)
	if err != nil {
		return "", "", "", fmt.Errorf("error generating %s: %w", f.path, err)
	}
	if layout != nil {
		p.needsHead = true
	}
	return src, reg, s.profilePattern(pattern, folderProfile), nil
}

func (s *routeScan) emitError(p *routePkg, file *File, f routeFileAt, pattern, baseName, scopeHash string) (string, string, error) {
	errCode := 404
	if baseName == "500" {
		errCode = 500
	}
	catchPattern := errorCatchPattern(pattern)
	if errCode == 404 {
		catchKey := "GET " + catchPattern
		if prev, dup := s.routeSources[catchKey]; dup {
			return "", "", fmt.Errorf("duplicate catch-all %s: %s and %s", catchPattern, prev, f.path)
		}
		s.routeSources[catchKey] = f.path
	}
	src, reg, err := GenerateErrorHandler(s.gen, file, p.pkg, errCode, catchPattern, scopeHash)
	if err != nil {
		return "", "", fmt.Errorf("error generating error page %s: %w", f.path, err)
	}
	return src, reg, nil
}

func (s *routeScan) profilePattern(pattern, folderProfile string) string {
	if folderProfile == "" || s.appliedProfile[pattern] {
		return ""
	}
	s.appliedProfile[pattern] = true
	return pattern
}

func normalizeMethodSections(file *File) {
	if len(file.Server) == 0 {
		file.Server = []ServerSection{{Method: "GET"}}
	}
	for i := range file.Server {
		if !file.Server[i].MethodExplicit {
			file.Server[i].Method = "GET"
		}
	}
	for i := range file.Bodies {
		if !file.Bodies[i].MethodExplicit {
			file.Bodies[i].Method = "GET"
		}
	}
}

func parentRoutePkg(pkgs map[string]*routePkg, rel string) *routePkg {
	for {
		idx := strings.LastIndex(rel, "/")
		if idx < 0 {
			return pkgs[""]
		}
		rel = rel[:idx]
		if p, ok := pkgs[rel]; ok {
			return p
		}
	}
}

func routeFileRel(root, dir, name string) string {
	rel := routeDirRel(root, dir)
	base := strings.TrimSuffix(name, ".dreego")
	if base == "+page" || base == "index" || base == "404" || base == "500" {
		return rel
	}
	if rel == "" {
		return base
	}
	return filepath.ToSlash(filepath.Join(rel, base))
}
