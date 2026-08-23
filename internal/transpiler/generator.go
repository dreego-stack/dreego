package transpiler

type Generator struct {
	defs      map[string]*ComponentDef
	src       string
	pkg       string
	module    string
	rootRel   string
	compPkgs  map[string]string
	compPaths map[string]string
	imports   map[string]map[string]string
}

func NewGenerator() *Generator {
	return &Generator{
		defs:      map[string]*ComponentDef{},
		compPkgs:  map[string]string{},
		compPaths: map[string]string{},
		imports:   map[string]map[string]string{},
	}
}

func (g *Generator) registerDef(name string, def *ComponentDef) {
	g.defs[name] = def
}

func (g *Generator) lookupDef(name string) *ComponentDef {
	return g.defs[name]
}

func (g *Generator) registerCompPkg(name, pkg, relDir string) {
	g.compPkgs[name] = pkg
	g.compPaths[pkg] = relDir
}

func (g *Generator) addImport(pkg, alias, path string) {
	if g.imports[pkg] == nil {
		g.imports[pkg] = map[string]string{}
	}
	g.imports[pkg][alias] = path
}

func (g *Generator) qualify(funcName string) string {
	pkg := g.compPkgs[funcName]
	if pkg == "" || pkg == g.pkg {
		return funcName
	}
	rel := g.compPaths[pkg]
	path := g.module + "/" + g.rootRel + "/" + rel
	g.addImport(g.pkg, pkg, path)
	return pkg + "." + funcName
}

type generator = Generator
