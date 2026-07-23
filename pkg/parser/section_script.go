package parser

import (
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
)

func (p *Parser) parseScriptSection() (*ast.ScriptSection, error) {
	content, err := p.readUntilClose("script")
	if err != nil {
		return nil, err
	}
	return &ast.ScriptSection{Code: strings.TrimSpace(content)}, nil
}
