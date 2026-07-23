package parser

import (
	"fmt"
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/lexer"
)

func (p *Parser) parseHeadSection() (*ast.HeadSection, error) {
	content, err := p.readUntilClose("head")
	if err != nil {
		return nil, err
	}
	return &ast.HeadSection{Content: strings.TrimSpace(content)}, nil
}

func (p *Parser) readUntilClose(tag string) (string, error) {
	p.advance()
	var content strings.Builder

	for {
		tok := p.current()
		if tok.Type == lexer.TokenEOF {
			return "", fmt.Errorf("unclosed <%s> at position %d", tag, tok.Pos)
		}
		if tok.Type == lexer.TokenTagClose && tok.Tag == tag {
			p.advance()
			return content.String(), nil
		}
		if tok.Type == lexer.TokenText {
			content.WriteString(tok.Value)
		}
		p.advance()
	}
}
