package transpiler

import (
	"fmt"
	"strings"
)

// componentParams computes the generated function signatures for a component.
// A trailing string prop with a default becomes variadic (…string) so the
// call site may omit it; the wrapper forwards the first element to the impl.
func componentParams(comp *ComponentDef) (decl, impl, call string, variadic string) {
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

// writePropDefaultFallbacks emits `if x == "" { x = default }` statements so a
// call site that omits the prop still renders the signature default. Only
// string props are covered: the variadic wrapper makes "omitted" (x0 == "")
// distinguishable from "passed" (x[0]), and an empty string is a valid value
// that the default may replace. bool/int props get NO fallback — an explicit
// false/0 is a valid value and must never be overwritten, so defaults for
// non-string types are not supported (the caller must pass the value).
func writePropDefaultFallbacks(buf *strings.Builder, comp *ComponentDef) {
	for _, p := range comp.Props {
		if p.Default == "" || p.Type != "string" {
			continue
		}
		buf.WriteString(fmt.Sprintf("\t\tif %s == \"\" {\n\t\t\t%s = %s\n\t\t}\n", p.Name, p.Name, p.Default))
	}
}
