package core

import (
	"fmt"
	"strings"

)

func (p *Parser) parseIfNodes() ([]TemplateNode, error) {
	var nodes []TemplateNode

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed {#if}")
		}
		if tok.Type == TokenIfClose {
			p.advance()
			return nodes, nil
		}
		if tok.Type == TokenTagClose {
			return nil, fmt.Errorf("unexpected </div> inside {#if}")
		}

		node, err := p.parseTemplateNode("if")
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
}

func (p *Parser) parseEachNodes() ([]TemplateNode, error) {
	var nodes []TemplateNode

	for {
		tok := p.current()
		if tok.Type == TokenEOF {
			return nil, fmt.Errorf("unclosed {#each}")
		}
		if tok.Type == TokenEachClose {
			p.advance()
			return nodes, nil
		}
		if tok.Type == TokenTagClose {
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
