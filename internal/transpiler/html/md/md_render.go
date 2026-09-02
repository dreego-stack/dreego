package md

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

var (
	imageRe       = regexp.MustCompile(`!\[([^\]]*)\]\(([^)]*)\)`)
	footnoteRefRe = regexp.MustCompile(`\[\^(\d+)\]`)
	footnoteDefRe = regexp.MustCompile(`^\[\^(\d+)\]:\s*(.*)$`)
)

var activeR *mdRenderer

func renderInline(s string) string {
	if activeR == nil {
		activeR = newRenderer()
	}
	return activeR.renderInline(s)
}

type mdRenderer struct {
	defs      map[string]string
	refCounts map[string]int
	refOrder  []string
}

func newRenderer() *mdRenderer {
	return &mdRenderer{defs: map[string]string{}, refCounts: map[string]int{}}
}

func (r *mdRenderer) addDef(num, content string) {
	if _, ok := r.defs[num]; !ok {
		content = strings.TrimSuffix(content, ".")
		r.defs[num] = r.renderInline(content)
	}
}

func (r *mdRenderer) hasDefs() bool {
	return len(r.defs) > 0
}

func (r *mdRenderer) renderInline(s string) string {
	var b strings.Builder
	for len(s) > 0 {
		i := strings.IndexByte(s, '<')
		if i < 0 || i+1 >= len(s) || !(isHTMLNameStart(s[i+1]) || s[i+1] == '/') {
			b.WriteString(r.renderInlineText(s))
			break
		}
		b.WriteString(r.renderInlineText(s[:i]))
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

func (r *mdRenderer) renderInlineText(s string) string {
	s = html.EscapeString(s)
	s = codeRe.ReplaceAllString(s, "<code>$1</code>")
	s = imageRe.ReplaceAllString(s, `<img src="$2" alt="$1">`)
	s = linkRe.ReplaceAllString(s, `<a href="$2">$1</a>`)
	s = strongRe.ReplaceAllString(s, "<strong>$1</strong>")
	s = emRe.ReplaceAllString(s, "<em>$1</em>")
	s = r.replaceFootnoteRefs(s)
	return s
}

func (r *mdRenderer) replaceFootnoteRefs(s string) string {
	return footnoteRefRe.ReplaceAllStringFunc(s, func(m string) string {
		num := footnoteRefRe.FindStringSubmatch(m)[1]
		r.refCounts[num]++
		count := r.refCounts[num]
		if count == 1 {
			r.refOrder = append(r.refOrder, num)
		}
		return `<sup class="footnote-ref"><a href="#fn-` + num + `" id="fnref-` + num + `">` + strconv.Itoa(count) + `</a></sup>`
	})
}

func (r *mdRenderer) footnotesSection() string {
	if len(r.refOrder) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<section class="footnotes"><ol>`)
	for _, num := range r.refOrder {
		content := r.defs[num]
		b.WriteString(`<li id="fn-` + num + `">` + content + ` <a href="#fnref-` + num + `" class="footnote-backref">↩</a></li>`)
	}
	b.WriteString(`</ol></section>`)
	return b.String()
}
