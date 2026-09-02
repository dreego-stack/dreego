package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/html/md"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
)

var mdOpenRe = regexp.MustCompile(`^<md(\s[^>]*)?>$`)

// transformInlineMd rewrites <md ...>...</md> regions inside a lang="html"
// body into a <div> wrapper around the markdown-converted node stream. The
// lexer emits the open and close tags as their own NodeText nodes, so the
// pre-pass walks the flat node stream and pairs them up.
func transformInlineMd(nodes []ir.TemplateNode) ([]ir.TemplateNode, error) {
	var out []ir.TemplateNode
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		if n.Type != ir.NodeText || !mdOpenRe.MatchString(n.Content) {
			out = append(out, n)
			continue
		}
		region, next, err := collectMdRegion(nodes, i)
		if err != nil {
			return nil, err
		}
		converted, err := md.TransformNodes(region)
		if err != nil {
			return nil, err
		}
		out = append(out, ir.TemplateNode{Type: ir.NodeText, Content: mdDivOpen(n.Content), Pos: n.Pos})
		out = append(out, converted...)
		out = append(out, ir.TemplateNode{Type: ir.NodeText, Content: "</div>", Pos: n.Pos})
		i = next
	}
	return out, nil
}

// collectMdRegion gathers the node stream between an <md> open tag and its
// matching </md> close tag, returning the region nodes and the index of the
// close tag. It rejects nested <md>, control flow inside the region, and an
// unclosed region.
func collectMdRegion(nodes []ir.TemplateNode, start int) ([]ir.TemplateNode, int, error) {
	open := nodes[start]
	var region []ir.TemplateNode
	if rest := open.Content[mdOpenTagLen(open.Content):]; rest != "" {
		region = append(region, ir.TemplateNode{Type: ir.NodeText, Content: rest, Pos: open.Pos})
	}
	for i := start + 1; i < len(nodes); i++ {
		n := nodes[i]
		switch {
		case n.Type == ir.NodeText && strings.Contains(n.Content, "</md>"):
			idx := strings.Index(n.Content, "</md>")
			if before := n.Content[:idx]; before != "" {
				region = append(region, ir.TemplateNode{Type: ir.NodeText, Content: before, Pos: n.Pos})
			}
			return region, i, nil
		case n.Type == ir.NodeIf || n.Type == ir.NodeEach:
			return nil, 0, fmt.Errorf("md blocks cannot span dreego control flow; close the md tag first at position %d", n.Pos)
		case n.Type == ir.NodeText && mdOpenRe.MatchString(n.Content):
			return nil, 0, fmt.Errorf("nested <md> tag at position %d is not supported", n.Pos)
		default:
			region = append(region, n)
		}
	}
	return nil, 0, fmt.Errorf("unclosed <md> tag at position %d", open.Pos)
}

func mdOpenTagLen(content string) int {
	m := mdOpenRe.FindStringIndex(content)
	if m == nil {
		return 0
	}
	return m[1]
}

// mdDivOpen builds the <div> open tag for an <md> open tag, moving the class
// attribute to the div and passing every other attribute through verbatim.
func mdDivOpen(openTag string) string {
	inner := strings.TrimSuffix(strings.TrimPrefix(openTag, "<md"), ">")
	inner = strings.TrimSpace(inner)
	classVal, rest := extractClass(inner)
	var b strings.Builder
	b.WriteString("<div")
	if classVal != "" {
		b.WriteString(` class="` + classVal + `"`)
	}
	if rest != "" {
		b.WriteString(" " + rest)
	}
	b.WriteString(">")
	return b.String()
}

// extractClass pulls the class attribute value out of an attribute string and
// returns it together with the remaining attributes, preserving quotes.
func extractClass(attrs string) (classVal, rest string) {
	i := 0
	for i < len(attrs) {
		for i < len(attrs) && attrs[i] == ' ' {
			i++
		}
		if i >= len(attrs) {
			break
		}
		start := i
		for i < len(attrs) && attrs[i] != ' ' && attrs[i] != '=' {
			i++
		}
		name := attrs[start:i]
		for i < len(attrs) && attrs[i] == ' ' {
			i++
		}
		if i < len(attrs) && attrs[i] == '=' {
			i++
			for i < len(attrs) && attrs[i] == ' ' {
				i++
			}
			if i < len(attrs) && (attrs[i] == '"' || attrs[i] == '\'') {
				q := attrs[i]
				i++
				vstart := i
				for i < len(attrs) && attrs[i] != q {
					i++
				}
				val := attrs[vstart:i]
				if i < len(attrs) {
					i++
				}
				if name == "class" {
					classVal = val
				} else {
					rest += " " + name + "=" + string(q) + val + string(q)
				}
			} else {
				vstart := i
				for i < len(attrs) && attrs[i] != ' ' {
					i++
				}
				val := attrs[vstart:i]
				if name == "class" {
					classVal = val
				} else {
					rest += " " + name + "=" + val
				}
			}
		} else if name != "class" {
			rest += " " + name
		}
	}
	return strings.TrimSpace(classVal), strings.TrimSpace(rest)
}

// rejectMdTagInMdBody reports an error when an <md> tag appears inside a body
// whose language is already md.
func rejectMdTagInMdBody(nodes []ir.TemplateNode) error {
	for _, n := range nodes {
		if n.Type == ir.NodeText && mdOpenRe.MatchString(n.Content) {
			return fmt.Errorf("<md> inside <body lang=\"md\"> at position %d is not supported; md is already the body language", n.Pos)
		}
	}
	return nil
}
