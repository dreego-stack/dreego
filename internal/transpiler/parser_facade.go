package transpiler

import parserpkg "github.com/dreego-stack/dreego/internal/transpiler/parser"

type Parser = parserpkg.Parser

func NewParser(tokens []Token) *Parser {
	return parserpkg.NewParser(tokens)
}
