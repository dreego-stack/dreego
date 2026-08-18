package core

import (
	"fmt"
	"regexp"
	"strings"
)

type Diagnostic struct {
	File  string
	Line  int
	Col   int
	Cause string
	Fix   string
}

func (d Diagnostic) String() string {
	return fmt.Sprintf("%s:%d:%d: %s Fix: %s", d.File, d.Line, d.Col, d.Cause, d.Fix)
}

func a11yCheck(nodes []TemplateNode) []Diagnostic {
	var diags []Diagnostic
	labelForIDs := collectLabelForIDs(nodes)
	scanA11y(nodes, labelForIDs, &diags)
	return diags
}

func a11yDiagnostics(nodes []TemplateNode) []string {
	ds := a11yCheck(nodes)
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.String()
	}
	return out
}

func collectLabelForIDs(nodes []TemplateNode) map[string]bool {
	ids := map[string]bool{}
	for i := range nodes {
		collectLabelForIDsIn(&nodes[i], ids)
	}
	return ids
}

func collectLabelForIDsIn(n *TemplateNode, ids map[string]bool) {
	if n.Type == NodeText {
		for _, label := range labelRe.FindAllString(n.Content, -1) {
			if forID := strings.TrimSpace(a11yAttrValue(label, "for")); forID != "" {
				ids[forID] = true
			}
		}
	}
	for i := range n.Children {
		collectLabelForIDsIn(&n.Children[i], ids)
	}
	for i := range n.ElseChildren {
		collectLabelForIDsIn(&n.ElseChildren[i], ids)
	}
}

func scanA11y(nodes []TemplateNode, labelForIDs map[string]bool, diags *[]Diagnostic) {
	for i := range nodes {
		walkA11y(&nodes[i], labelForIDs, diags)
	}
}

func walkA11y(n *TemplateNode, labelForIDs map[string]bool, diags *[]Diagnostic) {
	if n.Type == NodeText {
		checkImgAlt(n, diags)
		checkInputLabel(n, labelForIDs, diags)
	}
	for i := range n.Children {
		walkA11y(&n.Children[i], labelForIDs, diags)
	}
	for i := range n.ElseChildren {
		walkA11y(&n.ElseChildren[i], labelForIDs, diags)
	}
}

var imgRe = regexp.MustCompile(`(?i)<img\b([^>]*)>`)
var inputRe = regexp.MustCompile(`(?i)<input\b([^>]*)>`)
var labelRe = regexp.MustCompile(`(?i)<label\b[^>]*>`)
var labelOpenRe = regexp.MustCompile(`(?i)<label\b[^>]*>`)
var labelCloseRe = regexp.MustCompile(`(?i)</label\s*>`)
var attrDoubleRe = regexp.MustCompile(`(?i)\b(\w[\w-]*)\s*=\s*"([^"]*)"`)
var attrSingleRe = regexp.MustCompile(`(?i)\b(\w[\w-]*)\s*=\s*'([^']*)'`)

func checkImgAlt(n *TemplateNode, diags *[]Diagnostic) {
	for _, m := range imgRe.FindAllStringSubmatchIndex(n.Content, -1) {
		full := n.Content[m[0]:m[1]]
		attrs := n.Content[m[2]:m[3]]
		if !a11yHasAttr(attrs, "alt") {
			pos := n.Pos + m[0]
			line, col := posToLineCol(n.SourceText, pos)
			*diags = append(*diags, Diagnostic{
				File:  n.Source,
				Line:  line,
				Col:   col,
				Cause: fmt.Sprintf("<img> missing alt attribute (src=%q)", a11yAttrValue(full, "src")),
				Fix:   "add an alt attribute describing the image, or alt=\"\" for decorative images",
			})
		}
	}
}

func checkInputLabel(n *TemplateNode, labelForIDs map[string]bool, diags *[]Diagnostic) {
	for _, m := range inputRe.FindAllStringSubmatchIndex(n.Content, -1) {
		full := n.Content[m[0]:m[1]]
		typ := strings.ToLower(a11yAttrValue(full, "type"))
		if typ == "hidden" || typ == "submit" || typ == "button" || typ == "reset" {
			continue
		}
		id := a11yAttrValue(full, "id")
		if id != "" && labelForIDs[id] {
			continue
		}
		if a11yWrappedByLabel(n.SourceText, n.Pos+m[0], n.Pos+m[1]) {
			continue
		}
		pos := n.Pos + m[0]
		line, col := posToLineCol(n.SourceText, pos)
		name := a11yAttrValue(full, "name")
		*diags = append(*diags, Diagnostic{
			File:  n.Source,
			Line:  line,
			Col:   col,
			Cause: fmt.Sprintf("<input> without associated label (name=%q)", name),
			Fix:   "add a <label for=\"id\"> matching the input id",
		})
	}
}

func a11yHasAttr(attrs, name string) bool {
	for _, m := range attrDoubleRe.FindAllStringSubmatch(attrs, -1) {
		if strings.EqualFold(m[1], name) {
			return true
		}
	}
	for _, m := range attrSingleRe.FindAllStringSubmatch(attrs, -1) {
		if strings.EqualFold(m[1], name) {
			return true
		}
	}
	return false
}

func a11yAttrValue(tag, name string) string {
	for _, m := range attrDoubleRe.FindAllStringSubmatch(tag, -1) {
		if strings.EqualFold(m[1], name) {
			return m[2]
		}
	}
	for _, m := range attrSingleRe.FindAllStringSubmatch(tag, -1) {
		if strings.EqualFold(m[1], name) {
			return m[2]
		}
	}
	return ""
}

func a11yWrappedByLabel(content string, start, end int) bool {
	opens := labelOpenRe.FindAllStringIndex(content[:start], -1)
	if len(opens) == 0 {
		return false
	}
	closes := labelCloseRe.FindAllStringIndex(content[:start], -1)
	lastOpen := opens[len(opens)-1][0]
	if len(closes) > 0 && closes[len(closes)-1][0] > lastOpen {
		return false
	}
	return labelCloseRe.FindStringIndex(content[end:]) != nil
}
