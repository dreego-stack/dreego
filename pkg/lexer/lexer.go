package lexer

import (
	"fmt"
	"strings"
)

func Lex(input string) ([]Token, error) {
	var tokens []Token
	pos := 0
	stack := []string{}

	for pos < len(input) {
		insideDiv := len(stack) > 0 && stack[len(stack)-1] == "div"

		nextPos := -1
		nextCh := byte(0)

		for i := pos; i < len(input); i++ {
			if input[i] == '<' {
				nextPos = i
				nextCh = '<'
				break
			}
			if insideDiv && input[i] == '{' {
				nextPos = i
				nextCh = '{'
				break
			}
		}

		if nextPos < 0 {
			nextPos = len(input)
		}

		if nextPos > pos {
			tokens = append(tokens, Token{
				Type:  TokenText,
				Value: input[pos:nextPos],
				Pos:   pos,
			})
		}

		if nextPos >= len(input) {
			break
		}

		pos = nextPos

		if nextCh == '<' {
			tok := scanTag(input, &pos)
			tokens = append(tokens, tok)

			switch tok.Type {
			case TokenTagOpen:
				stack = append(stack, tok.Tag)
			case TokenTagClose:
				if len(stack) == 0 {
					return nil, fmt.Errorf("unexpected closing tag </%s> at position %d", tok.Tag, tok.Pos)
				}
				if stack[len(stack)-1] != tok.Tag {
					return nil, fmt.Errorf("mismatched closing tag </%s>, expected </%s> at position %d",
						tok.Tag, stack[len(stack)-1], tok.Pos)
				}
				stack = stack[:len(stack)-1]
			}
		} else if nextCh == '{' {
			tok, err := scanBrace(input, &pos)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
		}
	}

	if len(stack) > 0 {
		return nil, fmt.Errorf("unclosed tag <%s>", stack[len(stack)-1])
	}

	tokens = append(tokens, Token{Type: TokenEOF, Pos: pos})
	return tokens, nil
}

func scanTag(input string, pos *int) Token {
	start := *pos
	remaining := input[start:]

	tags := []string{"go", "div", "head", "script", "style"}

	for _, tag := range tags {
		closer := "</" + tag + ">"
		if strings.HasPrefix(remaining, closer) {
			*pos += len(closer)
			return Token{Type: TokenTagClose, Tag: tag, Pos: start}
		}
		opener := "<" + tag
		if strings.HasPrefix(remaining, opener) {
			end := strings.IndexByte(remaining, '>')
			if end < 0 {
				*pos += len(remaining)
				return Token{Type: TokenText, Value: remaining, Pos: start}
			}
			attrs := strings.TrimSpace(remaining[len(opener):end])
			attr := ""
			if strings.HasPrefix(attrs, "method=") {
				attr = strings.Trim(attrs[7:], "\"'")
			}
			*pos += end + 1
			return Token{Type: TokenTagOpen, Tag: tag, Attr: attr, Pos: start}
		}
	}

	end := strings.IndexByte(remaining, '>')
	if end >= 0 {
		*pos += end + 1
		return Token{Type: TokenText, Value: remaining[:end+1], Pos: start}
	}
	*pos += len(remaining)
	return Token{Type: TokenText, Value: remaining, Pos: start}
}
