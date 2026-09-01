package output

import (
	"fmt"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func slotExists(def *ir.ComponentDef, name string) bool {
	if def == nil {
		return false
	}
	for _, s := range def.Slots {
		if s == name {
			return true
		}
	}
	return false
}

func ValidateSlotName(def *ir.ComponentDef, name, filename, src string, pos int) error {
	if def == nil || slotExists(def, name) {
		return nil
	}
	loc := ir.SourceLocation(src, pos)
	if src == "" {
		loc = "?:?"
	}
	component := "?"
	if def != nil {
		component = def.Name
	}
	return fmt.Errorf("%s:%s: %s: unknown slot %q", filename, loc, component, name)
}

func NestedSlotError(call ir.TemplateNode, def *ir.ComponentDef, nested *ir.TemplateNode, src string) error {
	component := call.Tag
	if def != nil && def.Name != "" {
		component = def.Name
	}
	loc := ir.SourceLocation(src, nested.Pos)
	if src == "" {
		loc = "?:?"
	}
	return fmt.Errorf("%s:%s: %s: nested slot declaration is not allowed", call.Source, loc, component)
}

func FindNestedSlot(nodes []ir.TemplateNode) *ir.TemplateNode {
	for i := range nodes {
		n := &nodes[i]
		if n.Type == ir.NodeSlot {
			return n
		}
		if n.Type != ir.NodeComponentCall {
			if found := FindNestedSlot(n.Children); found != nil {
				return found
			}
			if found := FindNestedSlot(n.ElseChildren); found != nil {
				return found
			}
		}
	}
	return nil
}
