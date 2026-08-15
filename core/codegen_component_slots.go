package core

import (
	"fmt"
)

func slotExists(def *ComponentDef, name string) bool {
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

func validateSlotName(def *ComponentDef, name, filename, src string, pos int) error {
	if def == nil || slotExists(def, name) {
		return nil
	}
	loc := sourceLocation(src, pos)
	if src == "" {
		loc = "?:?"
	}
	component := "?"
	if def != nil {
		component = def.Name
	}
	return fmt.Errorf("%s:%s: %s: unknown slot %q", filename, loc, component, name)
}

func nestedSlotError(call TemplateNode, def *ComponentDef, nested *TemplateNode, src string) error {
	component := call.Tag
	if def != nil && def.Name != "" {
		component = def.Name
	}
	loc := sourceLocation(src, nested.Pos)
	if src == "" {
		loc = "?:?"
	}
	return fmt.Errorf("%s:%s: %s: nested slot declaration is not allowed", call.Source, loc, component)
}
