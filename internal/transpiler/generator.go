package transpiler

import "github.com/dreego-stack/dreego/internal/transpiler/ir"

type Generator = ir.Generator

func NewGenerator() *Generator {
	return ir.NewGenerator()
}

type generator = ir.Generator
