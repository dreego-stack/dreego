package html

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/codegen"
	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
	"github.com/dreego-stack/dreego/internal/dreefile/ir"
	jsoutput "github.com/dreego-stack/dreego/internal/dreefile/jsoutput"
	clientpkg "github.com/dreego-stack/dreego/internal/dreefile/sections/client"
	"github.com/dreego-stack/dreego/internal/dreefile/sections/head"
	"github.com/dreego-stack/dreego/internal/dreefile/sections/style"
)

func GenTempl(gen *codegen.State, file *ir.File, layout *codegen.Layout, scopeHash string, isGET bool) (string, error) {
	var buf strings.Builder

	if layout == nil && file.Head != nil && isGET {
		inSection := false
		headCode, err := head.GenWithMessages(gen, file.Head.Content, "b", "c")
		if err != nil {
			return "", err
		}
		headPending := false
		if len(file.Body.Nodes) > 0 && file.Body.Nodes[0].Type == ir.NodeText &&
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
			code, err := genTemplateNodeToState(gen, n, 1, "b", &inSection, file.Server)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
			if headPending && n.Type == ir.NodeText && strings.HasPrefix(n.Content, "<!") {
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
			code, err := genTemplateNodeToState(gen, n, 1, "b", &inSection, file.Server)
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
		clientArtifact, err := clientpkg.Client(gen, file)
		if err != nil {
			return "", err
		}
		buf.WriteString(jsoutput.GenClient(clientArtifact))
	}
	if file.Style != nil {
		scoped := style.ScopeCSS(file.Style.Code, scopeHash)
		buf.WriteString("\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", gogen.GoLiteral(scoped)))
		buf.WriteString("\tb.WriteString(\"</style>\")\n")
	}

	if layout != nil && isGET {
		buf.WriteString("\tpageContent := b.String()\n")
		buf.WriteString("\tb.Reset()\n")

		if file.Head != nil {
			headCode, err := head.GenWithMessages(gen, file.Head.Content, "b", "c")
			if err != nil {
				return "", err
			}
			buf.WriteString(headCode)
		}
		buf.WriteString("\tpageHead := b.String()\n")
		buf.WriteString("\tb.Reset()\n")

		layoutHead := ""
		if layout.File.Head != nil {
			layoutHead = layout.File.Head.Content
		}
		headPrefix, headSuffix := SplitHeadPlaceholder(layoutHead)
		if headPrefix != "" {
			buf.WriteString(fmt.Sprintf("\tlayoutHead := %s\n", gogen.GoLiteral(headPrefix)))
			buf.WriteString("\tlayoutHead = dedupeLayoutHead(layoutHead, pageHead)\n")
			buf.WriteString("\tb.WriteString(layoutHead)\n")
		}
		buf.WriteString("\tb.WriteString(pageHead)\n")
		if headSuffix != "" {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", gogen.GoLiteral(headSuffix)))
		}

		buf.WriteString("\thead := b.String()\n")
		buf.WriteString("\tb.Reset()\n")
		buf.WriteString("\tc.Set(\"slot\", pageContent)\n")

		layoutPkg := "layouts"
		layoutPath := gen.Module + "/" + gen.RootRel + "/layouts"
		gen.AddImportForCurrent(layoutPkg, layoutPath)
		buf.WriteString(fmt.Sprintf("\thtml, err := %s.%s(c, pageContent, head)\n", layoutPkg, layout.Name))
		buf.WriteString("\tif err != nil { return \"\", err }\n")
		buf.WriteString("\tb.WriteString(html)\n")
	}

	return buf.String(), nil
}

func SplitHeadPlaceholder(head string) (prefix, suffix string) {
	if !strings.Contains(head, ir.HeadPlaceholder) {
		return head, ""
	}
	parts := strings.SplitN(head, ir.HeadPlaceholder, 2)
	return parts[0], parts[1]
}

func HeadMergeHelpers() string {
	return `
func dedupeLayoutHead(layoutHead, pageHead string) string {
	lower := strings.ToLower(pageHead)
	if strings.Contains(lower, "<title") {
		layoutHead = stripTitleTag(layoutHead)
	}
	if strings.Contains(lower, ` + "`" + `name="description"` + "`" + `) || strings.Contains(lower, ` + "`" + `name='description'` + "`" + `) {
		layoutHead = stripMetaTag(layoutHead, "name=\"description\"", "name='description'")
	}
	if strings.Contains(lower, ` + "`" + `name="viewport"` + "`" + `) || strings.Contains(lower, ` + "`" + `name='viewport'` + "`" + `) {
		layoutHead = stripMetaTag(layoutHead, "name=\"viewport\"", "name='viewport'")
	}
	if strings.Contains(lower, "charset") {
		layoutHead = stripMetaTag(layoutHead, "charset")
	}
	return layoutHead
}

func stripTitleTag(s string) string {
	for {
		open := strings.Index(strings.ToLower(s), "<title")
		if open < 0 {
			return s
		}
		closeIdx := strings.Index(strings.ToLower(s[open:]), "</title>")
		if closeIdx < 0 {
			return s
		}
		end := open + closeIdx + len("</title>")
		s = s[:open] + s[end:]
	}
}

func stripMetaDescriptionTag(s string) string {
	return stripMetaTag(s, "name=\"description\"", "name='description'")
}

func stripMetaTag(s string, markers ...string) string {
	offset := 0
	for {
		open := strings.Index(strings.ToLower(s[offset:]), "<meta")
		if open < 0 {
			return s
		}
		open += offset
		end := strings.IndexByte(s[open:], '>')
		if end < 0 {
			return s
		}
		tag := strings.ToLower(s[open : open+end+1])
		matched := false
		for _, marker := range markers {
			if strings.Contains(tag, marker) {
				matched = true
				break
			}
		}
		if matched {
			s = s[:open] + s[open+end+1:]
			offset = open
			continue
		}
		offset = open + end + 1
	}
}
`
}
