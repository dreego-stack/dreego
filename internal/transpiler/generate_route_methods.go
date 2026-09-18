package transpiler

import (
	"fmt"
	"os"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
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
		ir.SetNodeSource(file.Body.Nodes, fpath, bodyOffset)
		ir.SetSourceText(file.Body.Nodes, raw)
		file.FormActions = scanFormActions(file.Body.Nodes)
		for _, diagnostic := range a11yDiagnostics(file.Body.Nodes) {
			fmt.Fprintf(os.Stderr, "warning: %s\n", diagnostic)
		}
	}
	return file, raw, nil
}
