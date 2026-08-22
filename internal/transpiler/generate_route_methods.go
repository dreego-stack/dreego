package transpiler

import (
	"fmt"
	"os"
)

func fileRegisteredMethods(file *File) []string {
	if len(file.FormActions) > 0 {
		action := file.FormActions[0]
		if findFormStruct(file.Go, action) != "" && findFormHandler(file.Go, action) {
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
	for _, section := range file.Go {
		add(section.Method)
	}
	for _, template := range file.Templates {
		add(template.Method)
	}
	if len(methods) == 0 {
		add("GET")
	}
	return methods
}

func parseRouteFile(gen *Generator, fpath string, data []byte) (*File, string, error) {
	raw := string(data)
	_, imports, body := ParseHeader(raw)
	tokens, err := Lex(body)
	if err != nil {
		return nil, "", fmt.Errorf("error lexing %s: %w", fpath, err)
	}
	p := NewParser(tokens)
	file, err := p.Parse()
	if err != nil {
		return nil, "", fmt.Errorf("error parsing %s: %w", fpath, err)
	}
	file.Imports = imports
	file.SourceContent = raw
	bodyOffset := len(raw) - len(body)
	if file.Template != nil {
		setNodeSource(file.Template.Nodes, fpath, bodyOffset)
		setSourceText(file.Template.Nodes, raw)
		file.FormActions = scanFormActions(file.Template.Nodes)
		for _, diagnostic := range a11yDiagnostics(file.Template.Nodes) {
			fmt.Fprintf(os.Stderr, "warning: %s\n", diagnostic)
		}
	}
	return file, raw, nil
}
