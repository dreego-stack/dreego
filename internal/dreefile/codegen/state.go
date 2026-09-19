package codegen

import (
	"fmt"

	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

type State struct {
	Defs             map[string]*ir.ComponentDef
	Src              string
	Pkg              string
	Module           string
	RootRel          string
	CompPkgs         map[string]string
	CompPaths        map[string]string
	CompAliases      map[string]string
	Imports          map[string]map[string]string
	GoImportPaths    map[string][]string
	Lua              map[string]bool
	MessageUses      []MessageUse
	MessageArguments map[string]map[string]string
}

type MessageUse struct {
	Key       string
	Arguments []string
}

func NewState() *State {
	return &State{
		Defs:             map[string]*ir.ComponentDef{},
		CompPkgs:         map[string]string{},
		CompPaths:        map[string]string{},
		CompAliases:      map[string]string{},
		Imports:          map[string]map[string]string{},
		GoImportPaths:    map[string][]string{},
		Lua:              map[string]bool{},
		MessageUses:      nil,
		MessageArguments: map[string]map[string]string{},
	}
}

func (g *State) MessageArgumentKind(key, name string) string {
	return g.MessageArguments[key][name]
}

func (g *State) RegisterMessageUse(key string, arguments []string) {
	g.MessageUses = append(g.MessageUses, MessageUse{Key: key, Arguments: append([]string(nil), arguments...)})
}

func (g *State) AddLuaFeatures(features []string) {
	for _, feature := range features {
		g.Lua[feature] = true
	}
}

func (g *State) RegisterDef(name string, def *ir.ComponentDef) {
	g.Defs[name] = def
}

func (g *State) LookupDef(name string) *ir.ComponentDef {
	if alias := g.CompAliases[name]; alias != "" {
		name = alias
	}
	return g.Defs[name]
}

func (g *State) RegisterCompPkg(name, pkg, relDir string) {
	g.CompPkgs[name] = pkg
	g.CompPaths[pkg] = relDir
}

func (g *State) RegisterCompAlias(alias, target string) error {
	if alias == "" || alias == target {
		return nil
	}
	if existing, ok := g.CompAliases[alias]; ok && existing != target {
		return fmt.Errorf("duplicate component alias %s: %s and %s", alias, existing, target)
	}
	g.CompAliases[alias] = target
	return nil
}

func (g *State) AddImport(pkg, alias, path string) {
	if g.Imports[pkg] == nil {
		g.Imports[pkg] = map[string]string{}
	}
	g.Imports[pkg][alias] = path
}

func (g *State) AddGoImportPath(pkg, path string) {
	if g.GoImportPaths == nil {
		g.GoImportPaths = map[string][]string{}
	}
	for _, existing := range g.GoImportPaths[pkg] {
		if existing == path {
			return
		}
	}
	g.GoImportPaths[pkg] = append(g.GoImportPaths[pkg], path)
}

func (g *State) Qualify(funcName string) string {
	if alias := g.CompAliases[funcName]; alias != "" {
		funcName = alias
	}
	pkg := g.CompPkgs[funcName]
	if pkg == "" || pkg == g.Pkg {
		return funcName
	}
	rel := g.CompPaths[pkg]
	path := g.Module + "/" + g.RootRel + "/" + rel
	g.AddImport(g.Pkg, pkg, path)
	return pkg + "." + funcName
}
