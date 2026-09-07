package js

import jsoutput "github.com/dreego-stack/dreego/internal/transpiler/js/output"

func Process(code string) jsoutput.Artifact {
	return jsoutput.Artifact{Code: code}
}
