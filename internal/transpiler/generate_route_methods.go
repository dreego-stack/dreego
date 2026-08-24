package transpiler

import (
	"fmt"
	"os"
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
	if file.Body != nil {
		setNodeSource(file.Body.Nodes, fpath, bodyOffset)
		setSourceText(file.Body.Nodes, raw)
		file.FormActions = scanFormActions(file.Body.Nodes)
		for _, diagnostic := range a11yDiagnostics(file.Body.Nodes) {
			fmt.Fprintf(os.Stderr, "warning: %s\n", diagnostic)
		}
	}
	return file, raw, nil
}
