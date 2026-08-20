package transpiler

import (
	"fmt"
	"strings"
)

func scanBrace(input string, pos *int) (Token, error) {
	start := *pos
	remaining := input[start:]

	if strings.HasPrefix(remaining, "{#else if ") {
		end := strings.IndexByte(remaining[10:], '}')
		if end < 0 {
			return Token{}, fmt.Errorf("unclosed {#else if at position %d", start)
		}
		*pos += 10 + end + 1
		return Token{Type: TokenElseIf, Value: strings.TrimSpace(remaining[10 : 10+end]), Pos: start}, nil
	}
	if strings.HasPrefix(remaining, "{#else}") {
		*pos += 7
		return Token{Type: TokenElse, Pos: start}, nil
	}
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
	if strings.HasPrefix(remaining, "{#each else}") {
		*pos += 12
		return Token{Type: TokenEachElse, Pos: start}, nil
	}
	if strings.HasPrefix(remaining, "{#each ") {
		end := strings.IndexByte(remaining[7:], '}')
		if end < 0 {
			return Token{}, fmt.Errorf("unclosed {#each at position %d", start)
		}
		*pos += 7 + end + 1
		return Token{Type: TokenEachOpen, Value: strings.TrimSpace(remaining[7 : 7+end]), Pos: start}, nil
	}

	if strings.HasPrefix(remaining, "{/slot}") {
		*pos += 7
		return Token{Type: TokenSlotClose, Pos: start}, nil
	}

	if strings.HasPrefix(remaining, "{#slot ") {
		end := strings.IndexByte(remaining[7:], '}')
		if end < 0 {
			return Token{}, fmt.Errorf("unclosed {#slot at position %d", start)
		}
		*pos += 7 + end + 1
		return Token{Type: TokenSlotOpen, Value: strings.TrimSpace(remaining[7 : 7+end]), Pos: start}, nil
	}

	if strings.HasPrefix(remaining, "{#slot}") {
		*pos += 7
		return Token{Type: TokenSlot, Pos: start}, nil
	}

	if strings.HasPrefix(remaining, "{#verbatim}") {
		*pos += 11
		closePos := strings.Index(input[*pos:], "{/verbatim}")
		if closePos < 0 {
			return Token{}, fmt.Errorf("unclosed {#verbatim} at position %d", start)
		}
		*pos += closePos + 11
		return Token{Type: TokenVerbatim, Value: input[start+11 : start+11+closePos], Pos: start}, nil
	}

	if !strings.HasPrefix(remaining, "{{") {
		return Token{}, fmt.Errorf("invalid template expression at position %d", start)
	}
	end := findExprEnd(remaining[2:])
	if end < 0 {
		return Token{}, fmt.Errorf("unclosed expression at position %d", start)
	}
	*pos += 2 + end + 2
	return Token{Type: TokenExpression, Value: strings.TrimSpace(remaining[2 : 2+end]), Pos: start}, nil
}

// findExprEnd returns the index of the closing "}}" in s, where s is the
// content after "{{", skipping "}}" inside string literals. Returns -1 when
// no closing delimiter exists.
func findExprEnd(s string) int {
	var quote byte
	escaped := false
	for i := 0; i+1 < len(s); i++ {
		c := s[i]
		if quote != 0 {
			if quote == '`' {
				if c == '`' {
					quote = 0
				}
				continue
			}
			if escaped {
				escaped = false
				continue
			}
			if c == '\\' {
				escaped = true
				continue
			}
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '"', '\'', '`':
			quote = c
		case '}':
			if s[i+1] == '}' {
				return i
			}
		}
	}
	return -1
}
