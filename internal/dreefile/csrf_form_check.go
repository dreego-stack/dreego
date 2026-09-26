package dreefile

import (
	"regexp"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

var csrfFormOpenRe = regexp.MustCompile(`(?i)<form\b[^>]*\bg-action\s*=\s*"([^"]*)"[^>]*>`)
var csrfFormCloseRe = regexp.MustCompile(`(?i)</form\s*>`)
var csrfTokenFieldRe = regexp.MustCompile(`(?i)(<input\b[^>]*\bname\s*=\s*["']?csrf_token\b|CSRFInput\s*\()`)

// csrfFormDiagnostics warns for every form with a g-action that does not render
// a csrf_token field. Such a form generates a POST handler, and the default CSRF
// middleware rejects the request with 403 when the token is missing.
func csrfFormDiagnostics(nodes []ir.TemplateNode, fpath string) []string {
	full, positions, src := flattenFormText(nodes)
	var diagnostics []string
	for _, loc := range csrfFormOpenRe.FindAllStringSubmatchIndex(full, -1) {
		start, end := loc[0], loc[1]
		innerEnd := len(full)
		if rel := csrfFormCloseRe.FindStringIndex(full[end:]); rel != nil {
			innerEnd = end + rel[0]
		}
		if csrfTokenFieldRe.MatchString(full[end:innerEnd]) {
			continue
		}
		if start >= len(positions) {
			continue
		}
		line, col := gogen.PosToLineCol(src, positions[start])
		action := full[loc[2]:loc[3]]
		d := Diagnostic{
			File:  fpath,
			Line:  line,
			Col:   col,
			Cause: `form with g-action="` + action + `" has no csrf_token field; its POST will be rejected with 403 if CSRF is active`,
			Fix:   `render c.CSRFInput() inside the form, or add a hidden input named csrf_token`,
		}
		diagnostics = append(diagnostics, d.String())
	}
	return diagnostics
}

func flattenFormText(nodes []ir.TemplateNode) (string, []int, string) {
	var b strings.Builder
	var positions []int
	src := ""
	collectFormText(nodes, &b, &positions, &src)
	return b.String(), positions, src
}

func collectFormText(nodes []ir.TemplateNode, b *strings.Builder, positions *[]int, src *string) {
	for i := range nodes {
		collectFormTextIn(&nodes[i], b, positions, src)
	}
}

func collectFormTextIn(n *ir.TemplateNode, b *strings.Builder, positions *[]int, src *string) {
	if n.SourceText != "" && *src == "" {
		*src = n.SourceText
	}
	switch {
	case n.Type == ir.NodeText && n.Content != "":
		appendText(b, positions, n.Content, n.Pos)
	case n.Type == ir.NodeExpression && n.SourceText != "":
		raw := rawExpressionAt(n.SourceText, n.Pos)
		if raw == "" {
			raw = "{{" + n.Content + "}}"
		}
		appendText(b, positions, raw, n.Pos)
	}
	for i := range n.Children {
		collectFormTextIn(&n.Children[i], b, positions, src)
	}
	for i := range n.ElseChildren {
		collectFormTextIn(&n.ElseChildren[i], b, positions, src)
	}
}

func appendText(b *strings.Builder, positions *[]int, text string, pos int) {
	for i := 0; i < len(text); i++ {
		b.WriteByte(text[i])
		*positions = append(*positions, pos+i)
	}
}

func rawExpressionAt(srcText string, pos int) string {
	if pos < 0 || pos+2 > len(srcText) || srcText[pos:pos+2] != "{{" {
		return ""
	}
	rest := srcText[pos+2:]
	if end := strings.Index(rest, "}}"); end >= 0 {
		return srcText[pos : pos+2+end+2]
	}
	return ""
}
