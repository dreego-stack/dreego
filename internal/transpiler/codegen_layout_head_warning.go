package transpiler

import "github.com/dreego-stack/dreego/internal/transpiler/ir"

func layoutHeadDedupeWarning(file *File) (string, bool) {
	if file == nil || file.Head != nil || file.Body == nil {
		return "", false
	}
	idx, reason, found := file.BodyHeadConflict()
	if !found {
		return "", false
	}
	if !bodyHasHeadDedupeTag(file) {
		return "", false
	}
	node := file.Body.Nodes[idx]
	line, col := ir.PosToLineCol(node.SourceText, node.Pos)
	location := node.Source
	if location == "" {
		location = file.SourcePath
	}
	if location == "" {
		location = "layout"
	}
	diagnostic := Diagnostic{
		File:  location,
		Line:  line,
		Col:   col,
		Cause: "body-level " + ir.HeadPlaceholder + " cannot be captured as literal markup (" + reason + "); head dedupe is disabled and a route title or description may be duplicated",
		Fix:   "keep " + ir.HeadPlaceholder + " in a text node after the layout head markup and before any component, expression, or message",
	}
	return diagnostic.String(), true
}

func bodyHasHeadDedupeTag(file *File) bool {
	for i := range file.Body.Nodes {
		if ir.HasHeadDedupeTag(file.Body.Nodes[i].Content) {
			return true
		}
	}
	return false
}
