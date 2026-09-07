package transpiler

import "github.com/dreego-stack/dreego/internal/transpiler/codegen"

type Generator = codegen.State

func NewGenerator() *Generator {
	return codegen.NewState()
}

type generator = codegen.State
