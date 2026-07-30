package core

import (
	"fmt"
	"strings"
)

func (p *Parser) parseGoSection() (*GoSection, error) {
	p.advance()
	var content strings.Builder

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed <go> at position %d", tok.Pos)
		}
		if tok.Type == TokenTagClose && tok.Tag == "go" {
			p.advance()
			return &GoSection{Code: strings.TrimSpace(content.String())}, nil
		}
		if tok.Type == TokenText {
			content.WriteString(tok.Value)
		}
		p.advance()
	}
}

func (p *Parser) parseNonDivSection(tag string) (string, error) {
	p.advance()
	var content strings.Builder
	depth := 1

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return "", fmt.Errorf("unclosed <%s> at position %d", tag, tok.Pos)
		}
		if tok.Type == TokenTagOpen {
			if tok.Tag == tag {
				depth++
			}
			content.WriteString("<" + tok.Tag)
			if tok.Attr != "" {
				content.WriteString(" " + tok.Attr)
			}
			content.WriteString(">")
		} else if tok.Type == TokenTagClose {
			if tok.Tag == tag {
				depth--
				if depth == 0 {
					p.advance()
					return content.String(), nil
				}
			}
			content.WriteString("</" + tok.Tag + ">")
		} else if tok.Type == TokenText {
			content.WriteString(tok.Value)
		} else if tok.Type == TokenExpression {
			content.WriteString("{" + tok.Value + "}")
		}
		p.advance()
	}
}
