package parser

import (
	"fmt"
	"strings"

	"codeberg.org/dreego/dreego/pkg/ast"
	"codeberg.org/dreego/dreego/pkg/lexer"
)

func (p *Parser) parseIfNodes() ([]ast.TemplateNode, error) {
	var nodes []ast.TemplateNode

	for {
		tok := p.current()
		if tok.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unclosed {#if}")
		}
		if tok.Type == lexer.TokenIfClose {
			p.advance()
			return nodes, nil
		}
		if tok.Type == lexer.TokenTagClose {
			return nil, fmt.Errorf("unexpected </div> inside {#if}")
		}

		node, err := p.parseTemplateNode("if")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

func (p *Parser) parseEachNodes() ([]ast.TemplateNode, error) {
	var nodes []ast.TemplateNode

	for {
		tok := p.current()
		if tok.Type == lexer.TokenEOF {
			return nil, fmt.Errorf("unclosed {#each}")
		}
		if tok.Type == lexer.TokenEachClose {
			p.advance()
			return nodes, nil
		}
		if tok.Type == lexer.TokenTagClose {
			return nil, fmt.Errorf("unexpected </div> inside {#each}")
		}

		node, err := p.parseTemplateNode("each")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

func parseEachClause(clause string) (items, item string, err error) {
	parts := strings.SplitN(clause, " as ", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("expected 'items as item', got %q", clause)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}
