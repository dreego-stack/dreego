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
		if containsTemplateSyntax(prefix) {
			return -1, "", "", false
		}
		return i, prefix, after, true
	}
	return -1, "", "", false
}

// BodyHeadConflict reports a body-level {#head} placeholder whose preceding
// head markup cannot be captured as a literal, so runtime head dedupe is
// disabled for that markup. nodeIndex is the index of the text node carrying
// the placeholder and reason names the blocking shape. It reports false when no
// placeholder exists or when the placeholder is preceded by literal markup
// only.
func (f *File) BodyHeadConflict() (nodeIndex int, reason string, found bool) {
	if f == nil || f.Body == nil {
		return -1, "", false
	}
	nodes := f.Body.Nodes
	placeholderIdx := -1
	for i := range nodes {
		if nodes[i].Type == NodeText && strings.Contains(nodes[i].Content, HeadPlaceholder) {
			placeholderIdx = i
			break
		}
	}
	if placeholderIdx < 0 {
		return -1, "", false
	}
	if _, _, _, ok := f.StaticBodyHead(); ok {
		return -1, "", false
	}
	for i := 0; i < placeholderIdx; i++ {
		if nodes[i].Type != NodeText {
			return placeholderIdx, "a non-text node precedes the " + HeadPlaceholder + " placeholder", true
		}
	}
	return placeholderIdx, "the " + HeadPlaceholder + " prefix contains template syntax that cannot be captured as a literal", true
}

func containsTemplateSyntax(s string) bool {
	return strings.Contains(s, "{#") || strings.Contains(s, "{{") || strings.Contains(s, "[[")
}

// StaticBodyHeadTail returns the literal markup that follows the {#head}
// placeholder: the remainder of the placeholder node plus every following text
// node until a non-text node or embedded template syntax appears. end is the
// index of the first node that is not part of the captured markup, so callers
// can skip the consumed nodes. It reports false when the placeholder remainder
// itself contains template syntax that cannot be captured as a literal.
func (f *File) StaticBodyHeadTail(idx int, after string) (tail string, end int, ok bool) {
	if f == nil || f.Body == nil || idx < 0 {
		return "", -1, false
	}
	if strings.Contains(after, "{#") || strings.Contains(after, "{{") || strings.Contains(after, "[[") {
		return "", -1, false
	}
	var b strings.Builder
	b.WriteString(after)
	nodes := f.Body.Nodes
	i := idx + 1
	for ; i < len(nodes); i++ {
		if nodes[i].Type != NodeText {
			break
		}
		if strings.Contains(nodes[i].Content, "{#") || strings.Contains(nodes[i].Content, "{{") || strings.Contains(nodes[i].Content, "[[") {
			break
		}
		b.WriteString(nodes[i].Content)
	}
	return b.String(), i, true
}

// HasHeadDedupeTag reports whether the given markup contains a tag whose layout
// copy a route head may override: a <title>, a meta description, a meta
// viewport, or a meta charset. HTML is case-insensitive, so the match folds the
// input before comparing.
func HasHeadDedupeTag(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "<title") ||
		strings.Contains(lower, "charset") ||
		strings.Contains(lower, `name="description"`) ||
		strings.Contains(lower, "name='description'") ||
		strings.Contains(lower, `name="viewport"`) ||
		strings.Contains(lower, "name='viewport'")
}
