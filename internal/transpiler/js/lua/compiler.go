package lua

import "strings"

type Artifact struct {
	Code    string
	Runtime string
}

func Compile(source string) (Artifact, error) {
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
	return Artifact{Code: isolate(code), Runtime: strings.Join(emitter.sortedFeatures(), ",")}, nil
}

func isolate(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return ""
	}
	return "(() => {\n  " + strings.ReplaceAll(code, "\n", "\n  ") + "\n})();"
}
