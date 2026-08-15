package core

import (
	"fmt"
	"strings"
)

func validateSelfClosingCalls(gen *generator, nodes []TemplateNode, src string) error {
	for i, n := range nodes {
		if n.Type == NodeComponentCall && n.SelfClose {
			def := gen.lookupDef(n.Tag)
			if def != nil && (def.HasDefaultSlot || def.HasNamedSlot) {
				if hasNonWhitespaceAfter(nodes, i+1) {
					loc := sourceLocation(src, n.Pos)
					if src == "" {
						loc = "?:?"
					}
					return fmt.Errorf("%s:%s: %s: self-closing call must not contain children", n.Source, loc, n.Tag)
				}
			}
		}
		if err := validateSelfClosingCalls(gen, n.Children, src); err != nil {
			return err
		}
		if err := validateSelfClosingCalls(gen, n.ElseChildren, src); err != nil {
			return err
		}
	}
	return nil
}

func hasNonWhitespaceAfter(nodes []TemplateNode, start int) bool {
	for i := start; i < len(nodes); i++ {
		if nodes[i].Type == NodeText && strings.TrimSpace(nodes[i].Content) == "" {
			continue
		}
		return true
	}
	return false
}

func hasDefaultSlot(nodes []TemplateNode) bool {
	return hasSlotWithContent(nodes, "")
}

func hasNamedSlot(nodes []TemplateNode) bool {
	for _, n := range nodes {
		if n.Type == NodeSlot && n.Content != "" {
			return true
		}
		if hasNamedSlot(n.Children) || hasNamedSlot(n.ElseChildren) {
			return true
		}
	}
	return false
}

func hasSlotWithContent(nodes []TemplateNode, content string) bool {
	for _, n := range nodes {
		if n.Type == NodeSlot && n.Content == content {
			return true
		}
		if hasSlotWithContent(n.Children, content) || hasSlotWithContent(n.ElseChildren, content) {
			return true
		}
	}
	return false
}
