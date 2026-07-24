package parser

import (
	"fmt"
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/lexer"
)

func (p *Parser) parseGoSection() (*ast.GoSection, error) {
	p.advance()
	var content strings.Builder

	for {
		tok := p.current()
		if tok.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unclosed <go> at position %d", tok.Pos)
		}
		if tok.Type == lexer.TokenTagClose && tok.Tag == "go" {
			p.advance()
			return &ast.GoSection{Code: strings.TrimSpace(content.String())}, nil
		}
		if tok.Type == lexer.TokenText {
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
		if tok.Type == lexer.TokenEOF {
			return "", fmt.Errorf("unclosed <%s> at position %d", tag, tok.Pos)
		}
		if tok.Type == lexer.TokenTagOpen {
			if tok.Tag == tag {
				depth++
			}
			content.WriteString("<" + tok.Tag)
			if tok.Attr != "" {
				content.WriteString(" " + tok.Attr)
			}
			content.WriteString(">")
		} else if tok.Type == lexer.TokenTagClose {
			if tok.Tag == tag {
				depth--
				if depth == 0 {
					p.advance()
					return content.String(), nil
				}
			}
			content.WriteString("</" + tok.Tag + ">")
		} else if tok.Type == lexer.TokenText {
			content.WriteString(tok.Value)
		} else if tok.Type == lexer.TokenExpression {
			content.WriteString("{" + tok.Value + "}")
		}
		p.advance()
	}
}
