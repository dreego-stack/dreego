package js

import jsoutput "github.com/dreego-stack/dreego/internal/dreefile/jsoutput"

func Process(code string) jsoutput.Artifact {
	return jsoutput.Artifact{Code: code}
}
