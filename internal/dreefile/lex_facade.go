package dreefile

import lex "github.com/dreego-stack/dreego/internal/dreefile/lexer"

func Lex(input string) ([]Token, error) {
	return lex.Lex(input)
}

func ParseHeader(input string) (*ComponentDef, []Import, string) {
	return lex.ParseHeader(input)
}

func ParseFileHeader(input string) (*FileHeader, string) {
	return lex.ParseFileHeader(input)
}

func ParseFileHeaderStrict(input string) (*FileHeader, string, error) {
	return lex.ParseFileHeaderStrict(input)
}
