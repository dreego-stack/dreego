package core

import (
	"fmt"
	"strings"
)

func compTextWithAttrs(s string) string {
	code, _ := compTextSection(s, false)
	return code
}

// compTextSection renders a NodeText content segment to Go code, resolving {{ … }}
// placeholders inside quoted attribute values but leaving <script>/<style> section
// bodies literal (the lexer treats those as raw text where {{ … }} is not an expression).
// It returns the generated code and the section state after this segment so the
// caller can carry section tracking across sibling text nodes.
func compTextSection(content string, inSection bool) (string, bool) {
	var parts []string
	cur := inSection
	var quote byte
	tagStart := -1
	start := 0
	i := 0
	for i < len(content) {
		if cur {
			if closeLen := sectionCloseLen(content[i:]); closeLen > 0 {
				end := i + closeLen
				parts = append(parts, goLiteral(content[start:end]))
				start = end
				i = end
				cur = false
				continue
			}
			i++
			continue
		}
		if strings.HasPrefix(content[i:], "<script") || strings.HasPrefix(content[i:], "<style") {
			if start < i {
				parts = append(parts, goLiteral(content[start:i]))
			}
			start = i
			cur = true
			quote = 0
			tagStart = -1
			continue
		}
		if content[i] == '<' {
			if start < i {
				parts = append(parts, goLiteral(content[start:i]))
			}
			start = i
			tagStart = i
			quote = 0
			i++
			continue
		}
		if (content[i] == '"' || content[i] == '\'') && (i == 0 || content[i-1] != '\\') {
			if quote == 0 {
				quote = content[i]
			} else if quote == content[i] {
				quote = 0
			}
			i++
			continue
		}
		if quote != 0 && strings.HasPrefix(content[i:], "{{") {
			closeIdx := strings.Index(content[i+2:], "}}")
			if closeIdx < 0 {
				i++
				continue
			}
			if start < i {
				parts = append(parts, goLiteral(content[start:i]))
			}
			expr := strings.TrimSpace(content[i+2 : i+2+closeIdx])
			code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr)
			parts = append(parts, fmt.Sprintf("dreego.%s(%s)", attrSafeFunc(content, tagStart, i), code))
			i += 2 + closeIdx + 2
			start = i
			continue
		}
		i++
	}
	if start < len(content) {
		parts = append(parts, goLiteral(content[start:]))
	}
	return strings.Join(parts, " + "), cur
}

// attrSafeFunc picks the safe-value rule for an expression placeholder inside a
// quoted attribute value. tagStart is the index of the '<' that opens the tag
// containing the placeholder; i is the placeholder index in content.
func attrSafeFunc(content string, tagStart, i int) string {
	if tagStart < 0 || tagStart >= i {
		return "SafeAttr"
	}
	tagEnd := tagEnd(content[tagStart:])
	if tagEnd < 0 {
		return "SafeAttr"
	}
	tagEnd += tagStart
	name := attrNameAt(content[tagStart:tagEnd], i-tagStart)
	if name == "" {
		return "SafeAttr"
	}
	return attrContext(name)
}

func sectionCloseLen(s string) int {
	for _, tag := range []string{"</script>", "</style>"} {
		if strings.HasPrefix(s, tag) {
			return len(tag)
		}
	}
	return 0
}
