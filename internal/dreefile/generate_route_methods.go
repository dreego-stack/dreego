package dreefile

import (
	"fmt"
	"os"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/gogen"
)

func fileRegisteredMethods(file *File) []string {
	if len(file.FormActions) > 0 {
		action := file.FormActions[0]
		if findFormStruct(file.Server, action) != "" && findFormHandler(file.Server, action) {
			return []string{"GET", "POST"}
		}
		return []string{"GET"}
	}
	seen := map[string]bool{}
	methods := []string{}
	add := func(method string) {
		if method == "" {
			method = "GET"
		}
		if !seen[method] {
			seen[method] = true
			methods = append(methods, method)
		}
	}
	for _, section := range file.Server {
		add(section.Method)
	}
	for _, template := range file.Bodies {
		add(template.Method)
	}
	if len(methods) == 0 {
		add("GET")
	}
	return methods
}

func parseRouteFile(gen *Generator, fpath string, data []byte) (*File, string, error) {
	raw := string(data)
	header, body, err := ParseFileHeaderStrict(raw)
	if err != nil {
		return nil, "", fmt.Errorf("%s:%w", fpath, err)
	}
	tokens, err := Lex(body)
	if err != nil {
		return nil, "", fmt.Errorf("error lexing %s: %w", fpath, err)
	}
	p := NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return nil, "", fmt.Errorf("error parsing %s: %w", fpath, err)
	}
	file.Imports = header.Imports
	file.Kind = header.Kind
	file.Layout = header.Layout
	file.GoImports = header.GoImports
	file.SourceContent = raw
	file.SourcePath = fpath
	bodyOffset := len(raw) - len(body)
	if file.Client != nil {
		file.Client.Pos += bodyOffset
	}
	if file.Body != nil {
		gogen.SetNodeSource(file.Body.Nodes, fpath, bodyOffset)
		gogen.SetSourceText(file.Body.Nodes, raw)
		file.FormActions = scanFormActions(file.Body.Nodes)
		for _, diagnostic := range a11yDiagnostics(file.Body.Nodes) {
			fmt.Fprintf(os.Stderr, "warning: %s\n", diagnostic)
		}
		if diagnostic, ok := bodyAttrDiagnostic(file, fpath, bodyOffset); ok {
			fmt.Fprintf(os.Stderr, "warning: %s\n", diagnostic)
		}
		if diagnostic, ok := alpineCSPDiagnostic(file.Body.Nodes, fpath); ok {
			fmt.Fprintf(os.Stderr, "warning: %s\n", diagnostic)
		}
	}
	return file, raw, nil
}

// bodyAttrDiagnostic warns when the <body> section tag carries HTML attributes.
// Those attributes sit on the route's body wrapper, but a layout supplies the
// real document <body>, so the attributes never reach the rendered element.
// The framework would otherwise drop them silently, leaving Alpine/HTMX hooks
// inert.
func bodyAttrDiagnostic(file *File, fpath string, bodyOffset int) (string, bool) {
	if file.Body == nil || strings.TrimSpace(file.Body.Attrs) == "" {
		return "", false
	}
	line, col := gogen.PosToLineCol(file.SourceContent, file.Body.Pos+bodyOffset)
	d := Diagnostic{
		File:  fpath,
		Line:  line,
		Col:   col,
		Cause: fmt.Sprintf("attributes on the route <body> tag (%s) are not applied to the layout's <body>", file.Body.Attrs),
		Fix:   "move them onto an element inside <body>, or add them to the layout's body tag",
	}
	return d.String(), true
}
