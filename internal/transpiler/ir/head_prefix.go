package ir

import "strings"

const HeadPlaceholder = "{#head}"

// StaticBodyHead locates the first body-level {#head} placeholder and returns
// the index of the text node carrying it, the literal markup before it and the
// literal remainder after it on the same node. It reports false when no
// placeholder exists, a non-text node precedes it, or the prefix contains
// template syntax that cannot be captured as a single literal.
func (f *File) StaticBodyHead() (idx int, prefix, after string, ok bool) {
	if f == nil || f.Body == nil {
		return -1, "", "", false
	}
	nodes := f.Body.Nodes
	for i := range nodes {
		if nodes[i].Type != NodeText {
			return -1, "", "", false
		}
		pos := strings.Index(nodes[i].Content, HeadPlaceholder)
		if pos < 0 {
			continue
		}
		var b strings.Builder
		for j := 0; j < i; j++ {
			b.WriteString(nodes[j].Content)
		}
		prefix = b.String() + nodes[i].Content[:pos]
		after = nodes[i].Content[pos+len(HeadPlaceholder):]
		if strings.Contains(prefix, "{#") || strings.Contains(prefix, "{{") || strings.Contains(prefix, "[[") {
			return -1, "", "", false
		}
		return i, prefix, after, true
	}
	return -1, "", "", false
}

// HasHeadDedupeTag reports whether the given markup contains a tag whose layout
// copy a route head may override: a <title> or a meta description.
func HasHeadDedupeTag(s string) bool {
	return strings.Contains(s, "<title") ||
		strings.Contains(s, `name="description"`) ||
		strings.Contains(s, "name='description'")
}
