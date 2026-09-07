package process

import (
	"fmt"

	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	jsinput "github.com/dreego-stack/dreego/internal/transpiler/js/js"
	jsoutput "github.com/dreego-stack/dreego/internal/transpiler/js/output"
	typescript "github.com/dreego-stack/dreego/internal/transpiler/js/ts"
)

func Client(file *ir.File) (jsoutput.Artifact, error) {
	if file.Client == nil {
		return jsoutput.Artifact{}, nil
	}
	return compile(file.Client.Language, file.Client.Code, file.Component, file.Server, file.SourcePath, lineAt(file.SourceContent, file.Client.Pos))
}

func Inline(node ir.TemplateNode, component *ir.ComponentDef, sections []ir.ServerSection) (jsoutput.Artifact, error) {
	return compile(node.Language, node.Content, component, sections, node.Source, lineAt(node.SourceText, node.Pos))
}

func compile(language, code string, component *ir.ComponentDef, sections []ir.ServerSection, sourcePath string, line int) (jsoutput.Artifact, error) {
	switch language {
	case "", "js":
		return jsinput.Process(code), nil
	case "ts":
		declarations, err := typescript.Declarations(component, sections)
		if err != nil {
			return jsoutput.Artifact{}, err
		}
		return typescript.Process(code, declarations, sourcePath, line)
	default:
		return jsoutput.Artifact{}, fmt.Errorf("unsupported client language %q", language)
	}
}

func lineAt(source string, position int) int {
	line := 1
	for i := 0; i < position && i < len(source); i++ {
		if source[i] == '\n' {
			line++
		}
	}
	return line
}
