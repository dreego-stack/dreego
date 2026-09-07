package codegen

import "github.com/dreego-stack/dreego/internal/transpiler/ir"

type State struct {
	Defs      map[string]*ir.ComponentDef
	Src       string
	Pkg       string
	Module    string
	RootRel   string
	CompPkgs  map[string]string
	CompPaths map[string]string
	Imports   map[string]map[string]string
}

func NewState() *State {
	return &State{
		Defs:      map[string]*ir.ComponentDef{},
		CompPkgs:  map[string]string{},
		CompPaths: map[string]string{},
		Imports:   map[string]map[string]string{},
	}
}

func (g *State) RegisterDef(name string, def *ir.ComponentDef) {
	g.Defs[name] = def
}

func (g *State) LookupDef(name string) *ir.ComponentDef {
	return g.Defs[name]
}

func (g *State) RegisterCompPkg(name, pkg, relDir string) {
	g.CompPkgs[name] = pkg
	g.CompPaths[pkg] = relDir
}

func (g *State) AddImport(pkg, alias, path string) {
	if g.Imports[pkg] == nil {
		g.Imports[pkg] = map[string]string{}
	}
	g.Imports[pkg][alias] = path
}

func (g *State) Qualify(funcName string) string {
	pkg := g.CompPkgs[funcName]
	if pkg == "" || pkg == g.Pkg {
		return funcName
	}
	rel := g.CompPaths[pkg]
	path := g.Module + "/" + g.RootRel + "/" + rel
	g.AddImport(g.Pkg, pkg, path)
	return pkg + "." + funcName
}
