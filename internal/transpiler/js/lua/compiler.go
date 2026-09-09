package lua

import (
	"net/url"
	"strconv"
	"strings"
)

type Artifact struct {
	Code    string
	Runtime string
}

func Compile(source string) (Artifact, error) {
	return CompileWithSource(source, "", 1)
}

func CompileWithSource(source, sourcePath string, sourceLine int) (Artifact, error) {
	tokens, err := lex(source)
	if err != nil {
		return Artifact{}, err
	}
	program, err := parse(tokens)
	if err != nil {
		return Artifact{}, err
	}
	emitter := newEmitter()
	code, err := emitter.program(program)
	if err != nil {
		return Artifact{}, err
	}
	return Artifact{Code: isolate(code, sourcePath, sourceLine), Runtime: strings.Join(emitter.sortedFeatures(), ",")}, nil
}

func isolate(code, sourcePath string, sourceLine int) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	if sourcePath == "" {
		return "(() => {\n  " + strings.ReplaceAll(code, "\n", "\n  ") + "\n})();"
	}
	cleanPath := strings.ReplaceAll(strings.ReplaceAll(sourcePath, "\r", ""), "\n", "")
	cleanPath = strings.TrimPrefix(cleanPath, "./")
	location := (&url.URL{Scheme: "dreego", Path: "/" + strings.TrimPrefix(cleanPath, "/")}).String()
	if sourceLine < 1 {
		sourceLine = 1
	}
	prefix := strings.Repeat("\n", sourceLine-1)
	return prefix + "(() => { /* source line " + strconv.Itoa(sourceLine) + " */ " + strings.ReplaceAll(code, "\n", "\n  ") + "\n})();\n//# sourceURL=" + location
}
