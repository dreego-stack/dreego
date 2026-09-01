package lexer

import (
	"fmt"
	"strings"

	"github.com/dreego-stack/dreego/internal/transpiler/tokens"
)

func Lex(input string) ([]tokens.Token, error) {
	var toks []tokens.Token
	pos := 0
	sectionStack := []string{}
	sectionTags := map[string]bool{
		"body": true, "client": true, "head": true, "server": true, "style": true,
		"go": true, "script": true,
	}

	for pos < len(input) {
		inSection := len(sectionStack) > 0
		curSection := ""
		if inSection {
			curSection = sectionStack[len(sectionStack)-1]
		}

		if inSection && (curSection == "server" || curSection == "client" || curSection == "style" || curSection == "go" || curSection == "script") {
			closer := "</" + curSection + ">"
			closePos := strings.Index(input[pos:], closer)
			if closePos < 0 {
				toks = append(toks, tokens.Token{Type: tokens.TokenText, Value: input[pos:], Pos: pos})
				pos = len(input)
				break
			}
			if closePos > 0 {
				toks = append(toks, tokens.Token{Type: tokens.TokenText, Value: input[pos : pos+closePos], Pos: pos})
			}
			toks = append(toks, tokens.Token{Type: tokens.TokenTagClose, Tag: curSection, Pos: pos + closePos})
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
			if (!inSection || curSection == "body" || curSection == "head") && input[i] == '{' && isTemplateBrace(input[i:]) {
				nextPos = i
				nextCh = '{'
				break
			}
		}

		if nextPos < 0 {
			nextPos = len(input)
		}

		if nextPos > pos {
			toks = append(toks, tokens.Token{
				Type:  tokens.TokenText,
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
			toks = append(toks, tok)

			if tok.Type == tokens.TokenTagOpen && sectionTags[tok.Tag] {
				sectionStack = append(sectionStack, tok.Tag)
			}
			if tok.Type == tokens.TokenTagClose && sectionTags[tok.Tag] {
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
			toks = append(toks, tok)
		}
	}

	if len(sectionStack) > 0 {
		return nil, fmt.Errorf("unclosed tag <%s>", sectionStack[len(sectionStack)-1])
	}

	toks = append(toks, tokens.Token{Type: tokens.TokenEOF, Pos: pos})
	return toks, nil
}
