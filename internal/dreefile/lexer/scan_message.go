package lexer

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/dreefile/dreecode"
	"github.com/dreego-stack/dreego/internal/dreefile/tokens"
)

func scanMessage(input string, pos *int) (tokens.Token, error) {
	start := *pos
	relativeEnd := dreecode.FindMessageEnd(input[start+2:])
	if relativeEnd < 0 {
		return tokens.Token{}, fmt.Errorf("unclosed message expression at position %d", start)
	}
	end := start + 2 + relativeEnd
	*pos = end + 2
	return tokens.Token{
		Type:  tokens.TokenMessage,
		Value: strings.TrimSpace(input[start+2 : end]),
		Pos:   start,
	}, nil
}
