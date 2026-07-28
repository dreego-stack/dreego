package core

import (
	"fmt"
	"strings"

)

func GenerateMethodHandler(file *File, layout *File, pkgName string, baseName string, pattern string, g GoSection, scopeHash string) (string, error) {
	methodSuffix := ""
	if g.Method != "GET" {
		methodSuffix = strings.ToUpper(g.Method)
	}
	funcName := "render" + toPascalCase(baseName) + methodSuffix
	handlerName := "Handle" + toPascalCase(baseName) + methodSuffix

	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("func %s(c *core.SSRContext) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if g.Code != "" {
		for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
			buf.WriteString("\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if file.Template != nil {
		if layout == nil && file.Head != nil && g.Method == "GET" {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Head.Content)))
		}
		if g.Method == "GET" {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		}
		for _, n := range file.Template.Nodes {
			buf.WriteString(genTemplateNode(n, 1))
		}
		if g.Method == "GET" {
			buf.WriteString("\tb.WriteString(\"</div>\")\n")
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

		if layout != nil && g.Method == "GET" {
			buf.WriteString("\tpageContent := b.String()\n")
			buf.WriteString("\tb.Reset()\n")
			if file.Head != nil {
				buf.WriteString(fmt.Sprintf("\tc.Set(\"head\", %s)\n", goLiteral(file.Head.Content)))
			}
			buf.WriteString("\tc.Set(\"slot\", pageContent)\n")
			if layout.Template != nil {
				for _, n := range layout.Template.Nodes {
					buf.WriteString(genLayoutNode(n, 1))
				}
			}
		}
	} else if g.Code != "" {
		buf.WriteString("\tb.WriteString(\"OK\")\n")
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", handlerName))
	buf.WriteString("\tc := core.NewSSR(w, r)\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(c)\n", funcName))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func init() {\n")
	buf.WriteString(fmt.Sprintf("\tcore.Register(\"%s\", \"%s\", %s)\n", g.Method, pattern, handlerName))
	buf.WriteString("}\n")

	return buf.String(), nil
}

func genLayoutNode(n TemplateNode, depth int) string {
	if n.Type == NodeSlot {
		return strings.Repeat("\t", depth) + "b.WriteString(c.Get(\"slot\"))\n"
	}
	if n.Type == NodeText {
		if strings.Contains(n.Content, "{#head}") {
			parts := strings.SplitN(n.Content, "{#head}", 2)
			var out string
			if parts[0] != "" {
				out += strings.Repeat("\t", depth) + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(parts[0]))
			}
			out += strings.Repeat("\t", depth) + "b.WriteString(c.Get(\"head\"))\n"
			if len(parts) > 1 && parts[1] != "" {
				out += strings.Repeat("\t", depth) + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(parts[1]))
			}
			return out
		}
		if strings.Contains(n.Content, "{#slot}") {
			parts := strings.SplitN(n.Content, "{#slot}", 2)
			var out string
			if parts[0] != "" {
				out += strings.Repeat("\t", depth) + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(parts[0]))
			}
			out += strings.Repeat("\t", depth) + "b.WriteString(c.Get(\"slot\"))\n"
			if len(parts) > 1 && parts[1] != "" {
				out += strings.Repeat("\t", depth) + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(parts[1]))
			}
			return out
		}
	}
	return genTemplateNode(n, depth)
}

func GenerateErrorHandler(file *File, pkgName string, code int, catchPattern string, scopeHash string) (string, error) {
	safeName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, pkgName)
	funcName := "renderError" + safeName + fmt.Sprintf("%d", code)
	handlerName := "HandleError" + safeName + fmt.Sprintf("%d", code)

	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("func %s(c *core.SSRContext) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if len(file.Go) > 0 && file.Go[0].Code != "" {
		for _, line := range strings.Split(strings.Trim(file.Go[0].Code, "\n"), "\n") {
			buf.WriteString("\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if file.Template != nil {
		if file.Head != nil {
			buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Head.Content)))
		}
		buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		for _, n := range file.Template.Nodes {
			buf.WriteString(genTemplateNode(n, 1))
		}
		buf.WriteString("\tb.WriteString(\"</div>\")\n")

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
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", handlerName))
	buf.WriteString("\tc := core.NewSSR(w, r)\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(c)\n", funcName))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	if code == 404 {
		buf.WriteString(fmt.Sprintf("\tw.WriteHeader(%d)\n", code))
	}
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func init() {\n")
	if code == 404 {
		buf.WriteString(fmt.Sprintf("\tcore.Register(\"\", \"%s\", %s)\n", catchPattern, handlerName))
	} else if code == 500 {
		buf.WriteString(fmt.Sprintf("\tcore.SetErrorHandler(%d, %s)\n", code, handlerName))
	}
	buf.WriteString("}\n")

	return buf.String(), nil
}

func GenerateComponent(file *File, scopeHash string) (string, error) {
	comp := file.Component
	if comp == nil {
		return "", fmt.Errorf("no component definition")
	}

	var buf strings.Builder

	params := ""
	for i, p := range comp.Props {
		if i > 0 {
			params += ", "
		}
		params += p.Name + " " + p.Type
	}

	buf.WriteString(fmt.Sprintf("func %s(%s) core.Component {\n", comp.Name, params))
	buf.WriteString("\treturn core.ComponentFunc(func(ctx *core.SSRContext) (string, error) {\n")
	buf.WriteString("\t\tvar b strings.Builder\n\n")

	if len(file.Go) > 0 && file.Go[0].Code != "" {
		for _, line := range strings.Split(strings.Trim(file.Go[0].Code, "\n"), "\n") {
			buf.WriteString("\t\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if file.Template != nil {
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		for _, n := range file.Template.Nodes {
			buf.WriteString("\t\t" + genTemplateNodeComp(n) + "\n")
		}
		buf.WriteString("\t\tb.WriteString(\"</div>\")\n")
	}

	if file.Style != nil {
		scoped := scopeCSS(file.Style.Code, scopeHash)
		buf.WriteString("\t\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(%s)\n", goLiteral(scoped)))
		buf.WriteString("\t\tb.WriteString(\"</style>\")\n")
	}

	buf.WriteString("\n\t\treturn b.String(), nil\n")
	buf.WriteString("\t})\n")
	buf.WriteString("}\n\n")

	return buf.String(), nil
}

func genTemplateNodeComp(n TemplateNode) string {
	switch n.Type {
	case NodeText:
		return fmt.Sprintf("b.WriteString(%s)", goLiteral(n.Content))
	case NodeExpression:
		return fmt.Sprintf("b.WriteString(html.EscapeString(fmt.Sprintf(\"%%v\", %s)))", n.Content)
	case NodeSlot:
		if n.Content != "" {
			return fmt.Sprintf("b.WriteString(ctx.Get(\"slot_%s\"))", n.Content)
		}
		return "b.WriteString(ctx.Get(\"slot\"))"
	case NodeComponentCall:
		return genComponentCall(n)
	case NodeVerbatim:
		return fmt.Sprintf("b.WriteString(%s)", goLiteral(n.Content))
	default:
		return ""
	}
}

func genComponentCall(n TemplateNode) string {
	parts := strings.SplitN(n.Tag, ".", 2)
	funcName := parts[len(parts)-1]
	if n.SelfClose {
		return fmt.Sprintf("%s(%s).Render(ctx)", funcName, n.Attrs)
	}
	return fmt.Sprintf("b.WriteString(\"<@%s>\")", n.Tag)
}
