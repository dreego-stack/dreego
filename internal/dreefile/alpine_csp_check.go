package dreefile

import (
	"regexp"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
	"github.com/dreego-stack/dreego/internal/dreefile/ir"
)

// alpineDirectiveRe matches Alpine.js directives in static markup. Alpine
// evaluates these values with new Function(), which needs `'unsafe-eval'` in
// the Content-Security-Policy. The framework's default CSP does not allow
// eval, so an Alpine page with the default policy is inert.
var alpineDirectiveRe = regexp.MustCompile(`(?i)(\sx-data\s*=|\sx-init\s*=|\sx-effect\s*=|\sx-on:[a-z-]+\s*=|\sx-bind:[a-z-]+\s*=|\sx-model\s*=|\sx-show\s*=|\sx-text\s*=|\sx-html\s*=|\sx-transition\b|\s@[a-z-]+\s*=)`)

// alpineCSPDiagnostic warns when a route uses Alpine directives under the
// default CSP, which blocks Alpine because it needs `'unsafe-eval'`.
func alpineCSPDiagnostic(nodes []ir.TemplateNode, fpath string) (string, bool) {
	node, pos, found := firstAlpineDirective(nodes)
	if !found {
		return "", false
	}
	line, col := gogen.PosToLineCol(node.SourceText, pos)
	d := Diagnostic{
		File:  fpath,
		Line:  line,
		Col:   col,
		Cause: "Alpine.js directives need 'unsafe-eval'; the default CSP does not allow it, so Alpine stays inert unless you widen the policy",
		Fix:   "call app.SetCSP with a script-src that includes 'unsafe-eval' (see the progressive-enhancement guide) or use plain JavaScript instead",
	}
	return d.String(), true
}

func firstAlpineDirective(nodes []ir.TemplateNode) (*ir.TemplateNode, int, bool) {
	for i := range nodes {
		if node, pos, found := alpineDirectiveInNode(&nodes[i]); found {
			return node, pos, true
		}
	}
	return nil, 0, false
}

func alpineDirectiveInNode(n *ir.TemplateNode) (*ir.TemplateNode, int, bool) {
	if n.Type == ir.NodeText {
		if m := alpineDirectiveRe.FindStringIndex(n.Content); m != nil {
			return n, n.Pos + m[0], true
		}
	}
	for i := range n.Children {
		if node, pos, found := alpineDirectiveInNode(&n.Children[i]); found {
			return node, pos, true
		}
	}
	for i := range n.ElseChildren {
		if node, pos, found := alpineDirectiveInNode(&n.ElseChildren[i]); found {
			return node, pos, true
		}
	}
	return nil, 0, false
}
