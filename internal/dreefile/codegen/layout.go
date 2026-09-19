package codegen

import "github.com/dreego-stack/dreego/internal/dreefile/ir"

type Layout struct {
	Rel    string
	Source string
	File   *ir.File
	Name   string
}
