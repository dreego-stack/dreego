package output

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

func CompTextWithAttrs(s string) string {
	code, _ := CompTextSection(s, false)
	return code
}

func CompTextSection(content string, inSection bool) (string, bool) {
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
				parts = append(parts, ir.GoLiteral(content[start:end]))
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
				parts = append(parts, ir.GoLiteral(content[start:i]))
			}
			start = i
			cur = true
			quote = 0
			tagStart = -1
			continue
		}
		if content[i] == '<' {
			if start < i {
				parts = append(parts, ir.GoLiteral(content[start:i]))
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
			closeIdx := ir.FindExprEnd(content[i+2:])
			if closeIdx < 0 {
				i++
				continue
			}
			if start < i {
				parts = append(parts, ir.GoLiteral(content[start:i]))
			}
			expr := strings.TrimSpace(content[i+2 : i+2+closeIdx])
			code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", expr)
			parts = append(parts, fmt.Sprintf("dreego.%s(%s)", AttrSafeFunc(content, tagStart, i), code))
			i += 2 + closeIdx + 2
			start = i
			continue
		}
		i++
	}
	if start < len(content) {
		parts = append(parts, ir.GoLiteral(content[start:]))
	}
	return strings.Join(parts, " + "), cur
}

func AttrSafeFunc(content string, tagStart, i int) string {
	if tagStart < 0 || tagStart >= i {
		return "SafeAttr"
	}
	tagEnd := ir.TagEnd(content[tagStart:])
	if tagEnd < 0 {
		return "SafeAttr"
	}
	tagEnd += tagStart
	name := ir.AttrNameAt(content[tagStart:tagEnd], i-tagStart)
	if name == "" {
		return "SafeAttr"
	}
	return ir.AttrContext(name)
}

func sectionCloseLen(s string) int {
	for _, tag := range []string{"</script>", "</style>"} {
		if strings.HasPrefix(s, tag) {
			return len(tag)
		}
	}
	return 0
}
