package core

import (
	"fmt"
	"strings"
)

func genTempl(gen *generator, file *File, layout *File, scopeHash string, isGET bool) (string, error) {
	var buf strings.Builder

	if layout == nil && file.Head != nil && isGET {
		inSection := false
		headCode, err := genHead(file.Head.Content, "b")
		if err != nil {
			return "", err
		}
		headPending := false
		if len(file.Template.Nodes) > 0 && file.Template.Nodes[0].Type == NodeText &&
			strings.HasPrefix(file.Template.Nodes[0].Content, "<!") {
			headPending = true
		}
		if !headPending {
			buf.WriteString(headCode)
		}
		if isGET {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		}
		for _, n := range file.Template.Nodes {
			code, err := genTemplateNodeToState(gen, n, 1, "b", &inSection)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
			if headPending && n.Type == NodeText && strings.HasPrefix(n.Content, "<!") {
				buf.WriteString(headCode)
				headPending = false
			}
		}
		if isGET {
			buf.WriteString("\tb.WriteString(\"</div>\")\n")
		}
	} else {
		inSection := false
		if isGET {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		}
		for _, n := range file.Template.Nodes {
			code, err := genTemplateNodeToState(gen, n, 1, "b", &inSection)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
		}
		if isGET {
			buf.WriteString("\tb.WriteString(\"</div>\")\n")
		}
	}

	if file.Script != nil {
		buf.WriteString("\tb.WriteString(\"<script>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Script.Code)))
		buf.WriteString("\tb.WriteString(\"</script>\")\n")
	}
	if file.Style != nil {
		scoped := scopeCSS(file.Style.Code, scopeHash)
		buf.WriteString("\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(scoped)))
		buf.WriteString("\tb.WriteString(\"</style>\")\n")
	}

	if layout != nil && isGET {
		buf.WriteString("\tpageContent := b.String()\n")
		buf.WriteString("\tb.Reset()\n")

		if layout.Head != nil && layout.Head.Content != "" {
			headContent := layout.Head.Content
			if strings.Contains(headContent, "{#head}") {
				parts := strings.SplitN(headContent, "{#head}", 2)
				prefix := parts[0]
				if file.Head != nil {
					prefix = dedupeHeadMerge(prefix, file.Head.Content)
				}
				if prefix != "" {
					buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(prefix)))
				}
				if file.Head != nil {
					headCode, err := genHead(file.Head.Content, "b")
					if err != nil {
						return "", err
					}
					buf.WriteString(headCode)
				} else {
					buf.WriteString("\tb.WriteString(\"\")\n")
				}
				if len(parts) > 1 && parts[1] != "" {
					buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(parts[1])))
				}
			} else {
				buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(headContent)))
			}
		}

		if file.Head != nil {
			buf.WriteString("\tvar headBuf strings.Builder\n")
			headCode, err := genHead(file.Head.Content, "headBuf")
			if err != nil {
				return "", err
			}
			buf.WriteString(headCode)
			buf.WriteString("\tc.Set(\"head\", headBuf.String())\n")
		}
		buf.WriteString("\tc.Set(\"slot\", pageContent)\n")
		if layout.Template != nil {
			inSection := false
			for _, n := range layout.Template.Nodes {
				code, err := genLayoutNodeState(gen, n, 1, &inSection)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
		}
	}

	return buf.String(), nil
}
