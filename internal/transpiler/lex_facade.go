package transpiler

import lex "github.com/dreego-stack/dreego/internal/transpiler/lexer"

func Lex(input string) ([]Token, error) {
	return lex.Lex(input)
}

func ParseHeader(input string) (*ComponentDef, []Import, string) {
	return lex.ParseHeader(input)
}
