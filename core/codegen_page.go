package core

import (
	"fmt"
	"strings"
)

func GenerateMethodHandler(file *File, layout *File, pkgName string, baseName string, pattern string, scopeHash string) (string, string, error) {
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

	buf.WriteString(fmt.Sprintf("\nfunc %s(c *dreego.SSRContext) (string, error) {\n", renderFunc))
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
			return "", "", err
		}
		buf.WriteString(typedCode)
	}

	if file.Template != nil && len(file.Template.Nodes) > 0 {
		templCode, err := genTempl(file, layout, scopeHash, true)
		if err != nil {
			return "", "", err
		}
		buf.WriteString(templCode)
	} else if !hasFormActions && firstMethod != "GET" {
		buf.WriteString("\tb.WriteString(\"OK\")\n")
	}

	buf.WriteString("\n\treturn b.String(), nil\n")
	buf.WriteString("}\n\n")

	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", getHandler))
	buf.WriteString("\tc := dreego.NewSSR(w, r)\n")
	buf.WriteString(fmt.Sprintf("\thtml, err := %s(c)\n", renderFunc))
	buf.WriteString("\tif err != nil {\n")
	buf.WriteString("\t\thttp.Error(w, err.Error(), http.StatusInternalServerError)\n")
	buf.WriteString("\t\treturn\n")
	buf.WriteString("\t}\n")
	buf.WriteString("\tw.Header().Set(\"Content-Type\", \"text/html; charset=utf-8\")\n")
	buf.WriteString("\tw.Write([]byte(html))\n")
	buf.WriteString("}\n")

	var postCode string
	if hasFormActions {
		var err error
		postCode, err = generateFormPostHandler(file, renderFunc, postHandler, pattern)
		if err != nil {
			return "", "", err
		}
		if !strings.HasPrefix(postCode, "//") {
			buf.WriteString(postCode)
		}
	}

	var reg strings.Builder
	if hasFormActions {
		reg.WriteString(fmt.Sprintf("\tapp.Register(\"GET\", \"%s\", %s)\n", pattern, getHandler))
		if postCode != "" && !strings.HasPrefix(postCode, "//") {
			reg.WriteString(fmt.Sprintf("\tapp.Register(\"POST\", \"%s\", %s)\n", pattern, postHandler))
		}
	} else {
		reg.WriteString(fmt.Sprintf("\tapp.Register(\"%s\", \"%s\", %s)\n", firstMethod, pattern, getHandler))
	}

	return buf.String(), reg.String(), nil
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
