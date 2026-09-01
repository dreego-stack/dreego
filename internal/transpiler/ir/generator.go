package ir

type Generator struct {
	Defs      map[string]*ComponentDef
	Src       string
	Pkg       string
	Module    string
	RootRel   string
	CompPkgs  map[string]string
	CompPaths map[string]string
	Imports   map[string]map[string]string
}

func NewGenerator() *Generator {
	return &Generator{
		Defs:      map[string]*ComponentDef{},
		CompPkgs:  map[string]string{},
		CompPaths: map[string]string{},
		Imports:   map[string]map[string]string{},
	}
}

func (g *Generator) RegisterDef(name string, def *ComponentDef) {
	g.Defs[name] = def
}

func (g *Generator) LookupDef(name string) *ComponentDef {
	return g.Defs[name]
}

func (g *Generator) RegisterCompPkg(name, pkg, relDir string) {
	g.CompPkgs[name] = pkg
	g.CompPaths[pkg] = relDir
}

func (g *Generator) AddImport(pkg, alias, path string) {
	if g.Imports[pkg] == nil {
		g.Imports[pkg] = map[string]string{}
	}
	g.Imports[pkg][alias] = path
}

func (g *Generator) Qualify(funcName string) string {
	pkg := g.CompPkgs[funcName]
	if pkg == "" || pkg == g.Pkg {
		return funcName
	}
	rel := g.CompPaths[pkg]
	path := g.Module + "/" + g.RootRel + "/" + rel
	g.AddImport(g.Pkg, pkg, path)
	return pkg + "." + funcName
}

func (g *Generator) registerDef(name string, def *ComponentDef) {
	g.RegisterDef(name, def)
}

func (g *Generator) lookupDef(name string) *ComponentDef {
	return g.LookupDef(name)
}

func (g *Generator) registerCompPkg(name, pkg, relDir string) {
	g.RegisterCompPkg(name, pkg, relDir)
}

func (g *Generator) addImport(pkg, alias, path string) {
	g.AddImport(pkg, alias, path)
}

func (g *Generator) qualify(funcName string) string {
	return g.Qualify(funcName)
}

type generator = Generator
