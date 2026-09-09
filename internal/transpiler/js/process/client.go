package process

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/dreego-stack/dreego/internal/transpiler/codegen"
	"github.com/dreego-stack/dreego/internal/transpiler/ir"
	jsinput "github.com/dreego-stack/dreego/internal/transpiler/js/js"
	luainput "github.com/dreego-stack/dreego/internal/transpiler/js/lua"
	jsoutput "github.com/dreego-stack/dreego/internal/transpiler/js/output"
	typescript "github.com/dreego-stack/dreego/internal/transpiler/js/ts"
)

func Client(gen *codegen.State, file *ir.File) (jsoutput.Artifact, error) {
	if file.Client == nil {
		return jsoutput.Artifact{}, nil
	}
	return compile(gen, file.Client.Language, file.Client.Code, file.Component, file.Server, file.SourcePath, lineAt(file.SourceContent, file.Client.Pos))
}

func Inline(gen *codegen.State, node ir.TemplateNode, component *ir.ComponentDef, sections []ir.ServerSection) (jsoutput.Artifact, error) {
	return compile(gen, node.Language, node.Content, component, sections, node.Source, lineAt(node.SourceText, node.Pos))
}

func compile(gen *codegen.State, language, code string, component *ir.ComponentDef, sections []ir.ServerSection, sourcePath string, line int) (jsoutput.Artifact, error) {
	switch language {
	case "", "js":
		return jsinput.Process(code), nil
	case "ts":
		declarations, err := typescript.Declarations(component, sections)
		if err != nil {
			return jsoutput.Artifact{}, err
		}
		return typescript.Process(code, declarations, sourcePath, line)
	case "lua":
		artifact, err := luainput.CompileWithSource(code, sourceName(sourcePath), line)
		if err != nil {
			return jsoutput.Artifact{}, luaDiagnostic(err, sourceName(sourcePath), line)
		}
		features := luainput.FeatureList(artifact.Runtime)
		gen.AddLuaFeatures(features)
		result := jsoutput.Artifact{Code: artifact.Code}
		if len(features) > 0 {
			result.RuntimePath = "/_dreego/lua.js"
		}
		return result, nil
	default:
		return jsoutput.Artifact{}, fmt.Errorf("unsupported client language %q", language)
	}
}

var luaLocation = regexp.MustCompile(`^Lua (\d+):(\d+): (.*)$`)

func luaDiagnostic(err error, sourcePath string, lineOffset int) error {
	match := luaLocation.FindStringSubmatch(err.Error())
	if len(match) != 4 {
		return fmt.Errorf("Lua error in %s near line %d: %w", sourcePath, lineOffset, err)
	}
	line, lineErr := strconv.Atoi(match[1])
	if lineErr != nil {
		return fmt.Errorf("Lua error in %s near line %d: %w", sourcePath, lineOffset, err)
	}
	return fmt.Errorf("%s(%d,%s): Lua: %s", sourcePath, line+lineOffset-1, match[2], match[3])
}

func sourceName(path string) string {
	if path == "" {
		return "input.dreego"
	}
	return path
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
