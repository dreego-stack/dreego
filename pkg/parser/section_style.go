package parser

import (
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
)

func (p *Parser) parseStyleSection() (*ast.StyleSection, error) {
	content, err := p.readUntilClose("style")
	if err != nil {
		return nil, err
	}
	return &ast.StyleSection{Code: strings.TrimSpace(content)}, nil
}
