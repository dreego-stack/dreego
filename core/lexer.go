package core

import (
	"fmt"
	"strings"
)

func Lex(input string) ([]Token, error) {
	var tokens []Token
	pos := 0
	sectionStack := []string{}
	sectionTags := map[string]bool{"go": true, "head": true, "script": true, "style": true}

	for pos < len(input) {
		inSection := len(sectionStack) > 0
		curSection := ""
		if inSection {
			curSection = sectionStack[len(sectionStack)-1]
		}

		if inSection && (curSection == "go" || curSection == "script" || curSection == "style") {
			closer := "</" + curSection + ">"
			closePos := strings.Index(input[pos:], closer)
			if closePos < 0 {
				tokens = append(tokens, Token{Type: TokenText, Value: input[pos:], Pos: pos})
				pos = len(input)
				break
			}
			if closePos > 0 {
				tokens = append(tokens, Token{Type: TokenText, Value: input[pos : pos+closePos], Pos: pos})
			}
			tokens = append(tokens, Token{Type: TokenTagClose, Tag: curSection, Pos: pos + closePos})
			pos += closePos + len(closer)
			sectionStack = sectionStack[:len(sectionStack)-1]
			continue
		}

		nextPos := -1
		nextCh := byte(0)

		for i := pos; i < len(input); i++ {
			if input[i] == '<' {
				nextPos = i
				nextCh = '<'
				break
			}
			if !inSection && input[i] == '{' && isTemplateBrace(input[i:]) {
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

			if tok.Type == TokenTagOpen && sectionTags[tok.Tag] {
				sectionStack = append(sectionStack, tok.Tag)
			}
			if tok.Type == TokenTagClose && sectionTags[tok.Tag] {
				if len(sectionStack) == 0 {
					return nil, fmt.Errorf("unexpected closing tag </%s> at position %d", tok.Tag, tok.Pos)
				}
				if sectionStack[len(sectionStack)-1] != tok.Tag {
					return nil, fmt.Errorf("mismatched closing tag </%s>, expected </%s> at position %d",
						tok.Tag, sectionStack[len(sectionStack)-1], tok.Pos)
				}
				sectionStack = sectionStack[:len(sectionStack)-1]
			}
		} else if nextCh == '{' {
			tok, err := scanBrace(input, &pos)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, tok)
		}
	}

	if len(sectionStack) > 0 {
		return nil, fmt.Errorf("unclosed tag <%s>", sectionStack[len(sectionStack)-1])
	}

	tokens = append(tokens, Token{Type: TokenEOF, Pos: pos})
	return tokens, nil
}

func scanTag(input string, pos *int) Token {
	start := *pos
	remaining := input[start:]

	if strings.HasPrefix(remaining, "<@") || strings.HasPrefix(remaining, "</@") {
		return scanComponentTag(input, pos)
	}

	knownTags := []string{"go", "div", "head", "script", "style"}

	for _, tag := range knownTags {
		closer := "</" + tag + ">"
		if strings.HasPrefix(remaining, closer) {
			*pos += len(closer)
			return Token{Type: TokenTagClose, Tag: tag, Pos: start}
		}
		opener := "<" + tag
		if strings.HasPrefix(remaining, opener) {
			next := byte(0)
			if len(remaining) > len(opener) {
				next = remaining[len(opener)]
			}
			if next != ' ' && next != '>' && next != '/' {
				continue
			}
			end := tagEnd(remaining)
			if end < 0 {
				*pos += len(remaining)
				return Token{Type: TokenText, Value: remaining, Pos: start}
			}
			attrs := strings.TrimSpace(remaining[len(opener):end])
			*pos += end + 1
			selfClose := strings.HasSuffix(attrs, "/")
			attrs = strings.TrimSpace(strings.TrimSuffix(attrs, "/"))
			return Token{Type: TokenTagOpen, Tag: tag, Attr: attrs, SelfClose: selfClose, Pos: start}
		}
	}

	if strings.HasPrefix(remaining, "</") {
		end := tagEnd(remaining)
		if end < 0 {
			*pos += len(remaining)
			return Token{Type: TokenText, Value: remaining, Pos: start}
		}
		tag := remaining[2:end]
		if idx := strings.IndexByte(tag, ' '); idx >= 0 {
			tag = tag[:idx]
		}
		*pos += end + 1
		return Token{Type: TokenTagClose, Tag: tag, Pos: start}
	}

	if remaining[0] == '<' {
		end := tagEnd(remaining)
		if end < 0 {
			*pos += len(remaining)
			return Token{Type: TokenText, Value: remaining, Pos: start}
		}
		body := remaining[1:end]
		tag := body
		if idx := strings.IndexByte(body, ' '); idx >= 0 {
			tag = body[:idx]
		}
		attrs := strings.TrimSpace(strings.TrimPrefix(body, tag))
		selfClose := strings.HasSuffix(attrs, "/")
		attrs = strings.TrimSpace(strings.TrimSuffix(attrs, "/"))
		*pos += end + 1
		return Token{Type: TokenTagOpen, Tag: strings.TrimSpace(tag), Attr: attrs, SelfClose: selfClose, Pos: start}
	}

	end := tagEnd(remaining)
	if end >= 0 {
		*pos += end + 1
		return Token{Type: TokenText, Value: remaining[:end+1], Pos: start}
	}
	*pos += len(remaining)
	return Token{Type: TokenText, Value: remaining, Pos: start}
}

func scanComponentTag(input string, pos *int) Token {
	start := *pos
	remaining := input[start:]

	prefix := "<@"
	if strings.HasPrefix(remaining, "</@") {
		prefix = "</@"
	}

	remaining = remaining[len(prefix):]

	end := tagEnd(remaining)
	if end < 0 {
		*pos += len(input) - start
		return Token{Type: TokenText, Value: input[start:], Pos: start}
	}

	body := strings.TrimSpace(remaining[:end])
	selfClose := strings.HasSuffix(body, "/")
	if selfClose {
		body = strings.TrimSuffix(body, "/")
		body = strings.TrimSpace(body)
	}

	parts := strings.Fields(body)
	tag := ""
	attrs := ""
	if len(parts) > 0 {
		tag = parts[0]
		attrs = strings.TrimSpace(strings.TrimPrefix(body, tag))
	}

	*pos = start + len(prefix) + end + 1

	if prefix == "</@" {
		return Token{Type: TokenComponentTagClose, Tag: tag, Pos: start}
	}
	if selfClose {
		return Token{Type: TokenComponentSelfClose, Tag: tag, Attr: attrs, Pos: start}
	}
	return Token{Type: TokenComponentTagOpen, Tag: tag, Attr: attrs, Pos: start}
}
