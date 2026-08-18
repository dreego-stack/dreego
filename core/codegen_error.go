package core

import (
	"fmt"
	"strings"
)

func GenerateErrorHandler(gen *generator, file *File, pkgName string, code int, catchPattern string, scopeHash string) (string, string, error) {

	safeName := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return -1
	}, pkgName)
	funcName := "renderError" + safeName + fmt.Sprintf("%d", code)
	handlerName := "HandleError" + safeName + fmt.Sprintf("%d", code)

	var buf strings.Builder

	buf.WriteString(fmt.Sprintf("func %s(c *dreego.SSRContext) (string, error) {\n", funcName))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if len(file.Go) > 0 && file.Go[0].Code != "" {
		for _, line := range strings.Split(strings.Trim(file.Go[0].Code, "\n"), "\n") {
			buf.WriteString("\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if file.Template != nil {
		suppressScope := false
		if len(file.Template.Nodes) > 0 && file.Template.Nodes[0].Type == NodeText {
			suppressScope = strings.HasPrefix(file.Template.Nodes[0].Content, "<!")
		}
		var headCode string
		if file.Head != nil {
			var err error
			headCode, err = genHead(file.Head.Content, "b")
			if err != nil {
				return "", "", err
			}
		}
		if !suppressScope {
			if headCode != "" {
				buf.WriteString(headCode)
			}
			buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		}
		headPending := suppressScope && headCode != ""
		for _, n := range file.Template.Nodes {
			code, err := genTemplateNode(gen, n, 1)
			if err != nil {
				return "", "", err
			}
			buf.WriteString(code)
			if headPending && n.Type == NodeText && strings.HasPrefix(n.Content, "<!") {
				buf.WriteString(headCode)
				headPending = false
			}
		}
		if !suppressScope {
			buf.WriteString("\tb.WriteString(\"</div>\")\n")
		}

		if file.Script != nil {
			buf.WriteString("\tb.WriteString(\"<script>\")\n")
			buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Script.Code)))
			buf.WriteString("\tb.WriteString(\"</script>\")\n")
		}
		if file.Style != nil {
			styleCode := file.Style.Code
			if !suppressScope {
				styleCode = scopeCSS(file.Style.Code, scopeHash)
			}
			buf.WriteString("\tb.WriteString(\"<style>\")\n")
			buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(styleCode)))
			buf.WriteString("\tb.WriteString(\"</style>\")\n")
		}
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", handlerName))
	buf.WriteString("\tc := dreego.NewSSR(w, r)\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(c)\n", funcName))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	if code == 404 {
		buf.WriteString(fmt.Sprintf("\tw.WriteHeader(%d)\n", code))
	}
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n\n")

	var reg strings.Builder
	if code == 404 {
		reg.WriteString(registrationStatement(fmt.Sprintf("app.Register(%q, %q, %s)", "", catchPattern, handlerName)))
	} else if code == 500 {
		reg.WriteString(registrationStatement(fmt.Sprintf("app.SetErrorHandler(%d, %s)", code, handlerName)))
	}

	return buf.String(), reg.String(), nil
}
