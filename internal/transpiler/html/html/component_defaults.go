package html

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func ComponentParams(comp *ir.ComponentDef) (decl, impl, call string, variadic string) {
	for i, p := range comp.Props {
		if i > 0 {
			decl += ", "
			impl += ", "
			call += ", "
		}
		decl += p.Name + " " + p.Type
		impl += p.Name + " " + p.Type
		call += p.Name
	}
	for i, p := range comp.Props {
		if p.Default != "" && p.Type == "string" && i == len(comp.Props)-1 {
			variadic = p.Name
			decl = strings.Replace(decl, variadic+" string", variadic+" ...string", 1)
			call = strings.Replace(call, variadic, variadic+"0", 1)
		}
	}
	return
}

func WritePropDefaultFallbacks(buf *strings.Builder, comp *ir.ComponentDef) {
	for _, p := range comp.Props {
		if p.Default == "" || p.Type != "string" {
			continue
		}
		buf.WriteString(fmt.Sprintf("\t\tif %s == \"\" {\n\t\t\t%s = %s\n\t\t}\n", p.Name, p.Name, p.Default))
	}
}
