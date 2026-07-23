package lexer

import (
	"fmt"
	"strings"
)

func scanBrace(input string, pos *int) (Token, error) {
	start := *pos
	remaining := input[start:]

	if strings.HasPrefix(remaining, "{/if}") {
		*pos += 5
		return Token{Type: TokenIfClose, Pos: start}, nil
	}
	if strings.HasPrefix(remaining, "{/each}") {
		*pos += 7
		return Token{Type: TokenEachClose, Pos: start}, nil
	}
	if strings.HasPrefix(remaining, "{#if ") {
		end := strings.IndexByte(remaining[5:], '}')
		if end < 0 {
			return Token{}, fmt.Errorf("unclosed {#if at position %d", start)
		}
		*pos += 5 + end + 1
		return Token{Type: TokenIfOpen, Value: strings.TrimSpace(remaining[5 : 5+end]), Pos: start}, nil
	}
	if strings.HasPrefix(remaining, "{#each ") {
		end := strings.IndexByte(remaining[7:], '}')
		if end < 0 {
			return Token{}, fmt.Errorf("unclosed {#each at position %d", start)
		}
		*pos += 7 + end + 1
		return Token{Type: TokenEachOpen, Value: strings.TrimSpace(remaining[7 : 7+end]), Pos: start}, nil
	}

	if strings.HasPrefix(remaining, "{#slot}") {
		*pos += 7
		return Token{Type: TokenSlot, Pos: start}, nil
	}

	end := strings.IndexByte(remaining[1:], '}')
	if end < 0 {
		return Token{}, fmt.Errorf("unclosed expression at position %d", start)
	}
	*pos += 1 + end + 1
	return Token{Type: TokenExpression, Value: strings.TrimSpace(remaining[1 : 1+end]), Pos: start}, nil
}
