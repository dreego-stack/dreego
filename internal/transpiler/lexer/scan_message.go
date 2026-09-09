package lexer

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func scanMessage(input string, pos *int) (tokens.Token, error) {
	start := *pos
	end := strings.Index(input[start+2:], "]]")
	if end < 0 {
		return tokens.Token{}, fmt.Errorf("unclosed message expression at position %d", start)
	}
	*pos = start + 2 + end + 2
	return tokens.Token{
		Type:  tokens.TokenMessage,
		Value: strings.TrimSpace(input[start+2 : start+2+end]),
		Pos:   start,
	}, nil
}
