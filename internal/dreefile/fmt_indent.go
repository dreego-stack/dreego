package dreefile

import "strings"

const indentUnit = "    "

// indentHTMLLines re-indents markup by element nesting depth, using four spaces
// per level. A line's indentation is the depth before its own tags are counted,
// so a closing tag aligns with its opening tag. Content inside <pre> or
// {#verbatim} is left verbatim. Single-line bodies are returned unchanged.
func indentHTMLLines(text string, base int) string {
	lines := strings.Split(text, "\n")
	if len(lines) < 2 {
		return text
	}
	raw := rawRegions(lines)
	out := make([]string, len(lines))
	depth := base
	for i, line := range lines {
		if raw[i] {
			out[i] = line
			continue
		}
		if strings.TrimSpace(line) == "" {
			out[i] = ""
		} else {
			d := depth - leadingDedentTags(line)
			if d < 0 {
				d = 0
			}
			out[i] = strings.Repeat(indentUnit, d) + strings.TrimLeft(line, " \t")
		}
		depth += htmlDepthDelta(line)
		if depth < 0 {
			depth = 0
		}
	}
	return strings.Join(out, "\n")
}

// rawRegions marks the lines strictly inside a <pre>…</pre> or
// {#verbatim}…{/verbatim} pair. The opening and closing lines themselves are
// not raw, so they take part in normal depth counting and indentation.
func rawRegions(lines []string) []bool {
	raw := make([]bool, len(lines))
	open := -1
	closer := ""
	for i, line := range lines {
		low := strings.ToLower(line)
		if open < 0 {
			if idx := strings.Index(low, "<pre"); idx >= 0 && !strings.Contains(low[idx:], "</pre") {
				open, closer = i, "</pre"
			} else if idx := strings.Index(low, "{#verbatim}"); idx >= 0 && !strings.Contains(low[idx:], "{/verbatim}") {
				open, closer = i, "{/verbatim}"
			}
			continue
		}
		if strings.Contains(low, closer) {
			for j := open + 1; j < i; j++ {
				raw[j] = true
			}
			open, closer = -1, ""
		}
	}
	if open >= 0 {
		for j := open + 1; j < len(lines); j++ {
			raw[j] = true
		}
	}
	return raw
}

// sectionLanguageFromTag reads the lang attribute from a root section tag. It
// is used to detect a markdown body, which must not be re-indented.
func sectionLanguageFromTag(section string) string {
	openEnd := strings.IndexByte(section, '>')
	if openEnd < 0 {
		return ""
	}
	tag := section[:openEnd+1]
	idx := strings.Index(tag, "lang=")
	if idx < 0 {
		return ""
	}
	value := strings.TrimLeft(tag[idx+len("lang="):], " \t")
	if value == "" {
		return ""
	}
	if value[0] == '"' || value[0] == '\'' {
		q := value[0]
		if end := strings.IndexByte(value[1:], q); end >= 0 {
			return strings.ToLower(value[1 : 1+end])
		}
		return ""
	}
	end := strings.IndexAny(value, " \t>")
	if end < 0 {
		return ""
	}
	return strings.ToLower(value[:end])
}

// htmlDepthDelta counts the element and block nesting change of one line. It
// ignores comments, markup-looking text inside text nodes, and the contents of
// {{ expression }} and [[ message ]] spans, so only real tags move the depth.
func htmlDepthDelta(line string) int {
	delta := 0
	i := 0
	for i < len(line) {
		switch line[i] {
		case '<':
			if strings.HasPrefix(line[i:], "<!--") {
				end := strings.Index(line[i:], "-->")
				if end < 0 {
					return delta
				}
				i += end + 3
				continue
			}
			nxt := byte(0)
			if i+1 < len(line) {
				nxt = line[i+1]
			}
			if nxt != '/' && nxt != '@' && !isASCIILetter(nxt) {
				i++
				continue
			}
			end := tagEndIn(line, i)
			if end < 0 {
				i++
				continue
			}
			inner := line[i+1 : end]
			switch {
			case strings.HasPrefix(inner, "/"):
				delta--
			case isVoidElement(inner), strings.HasSuffix(strings.TrimSpace(inner), "/"):
			default:
				delta++
			}
			i = end + 1
		case '{':
			if n := skipSpan(line[i:], "{{", "}}"); n > 0 {
				i += n
				continue
			}
			if d, n := ctrlDeltaAt(line[i:]); n > 0 {
				delta += d
				i += n
				continue
			}
			i++
		case '[':
			if n := skipSpan(line[i:], "[[", "]]"); n > 0 {
				i += n
				continue
			}
			i++
		default:
			i++
		}
	}
	return delta
}

// leadingDedentTags counts the closing tags and block terminators at the start
// of a line, so a line that closes an element or block aligns with the token
// that opened it. {#else} and friends dedent to their opening block.
func leadingDedentTags(line string) int {
	s := strings.TrimLeft(line, " \t")
	n := 0
	for {
		switch {
		case strings.HasPrefix(s, "</"):
			end := tagEndIn(s, 0)
			if end < 0 {
				return n
			}
			n++
			s = strings.TrimLeft(s[end+1:], " \t")
		case strings.HasPrefix(s, "{/if}"), strings.HasPrefix(s, "{/each}"), strings.HasPrefix(s, "{/slot}"):
			n++
			s = strings.TrimLeft(s[strings.IndexByte(s, '}')+1:], " \t")
		case isMidBlockAt(s):
			n++
			s = strings.TrimLeft(s[strings.IndexByte(s, '}')+1:], " \t")
		default:
			return n
		}
	}
}

// ctrlDeltaAt returns the nesting change of a control-flow token at the start
// of s, and the token length. Mid-markers such as {#else} report a zero delta:
// they dedent themselves but their content stays nested.
func ctrlDeltaAt(s string) (int, int) {
	switch {
	case isMidBlockAt(s):
		return 0, strings.IndexByte(s, '}') + 1
	case strings.HasPrefix(s, "{#if "), strings.HasPrefix(s, "{#each "), strings.HasPrefix(s, "{#slot "):
		return 1, strings.IndexByte(s, '}') + 1
	case strings.HasPrefix(s, "{/if}"), strings.HasPrefix(s, "{/each}"), strings.HasPrefix(s, "{/slot}"):
		return -1, strings.IndexByte(s, '}') + 1
	}
	return 0, 0
}

func isMidBlockAt(s string) bool {
	return strings.HasPrefix(s, "{#else}") ||
		strings.HasPrefix(s, "{#else if ") ||
		strings.HasPrefix(s, "{#each else}")
}

func tagEndIn(line string, open int) int {
	var quote byte
	for i := open + 1; i < len(line); i++ {
		c := line[i]
		if quote != 0 {
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '"', '\'':
			quote = c
		case '>':
			return i
		}
	}
	return -1
}

func skipSpan(s, open, close string) int {
	if !strings.HasPrefix(s, open) {
		return 0
	}
	if idx := strings.Index(s[len(open):], close); idx >= 0 {
		return len(open) + idx + len(close)
	}
	return 0
}

func isASCIILetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// voidElements are HTML elements that never have a closing tag, so they must
// not increase the nesting depth.
var voidElements = map[string]bool{
	"area": true, "base": true, "br": true, "col": true, "embed": true,
	"hr": true, "img": true, "input": true, "link": true, "meta": true,
	"param": true, "source": true, "track": true, "wbr": true,
}

func isVoidElement(inner string) bool {
	name := inner
	if idx := strings.IndexAny(name, " \t/"); idx >= 0 {
		name = name[:idx]
	}
	return voidElements[strings.ToLower(name)]
}
