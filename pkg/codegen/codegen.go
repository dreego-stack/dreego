package codegen

import (
	"fmt"
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
)

func GenerateMethodHandler(file *ast.File, layout *ast.File, pkgName string, baseName string, pattern string, g ast.GoSection, scopeHash string) (string, error) {
	methodSuffix := ""
	if g.Method != "GET" {
		methodSuffix = strings.ToUpper(g.Method)
	}
	funcName := "render" + toPascalCase(baseName) + methodSuffix
	handlerName := "Handle" + toPascalCase(baseName) + methodSuffix

	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("func %s(c *context.Context) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if g.Code != "" {
		for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
			buf.WriteString("\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if file.Template != nil {
		if layout == nil && file.Head != nil {
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
	buf.WriteString("\tc := &context.Context{W: w, R: r}\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(c)\n", funcName))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func init() {\n")
	buf.WriteString(fmt.Sprintf("\truntime.Register(\"%s\", \"%s\", %s)\n", g.Method, pattern, handlerName))
	buf.WriteString("}\n")

	return buf.String(), nil
}

func genLayoutNode(n ast.TemplateNode, depth int) string {
	if n.Type == ast.NodeSlot {
		return strings.Repeat("\t", depth) + "b.WriteString(c.Get(\"slot\"))\n"
	}
	if n.Type == ast.NodeText {
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
