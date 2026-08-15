package core

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

func collectSlotNames(nodes []TemplateNode) []string {
	seen := map[string]bool{}
	var names []string
	var walk func([]TemplateNode)
	walk = func(list []TemplateNode) {
		for _, n := range list {
			if n.Type == NodeSlot && n.Content != "" && !seen[n.Content] {
				seen[n.Content] = true
				names = append(names, n.Content)
			}
			walk(n.Children)
			walk(n.ElseChildren)
		}
	}
	walk(nodes)
	return names
}

func mergeUnique(a, b []string) []string {
	seen := map[string]bool{}
	for _, s := range a {
		seen[s] = true
	}
	for _, s := range b {
		seen[s] = true
	}
	var out []string
	for _, s := range a {
		out = append(out, s)
		seen[s] = false
	}
	for _, s := range b {
		if seen[s] {
			out = append(out, s)
		}
	}
	return out
}

func findNestedSlot(nodes []TemplateNode) *TemplateNode {
	for i := range nodes {
		n := &nodes[i]
		if n.Type == NodeSlot {
			return n
		}
		if n.Type != NodeComponentCall {
			if found := findNestedSlot(n.Children); found != nil {
				return found
			}
			if found := findNestedSlot(n.ElseChildren); found != nil {
				return found
			}
		}
	}
	return nil
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
