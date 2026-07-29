package core

import (
	"fmt"
	"strings"

)

func (p *Parser) parseHeadSection() (*HeadSection, error) {
	content, err := p.readUntilClose("head")
	if err != nil {
		return nil, err
	}
	return &HeadSection{Content: strings.TrimSpace(content)}, nil
}

func (p *Parser) readUntilClose(tag string) (string, error) {
	p.advance()
	var content strings.Builder

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return "", fmt.Errorf("unclosed <%s> at position %d", tag, tok.Pos)
		}
		if tok.Type == TokenTagClose && tok.Tag == tag {
			p.advance()
			return content.String(), nil
		}
		if tok.Type == TokenText {
			content.WriteString(tok.Value)
		}
		p.advance()
	}
}
