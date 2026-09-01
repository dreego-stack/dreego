package html

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/html/css"
	"github.com/dreego-stack/dreego/internal/transpiler/html/output"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func GenerateComponent(gen *ir.Generator, file *ir.File, scopeHash string) (string, error) {
	comp := file.Component
	if comp == nil {
		return "", fmt.Errorf("no component definition")
	}

	var buf strings.Builder

	declParams, implParams, callArgs, variadicName := ComponentParams(comp)

	if variadicName != "" {
		buf.WriteString(fmt.Sprintf("func %s(%s) dreego.Component {\n", comp.Name, declParams))
		buf.WriteString("\t" + variadicName + "0 := \"\"\n")
		buf.WriteString("\tif len(" + variadicName + ") > 0 {\n\t\t" + variadicName + "0 = " + variadicName + "[0]\n\t}\n")
		buf.WriteString("\treturn component" + comp.Name + "(" + callArgs + ")\n")
		buf.WriteString("}\n\n")
		buf.WriteString(fmt.Sprintf("func component%s(%s) dreego.Component {\n", comp.Name, implParams))
	} else {
		buf.WriteString(fmt.Sprintf("func %s(%s) dreego.Component {\n", comp.Name, declParams))
	}
	buf.WriteString("\treturn dreego.ComponentFunc(func(ctx *dreego.SSRContext) (string, error) {\n")
	WritePropDefaultFallbacks(&buf, comp)
	buf.WriteString("\t\tvar b strings.Builder\n\n")

	for _, g := range file.Server {
		if g.Code != "" {
			for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
				buf.WriteString("\t\t" + strings.TrimSpace(line) + "\n")
			}
			buf.WriteString("\n")
		}
	}

	if file.Body != nil {
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		g := &output.CompGen{Gen: gen, Builder: "b"}
		for _, n := range file.Body.Nodes {
			code, err := g.Node(n)
			if err != nil {
				return "", err
			}
			buf.WriteString("\t\t" + code + "\n")
		}
		buf.WriteString("\t\tb.WriteString(\"</div>\")\n")
	}

	if file.Style != nil {
		scoped := css.ScopeCSS(file.Style.Code, scopeHash)
		buf.WriteString("\t\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(%s)\n", ir.GoLiteral(scoped)))
		buf.WriteString("\t\tb.WriteString(\"</style>\")\n")
	}

	buf.WriteString("\n\t\treturn b.String(), nil\n")
	buf.WriteString("\t})\n")
	buf.WriteString("}\n\n")

	return buf.String(), nil
}
