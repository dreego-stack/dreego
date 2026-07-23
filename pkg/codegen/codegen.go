package codegen

import (
	"fmt"
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
)

func GenerateHandler(file *ast.File, baseName string) (string, error) {
	funcName := "render" + toPascalCase(baseName)
	handlerName := "Handle" + toPascalCase(baseName)

	var styleHash string
	var buf strings.Builder

	buf.WriteString("package main\n\n")
	buf.WriteString("import (\n")
	buf.WriteString("\t\"fmt\"\n")
	buf.WriteString("\t\"net/http\"\n")
	buf.WriteString("\t\"strings\"\n\n")
	buf.WriteString("\t\"codeberg.org/dreego/dreego/pkg/context\"\n")
	buf.WriteString(")\n\n")

	if file.Head != nil {
		buf.WriteString(fmt.Sprintf("const head_%s = %s\n\n", baseName, goLiteral(file.Head.Content)))
	}
	if file.Script != nil {
		buf.WriteString(fmt.Sprintf("const script_%s = %s\n\n", baseName, goLiteral(file.Script.Code)))
	}
	if file.Style != nil {
		styleHash = shortHash(file.Style.Code)
		buf.WriteString(fmt.Sprintf("const style_%s = %s\n\n", styleHash, goLiteral(file.Style.Code)))
	}

	buf.WriteString(fmt.Sprintf("func %s(ctx *context.Context) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n")

	if code := file.Go; code != nil && code.Code != "" {
		buf.WriteString("\n")
		for _, line := range strings.Split(strings.Trim(code.Code, "\n"), "\n") {
			buf.WriteString("\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if file.Head != nil {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(head_%s)\n", baseName))
	}

	if file.Template != nil {
		for _, n := range file.Template.Nodes {
			buf.WriteString(genTemplateNode(n, 1))
		}
	}

	if file.Script != nil {
		buf.WriteString("\tb.WriteString(\"<script>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(script_%s)\n", baseName))
		buf.WriteString("\tb.WriteString(\"</script>\")\n")
	}

	if file.Style != nil {
		buf.WriteString("\tb.WriteString(\"<style>\")\n")
		buf.WriteString(fmt.Sprintf("\tb.WriteString(style_%s)\n", styleHash))
		buf.WriteString("\tb.WriteString(\"</style>\")\n")
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", handlerName))
	buf.WriteString("\tctx := &context.Context{W: w, R: r}\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(ctx)\n", funcName))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n")

	if file.Style != nil {
		buf.WriteString(fmt.Sprintf("\nfunc %sCSS() string {\n", handlerName))
		buf.WriteString(fmt.Sprintf("\treturn style_%s\n", styleHash))
		buf.WriteString("}\n")
	}
	if file.Script != nil {
		buf.WriteString(fmt.Sprintf("\nfunc %sJS() string {\n", handlerName))
		buf.WriteString(fmt.Sprintf("\treturn script_%s\n", baseName))
		buf.WriteString("}\n")
	}

	return buf.String(), nil
}

func Generate(file *ast.File, baseName string) (string, error) {
	handler, err := GenerateHandler(file, baseName)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	buf.WriteString(handler)
	buf.WriteString("\nfunc main() {}\n")

	return buf.String(), nil
}

func GenerateMain() string {
	return `package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	http.ListenAndServe(":8080", mux)
}
`
}
