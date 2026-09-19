package dreefile

import "github.com/dreego-stack/dreego/internal/dreefile/codegen"

type Generator = codegen.State

func NewGenerator() *Generator {
	return codegen.NewState()
}

type generator = codegen.State
