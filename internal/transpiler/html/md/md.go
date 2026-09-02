package md

import (
	"fmt"
	"html"
	"regexp"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

var (
	atxHeading     = regexp.MustCompile(`^(#{1,6})\s+(.*)$`)
	ulItem         = regexp.MustCompile(`^[-*+]\s+(.*)$`)
	olItem         = regexp.MustCompile(`^(\d+)\.\s+(.*)$`)
	listMarker     = regexp.MustCompile(`^([-*+]|\d+\.)\s*$`)
	codeRe         = regexp.MustCompile("`([^`]*)`")
	linkRe         = regexp.MustCompile(`\[([^\]]*)\]\(([^)]*)\)`)
	strongRe       = regexp.MustCompile(`\*\*([^*]+)\*\*`)
	emRe           = regexp.MustCompile(`\*([^*]+)\*`)
	htmlBlockStart = regexp.MustCompile(`^</?([a-zA-Z][a-zA-Z0-9-]*)(\s|>|/)`)
)

func ToNodes(src string) ([]ir.TemplateNode, error) {
	return parseBlocks(src)
}

func parseBlocks(src string) ([]ir.TemplateNode, error) {
	lines := strings.Split(src, "\n")
	var nodes []ir.TemplateNode
	var para []string

	flushPara := func() {
		if len(para) > 0 {
			nodes = append(nodes, textNode("<p>"+renderInline(strings.Join(para, " "))+"</p>"))
			para = nil
		}
	}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "```") {
			flushPara()
			lang := strings.TrimSpace(strings.TrimPrefix(trimmed, "```"))
			var code []string
			i++
			for i < len(lines) && !strings.HasPrefix(strings.TrimSpace(lines[i]), "```") {
				code = append(code, lines[i])
				i++
			}
			open := "<pre><code"
			if lang != "" {
				open += ` class="language-` + lang + `"`
			}
			open += ">"
			nodes = append(nodes, textNode(open+strings.Join(code, "\n")+"</code></pre>"))
			continue
		}

		if htmlBlockStart.MatchString(trimmed) {
			flushPara()
			var raw []string
			for i < len(lines) && strings.TrimSpace(lines[i]) != "" {
				raw = append(raw, lines[i])
				i++
			}
			i--
			nodes = append(nodes, textNode(strings.Join(raw, "\n")))
			continue
		}

		if m := atxHeading.FindStringSubmatch(trimmed); m != nil {
			flushPara()
			level := len(m[1])
			nodes = append(nodes, textNode(fmt.Sprintf("<h%d>%s</h%d>", level, renderInline(m[2]), level)))
			continue
		}

		if isHR(trimmed) {
			flushPara()
			nodes = append(nodes, textNode("<hr>"))
			continue
		}

		if strings.HasPrefix(trimmed, ">") {
			flushPara()
			var q []string
			for i < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i]), ">") {
				q = append(q, strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(lines[i]), ">")))
				i++
			}
			i--
			nodes = append(nodes, textNode("<blockquote>"+renderInline(strings.Join(q, " "))+"</blockquote>"))
			continue
		}

		if listMarker.MatchString(trimmed) {
			return nil, fmt.Errorf("md: unparseable list marker on line %d", i+1)
		}

		if m := ulItem.FindStringSubmatch(trimmed); m != nil {
			flushPara()
			var items []string
			for i < len(lines) {
				lt := strings.TrimSpace(lines[i])
				if im := ulItem.FindStringSubmatch(lt); im != nil {
					items = append(items, "<li>"+renderInline(im[1])+"</li>")
					i++
				} else {
					break
				}
			}
			i--
			nodes = append(nodes, textNode("<ul>"+strings.Join(items, "")+"</ul>"))
			continue
		}

		if m := olItem.FindStringSubmatch(trimmed); m != nil {
			flushPara()
			var items []string
			for i < len(lines) {
				lt := strings.TrimSpace(lines[i])
				if im := olItem.FindStringSubmatch(lt); im != nil {
					items = append(items, "<li>"+renderInline(im[2])+"</li>")
					i++
				} else {
					break
				}
			}
			i--
			nodes = append(nodes, textNode("<ol>"+strings.Join(items, "")+"</ol>"))
			continue
		}

		if trimmed == "" {
			flushPara()
			continue
		}

		para = append(para, trimmed)
	}

	flushPara()
	return nodes, nil
}

func textNode(content string) ir.TemplateNode {
	return ir.TemplateNode{Type: ir.NodeText, Content: content}
}

func isHR(s string) bool {
	if !strings.HasPrefix(s, "-") {
		return false
	}
	for _, r := range s {
		if r != '-' {
			return false
		}
	}
	return len(s) >= 3
}

func renderInline(s string) string {
	var b strings.Builder
	for len(s) > 0 {
		i := strings.IndexByte(s, '<')
		if i < 0 || i+1 >= len(s) || !(isHTMLNameStart(s[i+1]) || s[i+1] == '/') {
			b.WriteString(renderInlineText(s))
			break
		}
		b.WriteString(renderInlineText(s[:i]))
		end := strings.IndexByte(s[i:], '>')
		if end < 0 {
			b.WriteString(s[i:])
			break
		}
		b.WriteString(s[i : i+end+1])
		s = s[i+end+1:]
	}
	return b.String()
}

func renderInlineText(s string) string {
	s = html.EscapeString(s)
	s = codeRe.ReplaceAllString(s, "<code>$1</code>")
	s = linkRe.ReplaceAllString(s, `<a href="$2">$1</a>`)
	s = strongRe.ReplaceAllString(s, "<strong>$1</strong>")
	s = emRe.ReplaceAllString(s, "<em>$1</em>")
	return s
}

func isHTMLNameStart(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}
