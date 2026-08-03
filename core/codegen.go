package core

import (
	"fmt"
	"strings"
)

func GenerateMethodHandler(file *File, layout *File, pkgName string, baseName string, pattern string, scopeHash string) (string, error) {
	hasTypedBlocks := false
	for _, g := range file.Go {
		if g.ContentType != "" && g.ContentType != "custom" {
			hasTypedBlocks = true
		}
	}
	hasFormActions := len(file.FormActions) > 0

	pascalBase := toPascalCase(baseName)
	renderFunc := "render" + pascalBase

	var firstMethod string
	var getHandler string
	var postHandler string

	if hasFormActions {
		getHandler = "Handle" + pascalBase + "Get"
		postHandler = "Handle" + pascalBase + "Post"
		firstMethod = ""
	} else {
		firstMethod = "GET"
		for _, g := range file.Go {
			if g.Method != "GET" {
				firstMethod = g.Method
			}
		}
		methodSuffix := ""
		if firstMethod != "GET" {
			methodSuffix = strings.ToUpper(firstMethod)
		}
		renderFunc = "render" + pascalBase + methodSuffix
		getHandler = "Handle" + pascalBase + methodSuffix
	}

	var buf strings.Builder

	pkgCode, inlineCode := splitGoSections(file.Go, hasFormActions)
	if pkgCode != "" {
		buf.WriteString(pkgCode)
	}

	buf.WriteString(fmt.Sprintf("\nfunc %s(c *core.SSRContext) (string, error) {\n", renderFunc))
	buf.WriteString("\tvar b strings.Builder\n\n")

	if inlineCode != "" {
		for _, line := range strings.Split(strings.Trim(inlineCode, "\n"), "\n") {
			buf.WriteString("\t" + strings.TrimSpace(line) + "\n")
		}
		buf.WriteString("\n")
	}

	if hasTypedBlocks {
		typedCode, err := genTypedBlocks(file)
		if err != nil {
			return "", err
		}
		buf.WriteString(typedCode)
	}

	if file.Template != nil && len(file.Template.Nodes) > 0 {
		templCode, err := genTempl(file, layout, scopeHash, true)
		if err != nil {
			return "", err
		}
		buf.WriteString(templCode)
	} else if !hasFormActions && firstMethod != "GET" {
		buf.WriteString("\tb.WriteString(\"OK\")\n")
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", getHandler))
	buf.WriteString("\tc := core.NewSSR(w, r)\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(c)\n", renderFunc))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n\n")

	var postCode string
	if hasFormActions {
		var err error
		postCode, err = generateFormPostHandler(file, renderFunc, postHandler, pattern)
		if err != nil {
			return "", err
		}
		if !strings.HasPrefix(postCode, "//") {
			buf.WriteString(postCode)
		}
	}

	buf.WriteString("func init() {\n")
	if hasFormActions {
		buf.WriteString(fmt.Sprintf("\tcore.Register(\"GET\", \"%s\", %s)\n", pattern, getHandler))
		if postCode != "" && !strings.HasPrefix(postCode, "//") {
			buf.WriteString(fmt.Sprintf("\tcore.Register(\"POST\", \"%s\", %s)\n", pattern, postHandler))
		}
	} else {
		buf.WriteString(fmt.Sprintf("\tcore.Register(\"%s\", \"%s\", %s)\n", firstMethod, pattern, getHandler))
	}
	buf.WriteString("}\n")

	return buf.String(), nil
}

func genTypedBlocks(file *File) (string, error) {
	var buf strings.Builder
	buf.WriteString("\tif true {\n")
	for _, g := range file.Go {
		if g.ContentType == "json" {
			buf.WriteString("\t\tif c.Wants(\"application/json\") {\n")
			buf.WriteString("\t\t\tc.W.Header().Set(\"Content-Type\", \"application/json; charset=utf-8\")\n")
			for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
				buf.WriteString("\t\t\t" + strings.TrimSpace(line) + "\n")
			}
			buf.WriteString("\t\t\treturn \"\", nil\n")
			buf.WriteString("\t\t}\n")
		}
		if g.ContentType == "xml" {
			buf.WriteString("\t\tif c.Wants(\"application/xml\") {\n")
			buf.WriteString("\t\t\tc.W.Header().Set(\"Content-Type\", \"application/xml; charset=utf-8\")\n")
			for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
				buf.WriteString("\t\t\t" + strings.TrimSpace(line) + "\n")
			}
			buf.WriteString("\t\t\treturn \"\", nil\n")
			buf.WriteString("\t\t}\n")
		}
	}
	buf.WriteString("\t}\n\n")
	return buf.String(), nil
}

func genTempl(file *File, layout *File, scopeHash string, isGET bool) (string, error) {
	var buf strings.Builder

	if layout == nil && file.Head != nil && isGET {
		headCode, err := genHead(file.Head.Content, "b")
		if err != nil {
			return "", err
		}
		buf.WriteString(headCode)
	}
	if isGET {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
	}
	for _, n := range file.Template.Nodes {
		code, err := genTemplateNode(n, 1)
		if err != nil {
			return "", err
		}
		buf.WriteString(code)
	}
	if isGET {
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

	if layout != nil && isGET {
		buf.WriteString("\tpageContent := b.String()\n")
		buf.WriteString("\tb.Reset()\n")

		if layout.Head != nil && layout.Head.Content != "" {
			headContent := layout.Head.Content
			if strings.Contains(headContent, "{#head}") {
				parts := strings.SplitN(headContent, "{#head}", 2)
				if parts[0] != "" {
					buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(parts[0])))
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
			for _, n := range layout.Template.Nodes {
				code, err := genLayoutNode(n, 1)
				if err != nil {
					return "", err
				}
				buf.WriteString(code)
			}
		}
	}

	return buf.String(), nil
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
			headCode, err := genHead(file.Head.Content, "b")
			if err != nil {
				return "", err
			}
			buf.WriteString(headCode)
		}
		buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		for _, n := range file.Template.Nodes {
			code, err := genTemplateNode(n, 1)
			if err != nil {
				return "", err
			}
			buf.WriteString(code)
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

	for _, g := range file.Go {
		if g.Code != "" {
			for _, line := range strings.Split(strings.Trim(g.Code, "\n"), "\n") {
				buf.WriteString("\t\t" + strings.TrimSpace(line) + "\n")
			}
			buf.WriteString("\n")
		}
	}

	if file.Template != nil {
		buf.WriteString(fmt.Sprintf("\t\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
		for _, n := range file.Template.Nodes {
			code, err := genTemplateNodeComp(n)
			if err != nil {
				return "", err
			}
			buf.WriteString("\t\t" + code + "\n")
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
