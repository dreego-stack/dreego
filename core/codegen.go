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
		buf.WriteString(genTypedBlocks(file))
	}

	if file.Template != nil && len(file.Template.Nodes) > 0 {
		buf.WriteString(genTempl(file, layout, scopeHash, true))
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
		postCode = generateFormPostHandler(file, renderFunc, postHandler, pattern)
		if !strings.HasPrefix(postCode, "//") {
			buf.WriteString(postCode)
		}
	}

	buf.WriteString("func init() {\n")
	if hasFormActions {
		buf.WriteString(fmt.Sprintf("\tcore.Register(\"GET\", \"%s\", %s)\n", pattern, getHandler))
		if !strings.HasPrefix(postCode, "//") {
			buf.WriteString(fmt.Sprintf("\tcore.Register(\"POST\", \"%s\", %s)\n", pattern, postHandler))
		}
	} else {
		buf.WriteString(fmt.Sprintf("\tcore.Register(\"%s\", \"%s\", %s)\n", firstMethod, pattern, getHandler))
	}
	buf.WriteString("}\n")

	return buf.String(), nil
}

func genTypedBlocks(file *File) string {
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
	return buf.String()
}

func genTempl(file *File, layout *File, scopeHash string, isGET bool) string {
	var buf strings.Builder

	if layout == nil && file.Head != nil && isGET {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Head.Content)))
	}
	if isGET {
		buf.WriteString(fmt.Sprintf("\tb.WriteString(\"<div data-scope=\\\"%s\\\">\")\n", scopeHash))
	}
	for _, n := range file.Template.Nodes {
		buf.WriteString(genTemplateNode(n, 1))
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
					buf.WriteString(fmt.Sprintf("\tb.WriteString(%s)\n", goLiteral(file.Head.Content)))
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
			buf.WriteString(fmt.Sprintf("\tc.Set(\"head\", %s)\n", goLiteral(file.Head.Content)))
		}
		buf.WriteString("\tc.Set(\"slot\", pageContent)\n")
		if layout.Template != nil {
			for _, n := range layout.Template.Nodes {
				buf.WriteString(genLayoutNode(n, 1))
			}
		}
	}

	return buf.String()
}

func splitGoSections(sections []GoSection, hasFormActions bool) (pkgCode string, inlineCode string) {
	var pkg []string
	var inl []string
	for _, g := range sections {
		if g.ContentType != "" && g.ContentType != "custom" {
			continue
		}
		trimmed := strings.TrimSpace(g.Code)
		if trimmed == "" {
			continue
		}
		firstLine := strings.TrimSpace(strings.SplitN(trimmed, "\n", 2)[0])
		isDeclaration := strings.HasPrefix(firstLine, "type ") || strings.HasPrefix(firstLine, "func ")
		if isDeclaration && hasFormActions {
			pkg = append(pkg, unindent(trimmed))
		} else {
			inl = append(inl, trimmed)
		}
	}
	result := strings.Join(pkg, "\n")
	if result != "" {
		result += "\n"
	}
	return result, strings.Join(inl, "\n")
}

func unindent(code string) string {
	lines := strings.Split(code, "\n")
	minIndent := -1
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}
		indent := len(l) - len(strings.TrimLeft(l, "\t"))
		if minIndent < 0 || indent < minIndent {
			minIndent = indent
		}
	}
	if minIndent <= 0 {
		return code
	}
	var out []string
	for _, l := range lines {
		if len(l) >= minIndent {
			out = append(out, l[minIndent:])
		} else {
			out = append(out, l)
		}
	}
	return strings.Join(out, "\n")
}

func generateFormPostHandler(file *File, renderFunc string, postHandler string, pattern string) string {
	action := file.FormActions[0]
	structName := findFormStruct(file.Go, action)
	if structName == "" {
		return fmt.Sprintf("// no form struct for action %s\n", action)
	}
	if !findFormHandler(file.Go, action) {
		return fmt.Sprintf("// no handler function for action %s\n", action)
	}
	hasValidate := hasValidateTag(file.Go, structName)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", postHandler))
	buf.WriteString("\tc := core.NewSSR(w, r)\n\n")
	buf.WriteString(fmt.Sprintf("\tvar form %s\n", structName))
	buf.WriteString("\tif err := core.BindForm(r, &form); err != nil {\n")
	buf.WriteString(fmt.Sprintf("\t\tc.Set(\"error__form\", err.Error())\n"))
	buf.WriteString(fmt.Sprintf("\t\thtml, _ := %s(c)\n", renderFunc))
	buf.WriteString("\t\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	buf.WriteString("\t\tw.Write([]byte(html))\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n\n")
	if hasValidate {
		buf.WriteString("\terrs := core.ValidateForm(form)\n")
		buf.WriteString("\tif len(errs) > 0 {\n")
		buf.WriteString("\t\tcore.SaveErrors(c, errs)\n")
		buf.WriteString("\t\tcore.SaveOld(c, form)\n")
		buf.WriteString(fmt.Sprintf("\t\thtml, _ := %s(c)\n", renderFunc))
		buf.WriteString("\t\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
		buf.WriteString("\t\tw.Write([]byte(html))\n")
		buf.WriteString("\t\treturn\n")
		buf.WriteString("\t}\n\n")
	}
	buf.WriteString(fmt.Sprintf("\tif err := %s(c, form); err != nil {\n", action))
	buf.WriteString("\t\tif err == core.ErrRedirect {\n")
	buf.WriteString("\t\t\treturn\n")
	buf.WriteString("\t\t}\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n\n")
	buf.WriteString("\thttp.Redirect(w, r, r.URL.Path, 303)\n")
	buf.WriteString("}\n\n")
	return buf.String()
}

func genLayoutNode(n TemplateNode, depth int) string {
	indent := strings.Repeat("\t", depth)
	if n.Type == NodeSlot {
		if n.Content != "" {
			return indent + fmt.Sprintf("b.WriteString(c.Get(\"slot_%s\"))\n", n.Content)
		}
		return indent + "b.WriteString(c.Get(\"slot\"))\n"
	}
	if n.Type == NodeText && (strings.Contains(n.Content, "{#head}") || strings.Contains(n.Content, "{#slot}")) {
		parts := splitLayoutText(n.Content)
		var out string
		for _, p := range parts {
			switch p {
			case "{#head}":
				out += indent + "b.WriteString(c.Get(\"head\"))\n"
			case "{#slot}":
				out += indent + "b.WriteString(c.Get(\"slot\"))\n"
			default:
				out += indent + fmt.Sprintf("b.WriteString(%s)\n", goLiteral(p))
			}
		}
		return out
	}
	return genTemplateNode(n, depth)
}

func splitLayoutText(s string) []string {
	var result []string
	for s != "" {
		headIdx := strings.Index(s, "{#head}")
		slotIdx := strings.Index(s, "{#slot}")

		next := -1
		nextLen := 0
		if headIdx >= 0 && (slotIdx < 0 || headIdx <= slotIdx) {
			next = headIdx
			nextLen = 7
		} else if slotIdx >= 0 {
			next = slotIdx
			nextLen = 7
		}

		if next < 0 {
			if s != "" {
				result = append(result, s)
			}
			break
		}
		if next > 0 {
			result = append(result, s[:next])
		}
		result = append(result, s[next:next+nextLen])
		s = s[next+nextLen:]
	}
	return result
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
		code := fmt.Sprintf("fmt.Sprintf(\"%%v\", %s)", n.Content)
		raw := false
		for _, f := range n.Filters {
			switch f {
			case "raw":
				raw = true
			case "upper":
				code = fmt.Sprintf("strings.ToUpper(%s)", code)
			}
		}
		if raw {
			return fmt.Sprintf("b.WriteString(%s)", code)
		}
		return fmt.Sprintf("b.WriteString(html.EscapeString(%s))", code)
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
