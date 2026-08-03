package core

import (
	"fmt"
	"strings"
)

func generateFormPostHandler(file *File, renderFunc string, postHandler string, pattern string) (string, error) {
	action := file.FormActions[0]
	structName := findFormStruct(file.Go, action)
	if structName == "" {
		return fmt.Sprintf("// no form struct for action %s\n", action), nil
	}
	if !findFormHandler(file.Go, action) {
		return fmt.Sprintf("// no handler function for action %s\n", action), nil
	}
	hasValidate := hasValidateTag(file.Go, structName)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("func %s(w http.ResponseWriter, r *http.Request) {\n", postHandler))
	buf.WriteString("\tc := core.NewSSR(w, r)\n\n")
	buf.WriteString(fmt.Sprintf("\tvar form %s\n", structName))
	buf.WriteString("\tif err := core.BindForm(r, \u0026form); err != nil {\n")
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
	return buf.String(), nil
}
