package transpiler

import (
	"fmt"
	"strings"
)

func genTempl(gen *Generator, file *File, layout *layoutEntry, scopeHash string, isGET bool) (string, error) {
	var buf strings.Builder

	if layout == nil && file.Head != nil && isGET {
		inSection := false
		headCode, err := genHead(file.Head.Content, "b")
		if err != nil {
			return "", err
		}
		headPending := false
		if len(file.Body.Nodes) > 0 && file.Body.Nodes[0].Type == NodeText &&
			strings.HasPrefix(file.Body.Nodes[0].Content, "<!") {
			headPending = true
		}
		if !headPending {
			buf.WriteString(headCode)
		}
		if isGET {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		}
		for _, n := range file.Body.Nodes {
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
		for _, n := range file.Body.Nodes {
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

	if file.Client != nil {
		buf.WriteString("\tb.WriteString(\"<script>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Client.Code)))
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

		if file.Head != nil {
			headCode, err := genHead(file.Head.Content, "b")
			if err != nil {
				return "", err
			}
			buf.WriteString(headCode)
		}
		buf.WriteString("\tpageHead := b.String()\n")
		buf.WriteString("\tb.Reset()\n")

		layoutHead := ""
		if layout.file.Head != nil {
			layoutHead = layout.file.Head.Content
		}
		headPrefix, headSuffix := splitHeadPlaceholder(layoutHead)
		if headPrefix != "" {
			buf.WriteString(fmt.Sprintf("\tlayoutHead := %s\n", goLiteral(headPrefix)))
			buf.WriteString("\tif strings.Contains(pageHead, \"<title\") {\n")
			buf.WriteString("\t\tlayoutHead = stripTitleTag(layoutHead)\n")
			buf.WriteString("\t}\n")
			buf.WriteString("\tif strings.Contains(pageHead, `name=\"description\"`) || strings.Contains(pageHead, `name='description'`) {\n")
			buf.WriteString("\t\tlayoutHead = stripMetaDescriptionTag(layoutHead)\n")
			buf.WriteString("\t}\n")
			buf.WriteString("\tb.WriteString(layoutHead)\n")
		}
		buf.WriteString("\tb.WriteString(pageHead)\n")
		if headSuffix != "" {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(headSuffix)))
		}

		buf.WriteString("\thead := b.String()\n")
		buf.WriteString("\tb.Reset()\n")
		buf.WriteString("\tc.Set(\"slot\", pageContent)\n")

		layoutPkg := "layouts"
		layoutPath := gen.module + "/" + gen.rootRel + "/layouts"
		gen.addImport(gen.pkg, layoutPkg, layoutPath)
		buf.WriteString(fmt.Sprintf("\thtml, err := %s.%s(c, pageContent, head)\n", layoutPkg, layout.name))
		buf.WriteString("\tif err != nil { return \"\", err }\n")
		buf.WriteString("\tb.WriteString(html)\n")
	}

	return buf.String(), nil
}

func headMergeHelpers() string {
	return `
func stripTitleTag(s string) string {
	for {
		open := strings.Index(s, "<title")
		if open < 0 {
			return s
		}
		closeIdx := strings.Index(s[open:], "</title>")
		if closeIdx < 0 {
			return s
		}
		end := open + closeIdx + len("</title>")
		s = s[:open] + s[end:]
	}
}

func stripMetaDescriptionTag(s string) string {
	offset := 0
	for {
		open := strings.Index(s[offset:], "<meta")
		if open < 0 {
			return s
		}
		open += offset
		end := strings.IndexByte(s[open:], '>')
		if end < 0 {
			return s
		}
		tag := s[open : open+end+1]
		if strings.Contains(tag, "name=\"description\"") || strings.Contains(tag, "name='description'") {
			s = s[:open] + s[open+end+1:]
			offset = open
			continue
		}
		offset = open + end + 1
	}
}
`
}
