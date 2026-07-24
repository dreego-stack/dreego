package codegen

import (
	"fmt"
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
)

func GenerateHandler(file *ast.File, layout *ast.File, pkgName string, baseName string, pattern string, scopeHash string) (string, error) {
	funcName := "render" + toPascalCase(baseName)
	handlerName := "Handle" + toPascalCase(baseName)

	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("package %s\n\n", pkgName))
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n")
	buf.WriteString("\t\"net/http\"\n")
	buf.WriteString("\t\"strings\"\n\n")
	buf.WriteString("\t\"codeberg.org/dreego/dreego/pkg/context\"\n")
	buf.WriteString("\t\"codeberg.org/dreego/dreego/pkg/runtime\"\n")
	buf.WriteString(")\n\n")

	if file.Head != nil {
		buf.WriteString(fmt.Sprintf("const head_%s = %s\n\n", baseName, goLiteral(file.Head.Content)))
	}
	if file.Script != nil {
		buf.WriteString(fmt.Sprintf("const script_%s = %s\n\n", baseName, goLiteral(file.Script.Code)))
	}
	if file.Style != nil {
		scoped := scopeCSS(file.Style.Code, scopeHash)
		buf.WriteString(fmt.Sprintf("const style_%s = %s\n\n", baseName, goLiteral(scoped)))
	}

	buf.WriteString(fmt.Sprintf("func renderPage%s(c *context.Context) (string, error) {\n", toPascalCase(baseName)))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if code := file.Go; code != nil && code.Code != "" {
		for _, line := range strings.Split(strings.Trim(code.Code, "\n"), "\n") {
			buf.WriteString("\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if layout == nil && file.Head != nil {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(head_%s)\n", baseName))
	}

	if file.Template != nil {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		for _, n := range file.Template.Nodes {
			buf.WriteString(genTemplateNode(n, 1))
		}
		buf.WriteString("\tb.WriteString(\"</div>\")\n")
	}

	if file.Script != nil {
		buf.WriteString("\tb.WriteString(\"<script>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(script_%s)\n", baseName))
		buf.WriteString("\tb.WriteString(\"</script>\")\n")
	}

	if file.Style != nil {
		buf.WriteString("\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(style_%s)\n", baseName))
		buf.WriteString("\tb.WriteString(\"</style>\")\n")
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")

	if layout != nil {
		buf.WriteString(fmt.Sprintf("func %s(ctx *context.Context) (string, error) {\n", funcName))
		buf.WriteString(fmt.Sprintf("\tpageContent, err := renderPage%s(ctx)\n", toPascalCase(baseName)))
		buf.WriteString("\tif err != nil { return \"\", err }\n")
		if file.Head != nil {
			buf.WriteString(fmt.Sprintf("\tctx.Set(\"head\", head_%s)\n", baseName))
		}
		buf.WriteString("\tctx.Set(\"slot\", pageContent)\n")
		buf.WriteString("\tvar b strings.Builder\n")

		if layout.Template != nil {
			for _, n := range layout.Template.Nodes {
				buf.WriteString(genLayoutNode(n, 1))
			}
		}

		buf.WriteString("\n\treturn b.String(), nil\n")
		buf.WriteString("}\n\n")
	} else {
		buf.WriteString(fmt.Sprintf("func %s(ctx *context.Context) (string, error) {\n", funcName))
		buf.WriteString(fmt.Sprintf("\treturn renderPage%s(ctx)\n", toPascalCase(baseName)))
		buf.WriteString("}\n\n")
	}

	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", handlerName))
	buf.WriteString("\tctx := &context.Context{W: w, R: r}\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(ctx)\n", funcName))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n\n")

	buf.WriteString("func init() {\n")
	buf.WriteString(fmt.Sprintf("\truntime.Register(\"GET\", \"%s\", %s)\n", pattern, handlerName))
	buf.WriteString("}\n")

	return buf.String(), nil
}

func genLayoutNode(n ast.TemplateNode, depth int) string {
	if n.Type == ast.NodeSlot {
		return strings.Repeat("\t", depth) + "b.WriteString(ctx.Get(\"slot\"))\n"
	}
	if n.Type == ast.NodeText {
		if strings.Contains(n.Content, "{#head}") {
			parts := strings.SplitN(n.Content, "{#head}", 2)
			var out string
			if parts[0] != "" {
				out += strings.Repeat("\t", depth) + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(parts[0]))
			}
			out += strings.Repeat("\t", depth) + "b.WriteString(ctx.Get(\"head\"))\n"
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
			out += strings.Repeat("\t", depth) + "b.WriteString(ctx.Get(\"slot\"))\n"
			if len(parts) > 1 && parts[1] != "" {
				out += strings.Repeat("\t", depth) + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(parts[1]))
			}
			return out
		}
	}
	return genTemplateNode(n, depth)
}
